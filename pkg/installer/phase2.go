package installer

import (
	"fmt"
	"strings"
	"time"

	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
	"github.com/kubermatic/kubermatic-installer/pkg/shared/dns"
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

	err = p.checkDNS()
	if err != nil {
		return fmt.Errorf("DNS check failed: %v", err)
	}

	return nil
}

func (p *phase2) installCharts() error {
	if err := p.helm.InstallChart("default", "certs", "charts/certs", p.valuesFile, nil, true); err != nil {
		return fmt.Errorf("could not install certs chart: %v", err)
	}

	if err := p.helm.InstallChart("iap", "iap", "charts/iap", p.valuesFile, nil, true); err != nil {
		return fmt.Errorf("could not install iap chart: %v", err)
	}

	// ensure that we do not check for CRD changes when installing Kubermatic
	kubermaticFlags := map[string]string{
		"kubermatic.checks.crd.disable": "true",
	}

	if err := p.helm.InstallChart(KubermaticNamespace, "kubermatic", "charts/kubermatic", p.valuesFile, kubermaticFlags, true); err != nil {
		return fmt.Errorf("could not install kubermatic chart: %v", err)
	}

	return nil
}

func (p *phase2) checkDNS() error {
	result := NewResult()

	if err := p.determineHostnames(&result); err != nil {
		return fmt.Errorf("could not determine hostnames: %v", err)
	}

	p.logger.Infof("Validating DNS settings…")

	validator := dns.NewValidator(15 * time.Minute)

	for _, record := range p.dnsRecords(result) {
		p.logger.Infof("Checking if %s points to %s…", record.Name, record.Target)

		if !validator.ValidateRecord(record) {
			return fmt.Errorf("could not resolve %s", record.Name)
		}
	}

	return nil
}
