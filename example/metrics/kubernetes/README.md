# Kubernetes Metrics

This guide deploys / configures:
- BindPlane collector
- BindPlane metrics least privileged user
- BindPlane Kubernetes metrics source

This guide assumes you have the following:
- A Kubernetes cluster with `kubectl` access
- A BindPlane metrics destination pre configured
- [bpcli installed](https://github.com/BlueMedoraPublic/bpcli)
- [terraform 0.13.x installed](https://www.terraform.io/downloads.html)
- [gomplate installed](https://github.com/hairyhenderson/gomplate)

## How it works

A script `deploy_collector.sh` is included that does the following:

1. Generate BindPlane collector Kubernetes yaml config
2. Deploy collector to Kubernetes
3. Generate and apply a Terraform module for BindPlane Kubernetes metrics

## Steps

Set your environment
```
export NAME="customer-c"
export COLLECTOR_NAME="${NAME}"
export COLLECTOR_SECRET_KEY=""
export BINDPLANE_API_KEY=""
```

Test (output should be empty, or show previously deployed collectors)
```
bpcli collector list
```

Deploy
```
./deploy_collector.sh
```

Verify
```
kubectl get pods
bpcli collector list
bpcli credential list
bpcli source list
```
