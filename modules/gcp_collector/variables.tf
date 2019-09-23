// variables that are required go here //

variable "project" {
  description = "The Google Cloud Project to deploy to"
}

variable "bindplane_secret_key" {
  description = "API key for bindplane"
}

variable "collector_name_prefix" {
  description = "The name to be assigned to the collector. Final name will have a random suffix"
}

variable "network_zone" {
  description = "The network zone to deploy to"
}

// variables with default values go here //

variable "compute_instance_size" {
  description = "Compute instance size: default to 1c4gb "
  default     = "n1-standard-1"
}

variable "compute_image" {
  description = "Compute image to deploy from"
  default     = "centos-cloud/centos-7"
}

variable "userdata_script" {
  description = "The userdata script to be used for collector installation"
  default     = "userdata.sh.tpl"
}

variable "network" {
  description = "The network to deploy to"
  default     = "default"
}

variable "disk_type" {
  description = "The persistent disk type to use"
  default     = "pd-standard"
}

variable "disk_size" {
  description = "The disk size to use, in Gigabytes"
  default     = "80"
}
