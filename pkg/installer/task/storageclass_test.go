package task

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kubermatic/kubermatic-installer/pkg/client/helm"
	"github.com/kubermatic/kubermatic-installer/pkg/installer/state"

	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestEnsureStorageClassTask(t *testing.T) {
	class := storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
		Provisioner: "example",
	}

	task := EnsureStorageClassTask{
		StorageClass: &class,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log := logrus.New()
	log.SetOutput(ioutil.Discard)

	state, err := state.NewInstallerState("../../../charts")
	if err != nil {
		t.Fatalf("Failed to create installer state: %v", err)
	}

	options := Options{}
	client := fake.NewFakeClient()

	helm, err := helm.NewCLI("", "", 30*time.Second, log)
	if err != nil {
		t.Fatalf("Failed to create Helm client: %v", err)
	}

	err = task.Run(ctx, &options, state, client, helm, log)
	if err != nil {
		t.Fatalf("Failed to run task: %v", err)
	}

	result := storagev1.StorageClass{}
	if err := client.Get(ctx, types.NamespacedName{Name: class.Name}, &result); err != nil {
		t.Fatalf("Expected storage class to exist, but failed to retrieve it: %v", err)
	}

	if result.Provisioner != class.Provisioner {
		t.Fatalf("Expected provisioner to be %q, but got %q.", class.Provisioner, result.Provisioner)
	}

	// running the task again should be safe
	err = task.Run(ctx, &options, state, client, helm, log)
	if err != nil {
		t.Fatalf("Failed to run task again: %v", err)
	}
}
