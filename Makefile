install-kubermatic: install-seed
	$(MAKE) -C installer install-kubermatic

install-seed: create-infrastructure
	$(MAKE) -C kubeadm-seed-installer install

create-infrastructure:
	$(MAKE) -C infrastructure/osk-cluster create

destroy-infrastructure:
	$(MAKE) -C infrastructure/osk-cluster destroy


