# https://raw.githubusercontent.com/terraform-aws-modules/terraform-aws-s3-bucket/31d76f933b05848be9aaf25befd43966e4065472/versions.tf

terraform {
  required_version = ">= 0.13.1"

  required_providers {
    aws = ">= 3.64"
  }
}
