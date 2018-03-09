<%
var.name_base = "seed-k8s"

# one or more SSH public keys to install on all machines.
# The first key in the list will be installed while the machine is created and then used for provisioning the machine,
# which means it must provide SSH access from the local machine (i.e. the corresponding private key must
# be have been installed locally)
# TODO these will only be uploaded once, at machine creation time.
#      They won't be updated if the list is changed and "make" rerun later.
var.ssh_keys = [
    "ssh-rsa AAAAB4.......... me@thismachine",
    "ssh-rsa AAAAB3.......... somebody@somewhere"
]

var.ssh_user_name = "ubuntu"

var.external_gateway = "caf8de33-1059-4473-a2c1-2a62d12294fa"  # network ext-net

var.pool = "ext-net"

var.region = "cbk"

# todo can come from environment
var.domain = "default"

# OpenStack network to use. If not specified, one will be created
#var.network_id = "64738e94-396f-49f9-89b1-eff6e541b72d"

var.sys11_office_vpn_nets = ['151.252.43.0/24', '176.74.56.128/25']

var.sys11_f5_backend_ips = ['37.49.152.98', '37.49.152.99', '37.49.152.100', '37.49.152.212', '37.49.152.213', '37.49.152.214']

var.sys11_cloud_nets = ['185.56.128.0/21'] # 185.56.128.0/22 + 185.56.132.0/22, according to netbox

# [30000,32767] are possible k8s NodePorts
var.internal_web_port_ranges = [[80,87],[30000,32767]]

# load balancer IPs for accessing the masters. If not specified, the first master IP will be used
var.master_lb_ips = None

var.master_image = "Ubuntu 16.04 sys11-cloudimg amd64"

var.master_flavor = "m1.small"

var.master_count = 3

# Floating IPs to assign to master nodes.
# Must contain a list of zero or more floating IPs. The IPs must haven been reserved in the tenant
# before this script is run.
# The IPs will be assigned to the master nodes in ascending order (master 0, 1 etc.), and will
# not be released from the tenant upon "make destroy". If there are more nodes than IPs in this list,
# the remaining nodes will receive dynamically reserved floating IPs
# (which will be released by "make destroy")
var.master_ips = []


var.etcd_image = "Ubuntu 16.04 sys11-cloudimg amd64"

var.etcd_flavor = "m1.small"

var.etcd_count = 3


var.worker_image = "Ubuntu 16.04 sys11-cloudimg amd64"

var.worker_flavor = "m1.small"

var.worker_count = 5

# Floating IPs to assign to worker nodes in ascending order.
# Same deal as with var.master_ips.
var.worker_ips = []

%>
