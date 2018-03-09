install: cloud.conf config.sh
	./install.sh

destroy: _cleanup destroy-localstate

_cleanup:
	./cleanup.sh

destroy-localstate:
	git ls-files -o | xargs rm -f
	rm -rf apiserver0pki

config.sh: ../infrastructure/osk-cluster/variables.makotemplate ../infrastructure/osk-cluster/terraform.tfstate

# delete locally generated (automatically recreated) files
clean:
	rm -f cloud.conf config.sh

include ../rules.make
