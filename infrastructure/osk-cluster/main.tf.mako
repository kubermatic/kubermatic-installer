<%! import os %>

<%
var.network_id = None
%>

<%include file="variables.mako" />

% if os.path.isfile("variables_override.mako"):
<%include file="variables_override.mako" />
% endif

provider "openstack" {
  domain_name = "${var.domain}"
}


## ####################### network #######################

% if var.network_id:

  <% net_id = var.network_id %>

% else:

  resource "openstack_networking_network_v2" "terraform" {
    name = "net.${var.name_base}"
    region = "${var.region}"
    admin_state_up = "true"
  }

  resource "openstack_networking_subnet_v2" "terraform" {
    name = "subnet.${var.name_base}"
    region = "${var.region}"
    network_id = "${"${"}openstack_networking_network_v2.terraform.id}"
    cidr = "10.0.4.0/24"
    ip_version = 4
    enable_dhcp = "true"
    dns_nameservers = ["37.123.105.116", "37.123.105.117"]
  }

  resource "openstack_networking_router_v2" "terraform" {
    name = "router2ext.${var.name_base}"
    region = "${var.region}"
    admin_state_up = "true"
    external_gateway = "${var.external_gateway}"
  }

  resource "openstack_networking_router_interface_v2" "terraform" {
    region = "${var.region}"
    router_id = "${"${"}openstack_networking_router_v2.terraform.id}"
    subnet_id = "${"${"}openstack_networking_subnet_v2.terraform.id}"
  }

  <% net_id = '${openstack_networking_network_v2.terraform.id}' %>

% endif


## ####################### security groups #######################

resource "openstack_compute_secgroup_v2" "ssh_access_secgroup" {
  name = "ssh_access.secgroup.${var.name_base}"
  region = "${var.region}"
  description = "Security group for SSH access from anywhere"
  rule {
    from_port = 22
    to_port = 22
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
  rule {
    ip_protocol = "icmp"
    from_port = "-1"
    to_port = "-1"
    cidr = "0.0.0.0/0"
  }
}

resource "openstack_compute_secgroup_v2" "web_access_secgroup" {
  name = "web_access.secgroup.${var.name_base}"
  region = "${var.region}"
  description = "Security group for http/https access from anywhere"
  rule {
    from_port = 80
    to_port = 80
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
  rule {
    from_port = 443
    to_port = 443
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}

resource "openstack_compute_secgroup_v2" "nodeport_access_secgroup" {
  name = "nodeport_access.secgroup.${var.name_base}"
  region = "${var.region}"
  description = "Security group for accessing K8s NodePorts from anywhere"
  rule {
    from_port = 30000
    to_port = 32767
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}

resource "openstack_compute_secgroup_v2" "apiserver_access_secgroup" {
  name = "apiserver_access.secgroup.${var.name_base}"
  region = "${var.region}"
  description = "Security group for https access to the K8s API server from anywhere"
  rule {
    from_port = 6443
    to_port = 6443
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}



resource "openstack_compute_secgroup_v2" "sys11internal_web_access_secgroup" {
  name = "sys11internal_web_access.secgroup.${var.name_base}"
  region = "${var.region}"
  description = "Security group for http/https access from the Sys11 Office and VPN networks"

  rule {
    from_port = -1
    to_port = -1
    ip_protocol = "icmp"
    cidr = "0.0.0.0/0"
  }

  % for net in var.sys11_office_vpn_nets + var.sys11_cloud_nets:
   % for pr in [[80,80],[443,443]]:
   rule {
     from_port = ${pr[0]}
     to_port = ${pr[1]}
     ip_protocol = "tcp"
     cidr = "${net}"
   }
   % endfor
  % endfor
}

resource "openstack_compute_secgroup_v2" "f5_access_secgroup" {
  name = "f5_access.secgroup.${var.name_base}"
  region = "${var.region}"
  description = "Security group for HTTP access from the F5 SSL terminator proxy/load balancer"

  rule {
    from_port = -1
    to_port = -1
    ip_protocol = "icmp"
    cidr = "0.0.0.0/0"
  }

  % for ip in var.sys11_f5_backend_ips:
   % for pr in var.internal_web_port_ranges:
   rule {
     from_port = ${pr[0]}
     to_port = ${pr[1]}
     ip_protocol = "tcp"
     cidr = "${ip}/32"
   }
   % endfor
  % endfor
}


## ####################### initial key pair #######################

resource "openstack_compute_keypair_v2" "keypair" {
  name = "keypair_${"${"}replace("${var.name_base}", ".", "_")}"
  public_key = "${var.ssh_keys[0]}"
  region = "${var.region}"
}



## ####################### masters #######################


% for i in range(var.master_count):

resource "openstack_compute_instance_v2" "k8s_master${i}" {
  name = "master${i}-${var.name_base}"
  region = "${var.region}"
  image_name = "${var.master_image}"
  flavor_name = "${var.master_flavor}"
  key_pair = "${"${"}openstack_compute_keypair_v2.keypair.name}"

  security_groups = [
      "default",
      "${"${"}openstack_compute_secgroup_v2.ssh_access_secgroup.name}",
      "${"${"}openstack_compute_secgroup_v2.apiserver_access_secgroup.name}"
  ]

  network {
    uuid = "${net_id}"
  }

  connection {
    user = "ubuntu"
  }

  ## TODO do this via scp or Ansible rather than Tf to support updating the list after machine creation
  provisioner "remote-exec" {
    inline = [
  % for key in var.ssh_keys[1:]:
      "echo '${key}' >>~/.ssh/authorized_keys",
  % endfor
    ]
  }

}

  % if i < len(var.master_ips):

resource "openstack_compute_floatingip_associate_v2" "master${i}_fip_assoc" {
  instance_id = "${"${"}openstack_compute_instance_v2.k8s_master${i}.id}"
  floating_ip = "${var.master_ips[i]}"
}

  % else:

resource "openstack_networking_floatingip_v2" "master_floating_ip_${i}" {
  region = "${var.region}"
  pool = "${var.pool}"
}

resource "openstack_compute_floatingip_associate_v2" "master${i}_fip_assoc" {
  instance_id = "${"${"}openstack_compute_instance_v2.k8s_master${i}.id}"
  floating_ip = "${"${"}openstack_networking_floatingip_v2.master_floating_ip_${i}.address}"
}

  %endif


% endfor




## ####################### workers #######################

% for i in range(var.worker_count):

resource "openstack_compute_instance_v2" "k8s_worker${i}" {
  name = "worker${i}-${var.name_base}"
  region = "${var.region}"
  image_name = "${var.worker_image}"
  flavor_name = "${var.worker_flavor}"
  key_pair = "${"${"}openstack_compute_keypair_v2.keypair.name}"

  security_groups = [
      "default",
      "${"${"}openstack_compute_secgroup_v2.ssh_access_secgroup.name}",
      "${"${"}openstack_compute_secgroup_v2.web_access_secgroup.name}",
      "${"${"}openstack_compute_secgroup_v2.nodeport_access_secgroup.name}"
  ]

  network {
    uuid = "${net_id}"
  }

  connection {
    user = "ubuntu"
  }

  ## TODO do this via scp or Ansible rather than Tf to support updating the list after machine creation
  provisioner "remote-exec" {
    inline = [
  % for key in var.ssh_keys[1:]:
      "echo '${key}' >>~/.ssh/authorized_keys",
  % endfor
    ]
  }

}

  % if i < len(var.worker_ips):

resource "openstack_compute_floatingip_associate_v2" "worker${i}_fip_assoc" {
  instance_id = "${"${"}openstack_compute_instance_v2.k8s_worker${i}.id}"
  floating_ip = "${var.worker_ips[i]}"
}

  % else:

resource "openstack_networking_floatingip_v2" "worker_floating_ip_${i}" {
  region = "${var.region}"
  pool = "${var.pool}"
}

resource "openstack_compute_floatingip_associate_v2" "worker${i}_fip_assoc" {
  instance_id = "${"${"}openstack_compute_instance_v2.k8s_worker${i}.id}"
  floating_ip = "${"${"}openstack_networking_floatingip_v2.worker_floating_ip_${i}.address}"
}

  % endif

% endfor




## ####################### etcds #######################

% for i in range(var.etcd_count):

resource "openstack_compute_instance_v2" "etcd${i}" {
  name = "etcd${i}-${var.name_base}"
  region = "${var.region}"
  image_name = "${var.etcd_image}"
  flavor_name = "${var.etcd_flavor}"
  key_pair = "${"${"}openstack_compute_keypair_v2.keypair.name}"

  security_groups = [ "default" ]

  network {
    uuid = "${net_id}"
  }

  connection {
    user = "ubuntu"
  }

  ## TODO do this via scp or Ansible rather than Tf to support updating the list after machine creation
  provisioner "remote-exec" {
    inline = [
  % for key in var.ssh_keys[1:]:
      "echo '${key}' >>~/.ssh/authorized_keys",
  % endfor
    ]
  }

}

% endfor
