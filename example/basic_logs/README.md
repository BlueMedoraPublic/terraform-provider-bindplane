# logs example

This example will:
- Create a log source config
- Create a log destination  config

## setup

Special care is required when handling the service account used to connect
BindPlane to Stackdriver.

1) use GCP Secret Manager to store your `service_account.json` file
2) edit `secrets.tf` to reflect your secret
```
data "google_secret_manager_secret_version" "bindplane_svc_act" {
  provider = google-beta
  project = "bpcli-dev"
  secret = "bindplane-service-account"
  version = 1
}
```
## usage

Deploy
```
terraform init
terraform apply
```

Cleanup
```
terraform destroy
```
