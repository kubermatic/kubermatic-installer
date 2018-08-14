resource "openstack_compute_instance_v2" "seed-installer-e2e" {
  count           = 6
  name            = "seed-e2e-test-${count.index}"
  image_name      = "Ubuntu 16.04 LTS 2018.03.26"
  flavor_name     = "m1.small"
  key_pair        = "seed-installer-e2e"
  security_groups = ["allow-ssh-all", "default"]

  network {
    name = "private"
  }
}

resource "openstack_networking_floatingip_v2" "e2e" {
  count = 6
  pool  = "ext-net"
}

resource "openstack_compute_floatingip_associate_v2" "e2e" {
  count       = 6
  floating_ip = "${element(openstack_networking_floatingip_v2.e2e.*.address, count.index)}"
  instance_id = "${element(openstack_compute_instance_v2.seed-installer-e2e.*.id, count.index)}"
  fixed_ip    = "${element(openstack_compute_instance_v2.seed-installer-e2e.*.network.0.fixed_ip_v4, count.index)}"
}
