output "apigw_url" {
  description = "API Gateway HTTP API URL (managed API with features like throttling and metrics)"
  value       = trimsuffix(aws_apigatewayv2_api.rest-api.api_endpoint, "/")
}


output "lambda_function_url" {
  description = "Direct HTTPS endpoint for the Lambda function (simplest way to expose the API)"
  value       = trimsuffix(aws_lambda_function_url.rest-api.function_url, "/")
}

output "alb_url" {
  description = "Application Load Balancer URL (good for complex routing and multiple functions)"
  value       = trimsuffix("http://${aws_lb.rest-api-alb.dns_name}", "/")
}

# output "app_runner_url" {
#   description = "AWS App Runner URL (fully managed service, good for containerized applications)"
#   value       = aws_apprunner_service.rest-api-app.service_url
# }

# output "cloudfront_url" {
#   description = "CloudFront distribution URL (global CDN, edge caching, DDoS protection)"
#   value       = aws_cloudfront_distribution.rest-api-cf.domain_name
# }

output "eventbridge_endpoint" {
  description = "EventBridge HTTP endpoint (good for event-driven architectures)"
  value       = trimsuffix(aws_cloudwatch_event_api_destination.rest-api-http.invocation_endpoint, "/")
}

output "all_endpoints" {
  description = "All available endpoints for the REST API"
  value = {
    lambda_function_url = trimsuffix(aws_lambda_function_url.rest-api.function_url, "/")
    apigw_url          = trimsuffix(aws_apigatewayv2_api.rest-api.api_endpoint, "/")
    alb_url            = trimsuffix("http://${aws_lb.rest-api-alb.dns_name}", "/")
    //app_runner_url     = aws_apprunner_service.rest-api-app.service_url
    //cloudfront_url     = aws_cloudfront_distribution.rest-api-cf.domain_name
    eventbridge_endpoint = trimsuffix(aws_cloudwatch_event_api_destination.rest-api-http.invocation_endpoint, "/")
  }
} 