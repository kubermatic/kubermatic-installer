package installer

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/kubermatic/kubermatic-installer/pkg/helm"
	"github.com/kubermatic/kubermatic-installer/pkg/kubectl"
	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
	"github.com/sirupsen/logrus"
)

const (
	HelmTillerNamespace      = "kube-system"
	HelmTillerServiceAccount = "tiller-sa"
	HelmTillerClusterRole    = "tiller-cluster-role"
)

type installer struct {
	manifest *manifest.Manifest
	logger   *logrus.Logger
}

func NewInstaller(manifest *manifest.Manifest, logger *logrus.Logger) *installer {
	return &installer{manifest, logger}
}

func (i *installer) Run() error {
	// create kubermatic's values.yaml
	values, err := LoadValuesFromFile("values.example.yaml")
	if err != nil {
		return err
	}

	err = values.ApplyManifest(i.manifest)
	if err != nil {
		return fmt.Errorf("failed to create Helm values.yaml: %v", err)
	}

	// prepare config files
	kubeconfigFile, err := i.dumpKubeconfig()
	if err != nil {
		return fmt.Errorf("failed to create kubeconfig: %v", err)
	}
	defer i.cleanupTempFile(kubeconfigFile)

	valuesFile, err := i.dumpHelmValues(values)
	if err != nil {
		return fmt.Errorf("failed to create values.yaml: %v", err)
	}
	defer i.cleanupTempFile(valuesFile)

	// create a Helm client
	helm, err := helm.NewHelm(kubeconfigFile, i.logger.WithField("backend", "helm"))
	if err != nil {
		return fmt.Errorf("failed to create Helm client: %v", err)
	}

	// create a kubectl client
	kubectl, err := kubectl.NewKubectl(kubeconfigFile, i.logger.WithField("backend", "kubectl"))
	if err != nil {
		return fmt.Errorf("failed to create kubectl client: %v", err)
	}

	return i.install(helm, kubectl, valuesFile)
}

func (i *installer) install(helm helm.Helm, kubectl kubectl.Kubectl, values string) error {
	if err := i.setupHelm(helm, kubectl); err != nil {
		return fmt.Errorf("failed to setup Helm: %v", err)
	}

	if err := i.installCharts(helm, kubectl, values); err != nil {
		return fmt.Errorf("failed to install charts: %v", err)
	}

	return nil
}

func (i *installer) setupHelm(helm helm.Helm, kubectl kubectl.Kubectl) error {
	if err := kubectl.CreateServiceAccount(HelmTillerNamespace, HelmTillerServiceAccount); err != nil {
		return fmt.Errorf("could not create tiller service account: %v", err)
	}

	if err := kubectl.CreateClusterRoleBinding(HelmTillerClusterRole, "cluster-admin", fmt.Sprintf("%s:%s", HelmTillerNamespace, HelmTillerServiceAccount)); err != nil {
		return fmt.Errorf("could not create tiller service account: %v", err)
	}

	// wait a bit for Kubernetes to settle
	time.Sleep(5 * time.Second)

	if err := helm.Init(HelmTillerServiceAccount, HelmTillerNamespace); err != nil {
		return fmt.Errorf("failed to init Helm: %v", err)
	}

	// wait for Helm to be ready
	time.Sleep(20 * time.Second)

	return nil
}

type helmChart struct {
	namespace string
	name      string
	directory string
}

func (i *installer) installCharts(helm helm.Helm, kubectl kubectl.Kubectl, values string) error {
	charts := []helmChart{
		{"nginx-ingress-controller", "nginx-ingress-controller", "charts/nginx-ingress-controller"},
		{"cert-manager", "cert-manager", "charts/cert-manager"},
		{"default", "certs", "charts/certs"},
		{"oauth", "oauth", "charts/oauth"},
		{"minio", "minio", "charts/minio"},
		{"kubermatic", "kubermatic", "charts/kubermatic"},
		{"nodeport-proxy", "nodeport-proxy", "charts/nodeport-proxy"},
	}

	if i.manifest.Monitoring.Enabled {
		charts = append(charts,
			helmChart{"monitoring", "prometheus", "charts/monitoring/prometheus"},
			helmChart{"monitoring", "node-exporter", "charts/monitoring/node-exporter"},
			helmChart{"monitoring", "kube-state-metrics", "charts/monitoring/kube-state-metrics"},
			helmChart{"monitoring", "grafana", "charts/monitoring/grafana"},
			helmChart{"monitoring", "alertmanager", "charts/monitoring/alertmanager"},
		)
	}

	charts = append(charts, helmChart{"iap", "iap", "charts/iap"})

	for _, chart := range charts {
		if err := helm.InstallChart(chart.namespace, chart.name, chart.directory, values); err != nil {
			return fmt.Errorf("could not install chart: %v", err)
		}
	}

	return nil
}

func (i *installer) dumpKubeconfig() (string, error) {
	return i.dumpTempFile("kubermatic.*.kubeconfig", i.manifest.Kubeconfig)
}

func (i *installer) dumpHelmValues(values KubermaticValues) (string, error) {
	return i.dumpTempFile("kubermatic.*.values.yaml", string(values.YAML()))
}

func (i *installer) dumpTempFile(fpattern string, contents string) (string, error) {
	tmpfile, err := ioutil.TempFile("", fpattern)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %v", err)
	}

	_, err = tmpfile.WriteString(contents)
	if err != nil {
		return "", fmt.Errorf("failed to write to file: %v", err)
	}

	err = tmpfile.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close file: %v", err)
	}

	return tmpfile.Name(), nil
}

func (i *installer) cleanupTempFile(filename string) {
	os.Remove(filename)
}
