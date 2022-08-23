# run npm build to build the dashboard
resource "null_resource" "build_blossom_dashboard" {
  provisioner "local-exec" {
    command = "cd ../blossom-dashboard/client; yarn; echo 0"
  }
}

module "s3_content_bucket" {
  source               = "../infrastructure/terraform/modules/aws/s3"
  bucket               = "${local.prefix}-content"
  tags                 = local.tags
  acl                  = "private"
  attach_public_policy = false
  block_public_acls    = false
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
  for_each = fileset("../blossom-dashboard/client/build", "**/*")
  key      = each.value
  source   = "../blossom-dashboard/client/build/${each.value}"
  etag     = filemd5("../blossom-dashboard/client/build/${each.value}")
  tags = {
    "Purpose" = "blossom-frontend"
  }
  # run npm build before uploading assets
  depends_on = [
    null_resource.build_blossom_dashboard
  ]
}
