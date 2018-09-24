package installer

import (
	"github.com/sirupsen/logrus"

	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
)

type installer struct {
	manifest *manifest.Manifest
	logger   *logrus.Logger
}

func NewInstaller(manifest *manifest.Manifest, logger *logrus.Logger) *installer {
	return &installer{manifest, logger}
}

func (i *installer) Run() error {
	return nil
}
