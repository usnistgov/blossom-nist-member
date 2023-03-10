resource "aws_api_gateway_rest_api" "gw" {
  name        = "${local.prefix}-gw"
  description = "The API Gateway for the ${module.vars.env.member_name} member"
  tags        = local.tags
}

resource "aws_api_gateway_deployment" "gw-deployment" {
  rest_api_id = aws_api_gateway_rest_api.gw.id

  triggers = {
    redeployment = sha1(jsonencode([
      # s3 integration
      aws_api_gateway_resource.s3,
      aws_api_gateway_method.s3,
      aws_api_gateway_integration.s3,
      aws_api_gateway_integration_response.s3,
      # s3 root integration
      aws_api_gateway_method.s3-root,
      aws_api_gateway_integration.s3-root,
      aws_api_gateway_integration_response.s3-root,
      # lambda integration
      aws_api_gateway_resource.lambda,
      aws_api_gateway_method.lambda,
      aws_api_gateway_integration.lambda
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }
}

output "gw_url" {
  value       = resource.aws_api_gateway_deployment.gw-deployment.invoke_url
  description = "The URL of the active APIGW deployment"
  sensitive   = true
}

output "base_url" {
  value       = "/${aws_api_gateway_stage.gw-stage.stage_name}/"
  description = "The base URL links should resolve against in the dashboard"
  sensitive   = false
}

resource "aws_api_gateway_stage" "gw-stage" {
  deployment_id = aws_api_gateway_deployment.gw-deployment.id
  rest_api_id   = aws_api_gateway_rest_api.gw.id
  stage_name    = "dev"
  tags          = local.tags
}

#############################
# S3 Integration Definition #
#############################

resource "aws_api_gateway_resource" "s3" {
  rest_api_id = aws_api_gateway_rest_api.gw.id
  parent_id   = aws_api_gateway_rest_api.gw.root_resource_id
  path_part   = "{proxy+}"
}

resource "aws_api_gateway_method" "s3" {
  rest_api_id   = aws_api_gateway_rest_api.gw.id
  resource_id   = aws_api_gateway_resource.s3.id
  http_method   = "GET"
  authorization = "NONE"

  request_parameters = {
    "method.request.path.proxy" = true
  }
}

resource "aws_api_gateway_method_response" "s3" {
  rest_api_id = aws_api_gateway_rest_api.gw.id
  resource_id = aws_api_gateway_resource.s3.id
  http_method = aws_api_gateway_method.s3.http_method

  status_code = 200
  response_parameters = {
    "method.response.header.Content-Type" = true
  }
}

resource "aws_api_gateway_integration" "s3" {
  rest_api_id = aws_api_gateway_rest_api.gw.id
  resource_id = aws_api_gateway_resource.s3.id
  http_method = aws_api_gateway_method.s3.http_method

  integration_http_method = "GET"
  type                    = "AWS"
  uri                     = "arn:aws:apigateway:us-east-1:s3:path/${module.s3_content_bucket.s3_bucket_id}/{proxy}"

  credentials          = data.aws_iam_role.apigw-proxy.arn
  passthrough_behavior = "WHEN_NO_MATCH"

  request_parameters = {
    "integration.request.path.proxy" = "method.request.path.proxy"
  }
}

resource "aws_api_gateway_integration_response" "s3" {
  rest_api_id = aws_api_gateway_rest_api.gw.id
  resource_id = aws_api_gateway_resource.s3.id
  http_method = aws_api_gateway_method.s3.http_method

  status_code = 200
  response_parameters = {
    "method.response.header.Content-Type" = "integration.response.header.Content-Type"
  }

  depends_on = [
    aws_api_gateway_integration.s3
  ]
}

##################################
# S3 Root Integration Definition #
##################################

resource "aws_api_gateway_method" "s3-root" {
  rest_api_id   = aws_api_gateway_rest_api.gw.id
  resource_id   = aws_api_gateway_rest_api.gw.root_resource_id
  http_method   = "GET"
  authorization = "NONE"
}

resource "aws_api_gateway_method_response" "s3-root" {
  rest_api_id = aws_api_gateway_rest_api.gw.id
  resource_id = aws_api_gateway_rest_api.gw.root_resource_id
  http_method = aws_api_gateway_method.s3-root.http_method

  status_code = 200
  response_parameters = {
    "method.response.header.Content-Type" = true
  }
}

resource "aws_api_gateway_integration" "s3-root" {
  rest_api_id = aws_api_gateway_rest_api.gw.id
  resource_id = aws_api_gateway_rest_api.gw.root_resource_id
  http_method = aws_api_gateway_method.s3-root.http_method

  integration_http_method = "GET"
  type                    = "AWS"
  uri                     = "arn:aws:apigateway:us-east-1:s3:path/${module.s3_content_bucket.s3_bucket_id}/index.html"

  credentials          = data.aws_iam_role.apigw-proxy.arn
  passthrough_behavior = "WHEN_NO_MATCH"
}

resource "aws_api_gateway_integration_response" "s3-root" {
  rest_api_id = aws_api_gateway_rest_api.gw.id
  resource_id = aws_api_gateway_rest_api.gw.root_resource_id
  http_method = aws_api_gateway_method.s3-root.http_method

  status_code = 200
  response_parameters = {
    "method.response.header.Content-Type" = "integration.response.header.Content-Type"
  }

  depends_on = [
    aws_api_gateway_integration.s3-root
  ]
}

#################################
# Lambda Integration Definition #
#################################

resource "aws_api_gateway_resource" "lambda" {
  rest_api_id = aws_api_gateway_rest_api.gw.id
  parent_id   = aws_api_gateway_rest_api.gw.root_resource_id

  path_part = "transaction"
}

resource "aws_api_gateway_method" "lambda" {
  rest_api_id = aws_api_gateway_rest_api.gw.id
  resource_id = aws_api_gateway_resource.lambda.id

  http_method          = "POST"
  authorization        = "COGNITO_USER_POOLS"
  authorizer_id        = aws_api_gateway_authorizer.cognito_integration.id
  authorization_scopes = ["openid", "email"]
}

resource "aws_api_gateway_integration" "lambda" {
  rest_api_id = aws_api_gateway_rest_api.gw.id
  resource_id = aws_api_gateway_resource.lambda.id
  http_method = aws_api_gateway_method.lambda.http_method

  uri                     = aws_lambda_function.query.invoke_arn
  type                    = "AWS_PROXY"
  integration_http_method = "POST"
}

resource "aws_lambda_permission" "lambda-permission" {
  statement_id  = "AllowMyDemoAPIInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.query.function_name
  principal     = "apigateway.amazonaws.com"

  # The /*/*/* part allows invocation from any stage, method and resource path
  # within API Gateway REST API.
  source_arn = "${aws_api_gateway_rest_api.gw.execution_arn}/*/*/*"
}

resource "aws_api_gateway_authorizer" "cognito_integration" {
  name        = "blossom_test-cognito_integration"
  type        = "COGNITO_USER_POOLS"
  rest_api_id = aws_api_gateway_rest_api.gw.id
  provider_arns = [
    tolist(data.aws_cognito_user_pools.identity.arns)[0]
  ]
}
