package lambdahttp

import (
	"encoding/base64"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// APIGatewayResponse represents the response format required by AWS API Gateway.
// It contains the status code, headers, and body of the HTTP response.
type APIGatewayResponse struct {
	events.APIGatewayProxyResponse
}

var _ Response = &APIGatewayResponse{}

// GetStatusCode returns the HTTP status code of the response
func (r APIGatewayResponse) GetStatusCode() int {
	return r.StatusCode
}

// GetHeaders returns the HTTP headers of the response
func (r APIGatewayResponse) GetHeaders() http.Header {
	headers := make(http.Header)
	for k, v := range r.Headers {
		headers.Set(k, v)
	}
	for k, v := range r.MultiValueHeaders {
		headers[k] = append(headers[k], v...)
	}
	return headers
}

// GetBody returns the response body as a string.
// If the body is base64 encoded, it will be decoded before returning.
func (r APIGatewayResponse) GetBody() string {
	if r.IsBase64Encoded {
		decodedBody, err := base64.StdEncoding.DecodeString(r.Body)
		if err != nil {
			return ""
		}
		return string(decodedBody)
	}
	return r.Body
}

func buildAPIGatewayResponse(w *lambdaResponseWriter) *APIGatewayResponse {
	response := &APIGatewayResponse{
		APIGatewayProxyResponse: events.APIGatewayProxyResponse{
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
