# logs example

## goal

1) define a source config
2) define a destination config
3) deploy a compute instance with and install the log agent
4) deploy source and destination configs to the log agent

## setup

Special care is required when handling the service account used to connect
BindPlane to Stackdriver.

Use GCP Secret Manager to store your `service_account.json` file.
Terraform retrieves the secret at runtime
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
export BINDPLANE_API_KEY=    <your api key>
export BINDPLANE_COMPANY_ID= <your company idd>
export BINDPLANE_SECRET_KEY= <your secret key>

export PROJECT=<your gcp project id>
```

Deploy
```
terraform init
terraform apply \
    -var "project=${PROJECT}" \
    -var "company_id=${BINDPLANE_COMPANY_ID}" \
    -var "secret_key=${BINDPLANE_SECRET_KEY}"
```

Cleanup
```
terraform destroy \
    -var "project=${PROJECT}" \
    -var "company_id=${BINDPLANE_COMPANY_ID}" \
    -var "secret_key=${BINDPLANE_SECRET_KEY}"
```
