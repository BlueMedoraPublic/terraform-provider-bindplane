data "template_file" "cassandra" {
  template = file("${path.module}/../../scripts/cassandra_install.sh.tpl")
}

resource "google_compute_instance" "cassandra" {
  count        = var.cassandra_instance_count
  name         = "cassandra-poc-${count.index}"
  project      = data.vault_generic_secret.poc_secrets.data["gcp_project"]
  machine_type = "n1-standard-1"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "centos-cloud/centos-7"
      size  = "40"
      type  = "pd-ssd"
    }
  }

  network_interface {
    network = "default"
    access_config {
    }
  }

  metadata_startup_script = "${data.template_file.cassandra.rendered};"
}

resource "bindplane_credential" "cassandra" {
  configuration = <<CONFIGURATION
{
  "name": "cassandra",
  "credential_type_id": "${var.bindplane_credential_types["cassandra"]}",
  "parameters": {
    "password": "cassandra",
    "username": "cassandra"
  }
}
CONFIGURATION

}


resource "bindplane_source" "cassandra" {
  provisioning_timeout = 260
  count                = "${var.cassandra_instance_count}"

  name                = "${element(google_compute_instance.cassandra.*.name, count.index)}"
  source_type         = "apache_cassandra"
  collector_id        = "${module.gcp_collector.id}"
  collection_interval = 300
  credential_id       = bindplane_credential.cassandra.id

  configuration = <<CONFIGURATION
{
  "host": "${element(google_compute_instance.cassandra.*.network_interface.0.network_ip, count.index)}",
  "port": 7199
}
CONFIGURATION
}
