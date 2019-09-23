// variables that are required go here //

variable "resource_group" {
  description = "The Azure resource group to deploy to"
}

variable "bindplane_secret_key" {
  description = "API key for bindplane"
}

variable "collector_name_prefix" {
  description = "The name to be assigned to the collector. Final name will have a random suffix"
}

variable "subnet_id" {
  description = "The subnet id to deploy to"
}

variable "admin_username" {
  description = "The sudo user"
}

// variables with default values go here //

variable "compute_instance_size" {
  description = "Compute instance size"
  default     = "Standard_DS1_v2"
}

variable "compute_image" {
  description = "Compute image to deploy from"

  default = {
    publisher = "OpenLogic"
    offer     = "CentOS"
    sku       = "7-CI"
    version   = "latest"
  }
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
  default     = "Standard_LRS"
}

variable "delete_os_disk_on_termination" {
  description = "Delete the OS disk when the instance is destroyed"
  default     = true
}

variable "location" {
  description = "The Azure region"
  default     = "Central US"
}
