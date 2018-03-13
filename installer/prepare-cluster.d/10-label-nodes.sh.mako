#!/bin/sh

<%!import os %>
<%include file="installer/variables.mako" />

% if os.path.isfile("../variables_override.mako"):
<%include file="installer/variables_override.mako" />
% endif

set -e

ZONE=${var.storage_zone}
KUBECTL="kubectl --kubeconfig=$(dirname $0)/../../kubeadm-seed-installer/kubeconfig"

$KUBECTL get node -o json | jq -r '.items[].metadata.name' | xargs -I__node $KUBECTL label node __node --overwrite=true failure-domain.beta.kubernetes.io/zone=$ZONE
