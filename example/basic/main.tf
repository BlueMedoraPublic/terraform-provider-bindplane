terraform {
  required_providers {
    bindplane = {
      source = "BlueMedoraPublic/bindplane"
      version = "0.2.5"
    }
  }
}

/*

Deploy a Google Compute Instance and configure Postgres

Vault is used to retrieve secret information such as project name

scripts/postgres_install.sh.tpl is used for configuration

*/
resource "google_compute_instance" "postgres" {
  name         = "postgres-bp-poc-${terraform.workspace}"
  project      = var.project
  machine_type = "f1-micro"
  zone         = "us-west1-b"

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

/*

Deploy a Google Compute Instance and install the Bindplane
collector using a metadata startup script.

This block uses the `gcp_collector` module to abstract the
instance deployment and collector installation.

*/
module "gcp_collector" {
  source = "../../modules/gcp_collector"

  // required parameters //
  project               = var.project
  bindplane_secret_key  = var.secret
  collector_name_prefix = "gcp-poc-${terraform.workspace}"
  network_zone          = "us-west1-a"
}

/*

Use the Bindplane provider to define a credential for postgres.

Username is pulled from Vault and password is generated randomly.

*/
resource "bindplane_credential" "postgres" {
  configuration = <<CONFIGURATION
{
  "name": "postgres",
  "credential_type_id": "6c2c63b5-d465-46ef-bfce-b0881066e43b",
  "parameters": {
    "username": "medora",
    "password": "password"
  }
}
CONFIGURATION

}

/*

Use the Bindplane provider to define a source for postgres.

The collector id, credential id, and postgres ip address are
all dynamically configured.

*/
resource "bindplane_source" "postgres_app_db_0" {
  provisioning_timeout = 120
  name = google_compute_instance.postgres.name
  source_type = "postgresql"
  collector_id = module.gcp_collector.id
  collection_interval = 60
  credential_id = bindplane_credential.postgres.id

  configuration = <<CONFIGURATION
{
    "collection_mode": "normal",
    "function_count": 20,
    "host": "${google_compute_instance.postgres.network_interface[0].network_ip}",
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

/*

This resource represents a rendered template file, used for
configuring Postgres

*/
data "template_file" "postgres" {
template = file("${path.module}/../../scripts/postgres_install.sh.tpl")

vars = {
username = "medora"
password = "password"
}
}

variable "project" {
}

variable "secret" {
}
