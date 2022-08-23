resource "aws_cognito_user_pool" "identity" {
  name = "${local.prefix}-identity"
  tags = local.tags
}
