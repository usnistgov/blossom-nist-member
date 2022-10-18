locals {
  # the output zip file containing the packaged lambda contents
  lambda_outpath = "lambda.zip"
  # the input directory to build
  lambda_srcdir = "lambda"
  # the output directory for the built lambda
  lambda_builddir = "${local.lambda_srcdir}/dist"
}

# This bucket stores the lambda's build artifacts
module "lambda_bucket" {
  source = "./module/s3"
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

resource "aws_lambda_function" "query" {
  runtime          = "nodejs16.x"
  function_name    = "${local.prefix}-lambda"
  memory_size      = 512
  publish          = true
  s3_bucket        = aws_s3_object.query_lambda.bucket
  s3_key           = aws_s3_object.query_lambda.key
  handler          = "index.handler"
  source_code_hash = data.archive_file.query_lambda.output_base64sha256
  role             = data.aws_iam_role.lambda_role.arn
  tags             = local.tags
  timeout          = 300
  vpc_config {
    subnet_ids = [
      "subnet-0e55c3a77dad7f698",
      "subnet-09e992c8f703e3662"
    ]
    security_group_ids = [
      "sg-0b94936c423b8e9ee",
      "sg-064e0191188232c57"
    ]
  }
  environment {
    variables = {
      CHANNEL_NAME    = module.vars.env.channel_name
      CONTRACT_NAME   = module.vars.env.contract_name
      PROFILE_ENCODED = filebase64("${path.module}/conn-profile-${module.vars.env.network_name}-${module.vars.env.member_name}.yaml")
    }
  }

}

data "aws_iam_role" "lambda_role" {
  name = "LambdaExecutionRole"
}

# this resource builds the lambda
resource "null_resource" "build-lambda" {
  triggers = {
    "package"      = sha256(file("${local.lambda_srcdir}/package.json"))
    "package-lock" = sha256(file("${local.lambda_srcdir}/package-lock.json"))
    "src"          = sha256(join("", [for f in fileset(local.lambda_srcdir, "src/**/*") : filesha256("${local.lambda_srcdir}/${f}")]))
  }
  provisioner "local-exec" {
    command = "pushd ${local.lambda_srcdir}; npm i; npm run build"
    interpreter = [
      "bash", "-c"
    ]
  }
}

data "archive_file" "query_lambda" {
  type        = "zip"
  source_dir  = local.lambda_builddir
  output_path = local.lambda_outpath

  depends_on = [
    null_resource.build-lambda
  ]
}

resource "aws_s3_object" "query_lambda" {
  bucket = module.lambda_bucket.s3_bucket_id

  key    = "lambda.zip"
  source = data.archive_file.query_lambda.output_path
  tags   = local.tags
  etag   = filesha256(data.archive_file.query_lambda.output_path)
}
