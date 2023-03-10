terraform {
  backend "s3" {
    # variables cannot be used here
    bucket = "nist-blossom-iac"
    key    = "terraform.tfstate"
    region = "us-east-1"
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.18.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
  # cannot be specified on our config due to permission issues :/
  # default_tags {
  #   tags = local.tags
  # }
}
