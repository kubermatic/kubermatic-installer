package command

import (
	"fmt"
	"sort"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/kubermatic/kubermatic-installer/pkg/client/helm"
	"github.com/kubermatic/kubermatic-installer/pkg/installer/state"
	"github.com/kubermatic/kubermatic-installer/pkg/shared"
)

func VersionCommand(logger *logrus.Logger) cli.Command {
	return cli.Command{
		Name:   "version",
		Usage:  "Prints the installer's version",
		Action: VersionAction(logger),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "short",
				Usage: "Omit git and chart information",
			},
		},
	}
}

func VersionAction(logger *logrus.Logger) cli.ActionFunc {
	return handleErrors(logger, setupLogger(logger, func(ctx *cli.Context) error {
		if ctx.Bool("short") {
			fmt.Printf("Kubermatic Installer %s\n", shared.INSTALLER_VERSION)
			return nil
		}

		installerState, err := state.NewInstallerState("charts")
		if err != nil {
			return fmt.Errorf("failed to determine installer chart state: %v", err)
		}

		fmt.Printf("Kubermatic Installer %s (git %s)\n", shared.INSTALLER_VERSION, shared.INSTALLER_GIT_HASH)

		var charts HelmCharts = installerState.HelmCharts
		sort.Sort(charts)

		for _, chart := range charts {
			fmt.Printf("%s %s (app version %s)\n", chart.Name, chart.Version, chart.AppVersion)
		}

		return nil
	}))
}

// HelmCharts is used to sort Helm charts by their name.
type HelmCharts []helm.Chart

func (a HelmCharts) Len() int           { return len(a) }
func (a HelmCharts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a HelmCharts) Less(i, j int) bool { return a[i].Name < a[j].Name }
