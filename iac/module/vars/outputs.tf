output "env" {
  value = local.environments[var.environment]
}

output "envs" {
  value = [
    "BLOSSON_NIST2"
  ]
}
