package task

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/kubermatic/kubermatic-installer/pkg/client/helm"
	"github.com/kubermatic/kubermatic-installer/pkg/client/kubernetes"
	"github.com/kubermatic/kubermatic-installer/pkg/installer/state"
)

type Task interface {
	Run(installerState *state.InstallerState, clusterState *state.ClusterState, kubeClient kubernetes.Client, helmClient helm.Client, log logrus.FieldLogger, dryRun bool) error
}

type GenericUpgradeTask struct {
	ChartName   string
	ReleaseName string
	Namespace   string
}

func (t *GenericUpgradeTask) Run(installerState *state.InstallerState, clusterState *state.ClusterState, kubeClient kubernetes.Client, helmClient helm.Client, log logrus.FieldLogger, dryRun bool) error {
	chart := installerState.GetChart(t.ChartName)
	if chart == nil {
		return fmt.Errorf("chart %s not found in installer bundle", t.ChartName)
	}

	release := clusterState.Release(t.ReleaseName, t.Namespace)
	if release != nil {
		log.WithFields(logrus.Fields{
			"from": release.AppVersion,
			"to":   chart.AppVersion,
		}).Infof("Updating %s chart…", t.ChartName)
	} else {
		log.WithField("version", chart.AppVersion).Infof("Installing %s chart…", t.ChartName)
	}

	if !dryRun {
		if release == nil {
			if err := kubeClient.CreateNamespace(t.Namespace); err != nil {
				return fmt.Errorf("failed to create %q namespace: %v", t.Namespace, err)
			}
		}

		if err := helmClient.InstallChart(t.Namespace, t.ReleaseName, chart.Directory, "", nil, true); err != nil {
			return fmt.Errorf("failed to install: %v", err)
		}
	}

	// always create the side-effect on the state, even in non-dry-run modes
	clusterState.UpdateRelease(t.ReleaseName, t.Namespace, chart)

	return nil
}

// InfoTask is used to just display some status information, but not actually do anything.
type InfoTask struct {
	Message string
}

func (t *InfoTask) Run(_ *state.InstallerState, _ *state.ClusterState, _ kubernetes.Client, _ helm.Client, log logrus.FieldLogger, _ bool) error {
	log.Info(t.Message)
	return nil
}

type EnsureStorageClassTask struct {
	StorageClass *kubernetes.StorageClass
}

func (t *EnsureStorageClassTask) Run(_ *state.InstallerState, clusterState *state.ClusterState, kubeClient kubernetes.Client, _ helm.Client, log logrus.FieldLogger, dryRun bool) error {
	if clusterState.HasStorageClass(t.StorageClass.Metadata.Name) {
		return nil
	}

	log.WithFields(logrus.Fields{
		"provisioner": t.StorageClass.Provisioner,
		"parameters":  t.StorageClass.Parameters,
	}).Infof("Creating StorageClass %s…", t.StorageClass.Metadata.Name)

	if !dryRun {
		err := kubeClient.CreateStorageClass(*t.StorageClass)
		if err != nil {
			return fmt.Errorf("StorageClass could not be created: %v", err)
		}

		log.Infof("StorageClass has been created successfully.")
	}

	// always create the side-effect on the state, even in non-dry-run modes
	clusterState.StorageClasses = append(clusterState.StorageClasses, *t.StorageClass)

	return nil
}
