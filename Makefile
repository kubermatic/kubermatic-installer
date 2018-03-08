include rules.make

create-infrastructure:
	$(MAKE) -C infrastructure/osk-cluster create

destroy-infrastructure:
	$(MAKE) -C infrastructure/osk-cluster destroy


values.yaml: variables.makotemplate variables_override.makotemplate datacenters.yaml kubeadm-seed-installer/kubeconfig values_more.yaml.makotemplate

values_more.yaml.makotemplate:
	if [ ! -f "$@" ]; then touch "$@"; fi

## TODO stuff below not adapted yet

install: kubeadm-seed-installer/cloud.conf kubeadm-seed-installer/config.sh
	true # TODO invoke kubeadm-seed-installer/install.sh

kubeadm-seed-installer/config.sh: create-infrastructure/osk-cluster kubeadm-seed-installer/config.sh.template infrastructure/osk-cluster/variables.template infrastructure/osk-cluster/terraform.tfstate
	./expand_template.py -i kubeadm-seed-installer/config.sh.template -o $@

kubeadm-seed-installer/cloud.conf: kubeadm-seed-installer/cloud.conf.template
	if [ -z "$${OS_AUTH_URL}" -o -z "$${OS_USERNAME}" ]; then \
	    echo "OpenStack access credentials not found in environment. Please set this up first."; \
	    exit 1; \
	fi
	eval "echo \"$$(cat "$<")\"" > $@
