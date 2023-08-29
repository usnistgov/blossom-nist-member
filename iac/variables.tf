module "vars" {
  source = "./module/vars"
  # Use the selected terraform workspace to select the environment
  environment       = terraform.workspace
  configuration_dir = "./configurations"
}

locals {
  prefix = "b-${module.vars.env.network_name}-${lower(module.vars.env.member_name)}"
  tags = {
    "Terraform"            = "true"
    "Blossom_Network_Name" = module.vars.env.network_name
    "Blossom_Member_Name"  = module.vars.env.member_name
  }
}

variable "hlf_debug" {
  type        = bool
  default     = false
  sensitive   = false
  description = "Enables HFC debug logging on the query lambda"
}

variable "cognito_debug" {
  type        = bool
  default     = true
  sensitive   = false
  description = "Enables localhost callback URLs in the Cognito user pool client"
}

variable "aws_region" {
  type        = string
  default     = "us-east-1"
  sensitive   = false
  description = "The AWS region to use"
}
