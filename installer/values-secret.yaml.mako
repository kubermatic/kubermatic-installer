<%! from base64 import b64encode %>
apiVersion: v1
kind: Secret
metadata:
  name: values
  namespace: kubermatic-installer
type: Opaque
data:
  values.yaml: "${read_file('values.yaml') | b64encode}"
