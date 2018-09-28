package helm

// Client describes the operations that the Helm client is providing to
// the installer. This is the minimum set of operations required to
// perform a Kubermatic installation.
type Client interface {
	Init(serviceAccount string) error
	InstallChart(namespace string, name string, directory string, values string) error
}
