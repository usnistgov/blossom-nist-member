data "aws_cognito_user_pools" "identity" {
  name = "blossom_test"
}

data "aws_cognito_user_pool_clients" "identity" {
  user_pool_id = data.aws_cognito_user_pools.identity.id
}
