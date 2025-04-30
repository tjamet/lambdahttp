resource "aws_cloudwatch_event_connection" "rest-api-connection" {
  name               = "rest-api-connection"
  description        = "Connection for rest-api HTTP endpoint"
  authorization_type = "BASIC"

  auth_parameters {
    basic {
      username = "dummy"
      password = "dummy"
    }
  }
}

resource "aws_cloudwatch_event_api_destination" "rest-api-http" {
  name                             = "rest-api-http"
  description                      = "HTTP endpoint for rest-api"
  invocation_endpoint              = aws_lambda_function_url.rest-api.function_url
  http_method                      = "POST"
  invocation_rate_limit_per_second = 20
  connection_arn                   = aws_cloudwatch_event_connection.rest-api-connection.arn
} 