package command

import (
	"fmt"

	"github.com/kubermatic/kubermatic-installer/pkg/shared"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func VersionCommand(logger *logrus.Logger) cli.Command {
	return cli.Command{
		Name:   "version",
		Usage:  "Prints the installer's version",
		Action: VersionAction(logger),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "short",
				Usage: "Omit git information",
			},
		},
	}
}

func VersionAction(logger *logrus.Logger) cli.ActionFunc {
	return handleErrors(logger, setupLogger(logger, func(ctx *cli.Context) error {
		if ctx.Bool("short") {
			fmt.Printf("Kubermatic Installer %s\n", shared.INSTALLER_VERSION)
		} else {
			fmt.Printf("Kubermatic Installer %s (git %s)\n", shared.INSTALLER_VERSION, shared.INSTALLER_GIT_HASH)
		}

		return nil
	}))
}
