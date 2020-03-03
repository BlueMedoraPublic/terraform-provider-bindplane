# logs example

## goal

1) define a source config
2) define a destination config
3) create a template that is paired with the source and destination config
4) deploy a compute instance with the log agent pre configured with the log template

## setup

Special care is required when handling the service account used to connect
BindPlane to Stackdriver.

1) use GCP Secret Manager to store your `service_account.json` file
2) edit `secrets.tf` to reflect your secret
```
data "google_secret_manager_secret_version" "bindplane_svc_act" {
  provider = google-beta
  project = var.project
  secret = "bindplane-service-account"
  version = 1
}
```

## usage

Set env
```
export BINDPLANE_API_KEY=<your api key>
export PROJECT=<your gcp project id>
```

Deploy
```
terraform init
terraform apply -var "project=${PROJECT}"
```

Cleanup
```
terraform destroy -var "project=${PROJECT}"
```
