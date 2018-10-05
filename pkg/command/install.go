package command

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"

	"github.com/kubermatic/kubermatic-installer/pkg/client/kubernetes"
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
				Usage: "Do not delete the temporary kubeconfig and values.yaml files",
			},
			cli.StringFlag{
				Name:   "values",
				Usage:  "Full path to where the Helm values.yaml should be written to (will never be deleted, regardless of --keep-files)",
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

		result, err := installer.NewInstaller(manifest, logger).Run(options)
		if err != nil {
			return err
		}

		fmt.Println("")
		fmt.Println("")
		fmt.Println("    Congratulations!")
		fmt.Println("")
		fmt.Println("    Kubermatic has been successfully installed to your Kubernetes")
		fmt.Println("    cluster. Please setup your DNS records to allow Kubermatic to")
		fmt.Println("    acquire its TLS certificates and enable inter-cluster")
		fmt.Println("    communication.")
		fmt.Println("")

		domain := manifest.Settings.BaseDomain

		if len(result.NginxIngresses) > 0 {
			target := dnsRecord(result.NginxIngresses[0])
			seedLength := len(manifest.SeedClusters[0])
			padding := strings.Repeat(" ", seedLength+1) // we need length+3, but in the format string we already put two spaces

			fmt.Printf("     %s  %s ➜ %s\n", padding, domain, target)
			fmt.Printf("     %s*.%s ➜ %s\n", padding, domain, target)
		}

		if len(result.NodeportIngresses) > 0 {
			target := dnsRecord(result.NodeportIngresses[0])

			fmt.Printf("     *.%s.%s ➜ %s\n", manifest.SeedClusters[0], domain, target)
		}

		fmt.Println("")
		fmt.Println("    Once the DNS changes have propagated, you can access the")
		fmt.Println("    Kubermatic dashboard using the following link:")
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
		return nil, fmt.Errorf("failed to decode file as YAML: %v", err)
	}

	return &manifest, nil
}

func dnsRecord(ingress kubernetes.Ingress) string {
	if ingress.Hostname != "" {
		return fmt.Sprintf("CNAME @ %s.", ingress.Hostname)
	} else {
		return fmt.Sprintf("A record @ %s", ingress.IP)
	}
}
