install: cloud.conf config.sh
	./install.sh

destroy:
	./cleanup.sh
	git ls-files -o | xargs rm -f
	rm -rf apiserver0pki

config.sh: ../infrastructure/osk-cluster/variables.makotemplate ../infrastructure/osk-cluster/terraform.tfstate

# delete locally generated (automatically recreated) files
clean:
	rm -f cloud.conf config.sh

include ../rules.make
