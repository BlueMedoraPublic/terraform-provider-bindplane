// Deploy a Google Compute Instance and install the Bindplane
// collector using a metadata startup script.
// This block uses the `gcp_collector` module to abstract the
// instance deployment and collector installation.
module "gcp_collector" {
  source = "../../modules/gcp_collector"

  // required parameters //
  project               = data.vault_generic_secret.poc_secrets.data["gcp_project"]
  bindplane_secret_key  = data.vault_generic_secret.poc_secrets.data["secret_key"]
  collector_name_prefix = "gcp-poc-${terraform.workspace}"
  network_zone          = "us-west1-a"

  // optional
  compute_instance_size = var.gcp_instance_type
  disk_type             = "pd-ssd"
}

module "gcp_collector2" {
  source = "../../modules/gcp_collector"

  // required parameters //
  project               = data.vault_generic_secret.poc_secrets.data["gcp_project"]
  bindplane_secret_key  = data.vault_generic_secret.poc_secrets.data["secret_key"]
  collector_name_prefix = "gcp-poc2-${terraform.workspace}"
  network_zone          = "us-east1-b"

  // optional
  compute_instance_size = var.gcp_instance_type
  disk_type             = "pd-ssd"
}

module "gcp_collector3" {
  source = "../../modules/gcp_collector"

  // required parameters //
  project               = data.vault_generic_secret.poc_secrets.data["gcp_project"]
  bindplane_secret_key  = data.vault_generic_secret.poc_secrets.data["secret_key"]
  collector_name_prefix = "gcp-poc3-${terraform.workspace}"
  network_zone          = "us-central1-a"

  // optional
  compute_instance_size = var.gcp_instance_type
  disk_type             = "pd-ssd"
}
