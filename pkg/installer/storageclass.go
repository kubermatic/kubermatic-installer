package installer

import (
	"github.com/kubermatic/kubermatic-installer/pkg/client/kubernetes"
	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
)

func NewEKSStorageClass() kubernetes.StorageClass {
	s := kubernetes.NewStorageClass()
	s.Metadata.Name = "kubermatic-fast"
	s.Provisioner = "kubernetes.io/aws-ebs"
	s.Parameters["type"] = "gp2"

	return s
}

func NewGKEStorageClass() kubernetes.StorageClass {
	s := kubernetes.NewStorageClass()
	s.Metadata.Name = "kubermatic-fast"
	s.Provisioner = "kubernetes.io/gce-pd"
	s.Parameters["type"] = "pd-ssd"

	return s
}

func NewAKSStorageClass() kubernetes.StorageClass {
	s := kubernetes.NewStorageClass()
	s.Metadata.Name = "kubermatic-fast"
	s.Provisioner = "kubernetes.io/azure-disk"
	s.Parameters["storageaccounttype"] = "Standard_LRS"
	s.Parameters["kind"] = "managed"

	return s
}

func StorageClassForProvider(p manifest.CloudProvider) *kubernetes.StorageClass {
	var sc kubernetes.StorageClass

	switch p {
	case manifest.ProviderAKS:
		sc = NewAKSStorageClass()
	case manifest.ProviderEKS:
		sc = NewEKSStorageClass()
	case manifest.ProviderGKE:
		sc = NewGKEStorageClass()
	default:
		return nil
	}

	return &sc
}
