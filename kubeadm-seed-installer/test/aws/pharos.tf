output "pharos_api" {
  value = {
    endpoint = "${aws_route53_record.master.fqdn}"
  }
}
 output "pharos_hosts" {
  value = {
    masters = {
      address         = "${aws_instance.master.*.public_ip}"
      private_address = "${aws_instance.master.*.private_ip}"
      role            = "master"
      user            = "ubuntu"
      # ssh_key_path    = "do_key"
    }
    #workers = {
      #address         = "${aws_instance.meta_worker.*.public_ip}"
      #private_address = "${aws_instance.meta_worker.*.private_ip}"
      #role            = "worker"
      #user            = "ubuntu"
      # ssh_key_path    = "do_key"
    #}
  }
}
