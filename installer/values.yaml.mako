<%! import os %>
<%! from base64 import b64encode %>

<%include file="dockercfgjson.mako" />

### Kubermatic
kubermatic:
  docker:
    secret: "${var._dockercfgjson}"
  auth:
    tokenIssuer: https://${var.base_domain}/dex
    clientID: ${var.kubermatic_client_id}
  datacenters: "${read_file('datacenters.yaml') | b64encode}"
  domain: ${var.base_domain}
  kubeconfig: "${read_file('kubeconfig') | b64encode}"

### Storage
storage:
  provider: ${var.storage_provider}
  zone: ${var.storage_zone}
  type: ${var.storage_type}

### Nginx definition
nginx:
  hostNetwork: true
  asDaemonSet: true

certificates:
  domains:
  - ${var.base_domain}
  - alertmanager.${var.base_domain}
  - grafana.${var.base_domain}
  - prometheus.${var.base_domain}

### Monitoring
prometheus:
  auth: '${b64encode(var.prometheus_username + ":" + var.prometheus_password)}'
  host: prometheus.${var.base_domain}

grafana:
  user: '${b64encode(var.grafana_username)}'
  password: '${b64encode(var.grafana_password)}'
  host: grafana.${var.base_domain}

dex:
  ingress:
    host: ${var.base_domain}
  clients:
  - id: ${var.kubermatic_client_id}
    name: Kubermatic
    secret: ${var.dex_secret | b64encode}
    RedirectURIs:
    - http://localhost:8000
    - http://localhost:8000/clusters
    - https://${var.base_domain}
    - https://${var.base_domain}/clusters
  connectors:
  - type: github
    id: github
    name: GitHub
    config:
      clientID: ${var.gh_client_id}
      clientSecret: ${var.gh_client_secret}
      redirectURI: https://${var.base_domain}/dex/callback
      orgs:
      - name: ${var.gh_orga_name}

<%include file="values_more.yaml.mako" />
