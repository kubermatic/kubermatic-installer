package installer

import (
	"fmt"
	"strings"

	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
	"github.com/sirupsen/logrus"
)

type phase2 struct {
	installer
}

func NewPhase2(options InstallerOptions, manifest *manifest.Manifest, logger *logrus.Logger) Installer {
	installer := NewInstaller(options, manifest, logger)

	return &phase2{
		*installer,
	}
}

func (p *phase2) Run() (Result, error) {
	defer p.cleanup()

	result := Result{}
	var err error

	// create a Helm client
	p.helm, err = p.HelmClient()
	if err != nil {
		return result, err
	}

	// create a Kubernetes client
	p.kubernetes, err = p.KubernetesClient()
	if err != nil {
		return result, err
	}

	// probe cluster and complete manifest
	err = p.probeCluster()
	if err != nil {
		return result, err
	}

	err = p.install(&result)

	return result, err
}

func (p *phase2) SuccessMessage(m *manifest.Manifest, r Result) string {
	return strings.Join([]string{
		"    Congratulations!",
		"",
		"    The final component has been installed and the cluster should",
		"    acquire its certificates in the next minutes. Once it is done",
		"    you can access Kubermatic using the follwing URL:",
		"",
		fmt.Sprintf("      %s", m.BaseURL()),
	}, "\n")
}

func (p *phase2) install(result *Result) error {
	// is this cluster good enough for us?
	if err := p.checkPrerequisites(); err != nil {
		return fmt.Errorf("failed to check prerequisites: %v", err)
	}

	// load Kubermatic's values.yaml
	values, err := p.prepareHelmValues()
	if err != nil {
		return err
	}

	result.HelmValues = values

	// install charts
	if err := p.installCharts(); err != nil {
		return fmt.Errorf("failed to install charts: %v", err)
	}

	p.logger.Info("Installation completed successfully!")

	return nil
}

func (p *phase2) checkPrerequisites() error {
	exists, err := p.kubernetes.HasService(HelmTillerNamespace, HelmTillerService)
	if err != nil {
		return fmt.Errorf("could not check for service: %v", err)
	}

	if !exists {
		return fmt.Errorf("Tiller service could not be found in namespace %s", HelmTillerNamespace)
	}

	return nil
}

func (p *phase2) installCharts() error {
	if err := p.helm.InstallChart("default", "certs", "charts/certs", p.valuesFile, true); err != nil {
		return fmt.Errorf("could not install certs chart: %v", err)
	}

	return nil
}
