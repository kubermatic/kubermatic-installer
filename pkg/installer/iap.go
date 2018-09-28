package installer

type IAPDeploymentIngress struct {
	Host string `yaml:"host"`
}

type IAPDeploymentConfig map[string]interface{}

type IAPDeployment struct {
	Name            string               `yaml:"name"`
	ClientID        string               `yaml:"client_id"`
	ClientSecret    string               `yaml:"client_secret"`
	EncryptionKey   string               `yaml:"encryption_key"`
	UpstreamService string               `yaml:"upstream_service"`
	UpstreamPort    int                  `yaml:"upstream_port"`
	Ingress         IAPDeploymentIngress `yaml:"ingress"`
	Config          IAPDeploymentConfig  `yaml:"config"`
}

type IAPResource struct {
	URI    string   `yaml:"uri"`
	Groups []string `yaml:"groups,omitempty"`
}

func NewNullIAPResource() IAPResource {
	return IAPResource{
		URI: "/*",
	}
}
