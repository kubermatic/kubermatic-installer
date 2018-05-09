#!/bin/bash

# See https://kubernetes.io/docs/setup/independent/install-kubeadm/

set -euo pipefail

DOCKER_VERSION="17.03.0~ce"

swapoff -a

export DEBIAN_FRONTEND=noninteractive

# Docker
apt-get update
apt-get -y install curl dnsutils iptables ebtables ethtool ca-certificates conntrack util-linux socat jq nfs-common lsb-release apt-transport-https gnupg2

# Kubernetes & Docker repo keys
curl -sSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
curl -sSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -

# kubernetes & Docker repos
OS_ID="$(lsb_release --id --short | tr '[:upper:]' '[:lower:]')"
OS_CODENAME="$(lsb_release --codename --short)"
# You'd think that kubernetes-${OS_CODENAME} belongs there instead, but the debian repo
# contains neither kubeadm nor kubelet, and the docs themselves suggest using xenial repo.
echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" > /etc/apt/sources.list.d/kubernetes.list
echo "deb [arch=amd64] https://download.docker.com/linux/${OS_ID} ${OS_CODENAME} stable" > /etc/apt/sources.list.d/docker.list

apt-get update
apt-get install -y kubelet kubeadm kubectl docker-ce=$DOCKER_VERSION*
