output "pharos_api" {
  value = {
    endpoint = "${aws_elb.master_elb.dns_name}"
  }
}

output "pharos_hosts" {
  value = {
    masters = {
      address         = "${aws_instance.master.*.public_ip}"
      private_address = "${aws_instance.master.*.private_ip}"
      role            = "master"
      user            = "ubuntu"

      ssh_key_path    = "../machine-key"
    }

    workers = {
      address         = "${aws_instance.worker.*.public_ip}"
      private_address = "${aws_instance.worker.*.private_ip}"
      role            = "worker"
      user            = "ubuntu"

      ssh_key_path    = "../machine-key"
    }
  }
}
