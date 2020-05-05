data "template_file" "fake_key" {
  template = file("fake_gcp_key.json")
}

resource "bindplane_log_destination_google" "default" {
  name = "google-default-${var.name}"
  credentials = data.template_file.fake_key.rendered
}

resource "bindplane_log_destination_google" "optional" {
  name = "google-optional-${var.name}"
  credentials = data.template_file.fake_key.rendered
  location = "us-east1"
  destination_version = "1.3.2" // latest is 1.3.3 as of 04/18/2020
}

resource "bindplane_log_destination" "generic" {
  name = "google-generic-${var.name}"
  destination_type_id = "stackdriver"
  destination_version = "1.3.2"
  configuration = <<CONFIGURATION
{
  "credentials": ${data.template_file.fake_key.rendered},
  "location": "us-west1"
}
CONFIGURATION

}
