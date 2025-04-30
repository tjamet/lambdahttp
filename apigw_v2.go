package lambdahttp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// newHTTPRequestFromAPIGWv2 creates a new http.Request from an API Gateway v2 event
func newHTTPRequestFromAPIGWv2(req map[string]interface{}) (*http.Request, error) {
	requestContext := req["requestContext"].(map[string]interface{})

	// Get the HTTP method from the request context
	httpMethod := requestContext["http"].(map[string]interface{})["method"].(string)

	// Get the raw path and query string
	rawPath := req["rawPath"].(string)
	rawQueryString := req["rawQueryString"].(string)

	// Construct the full URL
	fullURL := rawPath
	if rawQueryString != "" {
		fullURL = fmt.Sprintf("%s?%s", rawPath, rawQueryString)
	}

	// Get the request body if present
	var body io.Reader
	if rawBody, ok := req["body"].(string); ok {
		if isBase64Encoded, ok := req["isBase64Encoded"].(bool); ok && isBase64Encoded {
			decodedBody, err := base64.StdEncoding.DecodeString(rawBody)
			if err != nil {
				return nil, fmt.Errorf("failed to decode base64 body: %v", err)
			}
			body = bytes.NewReader(decodedBody)
		} else {
			body = bytes.NewReader([]byte(rawBody))
		}
	}

	// Create the HTTP request
	httpReq, err := http.NewRequest(httpMethod, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set the headers
	if headers, ok := req["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			if strValue, ok := value.(string); ok {
				fmt.Println("Setting header:", key, strValue)
				for _, v := range strings.Split(strValue, ",") {
					httpReq.Header.Add(key, v)
				}
			}
		}
	}

	return httpReq, nil
}
