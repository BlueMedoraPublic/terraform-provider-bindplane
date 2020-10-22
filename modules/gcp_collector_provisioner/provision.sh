#!/bin/bash

set -e

if [ -z ${agent_secret+x} ]; then echo "agent_secret is unset"; exit 1; fi
if [ -z ${agent_name+x} ]; then echo "agent_name is unset"; exit 1; fi
if [ -z ${zone+x} ]; then echo "zone is unset"; exit 1; fi
if [ -z ${project+x} ]; then echo "project is unset"; exit 1; fi
if [ -z ${instance_name+x} ]; then echo "instance_name is unset"; exit 1; fi
if [ -z ${strict_host_key_checking+x} ]; then echo "strict_host_key_checking is unset"; exit 1; fi

ssh_cmd_args="${instance_name} --project=${project} --zone=${zone} --strict-host-key-checking=${strict_host_key_checking} --force-key-file-overwrite"

ssh_port=22

wait_for_ssh () {
  ip=$(gcloud compute instances list --project=$project | awk '/'$instance_name'/ {print $5}')

  t=0
  max=12
  while :
  do
    nc -w 1 -z $ip $ssh_port && break

    if [ "$t" != "$max" ]
    then
      t=$((t+1))
      echo "Waiting for instance SSH. . ."
      sleep 5
    else
      echo "Timed out waiting for SSH on instance: ${instance_name}, ip ${ip}, port: ${ssh_port}"
      exit 1
    fi
  done
}

get_install_script () {
  gcloud compute ssh $ssh_cmd_args \
    --command="curl -s https://bindplane.bluemedora.com/collectors/unix_install.sh -o agent.sh && chmod +x agent.sh"
}

install_agent () {
  gcloud compute ssh $ssh_cmd_args \
    --command="if systemctl status bindplane-collector; then sudo systemctl stop bindplane-collector; fi"

  gcloud compute ssh $ssh_cmd_args \
    --command="sudo ./agent.sh --secret-key ${agent_secret} -n ${agent_name} -a y"
}

wait_for_ssh
get_install_script
install_agent
