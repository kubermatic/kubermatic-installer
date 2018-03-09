<%
var.kubermatic_installer_tag = 'v2.3.4'

var.base_domain = 'myk8s.de'

var.kubermatic_client_id = 'kubermatic'

var.gh_orga_name = 'my-orga'
var.gh_client_id = 'xxxxxxxxx'
var.gh_client_secret = 'xxxxxxxxxxxxxx'

var.docker_username = 'somebody'
var.docker_password = 'somepass'
var.docker_email = 'somebody@somewhere.com'

var.prometheus_username = "admin"
var.prometheus_password = "P6QVQLj8RmGRhxCLq74dVylXqyca1lu9"

var.grafana_username = "admin"
var.grafana_password = "P6QVQLj8RmGRhxCLq74dVylXqyca1lu9"

var.dex_secret = 'example-app-secret'

# paths of seed cluster kubeconfig files.
# must contain one entry for each seed cluster definition in datacenters.yaml,
# in the same order
var.seed_configs = [
  "../kubeadm-seed-installer/kubeconfig"
]

%>
