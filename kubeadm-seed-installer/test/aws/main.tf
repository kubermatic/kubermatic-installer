resource "aws_key_pair" "the_key" {
  key_name_prefix = "k8s-install-test"
  public_key      = "${file("~/.ssh/id_rsa.pub")}"
}

data "aws_ami" "coreos" {
  most_recent = true

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "name"
    values = ["CoreOS-stable-*"]
  }

  owners = ["595879546273"] # CoreOS
}

resource "aws_instance" "master" {
  count = 3

  ami                         = "${data.aws_ami.coreos.id}"
  instance_type               = "t3.small"
  key_name                    = "${aws_key_pair.the_key.key_name}"
  subnet_id                   = "${aws_subnet.main.id}"
  associate_public_ip_address = true
  availability_zone           = "${var.availability_zone}"
  vpc_security_group_ids      = ["${aws_security_group.masters.id}"]

  tags {
    Name = "install-test"
  }
}

output "master_public_ips" {
  value = "${join(" ", aws_instance.master.*.public_ip)}"
}

output "master_private_ips" {
  value = "${join(" ", aws_instance.master.*.private_ip)}"
}

output "worker_ips" {
  value = ""
}

output "loadbalancer_addr" {
  value = "${aws_elb.master_elb.dns_name}"
}
