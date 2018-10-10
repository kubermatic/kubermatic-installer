package installer

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kubermatic/kubermatic-installer/pkg/client/helm"
	"github.com/kubermatic/kubermatic-installer/pkg/client/kubernetes"
	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
	"github.com/sirupsen/logrus"
)

const (
	KubermaticNamespace      = "kubermatic"
	KubermaticStorageClass   = "kubermatic-fast"
	HelmTillerNamespace      = KubermaticNamespace
	HelmTillerServiceAccount = "tiller-sa"
	HelmTillerClusterRole    = "tiller-cluster-role"
)

type InstallerOptions struct {
	KeepFiles   bool
	HelmTimeout int
	ValuesFile  string
}

type installer struct {
	manifest *manifest.Manifest
	logger   *logrus.Logger
}

func NewInstaller(manifest *manifest.Manifest, logger *logrus.Logger) *installer {
	return &installer{manifest, logger}
}

func (i *installer) Run(opts InstallerOptions) (Result, error) {
	result := Result{}

	// create kubermatic's values.yaml
	values, err := LoadValuesFromFile("values.example.yaml")
	result.HelmValues = values
	if err != nil {
		return result, err
	}

	err = values.ApplyManifest(i.manifest)
	if err != nil {
		return result, fmt.Errorf("failed to create Helm values.yaml: %v", err)
	}

	// prepare config files
	kubeconfigFile, err := i.dumpKubeconfig()
	if err != nil {
		return result, fmt.Errorf("failed to create kubeconfig: %v", err)
	}
	if !opts.KeepFiles {
		defer i.cleanupTempFile(kubeconfigFile)
	}

	i.logger.Debugf("Dumped kubeconfig to %s.", kubeconfigFile)

	valuesFile, err := i.dumpHelmValues(values, opts.ValuesFile)
	if err != nil {
		return result, fmt.Errorf("failed to create values.yaml: %v", err)
	}
	if !opts.KeepFiles && opts.ValuesFile == "" {
		defer i.cleanupTempFile(valuesFile)
	}

	i.logger.Debugf("Created Helm values.yaml at %s.", valuesFile)

	// create a Helm client
	kubeContext := i.manifest.SeedClusters[0]
	helm, err := helm.NewCLI(kubeconfigFile, kubeContext, HelmTillerNamespace, opts.HelmTimeout, i.logger.WithField("backend", "helm"))
	if err != nil {
		return result, fmt.Errorf("failed to create Helm client: %v", err)
	}

	// create a kubectl client
	kubectl, err := kubernetes.NewKubectl(kubeconfigFile, kubeContext, i.logger.WithField("backend", "kubectl"))
	if err != nil {
		return result, fmt.Errorf("failed to create kubectl client: %v", err)
	}

	return result, i.install(helm, kubectl, &result, valuesFile)
}

func (i *installer) install(helm helm.Client, kubectl kubernetes.Client, result *Result, values string) error {
	if err := i.setupHelm(helm, kubectl, result); err != nil {
		return fmt.Errorf("failed to setup Helm: %v", err)
	}

	if err := i.checkPrerequisites(helm, kubectl); err != nil {
		return fmt.Errorf("failed to check prerequisites: %v", err)
	}

	if err := i.installCharts(helm, kubectl, result, values); err != nil {
		return fmt.Errorf("failed to install charts: %v", err)
	}

	if err := i.determineHostnames(helm, kubectl, result); err != nil {
		return fmt.Errorf("failed to determine hostnames: %v", err)
	}

	i.logger.Info("Installation completed successfully!")

	return nil
}

func (i *installer) setupHelm(helm helm.Client, kubectl kubernetes.Client, result *Result) error {
	if err := kubectl.CreateNamespace(KubermaticNamespace); err != nil {
		return fmt.Errorf("could not create namespace: %v", err)
	}

	if HelmTillerNamespace != KubermaticNamespace {
		if err := kubectl.CreateNamespace(HelmTillerNamespace); err != nil {
			return fmt.Errorf("could not create namespace: %v", err)
		}
	}

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

func (i *installer) checkPrerequisites(helm helm.Client, kubectl kubernetes.Client) error {
	exists, err := kubectl.HasStorageClass(KubermaticStorageClass)
	if err != nil {
		return fmt.Errorf("could not check for storage class: %v", err)
	}

	if !exists {
		sc := StorageClassForProvider(i.manifest.CloudProvider)
		if sc == nil {
			i.logger.Warnf("Storage class could not be found, please create it manually.", KubermaticStorageClass)
		} else {
			err := kubectl.CreateStorageClass(*sc)
			if err != nil {
				i.logger.Errorf("Storage class could not be found nor created: %v", KubermaticStorageClass, err)
			} else {
				i.logger.Infof("Automatically created storage class.", KubermaticStorageClass)
			}
		}
	}

	return nil
}

type helmChart struct {
	namespace string
	name      string
	directory string
	wait      bool
}

func (i *installer) installCharts(helm helm.Client, kubectl kubernetes.Client, result *Result, values string) error {
	charts := []helmChart{
		{"nginx-ingress-controller", "nginx-ingress-controller", "charts/nginx-ingress-controller", true},
		{"cert-manager", "cert-manager", "charts/cert-manager", true},
		{"default", "certs", "charts/certs", true},
		{"oauth", "oauth", "charts/oauth", true},
		{"minio", "minio", "charts/minio", true},
		{"kubermatic", KubermaticNamespace, "charts/kubermatic", true},
		{"nodeport-proxy", "nodeport-proxy", "charts/nodeport-proxy", true},

		// Do not wait for IAP to come up, because it depends on proper DNS names to be configured
		// and certificates to be acquired; this is something the user has to do *after* we tell
		// them the target IPs/hostnames for their DNS settings.
		{"iap", "iap", "charts/iap", false},
	}

	if i.manifest.Monitoring.Enabled {
		charts = append(charts,
			helmChart{"monitoring", "prometheus", "charts/monitoring/prometheus", true},
			helmChart{"monitoring", "node-exporter", "charts/monitoring/node-exporter", true},
			helmChart{"monitoring", "kube-state-metrics", "charts/monitoring/kube-state-metrics", true},
			helmChart{"monitoring", "grafana", "charts/monitoring/grafana", true},
			helmChart{"monitoring", "alertmanager", "charts/monitoring/alertmanager", true},
		)
	}

	for _, chart := range charts {
		if err := helm.InstallChart(chart.namespace, chart.name, chart.directory, values, chart.wait); err != nil {
			return fmt.Errorf("could not install chart: %v", err)
		}
	}

	return nil
}

func (i *installer) determineHostnames(helm helm.Client, kubectl kubernetes.Client, result *Result) error {
	ingresses, err := kubectl.ServiceIngresses("nginx-ingress-controller", "nginx-ingress-controller")
	if err != nil {
		return err
	}

	result.NginxIngresses = ingresses

	ingresses, err = kubectl.ServiceIngresses("nodeport-proxy", "nodeport-lb")
	if err != nil {
		return err
	}

	result.NodeportIngresses = ingresses

	return nil
}

func (i *installer) dumpKubeconfig() (string, error) {
	return i.dumpTempFile("kubermatic.*.kubeconfig", i.manifest.Kubeconfig)
}

func (i *installer) dumpHelmValues(values KubermaticValues, filename string) (string, error) {
	data := values.YAML()

	if len(filename) > 0 {
		return filename, ioutil.WriteFile(filename, data, 0644)
	}

	return i.dumpTempFile("kubermatic.*.values.yaml", string(data))
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
