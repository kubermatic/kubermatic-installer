package kubermatic

import (
	"github.com/kubermatic/kubermatic-installer/pkg/manifest"

	storagev1 "k8s.io/api/storage/v1"
)

func newEKSStorageClass() storagev1.StorageClass {
	s := storagev1.StorageClass{}
	s.Provisioner = "kubernetes.io/aws-ebs"
	s.Parameters["type"] = "gp2"

	return s
}

func newGKEStorageClass() storagev1.StorageClass {
	s := storagev1.StorageClass{}
	s.Provisioner = "kubernetes.io/gce-pd"
	s.Parameters["type"] = "pd-ssd"

	return s
}

func newAKSStorageClass() storagev1.StorageClass {
	s := storagev1.StorageClass{}
	s.Provisioner = "kubernetes.io/azure-disk"
	s.Parameters["storageaccounttype"] = "Standard_LRS"
	s.Parameters["kind"] = "managed"

	return s
}

func storageClassForProvider(name string, p manifest.CloudProvider) *storagev1.StorageClass {
	var sc storagev1.StorageClass

	switch p {
	case manifest.ProviderAKS:
		sc = newAKSStorageClass()
	case manifest.ProviderEKS:
		sc = newEKSStorageClass()
	case manifest.ProviderGKE:
		sc = newGKEStorageClass()
	default:
		return nil
	}

	sc.Name = name

	return &sc
}
