package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/kubermatic/kubermatic-installer/pkg/command"
	"github.com/kubermatic/kubermatic-installer/pkg/log"
	"github.com/kubermatic/kubermatic-installer/pkg/shared"
)

func main() {
	logger := log.NewLogrus()

	app := cli.NewApp()
	app.Name = "kubermatic-installer"
	app.Usage = "Helps configuring and setting up Kubermatic."
	app.Version = shared.INSTALLER_VERSION
	app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "enable more verbose output",
		},
	}

	app.Commands = []cli.Command{
		command.VersionCommand(logger),
		command.DeployCommand(logger),
	}

	app.Run(os.Args)
}
