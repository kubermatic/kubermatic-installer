package installer

import (
	"fmt"
	"strings"

	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
	"github.com/kubermatic/kubermatic-installer/pkg/shared/dns"
	"github.com/sirupsen/logrus"
)

type phase1 struct {
	installer
}

func NewPhase1(options InstallerOptions, manifest *manifest.Manifest, logger *logrus.Logger) Installer {
	installer := NewInstaller(options, manifest, logger)

	return &phase1{
		*installer,
	}
}

func (p *phase1) Run() (Result, error) {
	defer p.cleanup()

	result := NewResult()
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

func (p *phase1) SuccessMessage(m *manifest.Manifest, r Result) string {
	msg := []string{
		"    Congratulations!",
		"",
		"    Kubermatic has been successfully installed to your Kubernetes",
		"    cluster. Please setup your DNS records to allow Kubermatic to",
		"    acquire its TLS certificates and enable inter-cluster",
		"    communication.",
		"",
	}

	msg = append(msg, p.formatRecords(p.dnsRecords(r))...)
	msg = append(msg,
		"",
		"    Once the DNS changes have propagated, please perform the",
		"    final step of the installation by adding the --certificates",
		"    flag to the installer and running it again.",
	)

	return strings.Join(msg, "\n")
}

func (p *phase1) install(result *Result) error {
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

	// install Helm into cluster
	if err := p.installHelm(); err != nil {
		return fmt.Errorf("failed to setup Helm: %v", err)
	}

	// install CRDs
	if err := p.installCRDs(); err != nil {
		return fmt.Errorf("failed to install CRDs: %v", err)
	}

	// install charts
	if err := p.installCharts(); err != nil {
		return fmt.Errorf("failed to install charts: %v", err)
	}

	// determine hostnames/IPs
	if err := p.determineHostnames(result); err != nil {
		return fmt.Errorf("failed to determine hostnames: %v", err)
	}

	p.logger.Info("Installation completed successfully!")

	return nil
}

func (p *phase1) checkPrerequisites() error {
	exists, err := p.kubernetes.HasStorageClass(KubermaticStorageClass)
	if err != nil {
		return fmt.Errorf("could not check for storage class: %v", err)
	}

	if !exists {
		sc := StorageClassForProvider(p.manifest.CloudProvider)
		if sc == nil {
			p.logger.Warnf("Storage class could not be found, please create it manually.")
		} else {
			err := p.kubernetes.CreateStorageClass(*sc)
			if err != nil {
				return fmt.Errorf("storage class could not be found or created: %v", err)
			}

			p.logger.Infof("Automatically created storage class.")
		}
	}

	return nil
}

func (p *phase1) installCRDs() error {
	exists, err := p.kubernetes.HasCustomResourceDefinition("addons.kubermatic.k8s.io")
	if err != nil {
		return fmt.Errorf("could not check for CRDs: %v", err)
	}

	if !exists {
		err = p.kubernetes.ApplyManifests("charts/kubermatic/crd")
		if err != nil {
			return fmt.Errorf("could not create CRDs: %v", err)
		}
	}

	return nil
}

func (p *phase1) installHelm() error {
	if err := p.kubernetes.CreateNamespace(KubermaticNamespace); err != nil {
		return fmt.Errorf("could not create namespace: %v", err)
	}

	if HelmTillerNamespace != KubermaticNamespace {
		if err := p.kubernetes.CreateNamespace(HelmTillerNamespace); err != nil {
			return fmt.Errorf("could not create namespace: %v", err)
		}
	}

	if err := p.kubernetes.CreateServiceAccount(HelmTillerNamespace, HelmTillerServiceAccount); err != nil {
		return fmt.Errorf("could not create tiller service account: %v", err)
	}

	if err := p.kubernetes.CreateClusterRoleBinding(HelmTillerClusterRole, "cluster-admin", fmt.Sprintf("%s:%s", HelmTillerNamespace, HelmTillerServiceAccount)); err != nil {
		return fmt.Errorf("could not create tiller service account: %v", err)
	}

	if err := p.helm.Init(HelmTillerServiceAccount); err != nil {
		return fmt.Errorf("failed to init Helm: %v", err)
	}

	return nil
}

func (p *phase1) installCharts() error {
	charts := []helmChart{
		{"nginx-ingress-controller", "nginx-ingress-controller", "charts/nginx-ingress-controller", nil, true},
		{"cert-manager", "cert-manager", "charts/cert-manager", nil, true},
		{"oauth", "oauth", "charts/oauth", nil, true},
		{"minio", "minio", "charts/minio", nil, true},
		{"nodeport-proxy", "nodeport-proxy", "charts/nodeport-proxy", nil, true},
	}

	if p.manifest.Monitoring.Enabled {
		charts = append(charts,
			helmChart{"monitoring", "prometheus", "charts/monitoring/prometheus", nil, true},
			helmChart{"monitoring", "node-exporter", "charts/monitoring/node-exporter", nil, true},
			helmChart{"monitoring", "kube-state-metrics", "charts/monitoring/kube-state-metrics", nil, true},
			helmChart{"monitoring", "grafana", "charts/monitoring/grafana", nil, true},
			helmChart{"monitoring", "alertmanager", "charts/monitoring/alertmanager", nil, true},
		)
	}

	for _, chart := range charts {
		if err := p.helm.InstallChart(chart.namespace, chart.name, chart.directory, p.valuesFile, chart.flags, chart.wait); err != nil {
			return fmt.Errorf("could not install chart: %v", err)
		}
	}

	return nil
}

func (p *phase1) formatRecords(records []dns.Record) []string {
	width := 0
	for _, record := range records {
		if l := len(record.Name); l > width {
			width = l
		}
	}

	format := fmt.Sprintf("    %%%ds ➜ %%s", width)
	lines := make([]string, 0)

	for _, record := range records {
		var target string

		if record.Kind == dns.RecordKindCNAME {
			target = fmt.Sprintf("CNAME @ %s.", record.Target)
		} else {
			target = fmt.Sprintf("A record @ %s", record.Target)
		}

		lines = append(lines, fmt.Sprintf(format, record.Name, target))
	}

	return lines
}
