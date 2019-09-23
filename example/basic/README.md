# Bindplane Terraform Basic Example

## Description

This basic example will:
- deploy a single bindplane collector
- deploy a single postgresql compute instance
- configure a credential and source for the postgres instance

## Usage

Make sure you are authenticated to GCP with with a service
account or your personal account
```
gcloud auth
```

Set `BINDPLANE_API_KEY`
```
export BINDPLANE_API_KEY=<your api key>
```

Run `terraform` and enter your GCP projectID and Bindplane Secret Key
when prompted for them.

```
terraform init
terraform apply
```

Cleanup with:
```
terraform destroy
```
