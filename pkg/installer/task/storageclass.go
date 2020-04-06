package task

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/kubermatic/kubermatic-installer/pkg/client/helm"
	"github.com/kubermatic/kubermatic-installer/pkg/installer/state"

	storagev1 "k8s.io/api/storage/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type EnsureStorageClassTask struct {
	StorageClass *storagev1.StorageClass
}

func (t *EnsureStorageClassTask) String() string {
	return fmt.Sprintf("Ensure storage class %s exists", t.StorageClass.Name)
}

func (t *EnsureStorageClassTask) Run(ctx context.Context, _ *Options, _ *state.InstallerState, kubeClient ctrlruntimeclient.Client, _ helm.Client, log logrus.FieldLogger) error {
	class := storagev1.StorageClass{}
	err := kubeClient.Get(ctx, types.NamespacedName{Name: t.StorageClass.Name}, &class)
	if err != nil && !kerrors.IsNotFound(err) {
		return err
	}

	if err == nil {
		log.WithFields(logrus.Fields{
			"provisioner": class.Provisioner,
			"parameters":  class.Parameters,
		}).Info("Storage class already exists.")

		return nil
	}

	if err := kubeClient.Create(ctx, t.StorageClass); err != nil {
		return err
	}

	return nil
}
