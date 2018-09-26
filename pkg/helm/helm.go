package helm

import (
	"errors"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type Helm interface {
	Init(serviceAccount string, tillerNamespace string) error
	InstallChart(namespace string, name string, directory string, values string) error
}

type helm struct {
	kubeconfig string
	logger     logrus.FieldLogger
}

func NewHelm(kubeconfig string, logger logrus.FieldLogger) (Helm, error) {
	return &helm{
		kubeconfig: kubeconfig,
		logger:     logger,
	}, nil
}

func (h *helm) Init(serviceAccount string, tillerNamespace string) error {
	h.logger.Infof("Installing Helm using service account %s into tiller namespace %s...", serviceAccount, tillerNamespace)

	return h.run("init", "--service-account", serviceAccount, "--tiller-namespace", tillerNamespace)
}

func (h *helm) InstallChart(namespace string, name string, directory string, values string) error {
	h.logger.Infof("Installing chart %s into namespace %s...", name, namespace)

	command := []string{
		"upgrade",
		"--install",
		"--wait",
		"--timeout", "300",
		"--tiller-namespace", "kube-system",
		"--kube-context", "default",
		"--values", values,
		"--namespace", namespace,
		name,
		directory,
	}

	return h.run(command...)
}

func (h *helm) run(args ...string) error {
	cmd := exec.Command("helm", args...)
	cmd.Env = append(cmd.Env, "KUBECONFIG="+h.kubeconfig)

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(stdoutStderr))
	}

	return nil
}
