variable "project" {
  description = "GCP Project"
}

data "google_secret_manager_secret_version" "bindplane_svc_act" {
  provider = google-beta
  project = var.project
  secret = "bindplane-service-account"
  version = 1
}

resource "bindplane_log_destination" "stackdriver" {
  name = "stackdriver-tf-eks-demo"
  destination_type_id = "stackdriver"
  destination_version = "1.3.2"
  configuration = <<CONFIGURATION
{
  "credentials": ${data.google_secret_manager_secret_version.bindplane_svc_act.secret_data},
  "location": "us-west1"
}
CONFIGURATION

}

module "eks" {
  source = "../../modules/logs/eks_logs"
  name = "eks-demo"
  source_type_version = "2.0.0"
  destination_config_id = bindplane_log_destination.stackdriver.id
}

output "eks" {
  value = module.eks
}
