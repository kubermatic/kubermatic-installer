package state

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kubermatic/kubermatic-installer/pkg/client/helm"
)

type InstallerState struct {
	HelmCharts []helm.Chart
}

func (s *InstallerState) Clone() InstallerState {
	result := InstallerState{
		HelmCharts: []helm.Chart{},
	}

	for _, chart := range s.HelmCharts {
		result.HelmCharts = append(result.HelmCharts, chart.Clone())
	}

	return result
}

func NewInstallerState(chartDirectory string) (*InstallerState, error) {
	state := &InstallerState{
		HelmCharts: make([]helm.Chart, 0),
	}

	err := filepath.Walk(chartDirectory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			chartFile := filepath.Join(path, "Chart.yaml")

			if _, err := os.Stat(chartFile); err == nil {
				chart, err := helm.LoadChart(path)
				if err != nil {
					return fmt.Errorf("failed to read %s: %v", chartFile, err)
				}

				state.HelmCharts = append(state.HelmCharts, *chart)

				return filepath.SkipDir
			}
		}

		return nil
	})

	return state, err
}

func (s *InstallerState) GetChart(name string) *helm.Chart {
	for idx, c := range s.HelmCharts {
		if c.Name == name {
			return &s.HelmCharts[idx]
		}
	}

	return nil
}
