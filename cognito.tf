data "aws_cognito_user_pools" "identity" {
  name = "blossom_test"
}

data "aws_cognito_user_pool_clients" "identity" {
  user_pool_id = aws_cognito_user_pool.identity.id
}

# resource "aws_cognito_user_pool" "identity" {
#   name = "${local.prefix}-identity-userpool"
#   tags = local.tags
# }

# resource "aws_cognito_user_pool_client" "identity" {
#   name         = "${local.prefix}-identity-userpool-client"
#   user_pool_id = aws_cognito_user_pool.identity.id
# }
