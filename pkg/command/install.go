package command

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"

	"github.com/kubermatic/kubermatic-installer/pkg/installer"
	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
)

func InstallCommand(logger *logrus.Logger) cli.Command {
	return cli.Command{
		Name:      "install",
		Usage:     "Installs Kubernetes and Kubermatic using the pre-configured manifest",
		Action:    InstallAction(logger),
		ArgsUsage: "MANIFEST_FILE",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "keep-files",
				Usage: "do not delete generated kubeconfig and values.yaml in case of errors",
			},
			cli.IntFlag{
				Name:  "helm-timeout",
				Usage: "Number of seconds to wait for Helm operations to finish",
				Value: 300,
			},
		},
	}
}

func InstallAction(logger *logrus.Logger) cli.ActionFunc {
	return handleErrors(logger, setupLogger(logger, func(ctx *cli.Context) error {
		manifestFile := ctx.Args().First()
		if len(manifestFile) == 0 {
			return errors.New("no manifest file given")
		}

		manifest, err := loadManifest(manifestFile)
		if err != nil {
			return fmt.Errorf("failed to load manifest: %v", err)
		}

		options := installer.InstallerOptions{
			KeepFiles:   ctx.Bool("keep-files"),
			HelmTimeout: ctx.Int("helm-timeout"),
		}

		err = installer.NewInstaller(manifest, logger).Run(options)
		if err != nil {
			return err
		}

		logger.Info("Installation completed successfully!")

		fmt.Println("")
		fmt.Println("")
		fmt.Println("    Congratulations!")
		fmt.Println("")
		fmt.Println("    Kubermatic has been successfully installed to your Kubernetes")
		fmt.Println("    cluster. You can access the dashboard using the following URL")
		fmt.Println("    and start creating new clusters right now:")
		fmt.Println("")
		fmt.Printf("      %s", manifest.BaseURL())
		fmt.Println("")

		return nil
	}))
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
