package installer

import (
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
	return nil
}
