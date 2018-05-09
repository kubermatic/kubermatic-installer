#!/bin/bash

source ./config.sh

install_kubeadm() {
  local USERHOST="$1"

  local OS_ID=$(ssh "${USERHOST}" cat /etc/os-release | grep '^ID=' | sed s/^ID=//)

  case $OS_ID in
    ubuntu|debian)
      scp ./install-kubeadm-ubuntu.sh $USERHOST:~/etc/kubernetes/install-kubeadm-ubuntu.sh
      ssh ${USERHOST} "sudo mv ~/etc/kubernetes/install-kubeadm-ubuntu.sh /etc/kubernetes/install-kubeadm-ubuntu.sh"
      ssh ${USERHOST} "sudo bash /etc/kubernetes/install-kubeadm-ubuntu.sh"
    ;;
    coreos)
      scp ./install-kubeadm-coreos.sh ${USERHOST}:~/etc/kubernetes/install-kubeadm-coreos.sh
      ssh ${USERHOST} "sudo mv ~/etc/kubernetes/install-kubeadm-coreos.sh /etc/kubernetes/install-kubeadm-coreos.sh"
      ssh ${USERHOST} "sudo bash /etc/kubernetes/install-kubeadm-coreos.sh"
    ;;
    *)
      echo " ### Operating system '$OS_ID' is not supported."
      exit 1
    ;;
  esac
}

# Setup Workers
for ((i = 0; i < ${#WORKER_PUBLIC_IPS[@]}; i++)); do
        echo "Create Worker $i"
        TOKEN=$(ssh ${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[0]} "sudo kubeadm token create --print-join-command")

        ssh ${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]} "mkdir -p ~/etc/kubernetes"
        ssh ${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]} "sudo mkdir -p /etc/kubernetes"
        scp ./install-kubeadm-ubuntu.sh ${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]}:~/etc/kubernetes/install-kubeadm-ubuntu.sh
        scp ${CLOUD_CONFIG_FILE} ${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]}:~/etc/kubernetes/cloud-config
        ssh ${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]} "sudo cp ~/etc/kubernetes/* /etc/kubernetes/; sudo chown -R root:root /etc/kubernetes/"

        scp ./10-hostname.conf ${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]}:~/10-hostname.conf
        ssh ${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]} "sudo mkdir -p /etc/systemd/system/kubelet.service.d/"
        ssh ${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]} "sudo cp ~/10-hostname.conf /etc/systemd/system/kubelet.service.d/; sudo chown -R root:root /etc/systemd/system/kubelet.service.d/"

        install_kubeadm "${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]}"
        ssh ${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]} "sudo ${TOKEN}"
        # TODO(realfake) On all workers:
        for (( retry_count = 0; retry_count < 10; retry_count++ )); do ssh ${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]} "sudo sed -i 's#server:.*#server: https://'"${MASTER_LOAD_BALANCER_ADDRS[0]}"':6443#g' /etc/kubernetes/kubelet.conf" && break || sleep 10; done

        ssh ${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]} "sudo systemctl daemon-reload"
        ssh ${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]} "sudo systemctl restart kubelet"
done
