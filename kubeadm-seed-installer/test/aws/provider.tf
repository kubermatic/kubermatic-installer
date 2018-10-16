variable "region" {
  default = "eu-central-1"
}

variable "availability_zone" {
  default = "eu-central-1a"
}

variable "vpc_cidr" {
  default = "10.43.0.0/16"
}

provider "aws" {
  region = "${var.region}"
}

resource "random_id" "id" {
  byte_length = 8
}
