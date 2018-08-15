#!/usr/bin/env bash

set -e

export VAULT_TOKEN=$(vault write --format=json auth/approle/login \
  role_id=$APPROLE_ID \
  secret_id=$SECRET_ID|jq .auth.client_token -r)

export OS_AUTH_URL="$(vault read --field=OS_AUTH_URL $SECRET_BASEPATH)"
export OS_REGION_NAME="$(vault read --field=OS_REGION_NAME $SECRET_BASEPATH)"
export OS_USER_DOMAIN_NAME=Default
export OS_IDENTITY_API_VERSION=3
export OS_PASSWORD="$(vault read --field=password $SECRET_BASEPATH)"
export OS_USERNAME="$(vault read --field=username $SECRET_BASEPATH)"
export OS_PROJECT_ID="$(vault read --field=OS_PROJECT_ID $SECRET_BASEPATH)"
export OS_TENANT_NAME="$(vault read --field=OS_TENANT_NAME $SECRET_BASEPATH)"

mkdir -m 0700 -p ~/.ssh/
vault read --field=key dev/machine-controller-ssh-key > ~/.ssh/id_rsa
chmod 0600 ~/.ssh/id_rsa
