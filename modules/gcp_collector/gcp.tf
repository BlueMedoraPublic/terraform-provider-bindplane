terraform {
  required_providers {
    bindplane = {
      source = "BlueMedoraPublic/bindplane"
      version = "0.2.5"
    }
  }
}

resource "random_id" "suffix" {
  byte_length = 5
}

data "template_file" "user_data" {
  template = file("${path.module}/../../scripts/${var.userdata_script}")

  vars = {
    api_key        = var.bindplane_secret_key
    name = "${var.collector_name_prefix}-${random_id.suffix.hex}"
  }
}

resource "google_compute_instance" "collector" {
  name                      = "${var.collector_name_prefix}-${random_id.suffix.hex}"
  project                   = var.project
  machine_type              = var.compute_instance_size
  zone                      = var.network_zone
  allow_stopping_for_update = true

  boot_disk {
    initialize_params {
      image = var.compute_image
      size  = var.disk_size
      type  = var.disk_type
    }
  }

  network_interface {
    network = var.network

    access_config {
      // leave emtpy in order to get a public ip
    }
  }

  metadata_startup_script = "${data.template_file.user_data.rendered};"
}

resource "bindplane_collector" "collector" {
  name       = "${var.collector_name_prefix}-${random_id.suffix.hex}"
  depends_on = [google_compute_instance.collector]
}
