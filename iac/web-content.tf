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
  webcontent_srcdir   = "${path.module}/../dashboard"
  webcontent_builddir = "${local.webcontent_srcdir}/dist"
  webcontent_env = merge({
    VITE_CLIENT_ID     = data.aws_cognito_user_pool_client.main.id
    VITE_CLIENT_SECRET = data.aws_cognito_user_pool_client.main.client_secret
    VITE_AUTH_URL      = "https://blossomtest.auth.us-east-1.amazoncognito.com"
    },
    {
      BASE_URL = "${aws_api_gateway_stage.gw-stage.stage_name}/"
  })
}

# run npm build to build the dashboard
resource "null_resource" "build_blossom_dashboard" {
  triggers = {
    "package"      = sha256(file("${local.webcontent_srcdir}/package.json"))
    "package-lock" = sha256(file("${local.webcontent_srcdir}/package-lock.json"))
    "src"          = sha256(join("", [for f in fileset(local.webcontent_srcdir, "src/**/*") : filesha256("${local.webcontent_srcdir}/${f}")]))
    "env"          = jsonencode(local.webcontent_env)
  }
  provisioner "local-exec" {
    command     = "cd ${local.webcontent_srcdir}; npm i; npm run build"
    environment = local.webcontent_env
  }
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
  # run npm build before uploading assets
  depends_on = [
    null_resource.build_blossom_dashboard
  ]
}
