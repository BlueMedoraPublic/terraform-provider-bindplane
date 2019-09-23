#!/bin/bash

install () {
    yum install curl -y
    SECRET_KEY=${api_key} sh -c "$(curl -fsSl https://bindplane.bluemedora.com/collectors/unix_install.sh)" unix_install.sh -n ${name} -a y
}

[ ! -d "/opt/bluemedora" ] && install
