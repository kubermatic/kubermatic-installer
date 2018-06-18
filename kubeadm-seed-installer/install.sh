#!/usr/bin/env bash
set -xeu
set -o pipefail

source ./config.sh

export APISERVER_COUNT=${#MASTER_PUBLIC_IPS[*]}
export APISERVER_SANS_YAML=""
export ETCD_RING_SANS=""
export ETCD_RING_YAML=""
export ETCD_RING=""
export NEW_ETCD_CLUSTER_TOKEN=$(openssl rand -base64 32 | tr -dc 'a-zA-Z0-9' | fold -w 32 | head -n 1)

# constants
readonly POD_SUBNET="10.244.0.0/16"   # for flannel
readonly CNI_VERSION="v0.7.1"         # for coreos

export POD_SUBNET
export CNI_VERSION

export NODEPORT_RANGE=${NODEPORT_RANGE:-30000-32767}

SCRIPT_DIR="$(realpath "$(dirname "${BASH_SOURCE[0]}")")"
OFFLINE="false"

# sudo with local binary directories manually added to path. Needed because some
# dirstros don't correctly set up path in non-interactive sessions, e.g. RHEL
SUDO="sudo env PATH=\$PATH:/usr/local/bin:/opt/bin"

while [ $# -gt 0 ]; do
  case "$1" in
    --offline)
      OFFLINE="true"
    ;;
    *)
      echo "Unknown parameter \"$1\""
      exit 1
    ;;
  esac
  shift
done

kubeadm_install() {
    local SSHDEST=$1
    local OS_ID=$(ssh ${SSHDEST} "cat /etc/os-release" | grep '^ID=' | sed s/^ID=//)

    case $OS_ID in
        ubuntu|debian)
            kubeadm_install_deb ${SSHDEST}
        ;;
        coreos)
            kubeadm_install_coreos ${SSHDEST}
        ;;
        centos)
            kubeadm_install_centos ${SSHDEST}
        ;;
        *)
            echo " ### Operating system '$OS_ID' is not supported."
            exit 1
        ;;
    esac
}

kubeadm_install_deb() {
    local SSHDEST=$1

    ssh ${SSHDEST} <<SSHEOF
        set -xeu pipefail
        sudo swapoff -a

        source /etc/os-release

        sudo apt-get update
        sudo apt-get install -y --no-install-recommends \
            apt-transport-https \
            ca-certificates \
            curl \
            htop \
            lsb-release \
            rsync \
            tree

        curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
        curl -fsSL https://download.docker.com/linux/\${ID}/gpg | sudo apt-key add -

        echo "deb [arch=amd64] https://download.docker.com/linux/\${ID} \$(lsb_release -sc) stable" | \
            sudo tee /etc/apt/sources.list.d/docker.list

        # You'd think that kubernetes-\$(lsb_release -sc) belongs there instead, but the debian repo
        # contains neither kubeadm nor kubelet, and the docs themselves suggest using xenial repo.
        echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" | \
            sudo tee /etc/apt/sources.list.d/kubernetes.list
        sudo apt-get update

        docker_ver=\$(apt-cache madison docker-ce | grep ${DOCKER_VERSION} | head -1 | awk '{print \$3}')
        kube_ver=\$(apt-cache madison kubelet | grep ${KUBERNETES_VERSION} | head -1 | awk '{print \$3}')

        sudo apt-mark unhold docker-ce kubelet kubeadm kubectl
        sudo apt-get install -y --no-install-recommends \
            docker-ce=\${docker_ver} \
            kubeadm \
            kubectl=\${kube_ver} \
            kubelet=\${kube_ver}
        sudo apt-mark hold docker-ce kubelet kubeadm kubectl
        sudo systemctl daemon-reload
SSHEOF
}

kubeadm_install_coreos() {
    local SSHDEST=$1

    ssh ${SSHDEST} << SSHEOF
        set -xeu pipefail

        sudo mkdir -p /opt/cni/bin /etc/kubernetes/pki /etc/kubernetes/manifests
        curl -L \
            "https://github.com/containernetworking/plugins/releases/download/${CNI_VERSION}/cni-plugins-amd64-${CNI_VERSION}.tgz" | \
            sudo tar -C /opt/cni/bin -xz

        RELEASE="v${KUBERNETES_VERSION}"

        sudo mkdir -p /opt/bin
        cd /opt/bin
        sudo curl -L --remote-name-all \
            https://storage.googleapis.com/kubernetes-release/release/\${RELEASE}/bin/linux/amd64/{kubeadm,kubelet,kubectl}
        sudo chmod +x {kubeadm,kubelet,kubectl}

        curl -sSL "https://raw.githubusercontent.com/kubernetes/kubernetes/\${RELEASE}/build/debs/kubelet.service" | \
            sed "s:/usr/bin:/opt/bin:g" | \
            sudo tee /etc/systemd/system/kubelet.service
        sudo mkdir -p /etc/systemd/system/kubelet.service.d
        curl -sSL "https://raw.githubusercontent.com/kubernetes/kubernetes/\${RELEASE}/build/debs/10-kubeadm.conf" | \
            sed "s:/usr/bin:/opt/bin:g" | \
            sudo tee /etc/systemd/system/kubelet.service.d/10-kubeadm.conf

        sudo systemctl daemon-reload
        sudo systemctl enable docker.service kubelet.service
        sudo systemctl start docker.service kubelet.service
SSHEOF
}

kubeadm_install_centos() {
    echo " ### TODO support CentOS"
    exit 1
}

all_ips=(${MASTER_PUBLIC_IPS[*]} ${WORKER_PUBLIC_IPS[*]})
all_ips=($(printf "%s\n" "${all_ips[*]}" | sort -u))
all_master_ips=(${MASTER_LOAD_BALANCER_ADDRS[*]} ${MASTER_PUBLIC_IPS[*]} ${MASTER_PRIVATE_IPS[*]})
all_master_ips=($(printf "%s\n" "${all_master_ips[*]}" | sort -u))

for i in ${!MASTER_PRIVATE_IPS[*]}; do
    ETCD_RING+="etcd-${i}=https://${MASTER_PRIVATE_IPS[$i]}:2380,"
    ETCD_RING_YAML+="  - https://${MASTER_PRIVATE_IPS[$i]}:2379"$'\n'
    ETCD_RING_SANS+="  - ${MASTER_PRIVATE_IPS[$i]}"$'\n'
done
ETCD_RING="$(echo ${ETCD_RING} | sed 's/[,]*$//')"

for san in ${all_master_ips[*]}; do
    APISERVER_SANS_YAML+="- ${san}"$'\n'
done

kubeadm_config_template='
apiVersion: kubeadm.k8s.io/v1alpha1
kind: MasterConfiguration
cloudProvider: "${CLOUD_PROVIDER_FLAG}"
kubernetesVersion: v${KUBERNETES_VERSION}
api:
  advertiseAddress: "${advertiseAddress}"
  controlPlaneEndpoint: "${controlPlaneEndpoint}"
etcd:
  endpoints:
${ETCD_RING_YAML}
  caFile: /etc/kubernetes/pki/etcd/ca.crt
  certFile: /etc/kubernetes/pki/etcd/peer.crt
  keyFile: /etc/kubernetes/pki/etcd/peer.key
  serverCertSANs:
${ETCD_RING_SANS}
  peerCertSANs:
${ETCD_RING_SANS}
networking:
  podSubnet: ${POD_SUBNET}
apiServerCertSANs:
${APISERVER_SANS_YAML}
apiServerExtraArgs:
  endpoint-reconciler-type: lease
  service-node-port-range: ${NODEPORT_RANGE}
'

mkdir -p ./render/pki ./render/etcd ./render/cfg
touch ./render/cfg/cloud-config

if [ ! -z "${CLOUD_CONFIG_FILE}" ]; then
    cp "${CLOUD_CONFIG_FILE}" ./render/cfg/cloud-config
    kubeadm_config_template+='
  cloud-config: /etc/kubernetes/cloud-config
controllerManagerExtraArgs:
  cloud-config: /etc/kubernetes/cloud-config
'
fi

etcd_manifest_template='
apiVersion: v1
kind: Pod
metadata:
  annotations:
    scheduler.alpha.kubernetes.io/critical-pod: ""
  labels:
    component: etcd
    tier: control-plane
  name: etcd
  namespace: kube-system
spec:
  containers:
  - name: etcd
    command:
    - etcd
    - --advertise-client-urls=https://${etcd_ip}:2379
    - --cert-file=/etc/kubernetes/pki/etcd/server.crt
    - --client-cert-auth=true
    - --data-dir=/var/lib/etcd
    - --initial-advertise-peer-urls=https://${etcd_ip}:2380
    - --initial-cluster=${ETCD_RING}
    - --initial-cluster-state=new
    - --initial-cluster-token=${NEW_ETCD_CLUSTER_TOKEN}
    - --key-file=/etc/kubernetes/pki/etcd/server.key
    - --listen-client-urls=https://${etcd_ip}:2379
    - --listen-peer-urls=https://${etcd_ip}:2380
    - --name=${etcd_name}
    - --peer-cert-file=/etc/kubernetes/pki/etcd/peer.crt
    - --peer-client-cert-auth=true
    - --peer-key-file=/etc/kubernetes/pki/etcd/peer.key
    - --peer-trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt
    - --trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt
    image: k8s.gcr.io/etcd-amd64:${ETCD_VERSION}
    volumeMounts:
    - mountPath: /var/lib/etcd
      name: etcd-data
    - mountPath: /etc/kubernetes/pki/etcd
      name: etcd-certs
  hostNetwork: true
  volumes:
  - hostPath:
      path: /var/lib/etcd
      type: DirectoryOrCreate
    name: etcd-data
  - hostPath:
      path: /etc/kubernetes/pki/etcd
      type: DirectoryOrCreate
    name: etcd-certs
'

cat > render/cfg/20-cloudconfig-kubelet.conf <<EOF
[Service]
Environment="KUBELET_EXTRA_ARGS= --cloud-provider=${CLOUD_PROVIDER_FLAG} --cloud-config=/etc/kubernetes/cloud-config"
EOF

# render configs
export advertiseAddress="${MASTER_PRIVATE_IPS[0]}"
export controlPlaneEndpoint="${MASTER_LOAD_BALANCER_ADDRS[0]}"

echo "$kubeadm_config_template" | envsubst > render/cfg/master.yaml

for i in ${!MASTER_PRIVATE_IPS[*]}; do
    export etcd_ip=${MASTER_PRIVATE_IPS[$i]}
    export etcd_name="etcd-${i}"
    echo "$etcd_manifest_template" | envsubst > render/etcd/etcd_${i}.yaml
done

# install prerequisites on all nodes
for sshaddr in ${all_ips[*]}; do
    if [[ "$OFFLINE" != "true" ]]; then
      kubeadm_install "${SSH_LOGIN}@${sshaddr}"
    fi
    rsync -av ./render ${SSH_LOGIN}@${sshaddr}:

    ssh ${SSH_LOGIN}@${sshaddr} <<SSHEOF
        set -xeu pipefail

        sudo mkdir -p /etc/systemd/system/kubelet.service.d/ /etc/kubernetes
        sudo mv ./render/cfg/20-cloudconfig-kubelet.conf /etc/systemd/system/kubelet.service.d/
        sudo mv ./render/cfg/cloud-config /etc/kubernetes/cloud-config
        sudo chown root:root /etc/kubernetes/cloud-config
        sudo chmod 600 /etc/kubernetes/cloud-config
SSHEOF
done

rsync -av ./render ${SSH_LOGIN}@${MASTER_PUBLIC_IPS[0]}:

ssh ${SSH_LOGIN}@${MASTER_PUBLIC_IPS[0]} <<SSHEOF
    set -xeu pipefail

    $SUDO kubeadm alpha phase certs ca --config=./render/cfg/master.yaml
    $SUDO kubeadm alpha phase certs etcd-ca --config=./render/cfg/master.yaml
    $SUDO kubeadm alpha phase certs sa --config=./render/cfg/master.yaml
    sudo rsync -av /etc/kubernetes/pki/ ./render/pki/
    sudo chown -R $SSH_LOGIN ./render
SSHEOF

# download generated CAa
rsync -av ${SSH_LOGIN}@${MASTER_PUBLIC_IPS[0]}:render/pki/ ./render/pki/

# at first run: configure kubelet and establish ETCD ring
for i in ${!MASTER_PUBLIC_IPS[*]}; do
    rsync -av ./render ${SSH_LOGIN}@${MASTER_PUBLIC_IPS[$i]}:
    ssh ${SSH_LOGIN}@${MASTER_PUBLIC_IPS[$i]} <<SSHEOF
        set -xeu pipefail

        sudo rsync -av ./render/pki/ /etc/kubernetes/pki/
        rm -rf ./render/pki
        sudo chown -R root:root /etc/kubernetes/pki
        sudo mkdir -p /etc/kubernetes/manifests
        sudo cp ./render/etcd/etcd_${i}.yaml /etc/kubernetes/manifests/etcd.yaml
        $SUDO kubeadm alpha phase certs etcd-healthcheck-client --config=./render/cfg/master.yaml
        $SUDO kubeadm alpha phase certs etcd-peer --config=./render/cfg/master.yaml
        $SUDO kubeadm alpha phase certs etcd-server --config=./render/cfg/master.yaml
        $SUDO kubeadm alpha phase kubeconfig kubelet --config=./render/cfg/master.yaml
        sudo systemctl restart kubelet
SSHEOF
done

# establish everything else
for sshaddr in ${MASTER_PUBLIC_IPS[*]}; do
    ssh ${SSH_LOGIN}@${sshaddr} <<SSHEOF
        set -xeu
        $SUDO kubeadm init --config=./render/cfg/master.yaml \
          --ignore-preflight-errors=Port-10250,FileAvailable--etc-kubernetes-manifests-etcd.yaml,FileExisting-crictl
SSHEOF
done

sleep 30;

# put the value of QUAY_IO_MIRROR into the flannel YAML template
FLANNEL_YAML="$(sed 's/QUAY_IO_MIRROR/'"$QUAY_IO_MIRROR"'/' $SCRIPT_DIR/kube-flannel.yml)"

ssh ${SSH_LOGIN}@${MASTER_PUBLIC_IPS[0]} <<SSHEOF
    set -xeu pipefail

    mkdir -p ~/.kube
    sudo cp /etc/kubernetes/admin.conf ~/.kube/config
    sudo chown -R \$(id -u):\$(id -g) ~/.kube

    echo "$FLANNEL_YAML" | kubectl apply -f -

    kubectl -n kube-system get configmap kube-proxy -o yaml > kube-proxy-configmap.yaml
    sed -i -e 's#server:.*#server: https://'"${MASTER_LOAD_BALANCER_ADDRS[0]}"':6443#g' kube-proxy-configmap.yaml
    kubectl delete -f kube-proxy-configmap.yaml
    kubectl create -f kube-proxy-configmap.yaml
    kubectl -n kube-system delete pod -l k8s-app=kube-proxy
SSHEOF

sleep 10;

JOINTOKEN=$(ssh ${SSH_LOGIN}@${MASTER_PUBLIC_IPS[0]} "$SUDO kubeadm token create --print-join-command")

for sshaddr in ${WORKER_PUBLIC_IPS[*]}; do
    ssh ${SSH_LOGIN}@${sshaddr} "sudo ${JOINTOKEN}"
done
