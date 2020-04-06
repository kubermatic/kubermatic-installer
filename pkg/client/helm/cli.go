package helm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/Masterminds/semver"
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

func (c *cli) InstallChart(namespace string, releaseName string, chartDirectory string, valuesFile string, flags map[string]string) error {
	command := []string{
		"upgrade",
		"--install",
		"--values", valuesFile,
		"--atomic", // implies --wait
		"--timeout", c.timeout.String(),
	}

	set := make([]string, 0)

	for name, value := range flags {
		set = append(set, fmt.Sprintf("%s=%s", name, value))
	}

	if len(set) > 0 {
		command = append(command, "--set", strings.Join(set, ","))
	}

	command = append(command, releaseName, chartDirectory)

	_, err := c.run(namespace, command...)

	return err
}

func (c *cli) GetRelease(namespace string, name string) (*Release, error) {
	releases, err := c.ListReleases(namespace)
	if err != nil {
		return nil, err
	}

	for idx, r := range releases {
		if r.Namespace == namespace && r.Name == name {
			return &releases[idx], nil
		}
	}

	return nil, nil
}

func (c *cli) ListReleases(namespace string) ([]Release, error) {
	args := []string{"list", "--all", "-o", "json"}
	if namespace == "" {
		args = append(args, "--all-namespaces")
	}

	output, err := c.run(namespace, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list releases: %v", err)
	}

	releases := []Release{}
	if err := json.NewDecoder(bytes.NewReader(output)).Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to parse Helm output: %v", err)
	}

	for idx, release := range releases {
		nameParts := strings.Split(release.Chart, "-")
		tail := nameParts[len(nameParts)-1]

		version, err := semver.NewVersion(tail)
		if err == nil {
			releases[idx].Version = version
		}

		releases[idx].Chart = strings.Join(nameParts[:len(nameParts)-1], "-")
	}

	return releases, nil
}

func (c *cli) UninstallRelease(namespace string, name string) error {
	_, err := c.run(namespace, "uninstall", name)

	return err
}

func (c *cli) run(namespace string, args ...string) ([]byte, error) {
	globalArgs := []string{}

	if c.kubeContext != "" {
		globalArgs = append(globalArgs, "--kube-context", c.kubeContext)
	}

	if namespace != "" {
		globalArgs = append(globalArgs, "--namespace", namespace)
	}

	cmd := exec.Command("helm3", append(globalArgs, args...)...)
	cmd.Env = append(cmd.Env, "KUBECONFIG="+c.kubeconfig)

	c.logger.Debugf("$ KUBECONFIG=%s %s", c.kubeconfig, strings.Join(cmd.Args, " "))

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.New(strings.TrimSpace(string(stdoutStderr)))
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
		Status        ReleaseStatus `json:"status"`
	} `json:"info"`
}

func (c *cli) releaseStatus(namespace string, name string) ReleaseStatus {
	output, err := c.run(namespace, "status", name, "-o", "json")
	if err != nil {
		return ReleaseCheckFailed
	}

	status := helmStatus{}
	err = json.Unmarshal(output, &status)
	if err != nil {
		return ReleaseCheckFailed
	}

	return status.Info.Status
}

// isPurgeable determines whether a Helm release status indicates
// that we should delete the release before attempting to re-install
// it.
func (c *cli) isPurgeable(status ReleaseStatus) bool {
	return status == ReleaseStatusFailed
}
