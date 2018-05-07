#!/bin/bash

source ./config.sh

rm -rf apiserver0pki/

# Clean up Master
for ((i = 0; i < ${#MASTER_HOSTNAMES[@]}; i++)); do
  # `bash -l -c $COMMAND` will make sure that the PATH is properly initialized,
  # e.g. includes /opt/bin on CoreOS.
        ssh ${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[$i]} "sudo bash -l -c 'kubeadm reset' || true"
        ssh ${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[$i]} "sudo rm -rf /etc/kubernetes || true"
        ssh ${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[$i]} "rm -rf ~/etc/ || true"
done

# Clean up ETCD
for ((i = 0; i < ${#ETCD_HOSTNAMES[@]}; i++)); do
        ssh ${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]} "sudo systemctl stop etcd || true"
        ssh ${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]} "sudo systemctl disable etcd || true"
        ssh ${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]} "sudo rm -rf /var/lib/etcd || true"
  # The binaries are installed under /usr/local/bin on Ubuntu and /opt/bin on CoreOS
  # therefore we delete them from both locations.
        ssh ${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]} "sudo rm -f /usr/local/bin/etcd /opt/bin/etcd || true"
        ssh ${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]} "sudo rm -f /usr/local/bin/etcdctl /opt/bin/etcdctl || true"
        ssh ${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]} "sudo rm -rf /etc/kubernetes || true"
        ssh ${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]} "rm -rf ~/etc/ || true"
done

# Clean up Worker
for ((i = 0; i < ${#WORKER_PUBLIC_IPS[@]}; i++)); do
        ssh ${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]} "sudo kubeadm reset || true"
        ssh ${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]} "sudo rm -rf /etc/kubernetes || true"
        ssh ${DEFAULT_LOGIN_USER}@${WORKER_PUBLIC_IPS[$i]} "rm -rf ~/etc/ || true"
done
