package command

import (
	"net/http"
	"time"

	"github.com/kubermatic/kubermatic-installer/pkg/assets"
	cli "gopkg.in/urfave/cli.v1"
)

func WizardCommand(ctx *cli.Context) {
	s := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: http.FileServer(assets.Assets),

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  2 * time.Minute,
	}

	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
