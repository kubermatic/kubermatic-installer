package installer

import (
	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
)

const (
	KubermaticNamespace    = "kubermatic"
	KubermaticStorageClass = "kubermatic-fast"
)

type InstallerOptions struct {
	KeepFiles   bool
	HelmTimeout int
	ValuesFile  string
}

type helmChart struct {
	namespace string
	name      string
	directory string
	flags     map[string]string
	wait      bool
}

type Installer interface {
	Run() (Result, error)
	SuccessMessage(*manifest.Manifest, Result) string
	Manifest() *manifest.Manifest
}
