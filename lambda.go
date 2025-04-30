// Package lambdahttp provides functionality to run standard Go HTTP handlers in AWS Lambda
// with automatic integration for various AWS services like API Gateway, ALB, and EventBridge.
package lambdahttp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
)

type lambdaResponseWriter struct {
	statusCode int
	headers    http.Header
	body       []byte
}

var _ http.ResponseWriter = &lambdaResponseWriter{}

func (w *lambdaResponseWriter) Header() http.Header {
	if w == nil {
		return nil
	}
	if w.headers == nil {
		w.headers = make(http.Header)
	}
	return w.headers
}

func (w *lambdaResponseWriter) Write(b []byte) (int, error) {
	if w == nil {
		return 0, fmt.Errorf("lambdaResponseWriter is nil")
	}
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	w.body = append(w.body, b...)
	return len(b), nil
}

func (w *lambdaResponseWriter) WriteHeader(statusCode int) {
	if w == nil {
		return
	}
	if w.statusCode != 0 {
		fmt.Println("WriteHeader called with status code", statusCode, "but already set to", w.statusCode, "ignoring")
		return
	}
	w.statusCode = statusCode
}

func (w *lambdaResponseWriter) GetStatusCode() int {
	if w == nil {
		return 0
	}
	return w.statusCode
}

// Response represents a Lambda response that can be returned by the handler
type Response interface {
	GetStatusCode() int
	GetHeaders() http.Header
	GetBody() string
}

// IntegrationType represents the type of AWS integration
type IntegrationType int

// Integration types supported by the package
const (
	// IntegrationTypeUnknown represents an unknown or unsupported integration type
	IntegrationTypeUnknown IntegrationType = iota
	// IntegrationTypeAPIGWv1 represents API Gateway v1 integration
	IntegrationTypeAPIGWv1
	// IntegrationTypeAPIGWv2 represents API Gateway v2 integration
	IntegrationTypeAPIGWv2
	// IntegrationTypeALB represents Application Load Balancer integration
	IntegrationTypeALB
)

func (i IntegrationType) String() string {
	return []string{"APIGWv1", "APIGWv2", "ALB"}[i]
}

func getIntegrationType(req map[string]interface{}) (IntegrationType, error) {
	version, ok := req["version"]
	if !ok {
		return IntegrationTypeALB, nil
	}
	versionStr, ok := version.(string)
	if !ok {
		return 0, fmt.Errorf("version is not a string")
	}
	switch versionStr {
	case "1.0":
		return IntegrationTypeAPIGWv1, nil
	case "2.0":
		return IntegrationTypeAPIGWv2, nil
	default:
		return IntegrationTypeUnknown, fmt.Errorf("unknown version: %s", versionStr)
	}
}

func requestBuilder(version IntegrationType, req map[string]interface{}) (*http.Request, error) {
	switch version {
	case IntegrationTypeAPIGWv1:
		return newHTTPRequestFromAPIGWv1(req)
	case IntegrationTypeAPIGWv2:
		return newHTTPRequestFromAPIGWv2(req)
	case IntegrationTypeALB:
		return newHTTPRequestFromALB(req)
	default:
		return nil, fmt.Errorf("unknown integration type: %d", version)
	}
}

func buildResponse(version IntegrationType, w *lambdaResponseWriter) (Response, error) {
	switch version {
	case IntegrationTypeAPIGWv1:
		return buildAPIGatewayResponse(w), nil
	case IntegrationTypeAPIGWv2:
		return buildAPIGatewayResponse(w), nil
	case IntegrationTypeALB:
		return buildALBResponse(w), nil
	default:
		return nil, fmt.Errorf("unknown integration type: %d", version)
	}
}

// LambdaRequestContextKey is used to store and retrieve Lambda-specific request context values
type LambdaRequestContextKey string

const (
	// LambdaRequestContextKeyOriginalRequest is the context key for storing the original Lambda request
	LambdaRequestContextKeyOriginalRequest LambdaRequestContextKey = "originalRequest"
	// LambdaRequestContextKeyIntegrationType is the context key for storing the integration type
	LambdaRequestContextKeyIntegrationType LambdaRequestContextKey = "integrationType"
)

// GetOriginalRequest retrieves the original Lambda request from the HTTP request context.
// Returns nil if no original request is found in the context.
func GetOriginalRequest(httpRequest *http.Request) map[string]interface{} {
	originalRequest, ok := httpRequest.Context().Value(LambdaRequestContextKeyOriginalRequest).(map[string]interface{})
	if !ok {
		return nil
	}
	return originalRequest
}

// GetIntegrationType determines the AWS integration type from the HTTP request context.
// Returns IntegrationTypeUnknown if the integration type cannot be determined.
func GetIntegrationType(httpRequest *http.Request) IntegrationType {
	integrationType, ok := httpRequest.Context().Value(LambdaRequestContextKeyIntegrationType).(IntegrationType)
	if !ok {
		return IntegrationTypeALB
	}
	return integrationType
}

// NewAWSLambdaHTTPHandler creates a new AWS Lambda handler that can be used to handle HTTP requests
// It takes an http.Handler and returns a function that can be used as a Lambda handler
// The handler will automatically detect the api gateway version and build the request accordingly
// It also handles the case where the request is coming from an ALB
//
// Supported integration types:
// - APIGWv1: API Gateway v1
// - APIGWv2: API Gateway v2 (also used by eventBridge and lambda function url)
// - ALB: Application Load Balancer
func NewAWSLambdaHTTPHandler(h http.Handler) func(context.Context, map[string]interface{}) (Response, error) {
	return func(ctx context.Context, req map[string]interface{}) (Response, error) {
		integrationType, err := getIntegrationType(req)
		if err != nil {
			return nil, err
		}
		httpRequest, err := requestBuilder(integrationType, req)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, LambdaRequestContextKeyOriginalRequest, req)
		ctx = context.WithValue(ctx, LambdaRequestContextKeyIntegrationType, integrationType)

		w := &lambdaResponseWriter{}
		h.ServeHTTP(w, httpRequest.WithContext(ctx))

		response, err := buildResponse(integrationType, w)
		if err != nil {
			return nil, err
		}

		return response, nil
	}
}

// StartLambdaHandler starts the lambda handler
// It takes an http.Handler and runs it within a lambda
// The handler will automatically detect the api gateway version and build the request accordingly
// It also handles the case where the request is coming from an ALB
//
// Supported integration types:
// - APIGWv1: API Gateway v1
// - APIGWv2: API Gateway v2 (also used by eventBridge and lambda function url)
// - ALB: Application Load Balancer
func StartLambdaHandler(h http.Handler) {
	lambda.Start(NewAWSLambdaHTTPHandler(h))
}
