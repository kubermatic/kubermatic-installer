KUBECTL = kubectl --kubeconfig=../kubeadm-seed-installer/kubeconfig

install: values.yaml namespace docker-secret.yaml install-job.yaml
	$(KUBECTL) create -f serviceaccount.yaml || true
	$(KUBECTL) create -f clusterrolebinding.yaml || true
	$(KUBECTL) -n kubermatic-installer create secret generic values --from-file=values.yaml
	$(KUBECTL) apply -f docker-secret.yaml
	$(KUBECTL) create -f install-job.yaml

namespace:
	$(KUBECTL) apply -f namespace.yaml

values.yaml: variables.mako variables_override.mako datacenters.yaml kubeconfig values_more.yaml.mako

values_more.yaml.mako:
	if [ ! -f "$@" ]; then touch "$@"; fi

kubeconfig: ../kubeadm-seed-installer/kubeconfig

# delete locally generated (automatically recreated) files
clean:
	rm -f kubeconfig values.yaml docker-secret.yaml install-job.yaml

include ../rules.make
