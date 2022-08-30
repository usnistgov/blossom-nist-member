# run npm build to build the dashboard
resource "null_resource" "build_blossom_dashboard" {
  provisioner "local-exec" {
    command = "cd ../blossom-dashboard/client; yarn; PUBLIC_URL=/dev echo 0"
  }
}

module "s3_content_bucket" {
  source               = "../infrastructure/terraform/modules/aws/s3"
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

locals {
  content_type_map = {
    "html" = "text/html"
    "css"  = "text/css"
    "js"   = "text/javascript"
    "png"  = "image/png"
    "ico"  = "image/x-icon"
    "txt"  = "text/plain"
    "json" = "application/json"
    # idk what this is
    "map" = "application/json"
  }
}

resource "aws_s3_object" "web-content" {
  bucket   = module.s3_content_bucket.s3_bucket_id
  for_each = fileset("../blossom-dashboard/client/build", "**/*")
  key      = each.value
  source   = "../blossom-dashboard/client/build/${each.value}"
  etag     = filemd5("../blossom-dashboard/client/build/${each.value}")
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
