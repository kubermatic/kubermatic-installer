<%! import os %>
<%include file="variables.mako" />
% if os.path.isfile("variables_override.mako"):
<%include file="variables_override.mako" />
% endif

apiVersion: batch/v1
kind: Job
metadata:
  generateName: installer
  namespace: kubermatic-installer
spec:
  template:
    metadata:
      name: kubermatic-installer
    spec:
      serviceAccountName: kubermatic-installer
      containers:
      - name: kubermatic-installer
        image: kubermatic/installer:${var.kubermatic_installer_tag}
        command:
        - /kubermatic/deploy_helm_charts.sh
        args:
        - /kubermatic/values/values.yaml
        - /kubermatic
        - /kubermatic/versions-values.yaml
        imagePullPolicy: Always
        volumeMounts:
          - name: values
            mountPath: /kubermatic/values
            readOnly: true
      restartPolicy: Never
      imagePullSecrets:
        - name: dockercfg
      volumes:
      - name: values
        secret:
          secretName: values