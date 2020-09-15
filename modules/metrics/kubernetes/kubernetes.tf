terraform {
  required_providers {
    bindplane = {
      source = "BlueMedoraPublic/bindplane"
      version = "0.2.3"
    }
  }
}

variable "name" {
  description = "Name to use for the credential and source"
  type        = string
}

variable "bearer_token" {
  description = "The service account's bearer token"
  type        = string
}

variable "collector_id" {
  description = "The collector id to use for the source"
  type        = string
}

variable "api_server_address" {
  description = "Hostname or IP addresses of the API server"
  type        = string
}

variable "collection_interval" {
  description = "Collection interval (seconds)"
  type        = number
  default     = 60
}

variable "provisioning_timeout" {
  description = "Max time to wait before timing out the source creation"
  type        = number
  default     = 90
}

variable "collect_containers" {
  description = "Enable container collection"
  type        = bool
  default     = true
}

variable "collect_deployments" {
  description = "Enable deployment collection"
  type        = bool
  default     = true
}

variable "collect_jobs" {
  description = "Enable job collection"
  type        = bool
  default     = true
}

variable "collect_kubelet_api" {
  description = "Enable collection from the Kubelet API(s)"
  type        = bool
  default     = true
}

variable "collect_pods" {
  description = "Enable pod collection"
  type        = bool
  default     = true
}

variable "collect_volumes" {
  description = "Enable volume collection"
  type        = bool
  default     = true
}

variable "connection_timeout" {
  description = "API timeout (seconds)"
  type        = number
  default     = 30
}

variable "internal_external_ip_usage" {
  description = "Connect to kubelet using their internal or external address"
  type        = string
  default     = "internal_ip_addresses"
}

variable "kubelet_ssl_config" {
  description = "Enable TLS verification against the Kubelet API(s)"
  type        = string
  default     = "Verify"
}

variable "max_simultaneous_kubelet_requests" {
  description = "Max number of connections to the kubelet api"
  type        = number
  default     = 20
}

variable "ssl_config" {
  description = "Enable TLS verification against the API server"
  type        = string
  default     = "Verify"
}

variable "api_server_port" {
  description = "API server's port"
  type        = string
  default     = "443"
}

variable "kubelet_port" {
  description = "API server's port"
  type        = string
  default     = "10250"
}

locals {
  credential_type_id = "6945f85d-ae72-4012-8224-6af0784e0b42"
  source_type_id     = "kubernetes"
}

resource "bindplane_credential" "kubernetes" {
  configuration = <<CONFIGURATION
{
  "name": "${var.name}",
  "credential_type_id": "${local.credential_type_id}",
  "parameters": {
    "bearer_token": "${var.bearer_token}"
  }
}
CONFIGURATION

}

resource "bindplane_source" "kubernetes" {
  name                 = var.name
  source_type          = local.source_type_id
  collector_id         = var.collector_id
  collection_interval  = var.collection_interval
  credential_id        = bindplane_credential.kubernetes.id
  provisioning_timeout = var.provisioning_timeout

  configuration = <<CONFIGURATION
{
  "cluster_name": "${var.name}",
  "collect_containers": ${var.collect_containers},
  "collect_deployments": ${var.collect_deployments},
  "collect_jobs": ${var.collect_jobs},
  "collect_kubelet_api": ${var.collect_kubelet_api},
  "collect_pods": ${var.collect_pods},
  "collect_volumes": ${var.collect_volumes},
  "connection_timeout": ${var.connection_timeout},
  "host": "${var.api_server_address}",
  "internal_external_ip_usage": "${var.internal_external_ip_usage}",
  "kubelet_port": "${var.kubelet_port}",
  "kubelet_ssl_config": "${var.kubelet_ssl_config}",
  "max_simultaneous_kubelet_requests": ${var.max_simultaneous_kubelet_requests},
  "port": "${var.api_server_port}",
  "ssl_config": "${var.ssl_config}"
}
CONFIGURATION

}
