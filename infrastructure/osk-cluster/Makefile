install: main.tf
	$(call check_project, .)
	terraform init; \
	terraform apply -auto-approve -parallelism=2 .

destroy: main.tf
	$(call check_project, .)
	terraform destroy -force .
	rm -f .projectid
	#rm -f ../../kubeadm-seed-installer/config.sh

plan: main.tf
	$(call check_project, .)
	terraform init; \
	terraform plan .

# delete locally generated (automatically recreated) files
clean:
	rm -f main.tf

# helper function: $(call check_project <infradir>)
# record OSK project we're working under in <infradir>/.projectid and check it against $OS_TENANT_ID so the user later
# won't accidentally rerun the build under a different project
define check_project
	@if [ ! -f $1/.projectid ]; then \
	    if [ -z "$$OS_TENANT_ID" ]; then \
	        echo "OS_TENANT_ID not set in environment. Please source an OpenStack environment rc file" >&2; \
	        exit 1; \
	    fi; \
	    echo $$OS_TENANT_ID >$1/.projectid; \
	fi; \
	if [ `cat $1/.projectid ` != "$$OS_TENANT_ID" ]; then \
	    echo "Unexpected OpenStack project: Stack was built under project `cat $1/.projectid `, current project is $$OS_TENANT_ID" >&2; \
	    echo "Aborting. If you really want to proceed under this new project, remove .projectid and rerun the build." >&2; \
	    exit 1; \
	fi;
endef

include ../../rules.make

main.tf: variables.mako variables_override.mako
