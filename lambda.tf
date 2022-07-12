resource "aws_lambda_function" "auth" {
  runtime          = "nodejs16.x"
  function_name    = "handler"
  s3_bucket        = aws_s3_object.auth_lambda.bucket
  s3_key           = aws_s3_object.auth_lambda.key
  handler          = "handler.handle"
  source_code_hash = data.archive_file.auth_lambda.output_base64sha256
  role             = data.aws_iam_role.lambda_role.arn
}


data "aws_iam_role" "lambda_role" {
  name = "nistitlblossom-auto-tagging-lambda-role"
}
