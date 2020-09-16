#!/bin/bash

set -e

if [ -z ${COLLECTOR_NAME+x} ]; then echo "COLLECTOR_NAME is unset"; exit 1; fi
if [ -z ${COLLECTOR_SECRET_KEY+x} ]; then echo "COLLECTOR_SECRET_KEY is unset"; exit 1; fi
if [ -z ${BINDPLANE_API_KEY+x} ]; then echo "BINDPLANE_API_KEY is unset"; exit 1; fi

COLLECTOR_UUID="$(uuidgen | tr '[:upper:]' '[:lower:]')"
export COLLECTOR_UUID

API_SERVER_ADDRESS=$(kubectl config view --minify | grep server | cut -f 2- -d ":" | tr -d " " | rev | cut -c5- | rev | cut -c9-)
export API_SERVER_ADDRESS

collector_template="templates/collector.yaml.template"
terraform_template="templates/terraform.tf.template"

render_config() {
    gomplate \
        --file=$collector_template \
        --out="${COLLECTOR_NAME}.yaml"
}

deploy_collector() {
    kubectl apply -f "${COLLECTOR_NAME}.yaml"

    t=0
    max=15
    while :
    do
        echo "looking for collector with uuid ${COLLECTOR_UUID}"
        bpcli collector list | grep ${COLLECTOR_UUID} && break

        echo "waiting for collector. . ."
        if [ "$t" != "$max" ]; then
            t=$((t+1))
            echo "waiting for api to return the collector: ${COLLECTOR_UUID}"
            sleep 6
        else
            echo "could not find collector ${UNIX_TIME} with UUID ${COLLECTOR_UUID} after 60 seconds"
            kubectl delete -f "${COLLECTOR_NAME}.yaml"
            exit 1
        fi
    done
}

render_terraform() {
    BEARER_TOKEN=$(kubectl -n kube-system get secret -o json $(kubectl -n kube-system get secret | grep bluemedora-monitoring-token | awk '{print $1}') | jq -r .data.token | base64 -d)
    export BEARER_TOKEN

    gomplate \
        --file=$terraform_template \
        --out="${COLLECTOR_NAME}.tf"
}

apply_terraform() {
    terraform init
    terraform apply
}



render_config
deploy_collector
render_terraform
apply_terraform
