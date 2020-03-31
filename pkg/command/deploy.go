package command

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"

	"github.com/kubermatic/kubermatic-installer/pkg/client/helm"
	"github.com/kubermatic/kubermatic-installer/pkg/client/kubernetes"
	"github.com/kubermatic/kubermatic-installer/pkg/installer"
	"github.com/kubermatic/kubermatic-installer/pkg/installer/stack/kubermatic"
	"github.com/kubermatic/kubermatic-installer/pkg/installer/state"
	"github.com/kubermatic/kubermatic-installer/pkg/log"
	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
	"github.com/kubermatic/kubermatic-installer/pkg/shared"
)

func DeployCommand(logger *logrus.Logger) cli.Command {
	return cli.Command{
		Name:      "deploy",
		Usage:     "Installs or upgrades the current installation to the installer's built-in version",
		Action:    DeployAction(logger),
		ArgsUsage: "MANIFEST_FILE",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "confirm",
				Usage: "Perform the deployment instead of just listing the required steps",
			},
			cli.BoolFlag{
				Name:  "certificates",
				Usage: "Deploy Helm charts required for acquiring TLS certificates",
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
			cli.StringFlag{
				Name:   "kubeconfig",
				Usage:  "Full path to where a kubeconfig with cluster-admin permissions for the target cluster",
				EnvVar: "KUBECONFIG",
			},
			cli.StringFlag{
				Name:   "kube-context",
				Usage:  "Context to use from the given kubeconfig",
				EnvVar: "KUBE_CONTEXT",
			},
			cli.DurationFlag{
				Name:  "helm-timeout",
				Usage: "Time to wait for Helm operations to finish",
				Value: 5 * time.Minute,
			},
		},
	}
}

/*
1. load chart information from disk (names, versions, app versions)
2. load release info from kubernetes cluster
3. prepare deployers for each helm chart
4. let each deployer validate all preconditions
5. run each deployer
6. success!
*/

func DeployAction(logger *logrus.Logger) cli.ActionFunc {
	return handleErrors(logger, setupLogger(logger, func(ctx *cli.Context) error {
		initLogger := logger.WithField("version", shared.INSTALLER_VERSION)
		if ctx.GlobalBool("verbose") {
			initLogger = initLogger.WithField("git", shared.INSTALLER_GIT_HASH)
		}
		initLogger.Info("Initializing installer…")

		// load bundled chart information from disk
		installerState, err := state.NewInstallerState("charts")
		if err != nil {
			return fmt.Errorf("failed to determine installer chart state: %v", err)
		}

		// prepapre Kubernetes and Helm clients
		kubeconfig := ctx.String("kubeconfig")
		if len(kubeconfig) == 0 {
			return errors.New("no kubeconfig (--kubeconfig or $KUBECONFIG) given")
		}

		kubeContext := ctx.String("kube-context")
		helmTimeout := ctx.Duration("helm-timeout")

		kubeClient, err := kubernetes.NewKubectl(kubeconfig, kubeContext, logger.WithField("backend", "kubectl"))
		if err != nil {
			return fmt.Errorf("failed to create Kubernetes client: %v", err)
		}

		helmClient, err := helm.NewCLI(kubeconfig, kubeContext, helmTimeout, logger.WithField("backend", "helmcli"))
		if err != nil {
			return fmt.Errorf("failed to create Helm client: %v", err)
		}

		// inspect remote cluster and list helm releases
		clusterState, err := state.NewClusterState(kubeClient, helmClient)
		if err != nil {
			return fmt.Errorf("failed to determine cluster state: %v", err)
		}

		for _, r := range clusterState.HelmReleases {
			fmt.Printf("release: %+v\n", r)
		}

		// determine tasks required to upgrade the cluster to the installer state
		logger.Info("Planning kubermatic stack deployment…")
		tasks, err := kubermatic.DeploymentTasks(installerState, clusterState)
		if err != nil {
			return fmt.Errorf("failed to determine deployment steps: %v", err)
		}

		// inform the user about steps we're going to take
		if !ctx.Bool("confirm") {
			fakeClusterState := clusterState.Clone()
			fakeInstallerState := installerState.Clone()
			fakeLogger := log.NewPlan()

			for _, t := range tasks {
				err := t.Run(&fakeInstallerState, &fakeClusterState, kubeClient, helmClient, fakeLogger, true)
				if err != nil {
					return fmt.Errorf("failed: %v", err)
				}
			}

			logger.Warn("Run the installer with --confirm to actually perform the steps outlined above.")
			return nil
		}

		logger.Info("Deploying Kubermatic now!")
		os.Exit(0)

		// run the tasks!
		for _, t := range tasks {
			err := t.Run(installerState, clusterState, kubeClient, helmClient, logger, false)
			if err != nil {
				return fmt.Errorf("failed: %v", err)
			}
		}

		os.Exit(0)

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
			HelmTimeout: ctx.Duration("helm-timeout"),
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
