/*
Credential types are set here. These can be found using the
bindplane cli
bpcli source type get --id postgresql --json
bpcli source type get --id postgresql --json | jq .credential_types
*/
variable "bindplane_credential_types" {
  description = "credential type IDs"

  default = {
    postgres                        = "6c2c63b5-d465-46ef-bfce-b0881066e43b"
    google_compute                  = "3ba17d17-0380-4c41-9fbb-31fdb18a893a"
    microsoft_azure_virtualmachines = "515a6b7c-a570-4274-becf-ce651a13e281"
    pure_storage                    = "e96817da-9ef9-4fdd-a443-99fd00895ea0"
    kubernetes                      = "99c2617c-add5-46b4-817c-2c7b4c16907e"
    vcenter                         = "64ad46cf-e468-4978-8aee-1acff7fc7bfb"
    mysql                           = "03e7424b-83c5-41ca-8365-2abf528881d5"
    db2                             = "8157ff06-29f0-4f66-a9de-4e1ecb060ac1"
    mssql                           = "e6796295-67c1-4bae-8f38-6b9a2f404fa7"
    cassandra                       = "e781f081-2e65-48c7-81b2-a7ceb823b267"
  }
}

/*
We use Vault to store secrets that should not be stored within the
repository. Postgres username and Google Cloud project are stored
here and referenced by other resources.
*/
data "vault_generic_secret" "poc_secrets" {
  path = "secret/bindplane/bpcli/poc"
}

data "vault_generic_secret" "poc_onprem" {
  path = "secret/bindplane/bpcli/poc_onprem"
}

data "vault_generic_secret" "azure_account" {
  path = "secret/azure/account"
}

variable "postgres_instance_count" {
  default = 15
}

variable "mysql_instance_count" {
  default = 15
}

variable "docker_instance_count" {
  default = 15
}

variable "cassandra_instance_count" {
  default = 15
}

// Instance type for GCP collectors
variable "gcp_instance_type" {
  default = "n1-standard-4"
  //default = "n1-standard-2"
}
