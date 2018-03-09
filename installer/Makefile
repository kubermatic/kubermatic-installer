install-kubermatic: namespace values.yaml docker-secret.yaml install-job.yaml
	kubectl create -f serviceaccount.yaml || true
	kubectl create -f clusterrolebinding.yaml || true
	kubectl -n kubermatic-installer create secret generic values --from-file=values.yaml
	kubectl apply -f docker-secret.yaml
	kubectl create -f install-job.yaml

namespace:
	kubectl apply -f installer/namespace.yaml

values.yaml: variables.makotemplate variables_override.makotemplate datacenters.yaml kubeconfig values_more.yaml.makotemplate

values_more.yaml.makotemplate:
	if [ ! -f "$@" ]; then touch "$@"; fi

kubeconfig: ../kubeadm-seed-installer/kubeconfig

include ../rules.make
