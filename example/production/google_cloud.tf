resource "bindplane_credential" "google_cloud" {
  configuration = <<CONFIGURATION
{
  "name": "${data.vault_generic_secret.poc_secrets.data["gcp_project"]}-compute-engine",
  "credential_type_id": "${var.bindplane_credential_types["google_compute"]}",
  "parameters": {
    "private_key_json": "{\"auth_provider_x509_cert_url\":\"https://www.googleapis.com/oauth2/v1/certs\",\"auth_uri\":\"https://accounts.google.com/o/oauth2/auth\",\"client_email\":\"${data.vault_generic_secret.poc_secrets.data["client_email"]}\",\"client_id\":\"${data.vault_generic_secret.poc_secrets.data["client_id"]}\",\"client_x509_cert_url\":\"${data.vault_generic_secret.poc_secrets.data["client_x509_cert_url"]}\",\"private_key\":\"-----BEGIN PRIVATE KEY-----\\n${data.vault_generic_secret.poc_secrets.data["private_key"]}\\n-----END PRIVATE KEY-----\\n\",\"private_key_id\":\"${data.vault_generic_secret.poc_secrets.data["private_key_id"]}\",\"project_id\":\"${data.vault_generic_secret.poc_secrets.data["project_id"]}\",\"token_uri\":\"https://oauth2.googleapis.com/token\",\"type\":\"service_account\"}"
  }
}
CONFIGURATION

}

resource "bindplane_source" "google_cloud_storage" {
  provisioning_timeout = 90

  name                = "google-cloud-storage"
  source_type         = "google_cloud_storage"
  collector_id        = "${module.gcp_collector.id}"
  collection_interval = 300
  credential_id       = "${bindplane_credential.google_cloud.id}"

  configuration = <<CONFIGURATION
{
  "connection_timeout": 30,
  "max_single_request_retry_timeout": 64,
  "max_total_request_retry_timeout": 300,
  "metric_collection": "all",
  "port": 443,
  "projects": "*",
  "ssl_config": "Verify"
}
CONFIGURATION
}
