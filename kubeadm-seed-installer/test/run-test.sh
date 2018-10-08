#!/usr/bin/env bash
# vim: tw=500

set -eu
set -o pipefail

SCRIPTDIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null && pwd )"
cd "$SCRIPTDIR"

PROVIDER="$1"
if [ "$PROVIDER" != "openstack" ] && [ "$PROVIDER" != "aws" ] ; then
  echo "Unknown cloud provider $PROVIDER"
  exit 1
fi

source set_auth_vars.sh "$PROVIDER"

function cleanup {
  cd $STATEFILE_DIR
  kubectl delete pvc redis-datadir || true
  terraform destroy -auto-approve "./${PROVIDER}"
}
trap cleanup EXIT SIGINT

export STATEFILE_DIR=$PWD
terraform init "${PROVIDER}"
terraform apply --auto-approve "${PROVIDER}"

export MASTER_PUBLIC_IPS="$(terraform output master_public_ips)"
export MASTER_PRIVATE_IPS="$(terraform output master_private_ips)"
export WORKER_IPS="$(terraform output worker_ips)"

# This must be the first ip if its not a real loadbalancer that does healthchecking
LB_IP=$(terraform output loadbalancer_addr)

test -e config.sh || cp ../config-example.sh config.sh

sed -i "s#^MASTER_PUBLIC_IPS.*#MASTER_PUBLIC_IPS=($MASTER_PUBLIC_IPS)#g" config.sh
sed -i "s#^MASTER_PRIVATE_IPS.*#MASTER_PRIVATE_IPS=($MASTER_PRIVATE_IPS)#g" config.sh
sed -i "s#^WORKER_PUBLIC_IPS.*#WORKER_PUBLIC_IPS=($WORKER_IPS)#g" config.sh
sed -i "s#^MASTER_LOAD_BALANCER_ADDRS.*#MASTER_LOAD_BALANCER_ADDRS=($LB_IP)#g" config.sh
sed -i "s#^SSH_LOGIN.*#SSH_LOGIN=ubuntu#g" config.sh
sed -i "s#^export CLOUD_PROVIDER_FLAG.*#export CLOUD_PROVIDER_FLAG=${PROVIDER}#g" config.sh
sed -i "s#^export CLOUD_CONFIG_FILE.*#export CLOUD_CONFIG_FILE=/test/cloud.conf#g" config.sh

case ${PROVIDER} in
  openstack)
    cp ../cloudconfig-openstack.sample.conf cloud.conf

    sed -i "s#<< OS_AUTH_URL >>#${OS_AUTH_URL}#g" cloud.conf
    sed -i "s#<< OS_USERNAME >>#${OS_USERNAME}#g" cloud.conf
    sed -i "s#<< OS_PASSWORD >>#${OS_PASSWORD}#g" cloud.conf
    sed -i "s#<< OS_DOMAIN_NAME >>#${OS_USER_DOMAIN_NAME}#g" cloud.conf
    sed -i "s#<< OS_TENANT_NAME >>#${OS_TENANT_NAME}#g" cloud.conf
    sed -i "s#<< OS_REGION_NAME >>#${OS_REGION_NAME}#g" cloud.conf
  ;;
  aws)
    cp ../cloudconfig-aws.sample.conf cloud.conf

    sed -i "s#<< AWS_AVAILABILITY_ZONE >>#$(terraform output availability_zone)#g" cloud.conf
    sed -i "s#<< AWS_VPC >>#$(terraform output vpc)#g" cloud.conf
    sed -i "s#<< NAME_OF_YOUR_CLUSTER >>#installer-e2e-test-cluster#g" cloud.conf
    sed -i "s#<< AWS_SUBNET_ID >>#$(terraform output subnet)#g" cloud.conf
    sed -i "s#<< AWS_ROUTE_TABLE_ID >>#$(terraform output route_table)#g" cloud.conf
  ;;
  *)
    echo "Cloud provider ${PROVIDER} not yet implemented"
    exit 1
  ;;
esac

export CONFIG_FILE=$PWD/config.sh

echo "Successfully generated config, installing cluster"
cd ..

for MASTER_IP in "$MASTER_PUBLIC_IPS"; do
  timeout=0
  while ! ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no ubuntu@$MASTER_IP true; do
    if [ $(( timeout++ )) -gt 10 ]; then echo "Failed to connect via ssh!"; exit 1; fi
    sleep 5
  done
done

./install.sh

sleep 20

source $CONFIG_FILE

if [ -r ./generated-known_hosts ]; then
  export SSH_FLAGS="${SSH_FLAGS:-} -o UserKnownHostsFile=./generated-known_hosts"
else
  export SSH_FLAGS="${SSH_FLAGS:-} -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no"
fi

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
TIMEOUT=600
ELAPSED=0

while true
do
  if [ "$ELAPSED" -gt "$TIMEOUT" ]; then
    echo "ERROR: storage test deployment didnt come up in time, see:"
    echo "PVC Description:"
    echo "$(kubectl describe pvc redis-datadir)"
    echo "SC Description:"
    echo "$(kubectl describe sc foosc)"

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
./test/conformance.sh
