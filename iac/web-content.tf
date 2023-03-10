locals {
  content_type_map = {
    "html" = "text/html"
    "css"  = "text/css"
    "js"   = "text/javascript"
    "png"  = "image/png"
    "ico"  = "image/x-icon"
    "svg"  = "image/svg+xml"
    "txt"  = "text/plain"
    "json" = "application/json"
    # idk what this is
    "map" = "application/json"
  }
  webcontent_builddir = "${path.module}/../dashboard/dist"
}

output "vite_dev_env" {
  value       = <<-EOT
  VITE_CLIENT_ID=${data.aws_cognito_user_pool_client.main.id}
  VITE_CLIENT_SECRET=${data.aws_cognito_user_pool_client.main.client_secret}
  VITE_AUTH_URL=https://${data.aws_cognito_user_pools.identity.name}.auth.${var.aws_region}.amazoncognito.com
  PROXY_URL=${resource.aws_api_gateway_deployment.gw-deployment.invoke_url}
  EOT
  description = "The developer environment used by the dashboard"
  sensitive   = true
}

module "s3_content_bucket" {
  source               = "./module/s3"
  bucket               = "${local.prefix}-content"
  tags                 = local.tags
  acl                  = "private"
  attach_public_policy = false
  block_public_acls    = true
  # server_side_encryption_configuration = try(lookup(var.server_side_encryption_configuration, "rule"), {
  #   "rule" : {
  #     "apply_server_side_encryption_by_default" : {
  #       "sse_algorithm" : "aws:kms"
  #       "kms_master_key_id" : data.aws_kms_alias.s3.arn
  #     }
  #   }
  # })
}

resource "aws_s3_object" "web-content" {
  bucket   = module.s3_content_bucket.s3_bucket_id
  for_each = fileset(local.webcontent_builddir, "**/*")
  key      = each.value
  source   = "${local.webcontent_builddir}/${each.value}"
  etag     = filemd5("${local.webcontent_builddir}/${each.value}")
  tags = merge({
    "Purpose" = "blossom-frontend"
  }, local.tags)
  # extract the extension, apply it to the content_type_map
  content_type = local.content_type_map[split(".", each.value)[length(split(".", each.value)) - 1]]
}
