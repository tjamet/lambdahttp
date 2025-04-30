resource "aws_apigatewayv2_api" "rest-api" {
  name          = "rest-api"
  protocol_type = "HTTP"
  
  cors_configuration {
    allow_headers     = ["*"]
    allow_methods     = ["*"]
    allow_origins     = ["*"]
    expose_headers    = ["*"]
    max_age          = 300
  }
}

locals {
  versions = {
    "1" = "1.0"
    "2" = "2.0"
  }
}

resource "aws_apigatewayv2_integration" "rest-api" {
  for_each = local.versions
  api_id             = aws_apigatewayv2_api.rest-api.id
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.rest-api.invoke_arn
  payload_format_version = each.value
}

resource "aws_apigatewayv2_route" "all" {
  for_each = local.versions
  api_id    = aws_apigatewayv2_api.rest-api.id
  route_key = "ANY /v${each.key}/{path+}"
  target    = "integrations/${aws_apigatewayv2_integration.rest-api[each.key].id}"
}

resource "aws_apigatewayv2_stage" "rest-api" {
  api_id = aws_apigatewayv2_api.rest-api.id
  name   = "$default"
  auto_deploy = true

  default_route_settings {
    throttling_burst_limit = 1000
    throttling_rate_limit  = 500
    detailed_metrics_enabled = true
  }
}

resource "aws_lambda_permission" "apigw" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.rest-api.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.rest-api.execution_arn}/*/*/{path+}"
} 