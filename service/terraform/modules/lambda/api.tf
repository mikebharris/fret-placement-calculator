terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.99.1"
    }
  }
}
resource "aws_iam_role" "api_iam_role" {
  name                  = "${var.project}-api-iam-role"
  force_detach_policies = true
  assume_role_policy    = jsonencode({
    Version   = "2012-10-17"
    Statement = [
      {
        Action    = "sts:AssumeRole"
        Effect    = "Allow"
        Sid       = ""
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_cloudwatch_log_group" "api_cloudwatch_log_group" {
  name              = "/aws/lambda/${aws_lambda_function.api_lambda_function.function_name}"
  retention_in_days = 1
}

data "aws_iam_policy_document" "api_iam_policy_document" {
  statement {
    effect  = "Allow"
    actions = [
      "ec2:DescribeNetworkInterfaces",
    ]
  }
}

resource "aws_iam_role_policy_attachment" "api_policy_attachment" {
  role       = aws_iam_role.api_iam_role.name
  policy_arn = aws_iam_policy.api_iam_policy.arn
}

resource "aws_iam_role_policy_attachment" "api_policy_attachment_execution" {
  role       = aws_iam_role.api_iam_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_policy" "api_iam_policy" {
  name   = "${var.project}-api-iam-policy"
  path   = "/"
  policy = data.aws_iam_policy_document.api_iam_policy_document.json
}

data "archive_file" "api_lambda_function_distribution" {
  source_file = "../fret-placement-calculator-api/main"
  output_path = "../fret-placement-calculator-api/${var.project}-api.zip"
  type        = "zip"
}

resource "aws_s3_object" "api_lambda_function_distribution_bucket_object" {
  bucket = var.distribution_bucket
  key    = "lambdas/${var.project}-api/${var.project}-api.zip"
  source = data.archive_file.api_lambda_function_distribution.output_path
  etag   = filemd5(data.archive_file.api_lambda_function_distribution.output_path)
}

resource "aws_lambda_function" "api_lambda_function" {
  function_name    = "${var.project}-api"
  role             = aws_iam_role.api_iam_role.arn
  handler          = "main"
  runtime          = "go1.x"
  s3_bucket        = aws_s3_object.api_lambda_function_distribution_bucket_object.bucket
  s3_key           = aws_s3_object.api_lambda_function_distribution_bucket_object.key
  source_code_hash = data.archive_file.api_lambda_function_distribution.output_md5
  timeout          = 60
  memory_size      = 256

  environment {
    variables = {
      AWS_REGION           = var.region
    }
  }

  tags = {
    Name          = "${var.project}.lambda.api"
    Contact       = var.contact
    Project       = var.project
    Orchestration = var.orchestration
    Description   = "API for calculating and returning fret placements"
  }
}

resource "aws_lambda_function_url" "api_lambda_function_url" {
  authorization_type = "NONE"
  function_name      = aws_lambda_function.api_lambda_function.function_name
}


