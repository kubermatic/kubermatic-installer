package main

import (
	"os"

	"github.com/kubermatic/kubermatic-installer/pkg/command"
	"github.com/kubermatic/kubermatic-installer/pkg/shared"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	logger := setupLogging()

	app := cli.NewApp()
	app.Name = "kubermatic-installer"
	app.Usage = "Helps configuring and setting up Kubermatic."
	app.Version = shared.INSTALLER_VERSION

	app.Commands = []cli.Command{
		command.WizardCommand(logger),
		command.InstallCommand(logger),
	}

	app.Run(os.Args)
}

func setupLogging() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		FullTimestamp: true,
	}

	return logger
}
