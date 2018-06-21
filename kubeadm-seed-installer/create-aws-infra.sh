#!/bin/bash
#
# helper script to provision a suitable set of aws instances
# for seed cluster creation.
#
# features:
#   * extracts hostkey fingerprints from console output and
#     uses them to compile a known_hosts file for ssh access
#   * configurable number of master/worker
#   * generates configuration for install.sh (with IP addresses)
#
# pre-requisite:
#   * security group named "kubermatic-test-seed-clusters" allowing intra-sg
#     traffic and tcp:22,tcp:6443 from everywhere.
#   * ssh-keypair with name $SSHKEYNAME

tag=$LOGNAME-$( date "+tmp-%F-%s" )
ltmp=$(mktemp -d) ; trap "rm -rv $ltmp" EXIT
here=`dirname $0`

die() {
	echo "$*"
	exit 1
}

type jq >& /dev/null || die "jq command required."
type aws >& /dev/null || die "aws command required."

master_count=3
worker_count=1

[ -z "$SSHKEYNAME" ] && die "please specify the ssh keypair to use in \$SSHKEYNAME"

[ -e $here/render/ ] && die "$here/render/ directory already exists. Please cleanup first!"

set -e

aws ec2 run-instances \
	--image-id ami-c7e0c82c \
	--instance-type t2.small \
	--key-name $SSHKEYNAME \
	--client-token "tok-$tag" \
	--security-groups kubermatic-test-seed-clusters \
	--count $(( master_count + worker_count )) \
	--output json --no-paginate > $ltmp/instances.json

# get array of instance IDs
instances_string=$( cat $ltmp/instances.json |jq -r '.Instances[].InstanceId' )
read -d "" -r -a instances <<< "$instances_string" || true
master_instances=( ${instances[@]:0:$master_count} )
worker_instances=( ${instances[@]:$master_count} )
master_lb_instance=${master_instances[0]}

# wait for instances..s
echo -n "Waiting for: ${instances[@]} ... "
aws ec2 wait instance-running --filters "Name=client-token,Values=tok-$tag"
printf "Done.\n"

printf "Adding a name tag to all instances... "
# add name tag - master
idx=0
for instance in ${master_instances[@]}; do
	aws ec2 create-tags --resources "$instance" --tags "Key=Name,Value=${tag}-master-$(( idx++ ))" "Key=owner-logname,Value=$LOGNAME"
done
idx=0
for instance in ${worker_instances[@]}; do
	aws ec2 create-tags --resources "$instance" --tags "Key=Name,Value=${tag}-worker-$(( idx++ ))" "Key=owner-logname,Value=$LOGNAME"
done
printf "Done.\n"

aws ec2 describe-instances \
	--filters "Name=client-token,Values=tok-$tag" \
	--output json > $ltmp/describe-instance.json

pubip_by_instanceid() {
	iid=$1; shift
	cat $ltmp/describe-instance.json \
		| jq -r '.Reservations[].Instances[] | select(.InstanceId=="'$iid'") | .PublicIpAddress'
}

privip_by_instanceid() {
	iid=$1; shift
	cat $ltmp/describe-instance.json \
		| jq -r '.Reservations[].Instances[] | select(.InstanceId=="'$iid'") | .PrivateIpAddress'
}

# get all instances' host keys --> compile known_hosts file for ssh
for instance in ${instances[@]}; do
	printf "Waiting for console output for instance $instance "
	aws ec2 get-console-output --output text --instance-id $instance > $ltmp/console-$instance
	timeout=0
	while ! grep -q -- '-----END SSH HOST KEY KEYS-----' $ltmp/console-$instance; do
		printf "... "
		sleep 15
		aws ec2 get-console-output --output text --instance-id $instance > $ltmp/console-$instance
		[ $(( timeout++ )) -gt 40 ] && break
	done
	if [ $timeout -gt 40 ]; then
		printf "\nTimeout reached while waiting for console output on [$instance].\n"
	else
		printf "Success! ($timeout)\n"
		pubip=`pubip_by_instanceid $instance`
		cat $ltmp/console-$instance \
			| awk '/-----BEGIN SSH HOST KEY KEYS-----/{x=1;next} ; /-----END SSH HOST KEY KEYS-----/{x=0} ; x {print $0}' \
			| sed "s/^/$pubip &/" \
			>> $ltmp/known_hosts
	fi
done

# master
for instance in ${master_instances[@]}; do
	MASTER_PUBLIC_IPS="$MASTER_PUBLIC_IPS `pubip_by_instanceid $instance`"
	MASTER_PRIVATE_IPS="$MASTER_PRIVATE_IPS `privip_by_instanceid $instance`"
done
MASTER_LOAD_BALANCER_ADDRS=`pubip_by_instanceid $master_lb_instance`

# worker
for instance in ${worker_instances[@]}; do
	WORKER_PUBLIC_IPS="$WORKER_PUBLIC_IPS `pubip_by_instanceid $instance`"
	WORKER_PRIVATE_IPS="$WORKER_PRIVATE_IPS `privip_by_instanceid $instance`"
done

(
# WORKER_PRIVATE_IPS=($WORKER_PRIVATE_IPS )
cat <<EOF
MASTER_PUBLIC_IPS=($MASTER_PUBLIC_IPS )
MASTER_PRIVATE_IPS=($MASTER_PRIVATE_IPS )

WORKER_PUBLIC_IPS=($WORKER_PUBLIC_IPS )

ETCD_PUBLIC_IPS=($MASTER_PUBLIC_IPS )
ETCD_PRIVATE_IPS=($MASTER_PRIVATE_IPS )
MASTER_LOAD_BALANCER_ADDRS=( $MASTER_LOAD_BALANCER_ADDRS )

export CLOUD_CONFIG_FILE=""
SSH_LOGIN="ubuntu"
export CLOUD_PROVIDER_FLAG=""
EOF
) | tee $ltmp/seed-aws-config.sh

cp -vi $ltmp/seed-aws-config.sh $here/generated-config.sh
cp -vi $ltmp/known_hosts $here/generated-known_hosts

printf "\nDone.\n"

# printf "Waiting for phony user input before removing temporary directory $ltmp\n"
# read phony
