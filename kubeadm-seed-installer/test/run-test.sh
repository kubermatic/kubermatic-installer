#!/usr/bin/env bash
#set tw=500

function cleanup {
  OLD_EXIT_CODE=$?
  terraform destroy -auto-approve
  exit $OLD_EXIT_CODE
}
trap cleanup EXIT

set -e

terraform init
terraform apply --auto-approve


export MASTER_PUBLIC_IPS=""
for index in {0..2}; do
  IP=$(cat terraform.tfstate\
    |jq ".modules[0].resources.\"openstack_compute_floatingip_associate_v2.e2e.$index\".primary.attributes.floating_ip" -r)
  if [[ $MASTER_PUBLIC_IPS == "" ]]; then
    export MASTER_PUBLIC_IPS=$IP;
  else
    export MASTER_PUBLIC_IPS="$MASTER_PUBLIC_IPS $IP"
  fi
done
export MASTER_PRIVATE_IPS=""
for index in {0..2}; do
  IP=$(cat terraform.tfstate\
    |jq ".modules[0].resources.\"openstack_compute_floatingip_associate_v2.e2e.$index\".primary.attributes.fixed_ip" -r)
  if [[ $MASTER_PRIVATE_IPS == "" ]]; then
    export MASTER_PRIVATE_IPS=$IP;
  else
    export MASTER_PRIVATE_IPS="$MASTER_PRIVATE_IPS $IP"
  fi
done

export WORKER_IPS=""
unset IP
for index in {3..5}; do
  IP=$(cat terraform.tfstate\
    |jq ".modules[0].resources.\"openstack_compute_floatingip_associate_v2.e2e.$index\".primary.attributes.floating_ip" -r)
  if [[ $WORKER_IPS == "" ]]; then export WORKER_IPS=$IP; else export WORKER_IPS="$WORKER_IPS $IP"; fi
done

# This must be the first ip if its not a real loadbalancer that does healthchecking
LB_IP=$(cat terraform.tfstate\
    |jq ".modules[0].resources.\"openstack_compute_floatingip_associate_v2.e2e.0\".primary.attributes.floating_ip" -r)

if ! ls config.sh &>/dev/null; then cp ../config.sh .; fi

sed -i "s#MASTER_PUBLIC_IPS.*#MASTER_PUBLIC_IPS=($MASTER_PUBLIC_IPS)#g" config.sh
sed -i "s#MASTER_PRIVATE_IPS.*#MASTER_PRIVATE_IPS=($MASTER_PRIVATE_IPS)#g" config.sh
sed -i "s#WORKER_PUBLIC_IPS.*#WORKER_PUBLIC_IPS=($WORKER_IPS)#g" config.sh
sed -i "s#MASTER_LOAD_BALANCER_ADDRS.*#MASTER_LOAD_BALANCER_ADDRS=($LB_IP)#g" config.sh
sed -i "s#SSH_LOGIN.*#SSH_LOGIN=ubuntu#g" config.sh

export CONFIG_FILE=$PWD/config.sh

echo "Successfully generated config, installing cluster"
cd ..

./install.sh
