package helm

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type cli struct {
	kubeconfig      string
	kubeContext     string
	tillerNamespace string
	logger          logrus.FieldLogger
}

// NewCLI returns a new Client implementation that uses a local helm
// binary to perform chart installations.
func NewCLI(kubeconfig string, kubeContext string, tillerNamespace string, logger logrus.FieldLogger) (Client, error) {
	return &cli{
		kubeconfig:      kubeconfig,
		kubeContext:     kubeContext,
		tillerNamespace: tillerNamespace,
		logger:          logger,
	}, nil
}

func (c *cli) Init(serviceAccount string) error {
	c.logger.Infof("Installing Helm using service account %s into tiller namespace %s...", serviceAccount, c.tillerNamespace)

	return c.run("init", "--service-account", serviceAccount, "--tiller-namespace", c.tillerNamespace)
}

func (c *cli) InstallChart(namespace string, name string, directory string, values string) error {
	c.logger.Infof("Installing chart %s into namespace %s...", name, namespace)

	command := []string{
		"--tiller-namespace", c.tillerNamespace,
		"--kube-context", c.kubeContext,
		"upgrade",
		"--install",
		"--wait",
		"--timeout", "300",
		"--values", values,
		"--namespace", namespace,
		name,
		directory,
	}

	return c.run(command...)
}

func (c *cli) run(args ...string) error {
	cmd := exec.Command("helm", args...)
	cmd.Env = append(cmd.Env, "KUBECONFIG="+c.kubeconfig)

	c.logger.Debugf("$ KUBECONFIG=%s %s", c.kubeconfig, strings.Join(cmd.Args, " "))

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(stdoutStderr))
	}

	return nil
}
