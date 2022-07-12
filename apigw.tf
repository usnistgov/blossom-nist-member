resource "aws_apigatewayv2_api" "gw" {
  name          = "${local.prefix}-gw"
  protocol_type = "HTTP"

  # cors_configuration = {
  #   allow_headers = ["content-type", "x-amz-date", "authorization", "x-api-key", "x-amz-security-token", "x-amz-user-agent"]
  #   allow_methods = ["*"]
  #   allow_origins = ["*"]
  # }

  tags = local.tags
}

resource "aws_apigatewayv2_authorizer" "cognito" {
  api_id           = aws_apigatewayv2_api.gw.id
  name             = "${local.prefix}-identity-authorizer"
  authorizer_type  = "JWT"
  identity_sources = ["$request.header.Authorization"]

  jwt_configuration {
    # annoyingly at least one audience is required
    audience = ["transact"]
    issuer   = "https://${aws_cognito_user_pool.identity.endpoint}"
  }
}

# resource "aws_apigatewayv2_integration" "default" {
#   api_id           = aws_apigatewayv2_api.gw.id
#   integration_type = "HTTP_PROXY"
#   integration_uri  = module.vars.env.forward_url
# }

# resource "aws_apigatewayv2_route" "default" {
#   api_id    = aws_apigatewayv2_api.gw.id
#   route_key = "$default"
#   target    = "integrations/${aws_apigatewayv2_integration.default.id}"
# }

resource "aws_apigatewayv2_integration" "query" {
  api_id           = aws_apigatewayv2_api.gw.id
  integration_type = "AWS_PROXY"
  integration_uri  = aws_lambda_function.auth.invoke_arn
  # request_parameters = {
  #   FORWARD_URL = "${module.vars.env.forward_url}/transaction/query"
  # }
}

resource "aws_apigatewayv2_route" "query" {
  api_id             = aws_apigatewayv2_api.gw.id
  route_key          = "POST /transaction/query"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
  authorization_type = "JWT"
  target             = "integrations/${aws_apigatewayv2_integration.query.id}"
}

resource "aws_apigatewayv2_integration" "invoke" {
  api_id           = aws_apigatewayv2_api.gw.id
  integration_type = "AWS_PROXY"
  integration_uri  = aws_lambda_function.auth.invoke_arn
  # request_parameters = {
  #   FORWARD_URL = "${module.vars.env.forward_url}/transaction/invoke"
  # }
}

resource "aws_apigatewayv2_route" "invoke" {
  api_id             = aws_apigatewayv2_api.gw.id
  route_key          = "POST /transaction/invoke"
  authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
  authorization_type = "JWT"
  target             = "integrations/${aws_apigatewayv2_integration.invoke.id}"
}
