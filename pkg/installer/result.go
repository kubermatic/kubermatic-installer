package installer

import "github.com/kubermatic/kubermatic-installer/pkg/client/kubernetes"

// Result represents the various values and configurations
// that get created during the installation, like generated
// Helm values, passwords, IP addresses etc.
type Result struct {
	HelmValues        KubermaticValues
	NginxIngresses    []kubernetes.Ingress
	NodeportIngresses []kubernetes.Ingress
}
