package kubermatic

import (
	"github.com/kubermatic/kubermatic-installer/pkg/client/kubernetes"
	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
)

func newEKSStorageClass() kubernetes.StorageClass {
	s := kubernetes.NewStorageClass()
	s.Provisioner = "kubernetes.io/aws-ebs"
	s.Parameters["type"] = "gp2"

	return s
}

func newGKEStorageClass() kubernetes.StorageClass {
	s := kubernetes.NewStorageClass()
	s.Provisioner = "kubernetes.io/gce-pd"
	s.Parameters["type"] = "pd-ssd"

	return s
}

func newAKSStorageClass() kubernetes.StorageClass {
	s := kubernetes.NewStorageClass()
	s.Provisioner = "kubernetes.io/azure-disk"
	s.Parameters["storageaccounttype"] = "Standard_LRS"
	s.Parameters["kind"] = "managed"

	return s
}

func storageClassForProvider(name string, p manifest.CloudProvider) *kubernetes.StorageClass {
	var sc kubernetes.StorageClass

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

	sc.Metadata.Name = name

	return &sc
}
