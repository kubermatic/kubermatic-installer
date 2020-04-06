package task

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	storagev1 "k8s.io/api/storage/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
)

type EnsureStorageClassTask struct {
	StorageClass *storagev1.StorageClass
}

func (t *EnsureStorageClassTask) Required(_ *Config, state *State, _ *Options) (bool, error) {
	return !state.Cluster.HasStorageClass(t.StorageClass.Name), nil
}

func (t *EnsureStorageClassTask) Plan(_ *Config, _ *State, _ *Options, log logrus.FieldLogger) error {
	log.WithFields(logrus.Fields{
		"provisioner": t.StorageClass.Provisioner,
		"parameters":  t.StorageClass.Parameters,
	}).Infof("Create %s StorageClass.", t.StorageClass.Name)

	return nil
}

func (t *EnsureStorageClassTask) Run(ctx context.Context, _ *Config, state *State, clients *Clients, _ *Options, log logrus.FieldLogger) error {
	log.WithFields(logrus.Fields{
		"provisioner": t.StorageClass.Provisioner,
		"parameters":  t.StorageClass.Parameters,
	}).Infof("Creating StorageClass %s…", t.StorageClass.Name)

	err := clients.Kubernetes.Create(ctx, t.StorageClass)
	if err != nil && !kerrors.IsAlreadyExists(err) {
		return fmt.Errorf("StorageClass could not be created: %v", err)
	}

	log.Infof("StorageClass has been created successfully.")

	// update cluster state
	state.Cluster.StorageClasses = append(state.Cluster.StorageClasses, *t.StorageClass)

	return nil
}
