#!/bin/bash

# A script to setup a open judge environment.
# It expects to be run on Ubuntu 16.04 via 'sudo'

install_docker() {
    echo "--> Installing docker"
    apt update

    apt-get install -y \
        apt-transport-https \
        ca-certificates \
        curl \
        software-properties-common

    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -

    apt-key fingerprint 0EBFCD88

    add-apt-repository \
    "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
    $(lsb_release -cs) \
    stable"

    apt-get update

    apt-get install -y docker-ce
}

main() {
    CHECK_DOCKER_INSTALLATION=$(dpkg -l | grep -c docker-ce)

    if ! [ "$CHECK_DOCKER_INSTALLATION" -ne 0 ]; then
        install_docker
    else
        echo "--> Docker its already installed"
    fi

    systemctl daemon-reload
    service docker restart
}

main