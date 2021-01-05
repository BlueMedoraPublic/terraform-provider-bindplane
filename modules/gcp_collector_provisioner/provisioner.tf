terraform {
  required_providers {
    bindplane = {
      source = "BlueMedoraPublic/bindplane"
      version = "0.2.5"
    }
  }
}

variable "agent_secret" {
  description = "API key for bindplane"
}

variable "instance_name" {
  description = "The name to be assigned to the agent. Final name will have a random suffix"
}

variable "instance_id" {
  description = "The instance id"
}

variable "zone" {
  description = "The instance's compute zone"
}

variable "project" {
  description = "The insdtance's project"
}

variable "strict_host_key_checking" {
  description = "Strict host key checking for SSH commands. Set to 'no' to disable host key checking. https://cloud.google.com/sdk/gcloud/reference/compute/ssh"
  default = "yes"
}

resource "random_id" "suffix" {
  byte_length = 5
}

resource "null_resource" "agent" {
  // re-run this script if the instance id changes
  triggers = {
    instance_id = var.instance_id
  }

  provisioner "local-exec" {
    working_dir = path.module
    command = <<EOF
agent_secret=${var.agent_secret} \
instance_name=${var.instance_name} \
agent_name=${var.instance_name}-${random_id.suffix.hex} \
zone=${var.zone} \
project=${var.project} \
strict_host_key_checking=${var.strict_host_key_checking} \
./provision.sh
EOF
  }
}

resource "bindplane_collector" "agent" {
  name       = "${var.instance_name}-${random_id.suffix.hex}"
  depends_on = [null_resource.agent]
}
