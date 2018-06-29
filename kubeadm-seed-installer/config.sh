# Prereq, we need private keys for all machines in our posession.
export KUBERNETES_VERSION="1.10.2"
export DOCKER_VERSION="17.03"
export ETCD_VERSION="3.1.13"

# Must be one of "aws", "openstack", "vsphere" or "azure
export CLOUD_PROVIDER_FLAG=
# e.G. "./cloud-config.sh"
export CLOUD_CONFIG_FILE=

QUAY_IO_MIRROR="quay.io"
SSH_LOGIN="root"

MASTER_LOAD_BALANCER_ADDRS=()

# number of items should be ODD
MASTER_PUBLIC_IPS=()
MASTER_PRIVATE_IPS=()

# Additional Worker IP's (Don't enter APISERVER IP)
WORKER_PUBLIC_IPS=()

# A nodePort range to reserve for services.
# Leave it empty for default: 30000-32767
NODEPORT_RANGE=""

# The subnet used for pods (flannel)
# Leave it empty for default: 10.244.0.0/16
POD_SUBNET=""

# The subnet used for services
# Leave it empty for default: 10.96.0.0/12
SERVICE_SUBNET=""
