resource "aws_api_gateway_rest_api" "gw" {
  name        = "${local.prefix}-gw"
  description = "The API Gateway for the ${module.vars.env.member_name} member"
  body        = <<EOT
    {
      "swagger" : "2.0",
      "info" : {
        "version" : "2022-08-16T21:58:05Z",
        "title" : "apigw-proxy-test"
      },
      "host" : "ayk7nyknel.execute-api.us-east-1.amazonaws.com",
      "basePath" : "/prod",
      "schemes" : [ "https" ],
      "paths" : {
        "/{proxy+}" : {
          "x-amazon-apigateway-any-method" : {
            "produces" : [ "application/json" ],
            "parameters" : [ {
              "name" : "proxy",
              "in" : "path",
              "required" : true,
              "type" : "string"
            } ],
            "responses" : {
              "200" : {
                "description" : "200 response",
                "schema" : {
                  "$ref" : "#/definitions/Empty"
                },
                "headers" : {
                  "Content-Type" : {
                    "type" : "string"
                  }
                }
              }
            },
            "x-amazon-apigateway-integration" : {
              "credentials" : "${aws_iam_role.apigw-proxy.arn}",
              "httpMethod" : "GET",
              "uri" : "arn:aws:apigateway:us-east-1:s3:path/${module.s3_content_bucket.s3_bucket_id}/{proxy}",
              "responses" : {
                "default" : {
                  "statusCode" : "200",
                  "responseParameters" : {
                    "method.response.header.Content-Type" : "integration.response.header.Content-Type"
                  }
                }
              },
              "requestParameters" : {
                "integration.request.path.proxy" : "method.request.path.proxy"
              },
              "passthroughBehavior" : "when_no_match",
              "cacheNamespace" : "878fmp",
              "cacheKeyParameters" : [ "method.request.path.proxy", "integration.request.path.proxy" ],
              "type" : "aws"
            }
          }
        }
      },
      "definitions" : {
        "Empty" : {
          "type" : "object",
          "title" : "Empty Schema"
        }
      }
    }
  EOT
}

resource "aws_api_gateway_deployment" "gw-deployment" {
  rest_api_id = aws_api_gateway_rest_api.gw.id

  triggers = {
    redeployment = sha1(jsonencode(aws_api_gateway_rest_api.gw.body))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "gw-stage" {
  deployment_id = aws_api_gateway_deployment.gw-deployment.id
  rest_api_id   = aws_api_gateway_rest_api.gw.id
  stage_name    = "dev"
}

# resource "aws_apigatewayv2_api" "gw" {
#   name          = "${local.prefix}-gw"
#   protocol_type = "HTTP"

#   # cors_configuration = {
#   #   allow_headers = ["content-type", "x-amz-date", "authorization", "x-api-key", "x-amz-security-token", "x-amz-user-agent"]
#   #   allow_methods = ["*"]
#   #   allow_origins = ["*"]
#   # }

#   tags = local.tags
# }

# resource "aws_apigatewayv2_authorizer" "cognito" {
#   api_id           = aws_apigatewayv2_api.gw.id
#   name             = "${local.prefix}-identity-authorizer"
#   authorizer_type  = "JWT"
#   identity_sources = ["$request.header.Authorization"]

#   jwt_configuration {
#     # annoyingly at least one audience is required
#     audience = ["transact"]
#     issuer   = "https://${aws_cognito_user_pool.identity.endpoint}"
#   }
# }

# resource "aws_apigatewayv2_integration" "query" {
#   api_id           = aws_apigatewayv2_api.gw.id
#   integration_type = "AWS_PROXY"
#   integration_uri  = aws_lambda_function.auth.invoke_arn
#   # request_parameters = {
#   #   FORWARD_URL = "${module.vars.env.forward_url}/transaction/query"
#   # }
# }

# resource "aws_apigatewayv2_route" "query" {
#   api_id             = aws_apigatewayv2_api.gw.id
#   route_key          = "POST /transaction/query"
#   authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
#   authorization_type = "JWT"
#   target             = "integrations/${aws_apigatewayv2_integration.query.id}"
# }

# resource "aws_apigatewayv2_integration" "invoke" {
#   api_id           = aws_apigatewayv2_api.gw.id
#   integration_type = "AWS_PROXY"
#   integration_uri  = aws_lambda_function.auth.invoke_arn
#   # request_parameters = {
#   #   FORWARD_URL = "${module.vars.env.forward_url}/transaction/invoke"
#   # }
# }

# resource "aws_apigatewayv2_route" "invoke" {
#   api_id             = aws_apigatewayv2_api.gw.id
#   route_key          = "POST /transaction/invoke"
#   authorizer_id      = aws_apigatewayv2_authorizer.cognito.id
#   authorization_type = "JWT"
#   target             = "integrations/${aws_apigatewayv2_integration.invoke.id}"
# }
