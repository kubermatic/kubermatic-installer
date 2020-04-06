package task

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/kubermatic/kubermatic-installer/pkg/client/helm"
	"github.com/kubermatic/kubermatic-installer/pkg/installer/state"
	"github.com/kubermatic/kubermatic-installer/pkg/shared/operatorv1alpha1"
	"github.com/kubermatic/kubermatic-installer/pkg/yamled"

	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
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
	Kubernetes ctrlruntimeclient.Client
	Helm       helm.Client
}

type Options struct {
	DryRun           bool
	ForceHelmUpgrade bool
}

type Task interface {
	Required(config *Config, state *State, opt *Options) (bool, error)
	Plan(config *Config, state *State, opt *Options, log logrus.FieldLogger) error
	Run(ctx context.Context, config *Config, state *State, clients *Clients, opt *Options, log logrus.FieldLogger) error
}
