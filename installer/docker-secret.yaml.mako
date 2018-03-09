<%include file="dockercfgjson.mako" />

apiVersion: v1
kind: Secret
type: kubernetes.io/dockerconfigjson
metadata:
  name: dockercfg
  namespace: kubermatic-installer
data:
  .dockerconfigjson: "${var._dockercfgjson}"
