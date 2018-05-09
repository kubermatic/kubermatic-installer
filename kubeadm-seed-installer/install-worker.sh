#!/bin/bash

source ./config.sh
source ./functions.sh.include

# Setup Workers
echo "Creating workers"
for ((i = 0; i < ${#WORKER_PUBLIC_IPS[@]}; i++)); do
  echo "  ${i}: ${WORKER_PUBLIC_IPS[$i]}"
  TOKEN=$(ssh ${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[0]} "sudo kubeadm token create --print-join-command")

  ssh "${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]}" sudo mkdir -p /etc/kubernetes/ /etc/systemd/system/kubelet.service.d/
  rsync --rsync-path="sudo rsync" "${CLOUD_CONFIG_FILE}" "${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]}:/etc/kubernetes/cloud-config"
  rsync --rsync-path="sudo rsync" ./10-hostname.conf "${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]}:/etc/systemd/system/kubelet.service.d/10-hostname.conf"

  install_kubeadm "${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]}"
  ssh ${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]} "sudo ${TOKEN}"
  # TODO(realfake) On all workers:
  for (( retry_count = 0; retry_count < 10; retry_count++ )); do ssh ${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]} "sudo sed -i 's#server:.*#server: https://'"${MASTER_LOAD_BALANCER_ADDRS[0]}"':6443#g' /etc/kubernetes/kubelet.conf" && break || sleep 10; done

  ssh ${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]} "sudo systemctl daemon-reload && sudo systemctl restart kubelet"
done
