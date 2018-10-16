package installer

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kubermatic/kubermatic-installer/pkg/client/helm"
	"github.com/kubermatic/kubermatic-installer/pkg/client/kubernetes"
	helmvalues "github.com/kubermatic/kubermatic-installer/pkg/helm"
	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
	"github.com/sirupsen/logrus"
)

type installer struct {
	options  InstallerOptions
	manifest *manifest.Manifest
	logger   *logrus.Logger

	// runtime information
	kubeconfigFile string
	valuesFile     string
	helm           helm.Client
	kubernetes     kubernetes.Client
}

func NewInstaller(options InstallerOptions, manifest *manifest.Manifest, logger *logrus.Logger) *installer {
	return &installer{
		options:  options,
		manifest: manifest,
		logger:   logger,
	}
}

func (i *installer) kubeContext() string {
	return i.manifest.SeedClusters[0]
}

func (i *installer) kubeconfig() (string, error) {
	if i.kubeconfigFile == "" {
		var err error

		i.kubeconfigFile, err = i.dumpKubeconfig()
		if err != nil {
			return "", fmt.Errorf("failed to create kubeconfig: %v", err)
		}

		i.logger.Debugf("Dumped kubeconfig to %s.", i.kubeconfigFile)
	}

	return i.kubeconfigFile, nil
}

func (i *installer) Manifest() *manifest.Manifest {
	return i.manifest
}

func (i *installer) HelmClient() (helm.Client, error) {
	kubeconfig, err := i.kubeconfig()
	if err != nil {
		return nil, fmt.Errorf("failed to build Helm client: %v", err)
	}

	kubeContext := i.kubeContext()

	return helm.NewCLI(kubeconfig, kubeContext, HelmTillerNamespace, i.options.HelmTimeout, i.logger.WithField("backend", "helm"))
}

func (i *installer) KubernetesClient() (kubernetes.Client, error) {
	kubeconfig, err := i.kubeconfig()
	if err != nil {
		return nil, fmt.Errorf("failed to build Kubernetes client: %v", err)
	}

	kubeContext := i.kubeContext()

	return kubernetes.NewKubectl(kubeconfig, kubeContext, i.logger.WithField("backend", "kubectl"))
}

func (i *installer) dumpKubeconfig() (string, error) {
	return i.dumpTempFile("kubermatic.*.kubeconfig", i.manifest.Kubeconfig)
}

func (i *installer) probeCluster() error {
	class, err := i.kubernetes.DefaultStorageClass()
	if err != nil {
		return err
	}

	if class == nil {
		i.manifest.MinioStorageClass = KubermaticStorageClass
	}

	return nil
}

func (i *installer) prepareHelmValues() (helmvalues.Values, error) {
	// load Kubermatic's values.yaml
	values, err := helmvalues.LoadValuesFromFile("values.example.yaml")
	if err != nil {
		return helmvalues.Values{}, err
	}

	// apply manifest information to the values.yaml
	if err := values.ApplyManifest(i.manifest); err != nil {
		return values, fmt.Errorf("failed to create Helm values.yaml: %v", err)
	}

	// write values.yaml to file
	i.valuesFile, err = i.dumpHelmValues(values)
	if err != nil {
		return values, fmt.Errorf("failed to create values.yaml: %v", err)
	}

	i.logger.Debugf("Created Helm values.yaml at %s.", i.valuesFile)

	return values, nil
}

func (i *installer) dumpHelmValues(values helmvalues.Values) (string, error) {
	data := values.YAML()
	filename := i.options.ValuesFile

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

func (i *installer) cleanup() {
	if i.kubeconfigFile != "" && !i.options.KeepFiles {
		os.Remove(i.kubeconfigFile)
	}

	if i.valuesFile != "" && (!i.options.KeepFiles && i.options.ValuesFile == "") {
		os.Remove(i.valuesFile)
	}
}

func (i *installer) cleanupTempFile(filename string) {
	os.Remove(filename)
}
