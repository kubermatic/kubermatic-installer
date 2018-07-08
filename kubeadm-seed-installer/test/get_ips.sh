#!/usr/bin/env bash

set tw=500

set -e


export MASTER_IPS=""
for index in {0..2}; do
  IP=$(cat terraform.tfstate\
    |jq ".modules[0].resources.\"openstack_compute_floatingip_associate_v2.e2e.$index\".primary.attributes.floating_ip" -r)
  export MASTER_IPS="$MASTER_IPS $IP"
done

export WORKER_IPS=""
unset IP
for index in {0..2}; do
  IP=$(cat terraform.tfstate\
    |jq ".modules[0].resources.\"openstack_compute_floatingip_associate_v2.e2e.$index\".primary.attributes.floating_ip" -r)
  export WORKER_IPS="$WORKER_IPS $IP"
done

echo "MASTER IPS: $MASTER_IPS"
echo "WORKER_IPS: $WORKER_IPS"
