output "network_id" {
    value = module.vars.env.network_id
    sensitive = false
}

output "channel_name" {
    value = module.vars.env.channel_name
    sensitive = false
}

output "connection_profile_file" {
    value = "conn-profile-${module.vars.env.network_name}-${module.vars.env.member_name}.yaml"
    sensitive = false
}