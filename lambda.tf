module "remote_state_bucket" {
  source = "../infrastructure/terraform/modules/aws/s3"
  # attach_public_policy = var.attach_public_policy
  bucket               = "${local.prefix}-lambda"
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

resource "aws_lambda_function" "auth" {
  runtime          = "nodejs16.x"
  function_name    = "handler"
  s3_bucket        = aws_s3_object.auth_lambda.bucket
  s3_key           = aws_s3_object.auth_lambda.key
  handler          = "handler.handle"
  source_code_hash = data.archive_file.auth_lambda.output_base64sha256
  role             = data.aws_iam_role.lambda_role.arn
  # environment {
  #   variables = {

  #   }
  # }
}

data "aws_iam_role" "lambda_role" {
  name = "nistitlblossom-auto-tagging-lambda-role"
}

data "archive_file" "auth_lambda" {
  type        = "zip"
  source_dir  = "${path.module}/auth_lambda"
  output_path = "${path.module}/auth_lambda.zip"
}

resource "aws_s3_object" "auth_lambda" {
  bucket = module.remote_state_bucket.s3_bucket_id

  key    = "auth_lambda.zip"
  source = data.archive_file.auth_lambda.output_path

  etag = filesha1(data.archive_file.auth_lambda.output_path)
}
