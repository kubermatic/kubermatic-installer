KUBECTL = kubectl --kubeconfig=../kubeadm-seed-installer/kubeconfig

install-kubermatic: values.yaml namespace docker-secret.yaml install-job.yaml
	$(KUBECTL) create -f serviceaccount.yaml || true
	$(KUBECTL) create -f clusterrolebinding.yaml || true
	$(KUBECTL) -n kubermatic-installer create secret generic values --from-file=values.yaml
	$(KUBECTL) apply -f docker-secret.yaml
	$(KUBECTL) create -f install-job.yaml

namespace:
	$(KUBECTL) apply -f namespace.yaml

values.yaml: variables.makotemplate variables_override.makotemplate datacenters.yaml kubeconfig values_more.yaml.makotemplate

values_more.yaml.makotemplate:
	if [ ! -f "$@" ]; then touch "$@"; fi

kubeconfig: ../kubeadm-seed-installer/kubeconfig

include ../rules.make
