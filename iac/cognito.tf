data "aws_cognito_user_pools" "identity" {
  name = module.vars.env.cognito_user_pool_name
}

locals {
  cognito_user_pool_id = tolist(data.aws_cognito_user_pools.identity.ids)[0]
  debug_callback_url   = "http://localhost:5173/"
}

resource "aws_cognito_user_pool_client" "client" {
  name                                 = "${local.prefix}-user-pool-client"
  user_pool_id                         = local.cognito_user_pool_id
  generate_secret                      = true
  callback_urls                        = var.cognito_debug ? [local.apigw_url, local.debug_callback_url] : [local.apigw_url]
  logout_urls                          = var.cognito_debug ? [local.apigw_url, local.debug_callback_url] : [local.apigw_url]
  explicit_auth_flows                  = ["ALLOW_CUSTOM_AUTH", "ALLOW_REFRESH_TOKEN_AUTH", "ALLOW_USER_SRP_AUTH"]
  enable_token_revocation              = true
  prevent_user_existence_errors        = "ENABLED"
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["code"]
  allowed_oauth_scopes                 = ["email", "openid"]
  supported_identity_providers         = ["COGNITO"]
}

output "client_id" {
  value     = resource.aws_cognito_user_pool_client.client.id
  sensitive = false
}

output "client_secret" {
  value     = resource.aws_cognito_user_pool_client.client.client_secret
  sensitive = true
}

locals {
  cognito_domain_prefix = "${module.vars.env.network_name}-${lower(module.vars.env.member_name)}"
}

resource "aws_cognito_user_pool_domain" "domain" {
  # e.x. blosson-nist2.auth.us-east-1.amazoncognito.com
  domain       = local.cognito_domain_prefix
  user_pool_id = local.cognito_user_pool_id
}

output "auth_url" {
  value     = "https://${local.cognito_domain_prefix}.auth.${var.aws_region}.amazoncognito.com"
  sensitive = false
}
