package kubernetes

// Client describes the operations that are required to
// perform the Kubermatic installation.
type Client interface {
	CreateNamespace(name string) error
	CreateServiceAccount(namespace string, name string) error
	CreateClusterRoleBinding(name string, clusterRole string, serviceAccount string) error
	HasStorageClass(name string) (bool, error)
	ServiceIngresses(namespace string, serviceName string) ([]Ingress, error)
	StorageClasses() ([]StorageClass, error)
	CreateStorageClass(sc StorageClass) error
	HasService(namespace string, name string) (bool, error)
	HasCustomResourceDefinition(name string) (bool, error)
	ApplyManifests(source string) error
}
