module "vars" {
  source = "./module/vars"
  # Use the selected terraform workspace to select the environment
  environment = terraform.workspace
}

locals {
  prefix = "b-${module.vars.env.network_name}-${lower(module.vars.env.member_name)}"
  tags = {
    "Terraform"            = "true"
    "Blossom_Network_Name" = module.vars.env.network_name
    "Blossom_Member_Name"  = module.vars.env.member_name
  }
}
