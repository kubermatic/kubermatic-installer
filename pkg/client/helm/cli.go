package helm

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type cli struct {
	kubeconfig  string
	kubeContext string
	timeout     time.Duration
	logger      logrus.FieldLogger
}

// NewCLI returns a new Client implementation that uses a local helm
// binary to perform chart installations.
func NewCLI(kubeconfig string, kubeContext string, timeout time.Duration, logger logrus.FieldLogger) (Client, error) {
	if timeout.Seconds() < 10 {
		return nil, errors.New("timeout must be >= 10 seconds")
	}

	return &cli{
		kubeconfig:  kubeconfig,
		kubeContext: kubeContext,
		timeout:     timeout,
		logger:      logger,
	}, nil
}

func (c *cli) InstallChart(namespace string, name string, directory string, values string, flags map[string]string, wait bool) error {
	c.logger.Infof("Installing chart %s into namespace %s…", name, namespace)

	// Check if there is an existing release and it failed;
	// sometimes installations can fail because prerequisites were not setup properly,
	// like a missing storage class. In this case, we want to allow the user to just
	// run the installer again and pick up where they left. Unfortunately Helm does not
	// support "upgrade --install" on failed installations: https://github.com/helm/helm/issues/3353
	// To work around this, we check the release status and purge it manually if it's failed.
	status := c.releaseStatus(namespace, name)

	if c.isPurgeable(status) {
		c.logger.Warnf("Release status is %s, purging release before attempting to install it.", status)

		_, err := c.run(namespace, "uninstall", name)
		if err != nil {
			return fmt.Errorf("failed to uninstall existing release: %v", err)
		}
	} else {
		c.logger.Debugf("Release status is %s, attempting in-place upgrade.", status)
	}

	command := []string{
		"upgrade",
		"--install",
		"--values", values,
	}

	set := make([]string, 0)

	for name, value := range flags {
		set = append(set, fmt.Sprintf("%s=%s", name, value))
	}

	if len(set) > 0 {
		command = append(command, "--set", strings.Join(set, ","))
	}

	if wait {
		command = append(command, "--wait", "--timeout", c.timeout.String())
	}

	command = append(command, name, directory)

	_, err := c.run(namespace, command...)

	return err
}

func (c *cli) run(namespace string, args ...string) ([]byte, error) {
	args = append([]string{
		"--namespace", namespace,
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

type helmStatus struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Manifest  string `json:"manifest"`
	Info      struct {
		FirstDeployed time.Time     `json:"first_deployed"`
		LastDeployed  time.Time     `json:"last_deployed"`
		Status        releaseStatus `json:"status"`
	} `json:"info"`
}

func (c *cli) releaseStatus(namespace string, name string) releaseStatus {
	c.logger.Debugf("Checking release status…")

	output, err := c.run(namespace, "status", name, "-o", "json")
	if err != nil {
		return releaseCheckFailed
	}

	status := helmStatus{}
	err = json.Unmarshal(output, &status)
	if err != nil {
		return releaseCheckFailed
	}

	return status.Info.Status
}

// isPurgeable determines whether a Helm release status indicates
// that we should delete the release before attempting to re-install
// it.
func (c *cli) isPurgeable(status releaseStatus) bool {
	return status != releaseCheckFailed && status != releaseUnknown && status != releaseDeployed
}
