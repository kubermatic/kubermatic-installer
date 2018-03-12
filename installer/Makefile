KUBECTL = kubectl --kubeconfig=../kubeadm-seed-installer/kubeconfig

install: values-secret.yaml namespace docker-secret.yaml install-job.yaml
	$(KUBECTL) create -f serviceaccount.yaml || true
	$(KUBECTL) create -f clusterrolebinding.yaml || true
	$(KUBECTL) apply -f values-secret.yaml
	$(KUBECTL) apply -f docker-secret.yaml
	$(KUBECTL) create -f install-job.yaml

namespace:
	$(KUBECTL) apply -f namespace.yaml

values-secret.yaml: values.yaml

values.yaml: variables.mako variables_override.mako datacenters.yaml kubeconfig values_more.yaml.mako

values_more.yaml.mako:
	if [ ! -f "$@" ]; then touch "$@"; fi

kubeconfig: ../kubeadm-seed-installer/kubeconfig

# delete locally generated (automatically recreated) files
clean:
	rm -f kubeconfig values.yaml values-secret.yaml docker-secret.yaml install-job.yaml

include ../rules.make
