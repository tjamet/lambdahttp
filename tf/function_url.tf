resource "aws_lambda_function_url" "rest-api" {
  function_name      = aws_lambda_function.rest-api.function_name
  authorization_type = "NONE"

  cors {
    allow_credentials = true
    allow_headers     = ["*"]
    allow_methods     = ["*"]
    allow_origins     = ["*"]
    expose_headers    = ["*"]
    max_age           = 300
  }
} 