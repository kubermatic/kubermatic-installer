package installer

import (
	"fmt"

	"github.com/kubermatic/kubermatic-installer/pkg/helm"
	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
	"github.com/sirupsen/logrus"
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

	values.ApplyManifest(i.manifest)

	// create a Helm client
	helm, err := helm.NewHelm(i.manifest.Kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to create Helm client: %v", err)
	}
	defer helm.Close()

	return nil
}
