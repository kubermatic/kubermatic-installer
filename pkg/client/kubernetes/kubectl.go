package kubernetes

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type kubectl struct {
	kubeconfig  string
	kubeContext string
	logger      logrus.FieldLogger
}

// NewKubectl returns a Client implementation that uses a local
// kubectl binary to interact with a given Kubernetes cluster.
func NewKubectl(kubeconfig string, kubeContext string, logger logrus.FieldLogger) (Client, error) {
	return &kubectl{
		kubeconfig:  kubeconfig,
		kubeContext: kubeContext,
		logger:      logger,
	}, nil
}

func (k *kubectl) CreateServiceAccount(namespace string, name string) error {
	k.logger.Infof("Creating service account %s:%s...", namespace, name)

	k.logger.Debug("Checking if the servic eaccount already exists...")
	exists, err := k.exists(namespace, "serviceaccount", name)
	if err != nil {
		return fmt.Errorf("failed to check for service account: %v", err)
	}

	if exists {
		k.logger.Debug("Service account already exists, skipping creation.")
	} else {
		if _, err := k.run("create", "serviceaccount", "-n", namespace, name); err != nil {
			return fmt.Errorf("failed to create service account: %v", err)
		}
	}

	return nil
}

func (k *kubectl) CreateClusterRoleBinding(name string, clusterRole string, serviceAccount string) error {
	k.logger.Infof("Creating cluster role binding %s...", name)

	k.logger.Debug("Checking if the cluster role binding already exists...")
	exists, err := k.exists("", "clusterrolebinding", name)
	if err != nil {
		return fmt.Errorf("failed to check for cluster role binding: %v", err)
	}

	if exists {
		k.logger.Debug("Cluster role binding already exists, skipping creation.")
	} else {
		if _, err := k.run("create", "clusterrolebinding", name, "--clusterrole", clusterRole, "--serviceaccount", serviceAccount); err != nil {
			return fmt.Errorf("failed to create cluster role binding: %v", err)
		}
	}

	return nil
}

func (k *kubectl) run(args ...string) (string, error) {
	cmd := exec.Command("kubectl", append([]string{"--kubeconfig", k.kubeconfig, "--context", k.kubeContext}, args...)...)

	k.logger.Debugf("$ %s", strings.Join(cmd.Args, " "))

	stdoutStderr, err := cmd.CombinedOutput()
	output := string(stdoutStderr)

	if err != nil {
		err = errors.New(output)
	}

	return output, err
}

func (k *kubectl) exists(namespace string, kind string, name string) (bool, error) {
	var args []string

	if len(namespace) > 0 {
		args = []string{"get", "-n", namespace, "--ignore-not-found", kind, name}
	} else {
		args = []string{"get", "--ignore-not-found", kind, name}
	}

	output, err := k.run(args...)
	if err != nil {
		return false, err
	}

	return len(output) > 0, nil
}
