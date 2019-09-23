resource "bindplane_credential" "pure_storage" {
  configuration = <<CONFIGURATION
{
  "name": "pure_storage",
  "credential_type_id": "${var.bindplane_credential_types["pure_storage"]}",
  "parameters": {
    "username": "${data.vault_generic_secret.poc_onprem.data["pure_user"]}",
    "password": "${data.vault_generic_secret.poc_onprem.data["pure_pass"]}"
  }
}
CONFIGURATION

}

resource "bindplane_source" "purestorage_flasharray" {
  provisioning_timeout = 90
  name                 = "purestorage_flasharray"
  source_type          = "purestorage_flasharray"
  collector_id         = data.vault_generic_secret.poc_onprem.data["onprem_collector_id"]
  collection_interval  = 300
  credential_id        = bindplane_credential.pure_storage.id

  configuration = <<CONFIGURATION
{
    "host": "${data.vault_generic_secret.poc_onprem.data["pure_host"]}",
    "max_simultaneous_requests": 4,
    "port": 443,
    "ssl_config": "No Verify"
}
CONFIGURATION

}

resource "bindplane_credential" "kubernetes" {
  configuration = <<CONFIGURATION
  {
    "name": "kubernetes",
    "credential_type_id": "${var.bindplane_credential_types["kubernetes"]}",
    "parameters": {
      "bearer_token": "${data.vault_generic_secret.poc_onprem.data["k8s_token"]}"
    }
  }

CONFIGURATION

}

resource "bindplane_source" "kubernetes" {
  provisioning_timeout = 90
  name                 = "kubernetes"
  source_type          = "kubernetes"
  collector_id         = data.vault_generic_secret.poc_onprem.data["onprem_collector_id"]
  collection_interval  = 300
  credential_id        = bindplane_credential.kubernetes.id

  configuration = <<CONFIGURATION
{
    "collect_containers": true,
    "collect_deployments": true,
    "collect_events": false,
    "collect_jobs": true,
    "collect_kubelet_api": true,
    "collect_pods": true,
    "collect_volumes": true,
    "connection_timeout": 30,
    "host": "${data.vault_generic_secret.poc_onprem.data["k8s_host"]}",
    "internal_external_ip_usage": "internal_ip_addresses",
    "max_simultaneous_kubelet_requests": 20,
    "ssl_config": "No Verify"
}
CONFIGURATION

}

resource "bindplane_credential" "db2" {
  configuration = <<CONFIGURATION
{
  "name": "db2",
  "credential_type_id": "${var.bindplane_credential_types["db2"]}",
  "parameters": {
      "username": "${data.vault_generic_secret.poc_onprem.data["db2_user"]}",
      "password": "${data.vault_generic_secret.poc_onprem.data["db2_pass"]}"
  }
}
CONFIGURATION

}

resource "bindplane_source" "db2" {
  provisioning_timeout = 90

  name                = "db2"
  source_type         = "ibm_db2"
  collector_id        = data.vault_generic_secret.poc_onprem.data["onprem_collector_id"]
  collection_interval = 300
  credential_id       = bindplane_credential.db2.id

  configuration = <<CONFIGURATION
{
    "collect_events": true,
    "database_name": "SAMPLE",
    "host": "${data.vault_generic_secret.poc_onprem.data["db2_host"]}",
    "order_queries_by": "AVG_EXEC_TIME",
    "port": 50000,
    "query_count": 10

}
CONFIGURATION

}

