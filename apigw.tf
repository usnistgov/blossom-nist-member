module "apigateway-v2" {
  source  = "terraform-aws-modules/apigateway-v2/aws"
  version = "2.0.0"

  name          = "${local.prefix}-gw"
  protocol_type = "HTTP"

  cors_configuration = {
    allow_headers = ["content-type", "x-amz-date", "authorization", "x-api-key", "x-amz-security-token", "x-amz-user-agent"]
    allow_methods = ["*"]
    allow_origins = ["*"]
  }

  authorizers = {
    "cognito" = {
      name             = "${local.prefix}-identity-authorizer"
      authorizer_type  = "JWT"
      identity_sources = "$request.header.Authorization"
      issuer           = "https://${aws_cognito_user_pool.identity.endpoint}"
    }
  }

  integrations = {
    "POST /invoke" = {
      lambda_arn       = aws_lambda_function.auth
      integration_type = "HTTP_PROXY"
      authorizer_key   = "cognito"
      request_parameters = jsonencode({
        FORWARD_URL = "${module.vars.env.forward_url}/transaction/invoke"
      })
    }

    "POST /transaction/query" = {
      lambda_arn       = aws_lambda_function.auth
      integration_type = "HTTP_PROXY"
      authorizer_key   = "cognito"
      request_parameters = jsonencode({
        FORWARD_URL = "${module.vars.env.forward_url}/transaction/query"
      })
    }

    "$default" = {
      integration_type = "HTTP"
      integration_type = "${module.vars.env.forward_url}"
    }
  }

  tags = local.tags
}
