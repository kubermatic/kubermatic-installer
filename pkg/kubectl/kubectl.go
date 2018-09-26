package kubectl

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type Kubectl interface {
	CreateServiceAccount(namespace string, name string) error
	CreateClusterRoleBinding(name string, clusterRole string, serviceAccount string) error
}

type kubectl struct {
	kubeconfig string
	logger     logrus.FieldLogger
}

func NewKubectl(kubeconfig string, logger logrus.FieldLogger) (Kubectl, error) {
	return &kubectl{
		kubeconfig: kubeconfig,
		logger:     logger,
	}, nil
}

func (k *kubectl) CreateServiceAccount(namespace string, name string) error {
	k.logger.Infof("Creating serviceaccount %s:%s...", namespace, name)

	if err := k.run("create", "serviceaccount", "-n", namespace, name); err != nil {
		return fmt.Errorf("failed to create serviceaccount: %v", err)
	}

	return nil
}

func (k *kubectl) CreateClusterRoleBinding(name string, clusterRole string, serviceAccount string) error {
	k.logger.Infof("Creating clusterrolebinding %s...", name)

	if err := k.run("create", "clusterrolebinding", name, "--clusterrole", clusterRole, "--serviceaccount", serviceAccount); err != nil {
		return fmt.Errorf("failed to create clusterrolebinding: %v", err)
	}

	return nil
}

func (k *kubectl) run(args ...string) error {
	cmd := exec.Command("kubectl", args...)
	cmd.Env = append(cmd.Env, "KUBECONFIG="+k.kubeconfig)

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(stdoutStderr))
	}

	return nil
}
