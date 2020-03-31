package task

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kubermatic/kubermatic-installer/pkg/yamled"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type EnsureHelmReleaseTask struct {
	ChartName   string
	ReleaseName string
	Namespace   string
}

func (t *EnsureHelmReleaseTask) Required(_ *Config, state *State, opt *Options) (bool, error) {
	if opt.ForceHelmUpgrade {
		return true, nil
	}

	chart := state.Installer.GetChart(t.ChartName)
	if chart == nil {
		return false, fmt.Errorf("chart %s not found in installer bundle", t.ChartName)
	}

	release := state.Cluster.Release(t.ReleaseName, t.Namespace)

	return release == nil || !release.Version.Equal(chart.Version), nil
}

func (t *EnsureHelmReleaseTask) Plan(_ *Config, state *State, opt *Options, log logrus.FieldLogger) error {
	chart := state.Installer.GetChart(t.ChartName)
	if chart == nil {
		return fmt.Errorf("chart %s not found in installer bundle", t.ChartName)
	}

	log = log.WithField("namespace", t.Namespace)
	release := state.Cluster.Release(t.ReleaseName, t.Namespace)

	if release == nil {
		log.WithField("version", chart.Version).Infof("Install %s chart.", t.ChartName)
		return nil
	}

	if release.Version.Equal(chart.Version) && !opt.ForceHelmUpgrade {
		return nil
	}

	log.WithFields(logrus.Fields{
		"from": release.Version,
		"to":   chart.Version,
	}).Infof("Update %s chart.", t.ChartName)

	return nil
}

func (t *EnsureHelmReleaseTask) Run(config *Config, state *State, clients *Clients, opt *Options, log logrus.FieldLogger) error {
	chart := state.Installer.GetChart(t.ChartName)
	if chart == nil {
		return fmt.Errorf("chart %s not found in installer bundle", t.ChartName)
	}

	release := state.Cluster.Release(t.ReleaseName, t.Namespace)
	if release != nil {
		log = log.WithField("installed", release.Version)
	}

	log.WithField("version", chart.Version).Infof("Ensuring %s chart is installed…", t.ChartName)

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
