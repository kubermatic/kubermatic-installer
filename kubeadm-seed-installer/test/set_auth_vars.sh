#!/usr/bin/env bash

set -e

export VAULT_TOKEN=$(vault write --format=json auth/approle/login \
  role_id=$APPROLE_ID \
  secret_id=$SECRET_ID|jq .auth.client_token -r)

case ${PROVIDER} in
  openstack)
	SECRET_BASEPATH="dev/syseleven-openstack"
	export OS_AUTH_URL="$(vault read --field=OS_AUTH_URL $SECRET_BASEPATH)"
	export OS_REGION_NAME="$(vault read --field=OS_REGION_NAME $SECRET_BASEPATH)"
	export OS_USER_DOMAIN_NAME=Default
	export OS_IDENTITY_API_VERSION=3
	export OS_PASSWORD="$(vault read --field=password $SECRET_BASEPATH)"
	export OS_USERNAME="$(vault read --field=username $SECRET_BASEPATH)"
	export OS_PROJECT_ID="$(vault read --field=OS_PROJECT_ID $SECRET_BASEPATH)"
	export OS_TENANT_NAME="$(vault read --field=OS_TENANT_NAME $SECRET_BASEPATH)"
	;;
  aws)
	SECRET_BASEPATH="dev/e2e_testing_credentials"
	export AWS_ACCESS_KEY_ID="$(vault read --field=AWS_ACCESS_KEY_ID $SECRET_BASEPATH)"
	export AWS_SECRET_ACCESS_KEY="$(vault read --field=AWS_SECRET_ACCESS_KEY $SECRET_BASEPATH)"
	;;
  *)
	echo "Cloud provider ${PROVIDER} not yet implemented"
    exit 1
  ;;
esac

vault read --field=key dev/machine-controller-ssh-key > machine-key
chmod 0600 machine-key

# generate the pubkey
ssh-keygen  -y -f machine-key > machine-key.pub
