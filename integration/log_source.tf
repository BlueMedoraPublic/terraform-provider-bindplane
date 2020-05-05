resource "bindplane_log_source" "mysql" {
  name = "mysql-tf-${var.name}"
  source_type_id = "mysql"
  source_version = "2.0.0"
  configuration = <<CONFIGURATION
{
  "mysql_error": true,
  "error_log_path": "/var/log/mysql/error.log",
  "mysql_slow_query": true,
  "slow_query_log_path": "/var/log/mysql/slow.log",
  "mysql_general": false,
  "general_log_path": "/var/log/mysql/general.log",
  "read_from_head": false
}
CONFIGURATION
}
