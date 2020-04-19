#!/bin/bash

set -eE  # same as: `set -o errexit -o errtrace`

cd "$(dirname "$0")"

if [ -z ${BINDPLANE_API_KEY+x} ]; then echo "BINDPLANE_API_KEY is unset"; exit 1; fi
if [ -z ${BINDPLANE_COMPANY_ID+x} ]; then echo "BINDPLANE_COMPANY_ID is unset"; exit 1; fi
if [ -z ${BINDPLANE_SECRET_KEY+x} ]; then echo "BINDPLANE_SECRET_KEY is unset"; exit 1; fi
if [ -z ${COLLECTOR_SECRET_KEY+x} ]; then echo "COLLECTOR_SECRET_KEY is unset"; exit 1; fi

# globals
UNIX_TIME=$(date +%s)

# collector
COLLECTOR_NAME="intigration-test-${UNIX_TIME}"
COLLECTOR_UUID=$(uuidgen)
API_ADDRESS="https://production.api.bindplane.bluemedora.com"
BINDPLANE_HOME="/opt/bluemedora/bindplane-collector"

# log agent
AGENT_NAME="intigration-test-log-${UNIX_TIME}"

# postgres cred / source
POSTGRES_USER="postgres"
POSTGRES_PASSWORD="password"
SOURCE_NAME="${UNIX_TIME}-psql-integration"
SOURCE_FILE="${SOURCE_NAME}.source"
CRED_NAME="${UNIX_TIME}-psql-integration"
CRED_FILE="${CRED_NAME}.cred"

clean () {
    # cleanup generated files
    rm -f postgres_source.tf collector.tf

    bpcli source list --json | jq ".[] | select(.name | contains(\"${UNIX_TIME}\"))" | jq .id | xargs --no-run-if-empty -n1 bpcli source delete --id
    bpcli credential list | grep $UNIX_TIME | awk '{print $2}' | xargs --no-run-if-empty -n1 bpcli credential delete --id
    bpcli collector list | grep $UNIX_TIME | awk '{print $2}' | xargs -n1 --no-run-if-empty -n1 bpcli collector delete --id

    bpcli logs agent list | grep $UNIX_TIME | awk '{print $2}' | xargs --no-run-if-empty -n1 bpcli logs agent delete --agent-id
    bpcli logs template list | grep $UNIX_TIME | awk '{print $2}' | xargs --no-run-if-empty -n1 bpcli logs template delete --template-id
    bpcli logs destination config list | grep $UNIX_TIME | awk '{print $2}' | xargs --no-run-if-empty -n1 bpcli logs destination config delete --config-id
    bpcli logs source config list | grep $UNIX_TIME | awk '{print $2}' | xargs --no-run-if-empty -n1 bpcli logs source config delete --config-id

    docker ps | grep $UNIX_TIME | awk '{print $1}' | xargs --no-run-if-empty -I{} docker rm -f {} >> /dev/null

    rm -f terraform.tfstate terraform.tfstate.*
}
trap clean ERR

docker_psql () {
    docker run \
        -d \
        -p 5432:5432 \
        --name=$SOURCE_NAME \
        -e "POSTGRES_PASSWORD=${POSTGRES_PASSWORD}" \
        postgres:9.6
    sleep 15
}

docker_bindplane_collector () {
    docker run -d --name=$COLLECTOR_NAME \
        --entrypoint="/opt/bluemedora/bindplane-collector/scripts/run_collector_in_docker.sh" \
        -e "COLLECTOR_NAME=${UNIX_TIME}" \
        -e "COLLECTOR_UUID=${COLLECTOR_UUID}" \
        -e "COLLECTOR_SECRET_KEY=${COLLECTOR_SECRET_KEY}" \
        -e "API_ADDRESS=${API_ADDRESS}" \
        -e "BINDPLANE_HOME=${BINDPLANE_HOME}" \
        docker.io/bluemedora/bindplane-metrics-collector:latest

    t=0
    max=12
    while :
    do
        echo "looking for collector with uuid ${COLLECTOR_UUID}"
        bpcli collector list | grep ${COLLECTOR_UUID} && break

        echo "waiting for collector. . ."
        if [ "$t" != "$max" ]
        then
            t=$((t+1))
            echo "waiting for api to return the collector: ${COLLECTOR_UUID}"
            sleep 5
        else
            echo "could not find collector ${UNIX_TIME} with UUID ${COLLECTOR_UUID} after 60 seconds"
            clean
        fi
    done
}

docker_bindplane_log_agent () {
    docker run -d --name=$AGENT_NAME \
        -e "AGENT_NAME=${UNIX_TIME}" \
        -e "COMPANY_ID=${BINDPLANE_COMPANY_ID}" \
        -e "SECRET_KEY=${BINDPLANE_SECRET_KEY}" \
        docker.io/bluemedora/bindplane-log-agent:0.8.0

    AGENT_UUID=$(bpcli logs agent list | grep $AGENT_NAME | awk '{print $2}')
}

templates () {
    # set collector name dynamically
    sed "s/NAME/${UNIX_TIME}/g" collector.tf.tpl > collector.tf

    # set postgres address dynamically
    POSTGRES_ADDR=$(docker inspect --format '{{ .NetworkSettings.IPAddress }}' $SOURCE_NAME)
    sed "s/ADDR/${POSTGRES_ADDR}/g" postgres_source.tf.tpl > postgres_source.tf
}

test () {
    VAR_ARGS="-var name=$UNIX_TIME -var secret=${BINDPLANE_SECRET_KEY} -var collector-id=${COLLECTOR_UUID}"

    terraform init
    terraform validate

    # import the docker deployed collector in order to test deletion from the API
    # disabled for now but should be enabled later after this bug is fixed
    #  https://github.com/BlueMedoraPublic/terraform-provider-bindplane/issues/10
    # terraform import $VAR_ARGS bindplane_collector.collector ${COLLECTOR_UUID}

    terraform apply $VAR_ARGS -auto-approve
    terraform destroy $VAR_ARGS -auto-approve
}


clean
docker_psql
docker_bindplane_collector
docker_bindplane_log_agent
templates
test
clean
