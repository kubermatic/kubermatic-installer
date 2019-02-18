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
				Name:  "certificates",
				Usage: "Install Helm charts required for acquiring TLS certificates",
			},
			cli.BoolFlag{
				Name:  "keep-files",
				Usage: "Do not delete the temporary kubeconfig and values.yaml files",
			},
			cli.StringFlag{
				Name:   "values",
				Usage:  "Full path to where the Helm values.yaml should read from / be written to",
				EnvVar: "KUBERMATIC_VALUES_YAML",
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

		err = manifest.Validate()
		if err != nil {
			return fmt.Errorf("manifest is invalid: %v", err)
		}

		options := installer.InstallerOptions{
			KeepFiles:   ctx.Bool("keep-files"),
			HelmTimeout: ctx.Int("helm-timeout"),
			ValuesFile:  ctx.String("values"),
		}

		var phase installer.Installer

		if ctx.Bool("certificates") {
			phase = installer.NewPhase2(options, manifest, logger)
		} else {
			phase = installer.NewPhase1(options, manifest, logger)
		}

		result, err := phase.Run()
		if err != nil {
			return err
		}

		msg := phase.SuccessMessage(manifest, result)
		if len(msg) > 0 {
			fmt.Println("")
			fmt.Println("")
			fmt.Println(msg)
			fmt.Println("")
		}

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
		return nil, fmt.Errorf("failed to decode file as YAML: %v", err)
	}

	return &manifest, nil
}
