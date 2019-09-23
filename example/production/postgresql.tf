// Generate a random password and use it for postgres password
resource "random_id" "postgres_password" {
  byte_length = 8
}

// This resource represents a rendered template file, used for
// configuring Postgres
data "template_file" "postgres" {
  template = file("${path.module}/../../scripts/postgres_install.sh.tpl")

  vars = {
    username = data.vault_generic_secret.poc_secrets.data["postgres_username"]
    password = random_id.postgres_password.hex
  }
}

// Deploy a Google Compute Instance and configure Postgres
// Vault is used to retrieve secret information such as project name
// scripts/postgres_install.sh.tpl is used for configuration
resource "google_compute_instance" "postgres" {
  count        = var.postgres_instance_count
  name         = "pgsql-poc-${count.index}"
  project      = data.vault_generic_secret.poc_secrets.data["gcp_project"]
  machine_type = "f1-micro"
  zone         = "us-east1-b"

  boot_disk {
    initialize_params {
      image = "centos-cloud/centos-7"
      size  = "30"
      type  = "pd-ssd"
    }
  }

  network_interface {
    network = "default"
    access_config {
    }
  }

  metadata_startup_script = "${data.template_file.postgres.rendered};"
}

// Use the Bindplane provider to define a credential for postgres.
// Username is pulled from Vault and password is generated randomly.
resource "bindplane_credential" "postgres" {
  configuration = <<CONFIGURATION
{
  "name": "postgres",
  "credential_type_id": "${var.bindplane_credential_types["postgres"]}",
  "parameters": {
    "username": "${data.vault_generic_secret.poc_secrets.data["postgres_username"]}",
    "password": "${random_id.postgres_password.hex}"
  }
}
CONFIGURATION

}

// Use the Bindplane provider to define a source for postgres.
// The collector id, credential id, and postgres ip address are
// all dynamically configured.
resource "bindplane_source" "postgres_app_db" {
  provisioning_timeout = 180
  count                = var.postgres_instance_count

  name                = "${element(google_compute_instance.postgres.*.name, count.index)}"
  source_type         = "postgresql"
  collector_id        = module.gcp_collector2.id
  collection_interval = 300
  credential_id       = bindplane_credential.postgres.id

  configuration = <<CONFIGURATION
{
	"collection_mode": "normal",
	"function_count": 20,
    "host": "${element(google_compute_instance.postgres.*.network_interface.0.network_ip, count.index)}",
	"monitor_indexes": false,
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
