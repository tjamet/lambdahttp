# lambdahttp

[![Go Reference](https://pkg.go.dev/badge/github.com/tjamet/lambdahttp.svg)](https://pkg.go.dev/github.com/tjamet/lambdahttp)
[![Test](https://github.com/tjamet/lambdahttp/actions/workflows/ci.yml/badge.svg)](https://github.com/tjamet/lambdahttp/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/tjamet/lambdahttp)](https://goreportcard.com/report/github.com/tjamet/lambdahttp)

Run standard Go HTTP handlers in AWS Lambda with automatic integration for API Gateway, ALB, and EventBridge.

## Overview

This package allows you to run any standard Go `http.Handler` in AWS Lambda without modification. It automatically detects and handles requests from:
- API Gateway (v1 and v2)
- Application Load Balancer (ALB)
- EventBridge
- Lambda Function URLs

The main benefits are:
- **Low-Cost Migration**: Write your HTTP handlers using standard Go patterns and libraries, then deploy them to AWS Lambda without changing your code
- **Cloud Flexibility**: Easily switch between running your application on traditional servers and AWS Lambda without code changes
- **Future-Proof**: If you decide to move away from Lambda in the future, your code remains standard and portable

## Features

- 🔄 Use standard `http.Handler` interface
- 🔍 Automatic detection of AWS integration type
- 🌐 Support for multiple AWS services:
  - API Gateway v1
  - API Gateway v2
  - Application Load Balancer
  - EventBridge
  - Lambda Function URLs
- 📝 Preserves HTTP headers, query parameters, and body

## Installation

```bash
go get github.com/tjamet/lambdahttp
```

## Quick Start

Here's a simple example of how to use the package:

```go
package main

import (
    "net/http"
    "github.com/tjamet/lambdahttp"
)

func main() {
    // Create your HTTP handler as usual
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"message": "Hello from Lambda!"}`))
    })

    // Start the Lambda handler
    serverless.StartLambdaHandler(mux)
}
```

## Advanced Usage

### Custom Handler

If you need more control over the Lambda function, you can use `NewAWSLambdaHTTPHandler`:

```go
package main

import (
    "context"
    "net/http"
    "github.com/tjamet/lambdahttp"
)

func main() {
    handler := serverless.NewAWSLambdaHTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Access the original Lambda request if needed
        originalRequest := serverless.GetOriginalRequest(r)
        
        // Get the integration type
        integrationType := serverless.GetIntegrationType(r)
        
        // Your handler logic here
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"message": "Hello!"}`))
    }))

    lambda.Start(handler)
}
```

### Integration Types

The package automatically detects the following integration types:

- `IntegrationTypeAPIGWv1`: API Gateway REST API (v1)
- `IntegrationTypeAPIGWv2`: API Gateway HTTP API (v2), EventBridge, and Lambda Function URLs
- `IntegrationTypeALB`: Application Load Balancer

## AWS Service Configuration

### API Gateway

For API Gateway, you can use either REST API (v1) or HTTP API (v2). The package will automatically detect the version and handle the request appropriately.

### Application Load Balancer

When using an ALB trigger, make sure your target group protocol is set to HTTP.

### EventBridge

The package handles EventBridge events by converting them to HTTP requests. The request will be sent as a POST request to the handler.

### Lambda Function URLs

Function URLs are automatically detected and handled as API Gateway v2 requests.

### Infrastructure as Code Example

For a complete example of how to wire all these integrations together, check out our [Terraform module example](./tf). This module helps:

- Setting up Lambda functions with this package
- Configuring API Gateway integrations (v1 and v2)
- Setting up ALB targets
- Configuring EventBridge rules
- Managing Lambda Function URLs
- Handling permissions and IAM roles

## Contributing

Contributions are welcome! Please feel free to submit an Issue or a Pull Request.

## License

[Apache - 2.0] - See [LICENSE](LICENSE) file for details.
