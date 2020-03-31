package command

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/kubermatic/kubermatic-installer/pkg/client/helm"
	"github.com/kubermatic/kubermatic-installer/pkg/client/kubernetes"
	"github.com/kubermatic/kubermatic-installer/pkg/installer/stack/kubermatic"
	"github.com/kubermatic/kubermatic-installer/pkg/installer/state"
	"github.com/kubermatic/kubermatic-installer/pkg/installer/task"
	"github.com/kubermatic/kubermatic-installer/pkg/log"
	"github.com/kubermatic/kubermatic-installer/pkg/shared"
	"github.com/kubermatic/kubermatic-installer/pkg/shared/operatorv1alpha1"
	"github.com/kubermatic/kubermatic-installer/pkg/yamled"

	"sigs.k8s.io/yaml"
)

func DeployCommand(logger *logrus.Logger) cli.Command {
	return cli.Command{
		Name:   "deploy",
		Usage:  "Installs or upgrades the current installation to the installer's built-in version",
		Action: DeployAction(logger),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "confirm",
				Usage: "Perform the deployment instead of just listing the required steps",
			},
			cli.BoolFlag{
				Name:  "force",
				Usage: "Perform Helm upgrades even when the release is up-to-date",
			},
			cli.StringFlag{
				Name:   "config",
				Usage:  "Full path to the KubermaticConfiguration YAML file",
				EnvVar: "CONFIG_YAML",
			},
			cli.StringFlag{
				Name:   "helm-values",
				Usage:  "Full path to the Helm values.yaml used for customizing all charts",
				EnvVar: "VALUES_YAML",
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

		// load config files
		kubermaticConfig, err := loadKubermaticConfiguration(ctx.String("config"))
		if err != nil {
			return fmt.Errorf("failed to load KubermaticConfiguration: %v", err)
		}

		helmValues, err := loadHelmValues(ctx.String("helm-values"))
		if err != nil {
			return fmt.Errorf("failed to load Helm values: %v", err)
		}

		// validate the configuration
		logger.Info("Validating the provided configuration…")

		kubermaticConfig, helmValues, validationErrors := kubermatic.ValidateConfiguration(kubermaticConfig, helmValues, logger)
		if len(validationErrors) > 0 {
			logger.Error("The provided configuration files are invalid:")

			for _, e := range validationErrors {
				logger.Errorf("✘ %v", e)
			}

			return errors.New("please review your configuration and try again")
		}

		logger.Info("Provided configuration is valid.")

		// prepapre Kubernetes and Helm clients
		kubeconfig := ctx.String("kubeconfig")
		if len(kubeconfig) == 0 {
			return errors.New("no kubeconfig (--kubeconfig or $KUBECONFIG) given")
		}

		kubeContext := ctx.String("kube-context")
		helmTimeout := ctx.Duration("helm-timeout")

		kubeClient, err := kubernetes.NewKubectl(kubeconfig, kubeContext, logger)
		if err != nil {
			return fmt.Errorf("failed to create Kubernetes client: %v", err)
		}

		helmClient, err := helm.NewCLI(kubeconfig, kubeContext, helmTimeout, logger)
		if err != nil {
			return fmt.Errorf("failed to create Helm client: %v", err)
		}

		// inspect remote cluster and list helm releases
		clusterState, err := state.NewClusterState(kubeClient, helmClient)
		if err != nil {
			return fmt.Errorf("failed to determine cluster state: %v", err)
		}

		// determine tasks required to upgrade the cluster to the installer state
		logger.Info("Planning kubermatic stack deployment…")
		tasks, err := kubermatic.DeploymentTasks(installerState, clusterState)
		if err != nil {
			return fmt.Errorf("failed to determine deployment steps: %v", err)
		}

		// prepare tasks
		state := task.State{
			Cluster:   clusterState,
			Installer: installerState,
		}

		config := task.Config{
			Kubermatic: kubermaticConfig,
			Helm:       helmValues,
		}

		clients := task.Clients{
			Kubernetes: kubeClient,
			Helm:       helmClient,
		}

		options := task.Options{
			DryRun:           !ctx.Bool("confirm"),
			ForceHelmUpgrade: ctx.Bool("force"),
		}

		// check if any of the tasks are actually required to run
		requiredTasks := []task.Task{}

		for idx, t := range tasks {
			required, err := t.Required(&config, &state, &options)
			if err != nil {
				return err
			}

			if required {
				requiredTasks = append(requiredTasks, tasks[idx])
			}
		}

		if len(requiredTasks) == 0 {
			logger.Info("✓ Your installation is already up-to-date.")
			return nil
		}

		// print the task preview with a special fancy logger
		taskLogger := log.NewPlan()

		for _, t := range requiredTasks {
			err := t.Plan(&config, &state, &options, taskLogger)
			if err != nil {
				return err
			}
		}

		if options.DryRun {
			logger.Warn("Run the installer with --confirm to actually perform the steps outlined above.")
			return nil
		}

		// let the magic happen
		for _, t := range requiredTasks {
			err := t.Run(&config, &state, &clients, &options, logger)
			if err != nil {
				return err
			}
		}

		logger.Info("✓ Installation completed successfully.")

		return nil
	}))
}

func loadKubermaticConfiguration(filename string) (*operatorv1alpha1.KubermaticConfiguration, error) {
	if filename == "" {
		return nil, errors.New("no file specified via --config flag")
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := operatorv1alpha1.KubermaticConfiguration{}
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("failed to decode %s: %v", filename, err)
	}

	return &config, nil
}

func loadHelmValues(filename string) (*yamled.Document, error) {
	if filename == "" {
		return nil, errors.New("no file specified via --helm-values flag")
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	values, err := yamled.Load(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %s: %v", filename, err)
	}

	return values, nil
}
