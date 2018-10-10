resource "aws_key_pair" "the_key" {
  key_name_prefix = "k8s-install-test"
  public_key      = "${file("machine-key.pub")}"
}

data "aws_ami" "coreos" {
  most_recent = true

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-xenial-16.04-amd64-server-*"]
  }

  owners = ["099720109477"] # Canonical
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
  iam_instance_profile        = "${aws_iam_instance_profile.profile.name}"

  tags = "${map(
    "Name", "install-test",
    "kubernetes.io/cluster/${random_id.id.hex}", "shared"
  )}"
}

resource "aws_instance" "worker" {
  count = 3

  ami                         = "${data.aws_ami.coreos.id}"
  instance_type               = "t3.small"
  key_name                    = "${aws_key_pair.the_key.key_name}"
  subnet_id                   = "${aws_subnet.main.id}"
  associate_public_ip_address = true
  availability_zone           = "${var.availability_zone}"
  vpc_security_group_ids      = ["${aws_security_group.workers.id}"]
  iam_instance_profile        = "${aws_iam_instance_profile.profile.name}"

  tags = "${map(
    "Name", "install-test",
    "kubernetes.io/cluster/${random_id.id.hex}", "shared"
  )}"
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
  value = "${aws_route53_record.master.fqdn}"
}

# data for cloud-config
output "availability_zone" {
  value = "${var.availability_zone}"
}

output "vpc" {
  value = "${aws_vpc.main.id}"
}

output "subnet" {
  value = "${aws_subnet.main.id}"
}

output "route_table" {
  value = "${aws_route_table.main.id}"
}

output "cluster_name" {
  value = "${random_id.id.hex}"
}
