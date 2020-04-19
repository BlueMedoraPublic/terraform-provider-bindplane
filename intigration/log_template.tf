resource "bindplane_log_template" "eks" {
    name = "${var.name}-template"
    source_config_ids = [
        bindplane_log_source.mysql.id
    ]
    destination_config_id =  bindplane_log_destination_google.default.id
    agent_group = ""
}
