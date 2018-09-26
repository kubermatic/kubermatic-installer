package installer

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/kubermatic/kubermatic-installer/pkg/client/helm"
	"github.com/kubermatic/kubermatic-installer/pkg/client/kubernetes"
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

func (i *installer) Run(keepFiles bool) error {
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
	if !keepFiles {
		defer i.cleanupTempFile(kubeconfigFile)
	}

	i.logger.Debugf("Dumped kubeconfig to %s.", kubeconfigFile)

	valuesFile, err := i.dumpHelmValues(values)
	if err != nil {
		return fmt.Errorf("failed to create values.yaml: %v", err)
	}
	if !keepFiles {
		defer i.cleanupTempFile(valuesFile)
	}

	i.logger.Debugf("Created Helm values.yaml at %s.", valuesFile)

	// create a Helm client
	kubeContext := i.manifest.SeedClusters[0]
	helm, err := helm.NewCLI(kubeconfigFile, kubeContext, HelmTillerNamespace, i.logger.WithField("backend", "helm"))
	if err != nil {
		return fmt.Errorf("failed to create Helm client: %v", err)
	}

	// create a kubectl client
	kubectl, err := kubernetes.NewKubectl(kubeconfigFile, kubeContext, i.logger.WithField("backend", "kubectl"))
	if err != nil {
		return fmt.Errorf("failed to create kubectl client: %v", err)
	}

	return i.install(helm, kubectl, valuesFile)
}

func (i *installer) install(helm helm.Client, kubectl kubernetes.Client, values string) error {
	if err := i.setupHelm(helm, kubectl); err != nil {
		return fmt.Errorf("failed to setup Helm: %v", err)
	}

	if err := i.installCharts(helm, kubectl, values); err != nil {
		return fmt.Errorf("failed to install charts: %v", err)
	}

	return nil
}

func (i *installer) setupHelm(helm helm.Client, kubectl kubernetes.Client) error {
	if err := kubectl.CreateServiceAccount(HelmTillerNamespace, HelmTillerServiceAccount); err != nil {
		return fmt.Errorf("could not create tiller service account: %v", err)
	}

	if err := kubectl.CreateClusterRoleBinding(HelmTillerClusterRole, "cluster-admin", fmt.Sprintf("%s:%s", HelmTillerNamespace, HelmTillerServiceAccount)); err != nil {
		return fmt.Errorf("could not create tiller service account: %v", err)
	}

	if err := helm.Init(HelmTillerServiceAccount); err != nil {
		return fmt.Errorf("failed to init Helm: %v", err)
	}

	return nil
}

type helmChart struct {
	namespace string
	name      string
	directory string
}

func (i *installer) installCharts(helm helm.Client, kubectl kubernetes.Client, values string) error {
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
