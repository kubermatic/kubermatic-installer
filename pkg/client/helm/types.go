package helm

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver"
	"gopkg.in/yaml.v2"
)

type Release struct {
	Name      string          `json:"name"`
	Namespace string          `json:"namespace"`
	Chart     string          `json:"chart"`
	Revision  string          `json:"revision"`
	Version   *semver.Version `json:"-"`
	// AppVersion is not a semver, for example Minio has date-based versions.
	AppVersion string        `json:"app_version"`
	Status     releaseStatus `json:"status"`
}

func (r *Release) Clone() Release {
	copy := *r
	copy.Version = semver.MustParse(r.Version.Original())

	return copy
}

type Chart struct {
	Name       string          `yaml:"name"`
	Version    *semver.Version `yaml:"-"`
	VersionRaw string          `yaml:"version"`
	// AppVersion is not a semver, for example Minio has date-based versions.
	AppVersion string `yaml:"appVersion"`
	Directory  string
}

func (c *Chart) Clone() Chart {
	copy := *c
	copy.Version = semver.MustParse(c.Version.Original())

	return copy
}

func LoadChart(directory string) (*Chart, error) {
	f, err := os.Open(filepath.Join(directory, "Chart.yaml"))
	if err != nil {
		return nil, fmt.Errorf("failed to open Chart.yaml: %v", err)
	}
	defer f.Close()

	chart := &Chart{}
	if err := yaml.NewDecoder(f).Decode(chart); err != nil {
		return nil, fmt.Errorf("failed to read Chart.yaml: %v", err)
	}

	version, err := semver.NewVersion(chart.VersionRaw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version %q: %v", chart.VersionRaw, err)
	}

	chart.Version = version
	chart.Directory = directory

	return chart, nil
}
