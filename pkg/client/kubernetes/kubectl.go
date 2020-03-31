package kubernetes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
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

func (k *kubectl) CreateNamespace(name string) error {
	k.logger.Infof("Creating namespace %s…", name)

	k.logger.Debug("Checking if it already exists…")
	exists, err := k.exists("", "namespace", name)
	if err != nil {
		return fmt.Errorf("failed to check for namespace: %v", err)
	}

	if exists {
		k.logger.Debug("Namespace already exists, skipping creation.")
	} else {
		if _, err := k.run("create", "namespace", name); err != nil {
			return fmt.Errorf("failed to create namespace: %v", err)
		}
	}

	return nil
}

func (k *kubectl) CreateServiceAccount(namespace string, name string) error {
	k.logger.Infof("Creating service account %s:%s…", namespace, name)

	k.logger.Debug("Checking if the servic eaccount already exists…")
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
	k.logger.Infof("Creating cluster role binding %s…", name)

	k.logger.Debug("Checking if the cluster role binding already exists…")
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

func (k *kubectl) HasStorageClass(name string) (bool, error) {
	k.logger.Infof("Checking for storage class %s…", name)

	return k.exists("", "storageclass", name)
}

func (k *kubectl) ServiceIngresses(namespace string, serviceName string) ([]Ingress, error) {
	k.logger.Infof("Retrieving service %s/%s ingresses…", namespace, serviceName)

	output, err := k.run("-n", namespace, "get", "service", serviceName, "-o", "json")
	if err != nil {
		return nil, fmt.Errorf("failed to create namespace: %v", err)
	}

	type kubectlOutput struct {
		Status struct {
			LoadBalancer struct {
				Ingress []Ingress `json:"ingress"`
			} `json:"loadBalancer"`
		} `json:"status"`
	}

	out := kubectlOutput{}
	if err := json.Unmarshal([]byte(output), &out); err != nil {
		return nil, fmt.Errorf("failed to parse kubectl JSON: %v", err)
	}

	return out.Status.LoadBalancer.Ingress, nil
}

func (k *kubectl) CreateStorageClass(sc StorageClass) error {
	tmpFile, err := k.dumpResource(sc)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	_, err = k.run("create", "-f", tmpFile)
	if err != nil {
		return fmt.Errorf("failed to create StorageClass: %v", err)
	}

	return nil
}

func (k *kubectl) StorageClasses() ([]StorageClass, error) {
	k.logger.Info("Retrieving storage classes…")

	output, err := k.run("get", "storageclass", "-o", "yaml")
	if err != nil {
		return nil, fmt.Errorf("failed to list storage classes: %v", err)
	}

	type kubectlOutput struct {
		Items []StorageClass `yaml:"items"`
	}

	out := kubectlOutput{}
	if err := yaml.Unmarshal([]byte(output), &out); err != nil {
		return nil, fmt.Errorf("failed to parse kubectl JSON: %v", err)
	}

	return out.Items, nil
}

func (k *kubectl) HasService(namespace string, name string) (bool, error) {
	k.logger.Infof("Checking for service %s/%s…", namespace, name)

	return k.exists(namespace, "service", name)
}

func (k *kubectl) HasCustomResourceDefinition(name string) (bool, error) {
	k.logger.Infof("Checking for CRD %s…", name)

	return k.exists("", "customresourcedefinition", name)
}

func (k *kubectl) ApplyManifests(source string) error {
	k.logger.Infof("Applying manifests from %s…", source)

	_, err := k.run("apply", "-f", source)

	return err
}

func (k *kubectl) run(args ...string) (string, error) {
	cmd := exec.Command("kubectl", append([]string{"--context", k.kubeContext}, args...)...)
	cmd.Env = append(cmd.Env, "KUBECONFIG="+k.kubeconfig)

	k.logger.Debugf("$ KUBECONFIG=%s %s", k.kubeconfig, strings.Join(cmd.Args, " "))

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

func (k *kubectl) dumpResource(res interface{}) (string, error) {
	tmpfile, err := ioutil.TempFile("", "kubermatic")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %v", err)
	}

	err = yaml.NewEncoder(tmpfile).Encode(res)
	if err != nil {
		return "", fmt.Errorf("failed to encode resource as YAML: %v", err)
	}

	err = tmpfile.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close file: %v", err)
	}

	return tmpfile.Name(), nil
}
