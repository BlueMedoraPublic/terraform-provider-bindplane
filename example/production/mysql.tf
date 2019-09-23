resource "random_id" "mysql_password" {
  byte_length = 8
}

data "template_file" "mysql" {
  template = file("${path.module}/../../scripts/mysql_install.sh.tpl")

  vars = {
    database   = "mydb"
    mysql_user = data.vault_generic_secret.poc_onprem.data["mysql_user"]
    mysql_pass = data.vault_generic_secret.poc_onprem.data["mysql_pass"]
  }
}

resource "google_compute_instance" "mysql_gcp" {
  count        = var.mysql_instance_count
  name         = "mysql-poc-${count.index}"
  project      = data.vault_generic_secret.poc_secrets.data["gcp_project"]
  machine_type = "f1-micro"
  zone         = "us-west1-a"

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-minimal-1604-lts"
      size  = "30"
      type  = "pd-ssd"
    }
  }

  network_interface {
    network = "default"
    access_config {
    }
  }

  metadata_startup_script = "${data.template_file.mysql.rendered};"
}

resource "bindplane_credential" "mysql_gcp" {
  configuration = <<CONFIGURATION
{
  "name": "mysql",
  "credential_type_id": "${var.bindplane_credential_types["mysql"]}",
  "parameters": {
    "username": "${data.vault_generic_secret.poc_onprem.data["mysql_user"]}",
    "password": "${data.vault_generic_secret.poc_onprem.data["mysql_pass"]}"
  }
}
CONFIGURATION

}

resource "bindplane_source" "mysql_gcp" {
  provisioning_timeout = 180
  count                = var.mysql_instance_count

  name                = element(google_compute_instance.mysql_gcp.*.name, count.index)
  source_type         = "mysql"
  collector_id        = module.gcp_collector.id
  collection_interval = 300
  credential_id       = bindplane_credential.mysql_gcp.id

  configuration = <<CONFIGURATION
{
    "collection_mode": "normal",
    "connection_timeout": 15,
    "host": "${element(
  google_compute_instance.mysql_gcp.*.network_interface.0.network_ip,
  count.index,
)}",
    "monitor_databases": "user",
    "monitor_queries": true,
    "monitor_tables": true,
    "order_queries_by": "avg_latency",
    "order_tablespaces_by": "file_size",
    "port": 3306,
    "query_count": 10,
    "query_history_interval": 24,
    "query_timeout": 5,
    "ssl_config": "No SSL",
    "table_space_count": 200

}
CONFIGURATION

}
