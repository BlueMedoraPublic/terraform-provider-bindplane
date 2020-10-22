// required variables //

variable "project" {
  description = "The Google Cloud Project to deploy to"
}

variable "name" {
  description = "The collector VM name"
}

variable "secret_key" {
  description = "Secret key for agent"
}

// variables with default values //

variable "network_zone" {
  description = "The network zone to deploy to"
  default     = "us-east1-b"
}
