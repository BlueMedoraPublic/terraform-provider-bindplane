# GCP Collector Provisioner

This module is useful for situations where you want to install
a collector on a virtual machine managed outside of the module.

Alternative solution:
- Install the collector during VM deployment using a metadata startup script

## Requirements

Software
- [Terraform 0.13.x](https://www.terraform.io/downloads.html)
- [Google Cloud SDK](https://cloud.google.com/sdk/docs/install)

Environment:
- External SSH access to target VM

## Usage

### Required Parameters

| Parameter | Type | description |  
|-----------|------|-------------|
|project    |string|The Google Cloud Project ID (Required)|
|name       |string|The name of the compute instance that the BindPlane collector should be installed on (Required)|
|instance_id|string|The instance ID of the compute instance the BindPlane collector should be installed on (Required)|

### Example

Assuming the instance already exists, you only need to define
a `gcp_collector_provisioner` module. A full example can be
found in this repo (`example/gcp_collector_provisioner`).

```
variable "project" {}
variable "name" {}
variable "secret_key" {}
variable "instance_id" {}

module "agent" {
    source = "./../../modules/gcp_collector_provisioner"

    zone = us-east1-b
    project = var.project
    agent_secret = var.secret_key
    instance_id = google_compute_instance.collector.instance_id
    instance_name = var.name

}
```

Apply the module:
```
terraform apply \
    -var project=dev-project \
    -var name=dev-instance \
    - var instance_id=<instance id> \
    -var=secret_key=<bindplane_secret-key>
```
