Usage
==================

- [Requirements](#requirements)
- [Terraform Apply](#terraform-apply)
- [Example Configuration](#example-configuration)

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12.x
- `terraform-provider-bindplane` (see README.md Installation)
- [bpcli](https://github.com/BlueMedoraPublic/bpcli) (optional)

`bpcli` is optional, but useful for retrieving required information
such as configuration requirements or source type ids.

Terraform Apply
------------

Using the provider is very simple, the only requirement is
setting `BINDPLANE_API_KEY` in the execution environment

```sh
export BINDPLANE_API_KEY=<your api key>
terraform init
terraform plan
terraform apply
```


Example Configuration
------------

This example is not designed for production use. Detailed
examples can be found in the `examples/` directory of this repo.


Create a Postgres `bindplane_credential` resource
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

Create a Postgres `bindplane_source` resource
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

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>
