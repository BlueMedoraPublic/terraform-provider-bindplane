/*

bpcli logs source type parameters \
  --source-type-id mysql

*/
variable "error_log_path" {
  description = "example: you can use interpolation for configuraton params"
  default = "/var/log/mysql/mysqld.log"
}

resource "bindplane_log_source" "mysql" {
  name = "mysql-terraform"
  source_type_id = "mysql"
  source_version = "2.0.0"
  configuration = <<CONFIGURATION
{
  "mysql_error": true,
  "error_log_path": "${var.error_log_path}",
  "mysql_slow_query": true,
  "slow_query_log_path": "/var/log/mysql/slow.log",
  "mysql_general": false,
  "general_log_path": "/var/log/mysql/general.log",
  "read_from_head": false
}
CONFIGURATION
}

/*

bpcli logs destination type parameters \
  --destination-type-id stackdriver

*/
resource "bindplane_log_destination" "stackdriver" {
  name = "stackdriver-terraform"
  destination_type_id = "stackdriver"
  destination_version = "1.3.2"
  configuration = <<CONFIGURATION
{
  "credentials": ${data.google_secret_manager_secret_version.bindplane_svc_act.secret_data},
  "location": "us-west1"
}
CONFIGURATION

}

/*resource "bindplane_log_template" "mysql_prod" {
    name = "template-terraform"
    source_config_ids = [
        bindplane_log_source.mysql.id
    ]
    destination_config_id =  bindplane_log_destination.stackdriver.id
    agent_group = ""
}*/
