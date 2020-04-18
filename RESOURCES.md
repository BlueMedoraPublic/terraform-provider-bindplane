# Resources and Data Providers

- [Metrics](#metrics)
  * [bindplane_credential](#bindplane-credential)
  * [bindplane_source](#bindplane-source)
- [Logs](#logs)
  * [bindplane_log_source](#bindplane-log-source)
  * [bindplane_log_destination](#bindplane-log-destination)
  * [bindplane_log_destination_google](#bindplane-log-destination-google)
  * [bindplane_log_bind_source](#bindplane-log-bind-source)
  * [bindplane_log_bind_destination](#bindplane-log-bind-destination)
  * [bindplane_log_template](#bindplane-log-template)
  * [bindplane_log_agent_populate](#bindplane-log-agent-populate)
  * [bindplane_agent_install_cmd](#bindplane-agent-install-cmd)

## Metrics

### bindplane_credential

Metric Credentials represent a credential configuration that
can be paired with a Metric Source, for monitoring your environment.

Credential templates can be retrieved with bpcli:
```
CREDENTIAL_ID=$(bpcli source type get --id postgresql --json | jq '.credential_types[0].id')
bpcli credential type template --json --id $CREDENTIAL_ID
```

The resource takes a single json configuration because each
credential type is different:
```terraform
resource "bindplane_credential" "postgres" {
  configuration = <<CONFIGURATION
{
  "name": "postgres",
  "credential_type_id": "6c2c63b5-d465-46ef-bfce-b0881066e43b",
  "parameters": {
    "username": "medora",
    "password": "${var.password}"
  }
}
CONFIGURATION
}
```

Notice that the password field uses a variable to dynamically
insert a password. You can use a secret manager such as Vault
or GCP Secret Manager to handle your sensitive information. It
is not recommended to hardcode your passwords within the configuration,
but nothing is stopping you from doing so.

It is important to understand the nature of Terraform's state
management. All resources are stored in plain text. Make sure you
use a secure backend configuration when handling credentials.

An example of handling credentials properly can be found in the
examples directory of this repo: `examples/production/postgresq.tf`

[Terraform Backends](https://www.terraform.io/docs/backends/index.html)

### bindplane_source

Metric Sources represent a source configuration for monitoring
your environment. You should familiarize yourself with the process
using the web interface before moving onto Terraform. It is important
to understand that a Test Connection is performed everytime a
bindplane_source is created / modified. The operator should have
an understanding of how long their sources take to pass a test connection
in order to set the `provisioning_timeout` correctly.

Source templates can be retrieved with bpcli:
```
bpcli source type template --id postgresql --json
```

The resource takes a number of parameters as well as an arbitrary
json configuration:
```terraform
resource "bindplane_source" "postgres_app_db_0" {
  provisioning_timeout = 120
  name = google_compute_instance.postgres.name
  source_type = "postgresql"
  collector_id = 6c2c63b5-d465-46ef-bfce-b0881066e43b
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

Notice that we are inserting some parameters using Terraform:
- name
- credential_id
- configuration.host

You can see a full example in `examples/basic`


## Logs

### bindplane_log_source

Log source represents a log configuration that can be paired
with a logging agent. One log source can be used for many agents.

Log source parameters can be found with bpcli:
```
bpcli logs source type parameters --type-id mysql --json | jq
```

The resource takes a number of parameters as well as an arbitrary
json configuration:
```terraform
resource "bindplane_log_source" "mysql" {
  name = "mysql"
  source_type_id = "mysql"
  source_version = "2.0.0"
  configuration = <<CONFIGURATION
{
  "mysql_error": true,
  "error_log_path": "/var/log/mysql/mysqld.log",
  "mysql_slow_query": true,
  "slow_query_log_path": "/var/log/mysql/slow.log",
  "mysql_general": false,
  "general_log_path": "/var/log/mysql/general.log",
  "read_from_head": false
}
CONFIGURATION
}
```

### bindplane_log_destination

Log destination represents a log destination configuration that
can be paired with a logging agent. One log destination can be used
for many agents. This configuration is generic, and will work with all
current and future log destinations supported by BindPlane.

Log source parameters can be found with bpcli:
```
bpcli logs destination type parameters \
    --type-id stackdriver \
    --json | jq
```

```terraform
resource "bindplane_log_destination" "stackdriver" {
  name = "stackdriver"
  destination_type_id = "stackdriver"
  destination_version = "1.3.2"
  configuration = <<CONFIGURATION
{
  "credentials": ${data.google_secret_manager_secret_version.bindplane_svc_act.secret_data},
  "location": "us-west1"
}
CONFIGURATION
}
```

Notice that the configuration json dynamically inserts a
Google Cloud Service Account Json key using google_secret_manager_secret_version
resource from the Google provider.

It is important to understand the nature of Terraform's state
management. All resources are stored in plain text. Make sure you
use a secure backend configuration when handling credentials.

`google_secret_manager_secret_version` examples can be found here:
- `examples/basic_logs`
- `examples/basic_logs/eks`

[Terraform Backends](https://www.terraform.io/docs/backends/index.html)


### bindplane_log_destination_google

Google logging destination represents a simplified implementation of `bindplane-log-destination`.
It sets sane defaults and requires only two parameters.

```terraform
resource "bindplane_log_destination_google" "default" {
  name = "google"
  credentials = data.google_secret_manager_secret_version.bindplane_svc_act.secret_data
}
```

### bindplane_log_bind_source

Source configurations can be deployed to a log agent

```terraform
resource "bindplane_log_bind_source" "mysql" {
  source_config_id = bindplane_log_source.mysql.id
  agent_id         = bindplane_log_agent_populate.mysql.id
}
```

### bindplane_log_bind_destination

A Destination configuration can be deployed to a log agent

```terraform
resource "bindplane_log_bind_destination" "stackdriver" {
    destination_config_id = bindplane_log_destination.stackdriver.id
    agent_id              = bindplane_log_agent_populate.mysql.id
}
```

### bindplane_log_template

Log templates represent a log template configuration that can
be paired with a logging agent. One log template can be used
for many agents. Log templates combine one or more source configurations
and a single destination configuration. Templates are a powerful
way to handle automating your BindPlane Logs deployment.

The resource takes a name, source config id list, destination
config id, and an agent group (optional)
```terraform
resource "bindplane_log_template" "mysql" {
  name = "template-tf-${random_id.suffix.hex}"
  source_config_ids = [
    bindplane_log_source.mysql.id
  ]
  destination_config_id =  bindplane_log_destination.stackdriver.id
  agent_group = ""
}
```

NOTE: `agent_group` is a required parameter but can be left as
an empty string.

### bindplane_log_agent_populate

The `bindplane_log_agent_populate` can be used as a way to manage
cleaning up log agents when an instance is destroyed. This resource
is not a deployment method.

Log agents are deployed in several ways:
- Manually using the command retrieved from the Bindplane UI
- As a container in Docker or Kubernetes using a deployment yaml
- As part of a startup script

This example will show how to use the `bindplane_log_agent_populate`
on Google Cloud with a Compute instance running mysql. The script
in use can be found in this repo at `scripts/mysql_install_logs.sh.tpl`.
```terraform
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
  name         = "mysql"
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
  // wait for the mysql compute instance to run its metadata
  // startup script, which performs the agent install
  provisioning_timeout = 300
}
```

Terraform will do the following in order:
1) render the script by dynamically inserting a password,
company_id, secret_key, template_id
2) Deploy an Ubuntu instance with a startup script. This script
installs a BindPlane log agent with a template_id and mysql.
3) `bindplane_log_agent_populate` resource waits up to 300 seconds
for a log agent named `mysql` to populate in the API.

This might seem pointless, however, when it is time to destroy
the environment, Terraform will delete the log agent from BindPlane
after deleting the compute instance. Think of `bindplane_log_agent_populate`
resource as a post deletion hook for cleaning up old log agents.

NOTE: This resource will fail if it finds multiple log agents
with the same name. It is a good idea to deploy compute instances
with random hostname suffixes in order to guarantee uniquely named
log agents.

### bindplane_agent_install_cmd

This data resource can be used to retrieve a bindplane log
agent install command.

```
bpcli logs agent install-cmd --platform centos7
```
```terraform
resource "bindplane_agent_install_cmd" "centos7" {
  platform = "centos7"
}

output "centos7_agent_install_cmd" {
  value = bindplane_agent_install_cmd.centos7.command
}
```

<small><i><a href='http://ecotrust-canada.github.io/markdown-toc/'>Table of contents generated with markdown-toc</a></i></small>
