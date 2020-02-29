/*

Use the Bindplane provider to define a credential for postgres.

Username is pulled from Vault and password is generated randomly.

*/
resource "bindplane_log_source" "mssql" {
  configuration = <<CONFIGURATION
{
"name": "terraform_test",
"configuration": {
  "read_interval": "5",
  "max_reads": "1000",
  "read_from_head": false
},
"source_type_id": "microsoft_sqlserver",
"source_version": "3.2.0",
"custom_template": ""
}
CONFIGURATION

}
