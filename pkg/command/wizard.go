package command

import (
	"net/http"
	"time"

	"github.com/kubermatic/kubermatic-installer/pkg/assets"
)

func WizardCommand() error {
	s := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: http.FileServer(assets.Assets),

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  2 * time.Minute,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
