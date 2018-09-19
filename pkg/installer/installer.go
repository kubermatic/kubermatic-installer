package installer

import (
	"github.com/kubermatic/kubermatic-installer/pkg/shared"
	"github.com/sirupsen/logrus"
)

type installer struct {
	manifest *shared.Manifest
	logger   *logrus.Logger
}

func NewInstaller(manifest *shared.Manifest, logger *logrus.Logger) *installer {
	return &installer{manifest, logger}
}

func (i *installer) Run() error {
	return nil
}
