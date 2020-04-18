variable "project" {
  description = "gcp project id"
}

data "google_secret_manager_secret_version" "bindplane_svc_act" {
  provider = google-beta
  project = var.project
  secret = "bindplane-service-account"
  version = 1
}

// keep resource names random
resource "random_id" "suffix" {
  byte_length = 2
}

/*

use the google resource destination to define stackdriver with
only required params

*/
resource "bindplane_log_destination_google" "default" {
  name = "google-default-${random_id.suffix.hex}"
  credentials = data.google_secret_manager_secret_version.bindplane_svc_act.secret_data
}

/*

use the google resource destination to define stackdriver with
required and optional params

*/
resource "bindplane_log_destination_google" "optional" {
  name = "google-optional-${random_id.suffix.hex}"
  credentials = data.google_secret_manager_secret_version.bindplane_svc_act.secret_data
  location = "us-east1"
  destination_version = "1.3.2" // latest is 1.3.3 as of 04/18/2020
}

/*

use the genric destination to define stackdriver

*/
resource "bindplane_log_destination" "generic" {
  name = "google-generic-${random_id.suffix.hex}"
  destination_type_id = "stackdriver"
  destination_version = "1.3.2"
  configuration = <<CONFIGURATION
{
  "credentials": ${data.google_secret_manager_secret_version.bindplane_svc_act.secret_data},
  "location": "us-west1"
}
CONFIGURATION

}
