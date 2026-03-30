package lambdahttp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// newHTTPRequestFromAPIGWv1 creates a new http.Request from an API Gateway v1 event
func newHTTPRequestFromAPIGWv1(req *events.APIGatewayProxyRequest) (*http.Request, error) {
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

	// Create request URL
	reqURL := fmt.Sprintf("https://%s%s%s",
		req.Headers["Host"],
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

	for k, v := range req.Headers {
		if _, ok := req.MultiValueHeaders[k]; ok {
			continue
		}
		httpReq.Header.Add(k, v)
	}
	for k, v := range req.MultiValueHeaders {
		for _, v := range v {
			httpReq.Header.Add(k, v)
		}
	}

	protoComponents := strings.Split(req.RequestContext.Protocol, "/")
	if len(protoComponents) == 2 {
		httpReq.Proto = req.RequestContext.Protocol
		major, _ := strconv.Atoi(protoComponents[0])
		minor, _ := strconv.Atoi(protoComponents[1])
		httpReq.ProtoMajor = major
		httpReq.ProtoMinor = minor
	}

	return httpReq, nil
}
