package helm

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Helm interface {
	Close() error
}

type helm struct {
	kubeconfigFile string
}

func NewHelm(kubeconfig string) (Helm, error) {
	tmpfile, err := ioutil.TempFile("", "kubermatic.*.kubeconfig")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary kubeconfig file: %v", err)
	}

	_, err = tmpfile.WriteString(kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to write kubeconfig to file: %v", err)
	}

	err = tmpfile.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close kubeconfig file: %v", err)
	}

	return &helm{
		kubeconfigFile: tmpfile.Name(),
	}, nil
}

func (h *helm) Close() error {
	var err error

	if len(h.kubeconfigFile) > 0 {
		err = os.Remove(h.kubeconfigFile)
		h.kubeconfigFile = ""
	}

	return err
}
