install: cloud.conf config.sh
	./install.sh

config.sh: ../infrastructure/osk-cluster/variables.makotemplate ../infrastructure/osk-cluster/terraform.tfstate

include ../rules.make
