resource "bindplane_log_agent_populate" "agent" {
    name = var.name
    provisioning_timeout = 30 // the min allowed by the provider
}
