package task

import (
	"github.com/sirupsen/logrus"

	"github.com/kubermatic/kubermatic-installer/pkg/client/helm"
	"github.com/kubermatic/kubermatic-installer/pkg/client/kubernetes"
	"github.com/kubermatic/kubermatic-installer/pkg/installer/state"
	"github.com/kubermatic/kubermatic-installer/pkg/shared/operatorv1alpha1"
	"github.com/kubermatic/kubermatic-installer/pkg/yamled"
)

type State struct {
	Installer *state.InstallerState
	Cluster   *state.ClusterState
}

type Config struct {
	Kubermatic *operatorv1alpha1.KubermaticConfiguration
	Helm       *yamled.Document
}

type Clients struct {
	Kubernetes kubernetes.Client
	Helm       helm.Client
}

type Task interface {
	Run(config *Config, state *State, clients *Clients, log logrus.FieldLogger, dryRun bool) error
}
