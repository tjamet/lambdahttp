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
)

// newHTTPRequestFromAPIGWv1 creates a new http.Request from an API Gateway v1 event
func newHTTPRequestFromAPIGWv1(req map[string]interface{}) (*http.Request, error) {
	method := req["httpMethod"].(string)
	path := req["path"].(string)

	// Build query string
	queryParams := ""
	if qp, ok := req["queryStringParameters"].(map[string]interface{}); ok && len(qp) > 0 {
		values := url.Values{}
		for k, v := range qp {
			values.Add(k, v.(string))
		}
		queryParams = "?" + values.Encode()
	}

	// Create request URL
	reqURL := fmt.Sprintf("https://%s%s%s",
		req["headers"].(map[string]interface{})["Host"].(string),
		path,
		queryParams)

	var body io.Reader
	if b64Body, ok := req["body"].(string); ok && req["isBase64Encoded"].(bool) {
		decodedBody, err := base64.StdEncoding.DecodeString(b64Body)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 body: %v", err)
		}
		body = bytes.NewReader(decodedBody)
	}

	httpReq, err := http.NewRequest(method, reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	multiValueHeaders := req["multiValueHeaders"].(map[string]interface{})
	headers := req["headers"].(map[string]interface{})
	for k, v := range headers {
		if _, ok := multiValueHeaders[k]; ok {
			continue
		}
		httpReq.Header.Add(k, v.(string))
	}
	for k, v := range multiValueHeaders {
		for _, v := range v.([]interface{}) {
			httpReq.Header.Add(k, v.(string))
		}
	}

	// Set protocol version
	if proto, ok := req["requestContext"].(map[string]interface{})["protocol"].(string); ok {
		protoComponents := strings.Split(proto, "/")
		if len(protoComponents) == 2 {
			httpReq.Proto = proto
			major, _ := strconv.Atoi(protoComponents[0])
			minor, _ := strconv.Atoi(protoComponents[1])
			httpReq.ProtoMajor = major
			httpReq.ProtoMinor = minor
		}
	}

	return httpReq, nil
}
