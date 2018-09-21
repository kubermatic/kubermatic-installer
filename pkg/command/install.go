package command

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/kubermatic/kubermatic-installer/pkg/installer"
	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

func InstallCommand(logger *logrus.Logger) cli.Command {
	return cli.Command{
		Name:      "install",
		Usage:     "Installs Kubernetes and Kubermatic using the pre-configured manifest",
		Action:    InstallAction(logger),
		ArgsUsage: "MANIFEST_FILE",
	}
}

func InstallAction(logger *logrus.Logger) cli.ActionFunc {
	return handleErrors(logger, func(ctx *cli.Context) error {
		manifestFile := ctx.Args().First()
		if len(manifestFile) == 0 {
			return errors.New("no manifest file given")
		}

		manifest, err := loadManifest(manifestFile)
		if err != nil {
			return fmt.Errorf("failed to load manifest: %v", err)
		}

		installer := installer.NewInstaller(manifest, logger)

		return installer.Run()
	})
}

func loadManifest(filename string) (*manifest.Manifest, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	manifest := manifest.Manifest{}
	if err := yaml.Unmarshal(content, &manifest); err != nil {
		return nil, fmt.Errorf("failed to decode file as JSON: %v", err)
	}

	return &manifest, nil
}
