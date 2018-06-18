#!/usr/bin/env bash

set -xeu pipefail

source ./config.sh

all_public_ips=(${MASTER_PUBLIC_IPS[@]} ${WORKER_PUBLIC_IPS[@]})
all_public_ips=($(printf "%s\n" "${all_public_ips[@]}" | sort -u))

rm -rf ./render/

# sudo with local binary directories manually added to path. Needed because some
# dirstros don't correctly set up path in non-interactive sessions, e.g. RHEL
SUDO="sudo env PATH=\$PATH:/usr/local/bin:/opt/bin"

for i in ${!all_public_ips[*]}; do
    echo "cleanup: ${all_public_ips[$i]}"

    ssh ${SSH_LOGIN}@${all_public_ips[$i]} <<SSHEOF
        set -xeu pipefail

        $SUDO kubeadm reset
        sudo rm -rf ~/render/ /var/lib/etcd
SSHEOF
done
