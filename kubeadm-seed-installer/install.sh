#!/bin/bash

# See https://kubernetes.io/docs/setup/independent/high-availability/#install-cni-network
set -eu pipefail

SCRIPT_DIR="$(dirname "${BASH_SOURCE[0]}")"

source "$SCRIPT_DIR/config.sh"
source "$SCRIPT_DIR/functions.sh.include"

"$SCRIPT_DIR/install-prerequistes.sh"

echo "Ensuring that all hosts have rsync installed"
for hostname in ${ETCD_HOSTNAMES[@]} ${MASTER_HOSTNAMES[@]} ${WORKER_PUBLIC_IPS[@]}; do
  echo " * $hostname"
  ensure_host_has_rsync "${DEFAULT_LOGIN_USER}@$hostname"
done

# Generate PKI
[[ -f ca.pem ]] || cfssl gencert -initca ca-csr.json | cfssljson -bare ca -

# Generate etcd client CA.
[[ -f client.pem ]] || cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=client client.json | cfssljson -bare client

for ((i = 0; i < ${#ETCD_HOSTNAMES[@]}; i++)); do
  if [[ ! -f "config${i}.json" ]]; then
    echo "Generating PKI for etcd server ${i}"
    cfssl print-defaults csr > "config${i}.json"
    sed -i '0,/CN/{s/example\.net/'"${ETCD_HOSTNAMES[$i]}"'/}' "config${i}.json"
    sed -i 's/www\.example\.net/'"${ETCD_PRIVATE_IPS[$i]}"'/' "config${i}.json"
    sed -i 's/example\.net/'"${ETCD_HOSTNAMES[$i]}"'/' "config${i}.json"

    cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=server "config${i}.json" | cfssljson -bare "server${i}"
    cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=peer "config${i}.json" | cfssljson -bare "peer${i}"
  fi
done

echo "Copying PKI to etcd servers"
for ((i = 0; i < ${#ETCD_HOSTNAMES[@]}; i++)); do
  echo "  ${i}: ${ETCD_PUBLIC_IPS[$i]}"
  ssh "${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]}" sudo mkdir -p /etc/kubernetes/pki/etcd/
  rsync --rsync-path="sudo rsync" ca.pem ca-key.pem client.pem client-key.pem "${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]}:/etc/kubernetes/pki/etcd/"

  rsync --rsync-path="sudo rsync" "./config${i}.json" "${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]}:/etc/kubernetes/pki/etcd/config.json"
  rsync --rsync-path="sudo rsync" "./peer${i}.csr" "${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]}:/etc/kubernetes/pki/etcd/peer.csr"
  rsync --rsync-path="sudo rsync" "./peer${i}-key.pem" "${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]}:/etc/kubernetes/pki/etcd/peer-key.pem"
  rsync --rsync-path="sudo rsync" "./peer${i}.pem" "${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]}:/etc/kubernetes/pki/etcd/peer.pem"
  rsync --rsync-path="sudo rsync" "./server${i}.csr" "${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]}:/etc/kubernetes/pki/etcd/server.csr"
  rsync --rsync-path="sudo rsync" "./server${i}-key.pem" "${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]}:/etc/kubernetes/pki/etcd/server-key.pem"
  rsync --rsync-path="sudo rsync" "./server${i}.pem" "${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]}:/etc/kubernetes/pki/etcd/server.pem"
done

# Build etcd ring ie. etcd0=https://<etcd0-ip-address>:2380,etcd1=https://<etcd1-ip-address>:2380,etcd2=https://<etcd2-ip-address>:2380
ETCD_RING=""
for ((i = 0; i < ${#ETCD_HOSTNAMES[@]}; i++)); do
        ETCD_RING+=${ETCD_HOSTNAMES[$i]}"=https://"${ETCD_PRIVATE_IPS[$i]}":2380,"
done
ETCD_RING="$(echo ${ETCD_RING} | sed 's/[,]*$//')"
echo $ETCD_RING

NEW_ETCD_CLUSTER_TOKEN=$(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)
export ETCD_VERSION="v3.1.12" # Suggested version for Kubernetes 1.10
curl -sLO https://github.com/coreos/etcd/releases/download/${ETCD_VERSION}/etcd-${ETCD_VERSION}-linux-amd64.tar.gz
echo "Installing etcd on etcd servers"
for ((i = 0; i < ${#ETCD_HOSTNAMES[@]}; i++)); do
  echo "  ${i}: ${ETCD_PUBLIC_IPS[$i]}"
  scp -q "./etcd-${ETCD_VERSION}-linux-amd64.tar.gz" ${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]}:~/
  CMD="set -eu
  sudo mkdir -p /opt/bin/
  sudo tar -xzf ~/etcd-${ETCD_VERSION}-linux-amd64.tar.gz --strip-components=1 -C /opt/bin/
  sudo rm -rf ~/etcd-${ETCD_VERSION}-linux-amd64*
  sudo sh -c \"echo '' > /etc/etcd.env && echo PEER_NAME=${ETCD_HOSTNAMES[$i]} >> /etc/etcd.env && echo PRIVATE_IP=${ETCD_PRIVATE_IPS[$i]} >> /etc/etcd.env\""
  ssh ${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]} "$CMD"
  cat >etcd${i}.service <<EOL
[Unit]
Description=etcd
Documentation=https://github.com/coreos/etcd
Conflicts=etcd2.service

[Service]
EnvironmentFile=/etc/etcd.env
Type=notify
Restart=always
RestartSec=5s
LimitNOFILE=40000
TimeoutStartSec=0

ExecStart=/opt/bin/etcd --name ${ETCD_HOSTNAMES[$i]} \\
    --data-dir /var/lib/etcd \\
    --listen-client-urls https://${ETCD_PRIVATE_IPS[$i]}:2379 \\
    --advertise-client-urls https://${ETCD_PRIVATE_IPS[$i]}:2379 \\
    --listen-peer-urls https://${ETCD_PRIVATE_IPS[$i]}:2380 \\
    --initial-advertise-peer-urls https://${ETCD_PRIVATE_IPS[$i]}:2380 \\
    --cert-file=/etc/kubernetes/pki/etcd/server.pem \\
    --key-file=/etc/kubernetes/pki/etcd/server-key.pem \\
    --client-cert-auth \\
    --trusted-ca-file=/etc/kubernetes/pki/etcd/ca.pem \\
    --peer-cert-file=/etc/kubernetes/pki/etcd/peer.pem \\
    --peer-key-file=/etc/kubernetes/pki/etcd/peer-key.pem \\
    --peer-client-cert-auth \\
    --peer-trusted-ca-file=/etc/kubernetes/pki/etcd/ca.pem \\
    --initial-cluster ${ETCD_RING} \\
    --initial-cluster-token ${NEW_ETCD_CLUSTER_TOKEN} \\
    --initial-cluster-state new

[Install]
WantedBy=multi-user.target
EOL
        rsync --rsync-path="sudo rsync" "etcd${i}.service" "${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]}:/etc/systemd/system/etcd.service"
        # Stop and cleanup exsisting etcd Server
        ssh ${DEFAULT_LOGIN_USER}@${ETCD_PUBLIC_IPS[$i]} "sudo systemctl daemon-reload; sudo systemctl reset-failed; sudo systemctl start --no-block etcd; sudo systemctl enable etcd"
done

# Cloud Provider prerequistes
cat > ./10-hostname.conf <<EOF
[Service]
Environment="KUBELET_EXTRA_ARGS= --cloud-provider=${CLOUD_PROVIDER_FLAG} --cloud-config=/etc/kubernetes/cloud-config"
EOF

for ((i = 0; i < ${#MASTER_HOSTNAMES[@]}; i++)); do
        echo "Inject Kubelet Master Server ${i}"
        rsync --rsync-path="sudo rsync" ./10-hostname.conf "${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[$i]}:/etc/systemd/system/kubelet.service.d/10-hostname.conf"
        ssh ${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[$i]} "sudo systemctl daemon-reload"
done

# Master Server Setup
echo "Copying PKI to master servers"
for ((i = 0; i < ${#MASTER_HOSTNAMES[@]}; i++)); do
  echo "  ${i}: ${MASTER_PUBLIC_IPS[$i]}"
  ssh "${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[$i]}" sudo mkdir -p /etc/kubernetes/pki/etcd/
  rsync --rsync-path="sudo rsync" ca.pem client.pem client-key.pem "${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[$i]}:/etc/kubernetes/pki/etcd/"
done

ETCD_RING_YAML=""
for ((i = 0; i < ${#ETCD_HOSTNAMES[@]}; i++)); do
        ETCD_RING_YAML+="  - \"https://${ETCD_PRIVATE_IPS[$i]}:2379\""$'\n'
done
echo "$ETCD_RING_YAML"

SANS_RING_YAML=""
for ((i = 0; i < ${#MASTER_LOAD_BALANCER_ADDRS[@]}; i++)); do
        SANS_RING_YAML+="- \"${MASTER_LOAD_BALANCER_ADDRS[$i]}\""$'\n'
done
for ((i = 0; i < ${#MASTER_PUBLIC_IPS[@]}; i++)); do
        SANS_RING_YAML+="- \"${MASTER_PUBLIC_IPS[$i]}\""$'\n'
done
for ((i = 0; i < ${#MASTER_PRIVATE_IPS[@]}; i++)); do
        SANS_RING_YAML+="- \"${MASTER_PRIVATE_IPS[$i]}\""$'\n'
done
echo "$SANS_RING_YAML"

install_kubeadm "${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[0]}"
KUBEADM_TOKEN="$(ssh ${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[0]} "bash -l -c 'kubeadm token generate'")"

cat >kubeadm-config0.yaml <<EOL
apiVersion: kubeadm.k8s.io/v1alpha1
kind: MasterConfiguration
cloudProvider: "${CLOUD_PROVIDER_FLAG}"
kubernetesVersion: ${KUBERNETES_VERSION}
token: "${KUBEADM_TOKEN}"
tokenTTL: "0"
api:
  advertiseAddress: "${MASTER_PRIVATE_IPS[0]}"
etcd:
  endpoints:
${ETCD_RING_YAML}
  caFile: /etc/kubernetes/pki/etcd/ca.pem
  certFile: /etc/kubernetes/pki/etcd/client.pem
  keyFile: /etc/kubernetes/pki/etcd/client-key.pem
networking:
  podSubnet: "${POD_SUBNET}"
apiServerCertSANs:
${SANS_RING_YAML}
apiServerExtraArgs:
  #endpoint-reconciler-type=lease
  apiserver-count: "${#MASTER_HOSTNAMES[@]}"
  cloud-config: /etc/kubernetes/cloud-config
controllerManagerExtraArgs:
  cloud-config: /etc/kubernetes/cloud-config
EOL
rsync --rsync-path="sudo rsync" kubeadm-config0.yaml "${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[0]}:/etc/kubernetes/kubeadm-config.yaml"
rsync --rsync-path="sudo rsync" "${CLOUD_CONFIG_FILE}" "${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[0]}:/etc/kubernetes/cloud-config"
ssh ${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[0]} "sudo kubeadm init --config=/etc/kubernetes/kubeadm-config.yaml"

# Copy generated certificates back to our machine
mkdir -p apiserver0pki || true

rsync --recursive --rsync-path="sudo rsync" "${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[0]}:/etc/kubernetes/pki/" ./apiserver0pki/
rm -rf ./apiserver0pki/etcd
for ((i = 1; i < ${#MASTER_HOSTNAMES[@]}; i++)); do
        echo "Copy CA to new Master Servers ${i}"
        rsync --rsync-path="sudo rsync" ./apiserver0pki/{ca.crt,ca.key,sa.key,sa.pub} ${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[$i]}:/etc/kubernetes/pki/
        cat >kubeadm-config${i}.yaml <<EOL
apiVersion: kubeadm.k8s.io/v1alpha1
kind: MasterConfiguration
cloudProvider: "${CLOUD_PROVIDER_FLAG}"
kubernetesVersion: ${KUBERNETES_VERSION}
token: "${KUBEADM_TOKEN}"
tokenTTL: "0"
api:
  advertiseAddress: "${MASTER_PRIVATE_IPS[$i]}"
etcd:
  endpoints:
${ETCD_RING_YAML}
  caFile: /etc/kubernetes/pki/etcd/ca.pem
  certFile: /etc/kubernetes/pki/etcd/client.pem
  keyFile: /etc/kubernetes/pki/etcd/client-key.pem
networking:
  podSubnet: "${POD_SUBNET}"
apiServerCertSANs:
${SANS_RING_YAML}
apiServerExtraArgs:
  #endpoint-reconciler-type=lease
  apiserver-count: "${#MASTER_HOSTNAMES[@]}"
  cloud-config: /etc/kubernetes/cloud-config
controllerManagerExtraArgs:
  cloud-config: /etc/kubernetes/cloud-config
EOL
        rsync --rsync-path="sudo rsync" "./kubeadm-config${i}.yaml" "${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[$i]}:/etc/kubernetes/kubeadm-config.yaml"
        rsync --rsync-path="sudo rsync" "${CLOUD_CONFIG_FILE}" "${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[$i]}:/etc/kubernetes/cloud-config"
        install_kubeadm "${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[$i]}"

        ssh ${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[$i]} "sudo kubeadm init --config=/etc/kubernetes/kubeadm-config.yaml"
done

rsync --rsync-path="sudo rsync" "${DEFAULT_LOGIN_USER}@${MASTER_PUBLIC_IPS[0]}:/etc/kubernetes/admin.conf" kubeconfig
#sed -i -e 's/'"${MASTER_PRIVATE_IPS[0]}"'/'"${MASTER_LOAD_BALANCER_ADDRS[0]}"'/g' kubeconfig

# Wait for LB to be ready.
for (( i = 0; i < 10; i++ )); do
          kubectl --kubeconfig=kubeconfig apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml && break || sleep 20;
done

# Switch to LB
echo "Switch Kube-Proxy to LB addr"
kubectl --kubeconfig=kubeconfig -n kube-system get configmap kube-proxy -o yaml > kube-proxy-configmap.yaml
sed -i -e 's#server:.*#server: https://'"${MASTER_LOAD_BALANCER_ADDRS[0]}"':6443#g' kube-proxy-configmap.yaml
kubectl --kubeconfig=kubeconfig delete -f kube-proxy-configmap.yaml
kubectl --kubeconfig=kubeconfig create -f kube-proxy-configmap.yaml
kubectl --kubeconfig=kubeconfig -n kube-system delete pod -l k8s-app=kube-proxy

./install-worker.sh
