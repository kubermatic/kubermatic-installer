#!/usr/bin/env bash

set -xeu pipefail

CONFIG_FILE=${CONFIG_FILE:-./config.sh}

source $CONFIG_FILE

all_public_ips=(${MASTER_PUBLIC_IPS[@]} ${WORKER_PUBLIC_IPS[@]})
all_public_ips=($(printf "%s\n" "${all_public_ips[@]}" | sort -u))

rm -rf ./render/

# sudo with local binary directories manually added to path. Needed because some
# dirstros don't correctly set up path in non-interactive sessions, e.g. RHEL
SUDO="sudo env PATH=\$PATH:/usr/local/bin:/opt/bin"

# use generated known_hosts file if available
[ -r ./generated-known_hosts ] && export SSH_FLAGS="${SSH_FLAGS:-} -o UserKnownHostsFile=./generated-known_hosts"
! [ -r ./generated-known_hosts ] && export SSH_FLAGS="${SSH_FLAGS:-} -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no"

for i in ${!all_public_ips[*]}; do
    echo "cleanup: ${all_public_ips[$i]}"

    ssh $SSH_FLAGS ${SSH_LOGIN}@${all_public_ips[$i]} <<SSHEOF
        set -xeu pipefail

        yes|$SUDO kubeadm reset
        sudo rm -rf ~/render/ /var/lib/etcd
SSHEOF
done
