resource "bindplane_log_bind_source" "mysql" {
    source_config_id = bindplane_log_source.mysql.id
    agent_id         = bindplane_log_agent_populate.agent.id
}

resource "bindplane_log_bind_destination" "stackdriver" {
    destination_config_id = bindplane_log_destination_google.default.id
    agent_id              = bindplane_log_agent_populate.agent.id
}
