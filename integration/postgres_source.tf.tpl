resource "bindplane_credential" "postgres" {
  configuration = <<CONFIGURATION
{
  "name": ${var.name},
  "credential_type_id": "6c2c63b5-d465-46ef-bfce-b0881066e43b",
  "parameters": {
    "username": "postgres",
    "password": "password"
  }
}
CONFIGURATION

}

resource "bindplane_source" "postgres_app_db_0" {
  provisioning_timeout = 120
  name = var.name
  source_type = "postgresql"
  collector_id = var.collector-id
  collection_interval = 60
  credential_id = bindplane_credential.postgres.id

  configuration = <<CONFIGURATION
{
    "collection_mode": "normal",
    "function_count": 20,
    "host": "ADDR",
    "monitor_indexes": true,
    "monitor_sequences": false,
    "monitor_sessions": false,
    "monitor_tables": false,
    "monitor_triggers": false,
    "order_functions_by": "calls",
    "order_queries_by": "calls",
    "port": 5432,
    "query_count": 10,
    "show_query_text": true,
    "ssl_config": "No SSL"
}
CONFIGURATION

}
