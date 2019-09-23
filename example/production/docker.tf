data "template_file" "docker" {
  template = file("${path.module}/../../scripts/docker_install.sh.tpl")
}

resource "google_compute_instance" "docker" {
  count        = var.docker_instance_count
  name         = "docker-poc-${count.index}"
  project      = data.vault_generic_secret.poc_secrets.data["gcp_project"]
  machine_type = "f1-micro"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-minimal-1604-lts"
      size  = "10"
      type  = "pd-ssd"
    }
  }

  network_interface {
    network = "default"
    access_config {
    }
  }

  metadata_startup_script = "${data.template_file.docker.rendered};"
}

resource "bindplane_source" "docker" {
  provisioning_timeout = 180
  count                = "${var.docker_instance_count}"

  name                = "${element(google_compute_instance.docker.*.name, count.index)}"
  source_type         = "docker"
  collector_id        = "${module.gcp_collector3.id}"
  collection_interval = 300
  credential_id       = ""

  configuration = <<CONFIGURATION
{
    "collect_events": true,
    "connection_timeout": 10,
    "host_list": "${element(google_compute_instance.docker.*.network_interface.0.network_ip, count.index)}",
    "max_events": 10000,
    "maximum_concurrent_hosts": 10,
    "port_list": "8080"
}
CONFIGURATION
}
