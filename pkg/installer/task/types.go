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

type Options struct {
	DryRun           bool
	ForceHelmUpgrade bool
}

type Task interface {
	Required(config *Config, state *State, opt *Options) (bool, error)
	Plan(config *Config, state *State, opt *Options, log logrus.FieldLogger) error
	Run(config *Config, state *State, clients *Clients, opt *Options, log logrus.FieldLogger) error
}
