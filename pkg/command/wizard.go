package command

import (
	"net"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/kubermatic/kubermatic-installer/pkg/server"
)

func WizardCommand(logger *logrus.Logger) cli.Command {
	return cli.Command{
		Name:   "wizard",
		Usage:  "Launches a HTTP server that provides a web UI for configuring the manifest",
		Action: WizardAction(logger),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "port",
				Value: 8080,
				Usage: "HTTP port to listen on",
			},
			cli.StringFlag{
				Name:   "host",
				Value:  "127.0.0.1",
				Usage:  "HTTP host to listen on",
				EnvVar: "KUBERMATIC_LISTEN_HOST",
			},
		},
	}
}

func WizardAction(logger *logrus.Logger) cli.ActionFunc {
	return handleErrors(logger, setupLogger(logger, func(ctx *cli.Context) error {
		port := ctx.Int("port")
		host := ctx.String("host")
		addr := net.JoinHostPort(host, strconv.Itoa(port))

		s := server.NewServer(logger)

		linkAddr := addr
		if host == "0.0.0.0" {
			linkAddr = net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
		}

		logger.Infof("Starting webserver at %s.", host)
		logger.Infof("Depending on your setup you can access the installer via http://%s/…", linkAddr)

		return s.Start(addr)
	}))
}
