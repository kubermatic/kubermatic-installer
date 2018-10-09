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

terraform output -json > pharos_terraform.json

/usr/local/bundle/bin/pharos-cluster up -y -c "${PROVIDER}/cluster.yaml" --tf-json pharos_terraform.json
/usr/local/bundle/bin/pharos-cluster kubeconfig -y -c "${PROVIDER}/cluster.yaml" --tf-json pharos_terraform.json > generated-kubeconfig

export KUBECONFIG=$PWD/generated-kubeconfig

# -- test storage --------------------------------------------------------------
echo " *** Applying storage class"
if [[ -f "$PROVIDER/storage-class.yaml" ]]; then
  kubectl create -f "$PROVIDER/storage-class.yaml"
else
  echo "$PROVIDER/storage-class.yaml missing"
  exit 1
fi

echo " *** Applying deployment with a PVC"
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
./conformance.sh
