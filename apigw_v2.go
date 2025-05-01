package lambdahttp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// newHTTPRequestFromAPIGWv2 creates a new http.Request from an API Gateway v2 event
func newHTTPRequestFromAPIGWv2(req *events.APIGatewayV2HTTPRequest) (*http.Request, error) {
	requestContext := req.RequestContext

	// Get the HTTP method from the request context
	httpMethod := requestContext.HTTP.Method

	// Get the raw path and query string
	rawPath := req.RawPath
	rawQueryString := req.RawQueryString

	// Construct the full URL
	fullURL := rawPath
	if rawQueryString != "" {
		fullURL = fmt.Sprintf("%s?%s", rawPath, rawQueryString)
	}

	// Get the request body if present
	var body io.Reader
	if req.Body != "" {
		if req.IsBase64Encoded {
			decodedBody, err := base64.StdEncoding.DecodeString(req.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to decode base64 body: %v", err)
			}
			body = bytes.NewReader(decodedBody)
		} else {
			body = bytes.NewReader([]byte(req.Body))
		}
	}

	// Create the HTTP request
	httpReq, err := http.NewRequest(httpMethod, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set the headers
	for key, value := range req.Headers {
		for _, v := range strings.Split(value, ",") {
			httpReq.Header.Add(key, v)
		}
	}

	return httpReq, nil
}
