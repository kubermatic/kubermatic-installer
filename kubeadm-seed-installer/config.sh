# Prereq, we need private keys for all machines in our posession.
export KUBERNETES_VERSION="1.10.2"
export DOCKER_VERSION="17.03"
export ETCD_VERSION="3.1.13"

export CLOUD_PROVIDER_FLAG=
export CLOUD_CONFIG_FILE=

SSH_LOGIN="root"

MASTER_LOAD_BALANCER_ADDRS=()

# number of items should be ODD
MASTER_PUBLIC_IPS=()
MASTER_PRIVATE_IPS=()

# Additional Worker IP's (Don't enter APISERVER IP)
WORKER_PUBLIC_IPS=()
