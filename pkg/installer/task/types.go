package task

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kubermatic/kubermatic-installer/pkg/client/helm"
	"github.com/kubermatic/kubermatic-installer/pkg/installer/state"
	"github.com/kubermatic/kubermatic-installer/pkg/shared/operatorv1alpha1"
	"github.com/kubermatic/kubermatic-installer/pkg/yamled"

	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type TemporaryError struct {
	error
	RetryAfter time.Duration
}

func temporaryErrorf(retryAfter time.Duration, format string, args ...interface{}) *TemporaryError {
	return &TemporaryError{
		error:      fmt.Errorf(format, args...),
		RetryAfter: retryAfter,
	}
}

type Options struct {
	Kubermatic       *operatorv1alpha1.KubermaticConfiguration
	Helm             *yamled.Document
	ForceHelmUpgrade bool
}

type Task interface {
	String() string
	Run(ctx context.Context, opt *Options, installer *state.InstallerState, kubeClient ctrlruntimeclient.Client, helmClient helm.Client, log logrus.FieldLogger) error
}
