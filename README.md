terraform-provider-bindplane
==================

* [bindplane.bluemedora.com](https://bindplane.bluemedora.com)
* [terraform.io](https://www.terraform.io)
* [Bindplane API Documentation](https://docs.bindplane.bluemedora.com/reference#introduction)

[![Build Status](https://travis-ci.com/BlueMedoraPublic/terraform-provider-bindplane.svg?branch=master)](https://travis-ci.com/BlueMedoraPublic/terraform-provider-bindplane)
[![Go Report Card](https://goreportcard.com/badge/github.com/BlueMedoraPublic/terraform-provider-bindplane)](https://goreportcard.com/report/github.com/BlueMedoraPublic/terraform-provider-bindplane)

Usage Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x

Installation
------------

1) download the latest release for your platform
2) unzip the plugin
3) copy plugin to `~/.terraform.d/plugins` For Mac / Linux and `%APPDATA%\terraform.d\plugins` for Windows

Example Code
------------

See `USAGE.md` and `examples/` for detailed examples

Create a `bindplane_credential` resource
```terraform
resource "bindplane_credential" "postgres" {
  configuration = <<CONFIGURATION
{
  "name": "postgres",
  "credential_type_id": "6c2c63b5-d465-46ef-bfce-b0881066e43b",
  "parameters": {
    "username": "medora",
    "password": "password"
  }
}
CONFIGURATION

}
```

Create a `bindplane_source` resource
```terraform
resource "bindplane_source" "postgres_app_db_0" {
  provisioning_timeout = 120
  name = google_compute_instance.postgres.name
  source_type = "postgresql"
  collector_id = module.gcp_collector.id
  collection_interval = 60
  credential_id = bindplane_credential.postgres.id

  configuration = <<CONFIGURATION
{
    "collection_mode": "normal",
    "function_count": 20,
    "host": "${google_compute_instance.postgres.network_interface[0].network_ip}",
    "monitor_indexes": true,
    "monitor_sequences": false,
    "monitor_sessions": false,
    "monitor_tables": false,
    "monitor_triggers": false,
    "order_functions_by": "calls",
    "order_queries_by": "calls",
    "port": 5432,
    "query_count": 10,
    "show_query_text": true,
    "ssl_config": "No SSL"
}
CONFIGURATION

}
```

Usage
------------

Using the provider is very simple, the only requirement is
setting `BINDPLANE_API_KEY` in the execution environment

```sh
export BINDPLANE_API_KEY=<your api key>
terraform init
terraform plan
terraform apply
```

See `USAGE.md` and `examples/` for detailed examples

Building The Provider
---------------------

Install the following:
- [Make](https://www.gnu.org/software/make/)
- [Docker](https://www.docker.com/) 18.x (primary build method)
- [Go](https://golang.org/doc/install) 1.11+ (alternative build method)

Clone repository anywhere on your system (outside of your GOPATH),
this repository uses go modules, and does not need to be in the GOPATH

Enter the provider directory and build the provider with Docker

```sh
make test
make
```

Build artifacts can be found in the `artifacts/` directory

If you wish to build without Docker
```sh
make quick
```
