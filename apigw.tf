resource "aws_apigatewayv2_api" "blossom" {
  name          = "${local.prefix}-gw"
  protocol_type = "HTTP"

  cors_configuration = {
    allow_headers = ["content-type", "x-amz-date", "authorization", "x-api-key", "x-amz-security-token", "x-amz-user-agent"]
    allow_methods = ["*"]
    allow_origins = ["*"]
  }

  integrations = {
    "POST /invoke" = {
      integration_type = "HTTP_PROXY"
      integration_uri  = "https://example.com"
    }

    "POST /query" = {
      integration_type = "HTTP_PROXY"
      integration_uri  = "https://example.com"
    }
  }
}

resource "aws_apigatewayv2_stage" "blossom" {
  name = "${local.prefix}-gw-stage"

  api_id      = aws_apigatewayv2_api.blossom
  auto_deploy = true
}
