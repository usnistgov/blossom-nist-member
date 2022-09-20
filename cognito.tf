resource "aws_cognito_user_pool" "identity" {
  name = "${local.prefix}-identity-userpool"
  tags = local.tags
}

# resource "aws_cognito_user_pool_client" "identity" {
#   name         = "${local.prefix}-identity-userpool-client"
#   user_pool_id = aws_cognito_user_pool.identity.id
# }
