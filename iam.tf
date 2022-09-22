# This is managed externally by OISM for the NIST infrastructure

data "aws_iam_role" "apigw-proxy" {
  name = "frontend-apigw-s3-integration-role"
}

# # terraform import aws_iam_role.apigw-proxy frontend-apigw-s3-integration-role
# resource "aws_iam_role" "apigw-proxy" {
#   # we do not have permissions to change this role
#   name        = "frontend-apigw-s3-integration-role"
#   description = "see SCTASK0613704. added by Fisan 8/18/22"
#   tags        = {}
#   assume_role_policy = jsonencode({
#     Statement = [
#       {
#         Action = "sts:AssumeRole"
#         Effect = "Allow"
#         Principal = {
#           Service = "apigateway.amazonaws.com"
#         }
#         Sid = ""
#       }
#     ]
#     Version = "2012-10-17"
#   })
# }

# # terraform import aws_iam_policy.s3-integration arn:aws:iam::259202176582:policy/frontend-apigw-s3-integration-role_policy
# resource "aws_iam_policy" "s3-integration" {
#   name = "frontend-apigw-s3-integration-role_policy"
#   policy = jsonencode({
#     Version = "2012-10-17",
#     Statement = [
#       {
#         Effect = "Allow",
#         Action = [
#           "s3:Get*",
#           "s3:List*",
#           "s3-object-lambda:Get*",
#           "s3-object-lambda:List*"
#         ],
#         Resource = "*",
#         Condition = {
#           StringEquals = {
#             "s3:ExistingObjectTag/Purpose" = "blossom-frontend"
#           }
#         }
#       }
#     ]
#   })
# }

# # terraform import aws_iam_role_policy_attachment.apigw-proxy_s3-integration \
# #   frontend-apigw-s3-integration-role/arn:aws:iam::259202176582:policy/frontend-apigw-s3-integration-role_policy
# resource "aws_iam_role_policy_attachment" "apigw-proxy_s3-integration" {
#   role       = aws_iam_role.apigw-proxy.name
#   policy_arn = aws_iam_policy.s3-integration.arn
# }
