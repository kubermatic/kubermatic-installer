include ../rules.make

install: cloud.conf config.sh
	true # TODO invoke install.sh (it isn't idempotent)

config.sh: ../infrastructure/osk-cluster/variables.makotemplate ../infrastructure/osk-cluster/terraform.tfstate
