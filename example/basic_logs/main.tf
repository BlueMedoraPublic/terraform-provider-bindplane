resource "bindplane_log_source" "mysql" {
  configuration = <<CONFIGURATION
{
"name": "mysql-terraform",
"configuration": {
  "mysql_error": true,
  "error_log_path": "/var/log/mysql/mysqld.log",
  "mysql_slow_query": true,
  "slow_query_log_path": "/var/log/mysql/slow.log",
  "mysql_general": false,
  "general_log_path": "/var/log/mysql/general.log",
  "read_from_head": false
},
"source_type_id": "mysql",
"source_version": "2.0.0",
"custom_template": ""
}
CONFIGURATION

}

// service account json is pulled from gcp secret manager
// to avoid putting secrets in the repo
resource "bindplane_log_destination" "stackdriver" {
  configuration = <<CONFIGURATION
{
"name": "stackdriver-terraform",
"configuration": {
  "credentials": ${data.google_secret_manager_secret_version.bindplane_svc_act.secret_data},
  "location": "us-west1"
},
"destination_type_id": "stackdriver",
"destination_version": "1.3.2"
}
CONFIGURATION

}
