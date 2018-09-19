package command

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/kubermatic/kubermatic-installer/pkg/assets"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func WizardCommand(logger *logrus.Logger) cli.Command {
	return cli.Command{
		Name:   "wizard",
		Usage:  "Launches a HTTP server that provides a web UI for configuring the manifest",
		Action: WizardAction(logger),
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "port, p",
				Value: 8080,
				Usage: "HTTP port to listen on",
			},
			cli.StringFlag{
				Name:  "addr, a",
				Value: "127.0.0.1",
				Usage: "HTTP host to listen on",
			},
		},
	}
}

func WizardAction(logger *logrus.Logger) cli.ActionFunc {
	return handleErrors(logger, func(ctx *cli.Context) error {
		port := ctx.Int("port")
		addr := ctx.String("addr")
		host := net.JoinHostPort(addr, strconv.Itoa(port))

		s := http.Server{
			Addr:    host,
			Handler: http.FileServer(assets.Assets),

			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  2 * time.Minute,
		}

		logger.Infof("Starting webserver at http://%s/…", host)

		return s.ListenAndServe()
	})
}
