default: install

destroy: destroy-infrastructure/osk-cluster

install: loodse-scripts/cloud.conf loodse-scripts/config.sh
	true # TODO invoke loodse-scripts/install.sh

loodse-scripts/config.sh: create-infrastructure/osk-cluster loodse-scripts/config.sh.template infrastructure/osk-cluster/variables.template infrastructure/osk-cluster/terraform.tfstate
	./expand_template.py -i loodse-scripts/config.sh.template -o $@

loodse-scripts/cloud.conf: loodse-scripts/cloud.conf.template
	if [ -z "$${OS_AUTH_URL}" -o -z "$${OS_USERNAME}" ]; then \
	    echo "OpenStack access credentials not found in environment. Please set this up first."; \
	    exit 1; \
	fi
	eval "echo \"$$(cat "$<")\"" > $@


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

create-infrastructure/%: infrastructure/%/main.tf
	$(call check_project, infrastructure/$*/)
	cd infrastructure/$*/; \
	terraform init; \
	terraform apply -auto-approve -parallelism=2 .

destroy-infrastructure/%: infrastructure/%/main.tf
	$(call check_project, infrastructure/$*/)
	cd infrastructure/$*/; \
	terraform destroy -force .
	rm -f loodse-scripts/config.sh infrastructure/$*/.projectid

infrastructure/%/main.tf: infrastructure/%/main.tf.template infrastructure/%/variables.template
	./expand_template.py -i $< -o $@
