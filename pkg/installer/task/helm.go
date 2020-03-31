package task

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kubermatic/kubermatic-installer/pkg/yamled"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type GenericUpgradeTask struct {
	ChartName   string
	ReleaseName string
	Namespace   string
}

func (t *GenericUpgradeTask) Run(config *Config, state *State, clients *Clients, log logrus.FieldLogger, dryRun bool) error {
	chart := state.Installer.GetChart(t.ChartName)
	if chart == nil {
		return fmt.Errorf("chart %s not found in installer bundle", t.ChartName)
	}

	release := state.Cluster.Release(t.ReleaseName, t.Namespace)
	if release != nil {
		if release.Version.Equal(chart.Version) {
			log.WithField("version", chart.AppVersion).Debugf("Chart %s is up-to-date, nothing to do.", t.ChartName)
		} else {
			log.WithFields(logrus.Fields{
				"from": release.AppVersion,
				"to":   chart.AppVersion,
			}).Infof("Updating %s chart…", t.ChartName)
		}
	} else {
		log.WithField("version", chart.AppVersion).Infof("Installing %s chart…", t.ChartName)
	}

	if !dryRun {
		if release == nil {
			if err := clients.Kubernetes.CreateNamespace(t.Namespace); err != nil {
				return fmt.Errorf("failed to create %q namespace: %v", t.Namespace, err)
			}
		}

		helmValues, err := dumpHelmValues(config.Helm)
		if helmValues != "" {
			defer os.Remove(helmValues)
		}
		if err != nil {
			return err
		}

		if err := clients.Helm.InstallChart(t.Namespace, t.ReleaseName, chart.Directory, helmValues, nil, true); err != nil {
			return fmt.Errorf("failed to install: %v", err)
		}
	}

	// always create the side-effect on the state, even in non-dry-run modes
	state.Cluster.UpdateRelease(t.ReleaseName, t.Namespace, chart)

	return nil
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
