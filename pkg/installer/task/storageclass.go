package task

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/kubermatic/kubermatic-installer/pkg/client/kubernetes"
)

type EnsureStorageClassTask struct {
	StorageClass *kubernetes.StorageClass
}

func (t *EnsureStorageClassTask) Run(config *Config, state *State, clients *Clients, log logrus.FieldLogger, dryRun bool) error {
	if state.Cluster.HasStorageClass(t.StorageClass.Metadata.Name) {
		return nil
	}

	log.WithFields(logrus.Fields{
		"provisioner": t.StorageClass.Provisioner,
		"parameters":  t.StorageClass.Parameters,
	}).Infof("Creating StorageClass %s…", t.StorageClass.Metadata.Name)

	if !dryRun {
		err := clients.Kubernetes.CreateStorageClass(*t.StorageClass)
		if err != nil {
			return fmt.Errorf("StorageClass could not be created: %v", err)
		}

		log.Infof("StorageClass has been created successfully.")
	}

	// always create the side-effect on the state, even in non-dry-run modes
	state.Cluster.StorageClasses = append(state.Cluster.StorageClasses, *t.StorageClass)

	return nil
}
