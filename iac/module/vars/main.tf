variable "environment" {
  description = "The environment (hlf member) for which to fetch the configuration for"
  type        = string
}

variable "configuration_dir" {
  description = "The directory to fetch configuration files from"
  type        = string
}

locals {
  configuration_files = fileset(var.configuration_dir, "*.json")
  # environment => configuration keypairs
  configurations = {
    for file_name in local.configuration_files :
    trimsuffix(file_name, ".json") => jsondecode(file("${var.configuration_dir}/${file_name}"))
  }
}

output "env" {
  description = "The configuration for the given environment"
  value       = local.configurations[var.environment]
}

output "envs" {
  description = "All defined environments"
  value       = keys(local.configurations)
}
