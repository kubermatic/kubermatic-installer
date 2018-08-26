package main

import (
	"os"

	"github.com/kubermatic/kubermatic-installer/pkg/command"
	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "kubermatic-installer"
	app.Usage = "Helps configuring and setting up Kubermatic seed clusters."
	app.Version = "1.0.0"

	app.Commands = []cli.Command{
		cli.Command{
			Name:   "wizard",
			Usage:  "Launches a HTTP server that provides a web UI",
			Action: command.WizardCommand,
		},
	}

	app.Run(os.Args)
}
