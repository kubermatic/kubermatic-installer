package installer

import (
	"github.com/kubermatic/kubermatic-installer/pkg/client/kubernetes"
	"github.com/kubermatic/kubermatic-installer/pkg/helm"
)

// Result represents the various values and configurations
// that get created during the installation, like generated
// Helm values, passwords, IP addresses etc.
type Result struct {
	HelmValues        *helm.Values
	NginxIngresses    []kubernetes.Ingress
	NodeportIngresses []kubernetes.Ingress
}

func NewResult() Result {
	return Result{
		HelmValues: helm.NewValues(),
	}
}
