KUBECTL = kubectl --kubeconfig=../kubeadm-seed-installer/kubeconfig

SHELL = /bin/bash

install: prepare-cluster values-secret.yaml namespace docker-secret.yaml install-job.yaml
	$(KUBECTL) create -f serviceaccount.yaml || true
	$(KUBECTL) create -f clusterrolebinding.yaml || true
	$(KUBECTL) apply -f values-secret.yaml
	$(KUBECTL) apply -f docker-secret.yaml
	$(KUBECTL) create -f install-job.yaml

prepare-cluster:
	$(MAKE) -C prepare-cluster.d
	run-parts --regex='\.sh$$' prepare-cluster.d

namespace:
	$(KUBECTL) apply -f namespace.yaml

values-secret.yaml: values.yaml

values.yaml: variables.mako variables_override.mako datacenters.yaml kubeconfig values_more.yaml.mako

values_more.yaml.mako:
	if [ ! -f "$@" ]; then touch "$@"; fi

kubeconfig: ../kubeadm-seed-installer/kubeconfig

# delete locally generated (automatically recreated) files
clean:
	$(MAKE) -C prepare-cluster.d clean
	rm -f kubeconfig values.yaml values-secret.yaml docker-secret.yaml install-job.yaml

# diff currently installed artifacts (values/datacenters/kubeconfig) against what would be installed by "make install". Needs yq (pip install yq).
diff: values.yaml
	@values=$$(mktemp); \
	$(KUBECTL) --namespace=kubermatic-installer get secret values -o json | jq -r '.data."values.yaml"' | base64 -d >$$values; \
	echo ">>>values.yaml:"; \
	diff -u <(cat $$values | grep -v 'datacenters:' | grep -v 'kubeconfig:') <(cat values.yaml | grep -v 'datacenters:' | grep -v 'kubeconfig:'); \
	echo ">>>datacenters.yaml:"; \
	diff -u <(cat $$values | yq -r '.kubermatic.datacenters' | base64 -d) <(cat values.yaml | yq -r '.kubermatic.datacenters' | base64 -d); \
	echo ">>>kubeconfig:"; \
	diff -u <(cat $$values | yq -r '.kubermatic.kubeconfig' | base64 -d) <(cat values.yaml | yq -r '.kubermatic.kubeconfig' | base64 -d); \
	rm -f $$values

include ../rules.make
