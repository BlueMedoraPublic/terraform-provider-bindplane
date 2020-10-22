resource "google_compute_instance" "collector-centos" {
  name                      = "${var.name}-centos"
  project                   = var.project
  machine_type              = "g1-small"
  zone                      = var.network_zone
  allow_stopping_for_update = true

  boot_disk {
    initialize_params {
      image = "centos-cloud/centos-7"
      size  = "30"
      type  = "pd-standard"
    }
  }

  network_interface {
    network = "default"

    access_config {
      // leave emtpy in order to get a public ip
    }
  }
}

module "agent-centos" {
    source = "./../../modules/gcp_collector_provisioner"

    agent_secret = var.secret_key
    instance_name = "${var.name}-centos"
    instance_id = google_compute_instance.collector-centos.instance_id
    zone = var.network_zone
    project = var.project


    /* For example purposes, ssh host key checking is
       disabled (set to 'no') because we often get the
       same IP address when rapidly creating / destroying
       the VM resource, causing host key conflicts
    */
    strict_host_key_checking = "no"
}
