package kubernetes

type Ingress struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
}

type StorageClass struct {
	Kind        string               `yaml:"kind"`
	APIVersion  string               `yaml:"apiVersion"`
	Metadata    StorageClassMetadata `yaml:"metadata"`
	Provisioner string               `yaml:"provisioner"`
	Parameters  map[string]string    `yaml:"parameters"`
}

type StorageClassMetadata struct {
	Name        string            `yaml:"name"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

func NewStorageClass() StorageClass {
	return StorageClass{
		Kind:       "StorageClass",
		APIVersion: "storage.k8s.io/v1",
		Parameters: make(map[string]string),
	}
}
