package kubernetes

// Client describes the operations that are required to
// perform the Kubermatic installation.
type Client interface {
	CreateNamespace(name string) error
	CreateServiceAccount(namespace string, name string) error
	CreateClusterRoleBinding(name string, clusterRole string, serviceAccount string) error
}
