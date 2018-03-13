#!/bin/bash
<%! import os %>
\
<%include file="infrastructure/osk-cluster/variables.mako" />\
\
% if os.path.isfile("../infrastructure/osk-cluster/variables_override.mako"):
<%include file="infrastructure/osk-cluster/variables_override.mako" />\
% endif
\
<% infra = read_tfstate('../infrastructure/osk-cluster/terraform.tfstate') %>\
\
# Prereq, we need private keys for all machines in our posession.
KUBERNETES_VERSION="v1.9.2"

CLOUD_PROVIDER_FLAG=openstack
CLOUD_CONFIG_FILE=./cloud.conf
DEFAULT_PRIVATE_IP4_INTERFACE=eth0
DEFAULT_LOGIN_USER=ubuntu
ETCD_HOSTNAMES=(\
% for i in range(var.etcd_count):
 ${infra['etcd%i' % i]['name']} \
% endfor
)
ETCD_PRIVATE_IPS=(\
% for i in range(var.etcd_count):
 ${infra['etcd%i' % i]['network.0.fixed_ip_v4']} \
% endfor
)
ETCD_PUBLIC_IPS=(\
% for i in range(var.etcd_count):
 ${infra['etcd%i' % i]['network.0.fixed_ip_v4']} \
% endfor
)

POD_SUBNET="10.244.0.0/16" # Canal
MASTER_LOAD_BALANCER_ADDRS=(\
% if var.master_lb_ips:
 ${var.master_lb_ips}\
% else:
 ${infra['master0_fip_assoc']['floating_ip']}\
%endif
)
MASTER_HOSTNAMES=(\
% for i in range(var.master_count):
 ${infra['k8s_master%i' % i]['name']} \
% endfor
)
MASTER_PRIVATE_IPS=(\
% for i in range(var.master_count):
 ${infra['k8s_master%i' % i]['network.0.fixed_ip_v4']} \
% endfor
)
MASTER_PUBLIC_IPS=(\
% for i in range(var.master_count):
 ${infra['master%i_fip_assoc' % i]['floating_ip']} \
% endfor
)

# Additional Worker IP's (Don't enter APISERVER IP)
WORKER_PRIVATE_IPS=(\
% for i in range(var.worker_count):
 ${infra['k8s_worker%i' % i]['network.0.fixed_ip_v4']} \
% endfor
)
WORKER_PUBLIC_IPS=(\
% for i in range(var.worker_count):
 ${infra['worker%i_fip_assoc' % i]['floating_ip']} \
% endfor
)