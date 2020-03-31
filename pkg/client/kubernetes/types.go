package kubernetes

const (
	DefaultStorageClassAnnotation     = "storageclass.kubernetes.io/is-default-class"
	DefaultStorageClassAnnotationBeta = "storageclass.beta.kubernetes.io/is-default-class"
)

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

func (s *StorageClass) IsDefault() bool {
	return s.Metadata.Annotations[DefaultStorageClassAnnotation] == "true" || s.Metadata.Annotations[DefaultStorageClassAnnotationBeta] == "true"
}

func (s *StorageClass) Clone() StorageClass {
	copy := NewStorageClass()

	copy.Kind = s.Kind
	copy.APIVersion = s.APIVersion
	copy.Provisioner = s.Provisioner
	copy.Metadata = StorageClassMetadata{
		Name:        s.Metadata.Name,
		Annotations: make(map[string]string),
	}

	for k, v := range s.Parameters {
		copy.Parameters[k] = v
	}

	for k, v := range s.Metadata.Annotations {
		copy.Metadata.Annotations[k] = v
	}

	return copy
}
