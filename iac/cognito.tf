data "aws_cognito_user_pools" "identity" {
  name = "blossom_test"
}

data "aws_cognito_user_pool_clients" "identity" {
  user_pool_id = tolist(data.aws_cognito_user_pools.identity.ids)[0]
}

data "aws_cognito_user_pool_client" "main" {
  user_pool_id = tolist(data.aws_cognito_user_pools.identity.ids)[0]
  client_id = data.aws_cognito_user_pool_clients.identity.client_ids[0]
} 
