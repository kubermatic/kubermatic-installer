install-kubermatic: install-seed
	$(MAKE) -C installer install

install-seed: install-infrastructure
	$(MAKE) -C kubeadm-seed-installer install

install-infrastructure:
	$(MAKE) -C infrastructure/osk-cluster install

plan-infrastructure:
	$(MAKE) -C infrastructure/osk-cluster plan

destroy-infrastructure:
	$(MAKE) -C infrastructure/osk-cluster destroy
	$(MAKE) -C kubeadm-seed-installer destroy-localstate clean

destroy-seed:
	$(MAKE) -C kubeadm-seed-installer destroy

# delete locally generated (automatically recreated) files
clean:
	$(MAKE) -C infrastructure/osk-cluster clean
	$(MAKE) -C infrastructure/osk-cluster clean
	$(MAKE) -C kubeadm-seed-installer clean
