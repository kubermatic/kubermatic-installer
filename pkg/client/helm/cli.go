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

// helm3 --namespace test list -o json
// [{"name":"blackbox-exporter","namespace":"test","revision":"2","updated":"2020-03-16 18:29:14.704404538 +0100 CET","status":"deployed","chart":"blackbox-exporter-1.0.3","app_version":"v0.15.2"}]

// helm3 --namespace test history blackbox-exporter -o json | jq
// [
//   {
//     "revision": 1,
//     "updated": "2020-03-16T18:22:35.155907509+01:00",
//     "status": "deployed",
//     "chart": "blackbox-exporter-1.0.2",
//     "app_version": "v0.15.1",
//     "description": "Install complete"
//   }
// ]

func (c *cli) ListReleases(namespace string) ([]Release, error) {
	logger := c.logger.WithField("namespace", namespace)
	logger.Info("Listing installed releases…")

	args := []string{"list", "-o", "json"}
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
		if err != nil {
			logger.Warnf("Release %s/%s has no valid version number (%q): %v", release.Namespace, release.Name, tail, err)
		} else {
			releases[idx].Version = version
		}

		releases[idx].Chart = strings.Join(nameParts[:len(nameParts)-1], "-")
	}

	return releases, nil
}

func (c *cli) run(namespace string, args ...string) ([]byte, error) {
	globalArgs := []string{"--kube-context", c.kubeContext}

	if namespace != "" {
		globalArgs = append(globalArgs, "--namespace", namespace)
	}

	cmd := exec.Command("helm3", append(globalArgs, args...)...)
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
