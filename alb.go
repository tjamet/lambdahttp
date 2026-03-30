package lambdahttp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type albResponse struct {
	events.ALBTargetGroupResponse
}

var _ Response = &albResponse{}

func (r albResponse) GetStatusCode() int {
	return r.StatusCode
}

func (r albResponse) GetHeaders() http.Header {
	headers := make(http.Header)
	for k, v := range r.Headers {
		headers.Set(k, v)
	}
	for k, v := range r.MultiValueHeaders {
		headers[k] = append(headers[k], v...)
	}
	return headers
}

func (r albResponse) GetBody() string {
	if r.IsBase64Encoded {
		decodedBody, err := base64.StdEncoding.DecodeString(r.Body)
		if err != nil {
			return ""
		}
		return string(decodedBody)
	}
	return r.Body
}

// newHTTPRequestFromALB creates a new http.Request from an ALB event
func newHTTPRequestFromALB(req *events.ALBTargetGroupRequest) (*http.Request, error) {
	method := req.HTTPMethod
	path := req.Path

	// Build query string
	queryParams := ""
	if len(req.QueryStringParameters) > 0 {
		values := url.Values{}
		for k, v := range req.QueryStringParameters {
			values.Add(k, v)
		}
		queryParams = "?" + values.Encode()
	}

	proto := req.Headers["x-forwarded-proto"]
	if proto == "" {
		proto = "http"
	}

	// Create request URL
	reqURL := fmt.Sprintf("%s://%s%s%s",
		proto,
		req.Headers["host"],
		path,
		queryParams)

	var body io.Reader
	if req.Body != "" {
		if req.IsBase64Encoded {
			decodedBody, err := base64.StdEncoding.DecodeString(req.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to decode base64 body: %v", err)
			}
			body = bytes.NewReader(decodedBody)
		} else {
			body = strings.NewReader(req.Body)
		}
	}

	httpReq, err := http.NewRequest(method, reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add headers
	for k, v := range req.Headers {
		for _, v := range strings.Split(v, ",") {
			httpReq.Header.Add(k, v)
		}
	}

	return httpReq, nil
}

func buildALBResponse(w *lambdaResponseWriter) *albResponse {
	response := &albResponse{
		ALBTargetGroupResponse: events.ALBTargetGroupResponse{
			Headers:           make(map[string]string),
			MultiValueHeaders: make(map[string][]string),
		},
	}
	response.StatusCode = w.GetStatusCode()
	for k, v := range w.headers {
		if len(v) == 1 {
			response.Headers[k] = v[0]
		} else {
			response.MultiValueHeaders[k] = v
		}
	}
	encodedBody := base64.StdEncoding.EncodeToString(w.body)
	response.Body = encodedBody
	response.IsBase64Encoded = true
	return response
}
