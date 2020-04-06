package helm

// Client describes the operations that the Helm client is providing to
// the installer. This is the minimum set of operations required to
// perform a Kubermatic installation.
type Client interface {
	InstallChart(namespace string, releaseName string, chartDirectory, valuesFile string, flags map[string]string) error
	GetRelease(namespace string, name string) (*Release, error)
	ListReleases(namespace string) ([]Release, error)
	UninstallRelease(namespace string, name string) error
}
