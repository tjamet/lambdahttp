locals {
  lambda_compile_dir   = "${path.cwd}/.lambdas/rest-api"
  lambda_binary        = "${path.cwd}/.lambdas/rest-api/bootstrap"
  go_arch_mapping = {
    "x86_64" = "amd64"
  }
  architecture = "arm64"
}

resource "null_resource" "build-go-binary" {
  triggers = {
    always_run = "${timestamp()}"
  }
  provisioner "local-exec" {
    command = "mkdir -p ${local.lambda_compile_dir} && GOOS=linux GOARCH=${lookup(local.go_arch_mapping, local.architecture, local.architecture)} CGO_ENABLED=0 GOFLAGS=-trimpath go build -mod=readonly -ldflags='-s -w' -o ${local.lambda_binary} ${path.cwd}/lambdas/rest-api"
  }
}

data "archive_file" "build-go-lambda" {
  depends_on = [null_resource.build-go-binary]

  type        = "zip"
  source_file = local.lambda_binary
  output_path = "${local.lambda_compile_dir}/rest-api-${null_resource.build-go-binary.id}.zip"
}

resource "aws_lambda_function" "rest-api" {
  function_name = "rest-api"
  handler       = "bootstrap"
  filename      = data.archive_file.build-go-lambda.output_path
  runtime       = "provided.al2023"
  role          = aws_iam_role.lambda-exec-role.arn
  architectures = [local.architecture]
  depends_on    = [data.archive_file.build-go-lambda]
}

resource "aws_iam_role" "lambda-exec-role" {
  name = "lambda-exec-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Service = "lambda.amazonaws.com"
        },
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda-exec-role" {
  role       = aws_iam_role.lambda-exec-role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_role_policy_attachment" "lambda-alb" {
  role       = aws_iam_role.lambda-exec-role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
} 