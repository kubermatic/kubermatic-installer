resource "aws_vpc" "main" {
  cidr_block           = "${var.vpc_cidr}"
  enable_dns_hostnames = true

  tags {
    Name = "k8s-seed-tf"
  }
}

resource "aws_internet_gateway" "gw" {
  vpc_id = "${aws_vpc.main.id}"

  tags {
    Name = "k8s-seed-tf"
  }
}

resource "aws_subnet" "main" {
  vpc_id                  = "${aws_vpc.main.id}"
  cidr_block              = "${var.vpc_cidr}"
  availability_zone       = "${var.availability_zone}"
  map_public_ip_on_launch = true

  tags {
    Name = "k8s-seed-tf"
  }
}

resource "aws_route_table" "main" {
  vpc_id     = "${aws_vpc.main.id}"
  depends_on = ["aws_vpc.main", "aws_internet_gateway.gw"]

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = "${aws_internet_gateway.gw.id}"
  }

  tags {
    Name = "k8s-seed-tf"
  }
}

resource "aws_main_route_table_association" "a" {
  vpc_id         = "${aws_vpc.main.id}"
  route_table_id = "${aws_route_table.main.id}"
}

resource "aws_security_group" "kubernetes_api" {
  vpc_id = "${aws_vpc.main.id}"
  name   = "k8s-seed-tf-api-${random_id.id.hex}"

  # Allow inbound traffic to the port used by Kubernetes API HTTPS
  ingress {
    from_port   = 6443
    to_port     = 6443
    protocol    = "TCP"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow all outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "k8s-seed-tf"
  }
}

resource "aws_security_group" "masters" {
  vpc_id = "${aws_vpc.main.id}"
  name   = "k8s-seed-tf-masters-${random_id.id.hex}"

  # Allow all outbound
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow ICMP
  ingress {
    from_port   = 8
    to_port     = 0
    protocol    = "icmp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow all internal
  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["${var.vpc_cidr}"]
  }

  # Allow all traffic from the API ELB
  ingress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    security_groups = ["${aws_security_group.kubernetes_api.id}"]
  }

  # Allow SSH
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "k8s-seed-tf"
  }
}

resource "aws_elb" "master_elb" {
  depends_on      = ["aws_internet_gateway.gw"]
  name            = "k8s-seed-tf-api-${random_id.id.hex}"
  internal        = false
  instances       = ["${aws_instance.master.*.id}"]
  subnets         = ["${aws_subnet.main.id}"]
  security_groups = ["${aws_security_group.kubernetes_api.id}"]

  listener {
    lb_port           = 6443
    instance_port     = 6443
    lb_protocol       = "TCP"
    instance_protocol = "TCP"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 15
    target              = "HTTPS:6443/healthz"
    interval            = 30
  }
}

data "aws_route53_zone" "loodse" {
  name = "aws.loodse.com."
}

resource "aws_route53_record" "wildcard" {
  depends_on = ["aws_elb.master_elb"]
  zone_id    = "${data.aws_route53_zone.loodse.zone_id}"
  name       = "*.${random_id.id.hex}.aws.loodse.com"
  type       = "CNAME"
  ttl        = "300"
  records    = ["${aws_elb.master_elb.dns_name}"]
}

resource "aws_route53_record" "master" {
  depends_on = ["aws_elb.master_elb"]
  zone_id    = "${data.aws_route53_zone.loodse.zone_id}"
  name       = "${random_id.id.hex}.aws.loodse.com"
  type       = "CNAME"
  ttl        = "300"
  records    = ["${aws_elb.master_elb.dns_name}"]
}
