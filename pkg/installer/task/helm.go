package task

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/kubermatic/kubermatic-installer/pkg/client/helm"
	"github.com/kubermatic/kubermatic-installer/pkg/installer/state"
	"github.com/kubermatic/kubermatic-installer/pkg/yamled"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type EnsureHelmReleaseTask struct {
	Chart       *helm.Chart
	Namespace   string
	ReleaseName string
}

func (t *EnsureHelmReleaseTask) String() string {
	return fmt.Sprintf("Install %s %s", t.Chart.Name, t.Chart.Version)
}

func (t *EnsureHelmReleaseTask) Run(ctx context.Context, opt *Options, installer *state.InstallerState, kubeClient ctrlruntimeclient.Client, helmClient helm.Client, log logrus.FieldLogger) error {
	// if tasks are run concurrently, this would make sense
	// log = log.WithField("chart", t.ChartName)

	// ensure namespace exists
	log.WithField("namespace", t.Namespace).Info("Ensuring namespace…")

	ns := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: t.Namespace,
		},
	}

	if err := kubeClient.Create(ctx, &ns); err != nil && !kerrors.IsAlreadyExists(err) {
		return fmt.Errorf("failed to create namespace: %v", err)
	}

	// find possible pre-existing release
	log.WithField("name", t.ReleaseName).Info("Checking for release…")

	release, err := helmClient.GetRelease(t.Namespace, t.ReleaseName)
	if err != nil {
		return fmt.Errorf("failed to check for an existing release: %v", err)
	}

	// release exists already, check if it's valid
	if release != nil {
		log.WithFields(logrus.Fields{
			"version": release.Version,
			"status":  release.Status,
		}).Info("Existing release found.")

		if release.Status.IsPending() {
			return temporaryErrorf(2*time.Second, "release is in %s status", release.Status)
		}

		// Sometimes installations can fail because prerequisites were not setup properly,
		// like a missing storage class. In this case, we want to allow the user to just
		// run the installer again and pick up where they left. Unfortunately Helm does not
		// support "upgrade --install" on failed installations: https://github.com/helm/helm/issues/3353
		// To work around this, we check the release status and purge it manually if it's failed.
		if statusRequiresPurge(release.Status) {
			log.Warn("Uninstalling defunct release before a clean installation is attempted…")

			if err := helmClient.UninstallRelease(t.Namespace, t.ReleaseName); err != nil {
				return fmt.Errorf("failed to uninstall release %s: %v", t.ReleaseName, err)
			}

			release = nil
			log.Info("Release has been uninstalled.")
		}

		// Now we have either a stable release or nothing at all.
	}

	if release == nil {
		log.Info("Installing release…")
	} else if release.Version.GreaterThan(t.Chart.Version) {
		log.Infof("Downgrading release to %s…", t.Chart.Version)
	} else if release.Version.LessThan(t.Chart.Version) {
		log.Infof("Updating release to %s…", t.Chart.Version)
	} else if opt.ForceHelmUpgrade {
		log.Info("Re-installing release because --force is set…")
	} else {
		log.Info("Release is up-to-date, nothing to do.")
		return nil
	}

	helmValues, err := dumpHelmValues(opt.Helm)
	if helmValues != "" {
		defer os.Remove(helmValues)
	}
	if err != nil {
		return err
	}

	if err := helmClient.InstallChart(t.Namespace, t.ReleaseName, t.Chart.Directory, helmValues, nil); err != nil {
		return fmt.Errorf("failed to install: %v", err)
	}

	return nil
}

func statusRequiresPurge(status helm.ReleaseStatus) bool {
	return status == helm.ReleaseStatusFailed
}

func dumpHelmValues(values *yamled.Document) (string, error) {
	f, err := ioutil.TempFile("", "helmvalues.*")
	if err != nil {
		return "", err
	}
	defer f.Close()

	err = yaml.NewEncoder(f).Encode(values)
	if err != nil {
		err = fmt.Errorf("failed to write Helm values to file: %v", err)
	}

	return f.Name(), err
}
