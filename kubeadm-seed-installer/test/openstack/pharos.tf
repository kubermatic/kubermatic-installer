output "pharos_api" {
  value = {
    endpoint = "${openstack_networking_floatingip_v2.e2e.0.address}"
  }
}

output "pharos_hosts" {
  value = {
    masters = {
      address         = "${slice(openstack_networking_floatingip_v2.e2e.*.address, 0, 3)}"
      private_address = "${slice(openstack_compute_instance_v2.seed-installer-e2e.*.network.0.fixed_ip_v4, 0, 3)}"
      role            = "master"
      user            = "ubuntu"

      # ssh_key_path    = "do_key"
    }

    workers = {
      address         = "${slice(openstack_networking_floatingip_v2.e2e.*.address, 3, 6)}"
      private_address = "${slice(openstack_compute_instance_v2.seed-installer-e2e.*.network.0.fixed_ip_v4, 3, 6)}"
      role            = "worker"
      user            = "ubuntu"

      # ssh_key_path    = "do_key"
    }
  }
}
