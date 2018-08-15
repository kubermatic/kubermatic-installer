#!/usr/bin/env bash
# vim: tw=500

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
cd "$SCRIPTDIR"

source set_auth_vars.sh

function cleanup {
  cd $STATEFILE_DIR
  kubectl delete pvc redis-datadir || true
  terraform destroy -auto-approve
}
trap cleanup EXIT SIGINT SIGKILL

set -e

export STATEFILE_DIR=$PWD
terraform init
terraform apply --auto-approve

export MASTER_PUBLIC_IPS=""
for index in {0..2}; do
  IP=$(cat terraform.tfstate\
    |jq ".modules[0].resources.\"openstack_compute_floatingip_associate_v2.e2e.$index\".primary.attributes.floating_ip" -r)
  export MASTER_PUBLIC_IPS="$MASTER_PUBLIC_IPS $IP"
done
export MASTER_PRIVATE_IPS=""
for index in {0..2}; do
  IP=$(cat terraform.tfstate\
    |jq ".modules[0].resources.\"openstack_compute_floatingip_associate_v2.e2e.$index\".primary.attributes.fixed_ip" -r)
  export MASTER_PRIVATE_IPS="$MASTER_PRIVATE_IPS $IP"
done

export WORKER_IPS=""
unset IP
for index in {3..5}; do
  IP=$(cat terraform.tfstate\
    |jq ".modules[0].resources.\"openstack_compute_floatingip_associate_v2.e2e.$index\".primary.attributes.floating_ip" -r)
  export WORKER_IPS="$WORKER_IPS $IP"
done

# This must be the first ip if its not a real loadbalancer that does healthchecking
LB_IP=$(cat terraform.tfstate\
    |jq ".modules[0].resources.\"openstack_compute_floatingip_associate_v2.e2e.0\".primary.attributes.floating_ip" -r)

test -e config.sh || cp ../config-example.sh config.sh

sed -i "s#^MASTER_PUBLIC_IPS.*#MASTER_PUBLIC_IPS=($MASTER_PUBLIC_IPS)#g" config.sh
sed -i "s#^MASTER_PRIVATE_IPS.*#MASTER_PRIVATE_IPS=($MASTER_PRIVATE_IPS)#g" config.sh
sed -i "s#^WORKER_PUBLIC_IPS.*#WORKER_PUBLIC_IPS=($WORKER_IPS)#g" config.sh
sed -i "s#^MASTER_LOAD_BALANCER_ADDRS.*#MASTER_LOAD_BALANCER_ADDRS=($LB_IP)#g" config.sh
sed -i "s#^SSH_LOGIN.*#SSH_LOGIN=ubuntu#g" config.sh
sed -i "s#^export CLOUD_PROVIDER_FLAG.*#export CLOUD_PROVIDER_FLAG=openstack#g" config.sh
sed -i "s#^export CLOUD_CONFIG_FILE.*#export CLOUD_CONFIG_FILE=/test/cloud.conf#g" config.sh

cp ../cloudconfig-openstack.sample.conf cloud.conf

sed -i "s#<< OS_AUTH_URL >>#${OS_AUTH_URL}#g" cloud.conf
sed -i "s#<< OS_USERNAME >>#${OS_USERNAME}#g" cloud.conf
sed -i "s#<< OS_PASSWORD >>#${OS_PASSWORD}#g" cloud.conf
sed -i "s#<< OS_DOMAIN_NAME >>#${OS_USER_DOMAIN_NAME}#g" cloud.conf
sed -i "s#<< OS_TENANT_NAME >>#${OS_TENANT_NAME}#g" cloud.conf
sed -i "s#<< OS_REGION_NAME >>#${OS_REGION_NAME}#g" cloud.conf

export CONFIG_FILE=$PWD/config.sh

echo "Successfully generated config, installing cluster"
cd ..

timeout=0
while ! ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no ubuntu@$LB_IP true; do
  if [ $(( timeout++ )) -gt 10 ]; then echo "Failed to connect via ssh!"; exit 1; fi
  sleep 5
done

chmod +x install.sh
./install.sh

sleep 20

source $CONFIG_FILE

[ -r ./generated-known_hosts ] && export SSH_FLAGS="${SSH_FLAGS:-} -o UserKnownHostsFile=./generated-known_hosts"
! [ -r ./generated-known_hosts ] && export SSH_FLAGS="${SSH_FLAGS:-} -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no"

scp ${SSH_FLAGS} ${SSH_LOGIN}@${MASTER_PUBLIC_IPS[0]}:./.kube/config ./generated-kubeconfig
export KUBECONFIG=$PWD/generated-kubeconfig

# -- test storage --------------------------------------------------------------
cat <<EOF | kubectl create -f -
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: foosc
provisioner: kubernetes.io/cinder
reclaimPolicy: Delete
EOF

cat <<EOF | kubectl create -f -
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: redis-datadir
spec:
  storageClassName: foosc
  accessModes: [ "ReadWriteOnce" ]
  resources:
    requests:
      storage: "20Gi"
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: redis-master
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: redis
        role: master
        tier: backend
    spec:
      containers:
      - name: master
        image: redis:4
        ports:
        - containerPort: 6379
        volumeMounts:
        - name: datadir
          mountPath: /data
      volumes:
        - name: datadir
          persistentVolumeClaim:
            claimName: redis-datadir
EOF

INTERVAL=30
TIMEOUT=300 # 5 minutes
ELAPSED=0

while true
do
  if [ "$ELAPSED" -gt "$TIMEOUT" ]; then
    echo "ERROR: storage test deployment didnt come up in time."
    exit 1
  fi

  AVAILABLE=$(kubectl get deployment redis-master -o 'jsonpath={.status.availableReplicas}' 2>&1 || echo '0')

  if ((AVAILABLE > 0)); then
    echo "storage test succeeded, pvc bound after $ELAPSED seconds."
    break
  fi

  sleep $INTERVAL
  ELAPSED=$((ELAPSED + INTERVAL))
  echo "waiting for pvc to be bound since $ELAPSED seconds"
done

# -- run conformance testsuite (sonobuoy) --------------------------------------
chmod +x ./test/conformance.sh
./test/conformance.sh
