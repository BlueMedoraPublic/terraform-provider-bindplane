/*
required variables
*/

variable "name" {
  description = "eks log source name. will be used to derive the template name."
  type = string
}

variable "destination_config_id" {
  description = "destination configuration id, for the log template"
  type = string
}

variable "source_type_version" {
  description = "eks log source type version"
  type = string
}

/*
variables with default values
*/

variable "container_logs" {
  description = "enable container logs"
  default = true
  type = bool
}

variable "container_file_exp" {
  description = "container logs file path"
  default = "/var/log/containers/*.log"
  type = string
}

variable "kube_proxy_logs" {
  description = "enable kube proxy logs"
  default = true
  type = bool
}

variable "kubelet_logs" {
  description = "enable kubelet logs"
  default = true
  type = bool
}

variable "controller_manager_logs" {
  description = "enable controller logs"
  default = true
  type = bool
}

variable "scheduler_logs" {
  description = "enable scheduler logs"
  default = true
  type = bool
}

variable "apiserver_logs" {
  description = "enable api server logs"
  default = true
  type = bool
}

variable "source_type_id" {
  description = "source type id"
  default = "amazon_eks"
  type = string
}

/*
resources
*/

resource "bindplane_log_source" "eks" {
  name           = var.name
  source_type_id = var.source_type_id
  source_version = var.source_type_version
  configuration  = <<CONFIGURATION
{
    "container_logs": ${var.container_logs},
    "container_file_exp": "${var.container_file_exp}",
    "kube_proxy_logs": ${var.kube_proxy_logs},
    "kubelet_logs": ${var.kubelet_logs},
    "kube_controller_manager_logs": ${var.controller_manager_logs},
    "kube_scheduler_logs": ${var.apiserver_logs},
    "kube_apiserver_logs": ${var.apiserver_logs}
}
CONFIGURATION
}

resource "bindplane_log_template" "eks" {
    name = "${var.name}-template"
    source_config_ids = [
        bindplane_log_source.eks.id
    ]
    destination_config_id = var.destination_config_id
    agent_group = ""
}
