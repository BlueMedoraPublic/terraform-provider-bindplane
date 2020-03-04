variable "project" {
  description = "gcp project id"
}

variable "company_id" {
  description = "bindplane company id"
}

variable "secret_key" {
  description = "bindplane secret key"
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

bpcli logs destination type parameters \
  --destination-type-id stackdriver

*/
resource "bindplane_log_destination" "stackdriver" {
  name = "stackdriver-tf-${random_id.suffix.hex}"
  destination_type_id = "stackdriver"
  destination_version = "1.3.2"
  configuration = <<CONFIGURATION
{
  "credentials": ${data.google_secret_manager_secret_version.bindplane_svc_act.secret_data},
  "location": "us-west1"
}
CONFIGURATION

}

/*

bpcli logs source type parameters \
  --source-type-id mysql

*/
variable "error_log_path" {
  description = "example: you can use interpolation for configuraton params"
  default = "/var/log/mysql/mysqld.log"
}

resource "bindplane_log_source" "mysql" {
  name = "mysql-tf-${random_id.suffix.hex}"
  source_type_id = "mysql"
  source_version = "2.0.0"
  configuration = <<CONFIGURATION
{
  "mysql_error": true,
  "error_log_path": "${var.error_log_path}",
  "mysql_slow_query": true,
  "slow_query_log_path": "/var/log/mysql/slow.log",
  "mysql_general": false,
  "general_log_path": "/var/log/mysql/general.log",
  "read_from_head": false
}
CONFIGURATION
}

resource "bindplane_log_template" "mysql" {
    name = "template-tf-${random_id.suffix.hex}"
    source_config_ids = [
        bindplane_log_source.mysql.id
    ]
    destination_config_id =  bindplane_log_destination.stackdriver.id
    agent_group = ""
}

resource "random_id" "mysql_password" {
  byte_length = 8
}

data "template_file" "mysql" {
  template = file("${path.module}/../../scripts/mysql_install_logs.sh.tpl")
  vars = {
    database   = "demo"
    mysql_user = "demo"
    mysql_pass = random_id.mysql_password.hex
    company_id = var.company_id
    secret_key = var.secret_key
    template_id = bindplane_log_template.mysql.id
  }
}

resource "google_compute_instance" "mysql" {
  name         = "mysql-tf-${random_id.suffix.hex}"
  machine_type = "g1-small"
  zone         = "us-central1-a"
  project      = var.project

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-minimal-1804-lts"
    }
  }

  network_interface {
    network = "default"
    access_config {
      // Ephemeral IP
    }
  }

  metadata_startup_script = "${data.template_file.mysql.rendered};"
}

resource "bindplane_log_agent_populate" "mysql" {
    name = google_compute_instance.mysql.name
    // wait up to two minutes for the mysql compute
    // instance to run its metadata startup script, which
    // performs the install
    provisioning_timeout = 180
}
