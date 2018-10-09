package helm

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type cli struct {
	kubeconfig      string
	kubeContext     string
	tillerNamespace string
	timeout         int
	logger          logrus.FieldLogger
}

// NewCLI returns a new Client implementation that uses a local helm
// binary to perform chart installations.
func NewCLI(kubeconfig string, kubeContext string, tillerNamespace string, timeout int, logger logrus.FieldLogger) (Client, error) {
	if timeout < 10 {
		return nil, errors.New("timeout must be >= 10 seconds")
	}

	return &cli{
		kubeconfig:      kubeconfig,
		kubeContext:     kubeContext,
		tillerNamespace: tillerNamespace,
		timeout:         timeout,
		logger:          logger,
	}, nil
}

func (c *cli) Init(serviceAccount string) error {
	c.logger.Infof("Installing Helm using service account %s into tiller namespace %s...", serviceAccount, c.tillerNamespace)

	_, err := c.run("init", "--service-account", serviceAccount, "--wait")

	return err
}

func (c *cli) InstallChart(namespace string, name string, directory string, values string, wait bool) error {
	c.logger.Infof("Installing chart %s into namespace %s...", name, namespace)

	// Check if there is an existing release and it failed;
	// sometimes installations can fail because prerequisites were not setup properly,
	// like a missing storage class. In this case, we want to allow the user to just
	// run the installer again and pick up where they left. Unfortunately Helm does not
	// support "upgrade --install" on failed installations: https://github.com/helm/helm/issues/3353
	// To work around this, we check the release status and purge it manually if it's failed.
	status := c.releaseStatus(name)

	if c.isPurgeable(status) {
		c.logger.Warnf("Release status is %s, purging release before attempting to install it.", status)

		_, err := c.run("delete", "--purge", name)
		if err != nil {
			return fmt.Errorf("failed to clean-up existing release: %v", err)
		}
	} else {
		c.logger.Debugf("Release status is %s, attempting in-place upgrade.", status)
	}

	command := []string{
		"upgrade",
		"--install",
		"--values", values,
		"--namespace", namespace,
	}

	if wait {
		command = append(command, "--wait", "--timeout", strconv.Itoa(c.timeout))
	}

	command = append(command, name, directory)

	_, err := c.run(command...)

	return err
}

func (c *cli) run(args ...string) ([]byte, error) {
	args = append([]string{
		"--tiller-namespace", c.tillerNamespace,
		"--kube-context", c.kubeContext,
	}, args...)

	cmd := exec.Command("helm", args...)
	cmd.Env = append(cmd.Env, "KUBECONFIG="+c.kubeconfig)

	c.logger.Debugf("$ KUBECONFIG=%s %s", c.kubeconfig, strings.Join(cmd.Args, " "))

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.New(string(stdoutStderr))
	}

	return stdoutStderr, err
}

func (c *cli) releaseStatus(name string) releaseStatus {
	c.logger.Debugf("Checking release status...")

	output, err := c.run(
		"--tiller-namespace", c.tillerNamespace,
		"--kube-context", c.kubeContext,
		"status",
		name,
		"-o",
		"json",
	)
	if err != nil {
		return releaseCheckFailed
	}

	type helmOutput struct {
		Info struct {
			Status struct {
				Code releaseStatus `json:"code"`
			} `json:"status"`
		} `json:"info"`
	}

	status := helmOutput{}
	err = json.Unmarshal(output, &status)
	if err != nil {
		return releaseCheckFailed
	}

	return status.Info.Status.Code
}

// isPurgeable determines whether a Helm release status indicates
// that we should delete the release before attempting to re-install
// it.
func (c *cli) isPurgeable(status releaseStatus) bool {
	return status != releaseCheckFailed && status != releaseUnknown && status != releaseDeployed
}
