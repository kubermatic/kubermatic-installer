package manifest

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kubermatic/kubermatic/api/pkg/provider"
)

const VERSION = "1"

type Manifest struct {
	Version        string                        `yaml:"version"`
	Kubeconfig     string                        `yaml:"kubeconfig"`
	Secrets        SecretsManifest               `yaml:"secrets"`
	SeedClusters   []string                      `yaml:"seedClusters"`
	Datacenters    map[string]DatacenterManifest `yaml:"datacenters"`
	Monitoring     MonitoringManifest            `yaml:"monitoring"`
	Logging        LoggingManifest               `yaml:"logging"`
	Authentication AuthenticationManifest        `yaml:"authentication"`
	Settings       SettingsManifest              `yaml:"settings"`
}

func (m *Manifest) Validate() error {
	if m.Version != VERSION {
		return errors.New("unknown or invalid manifest version")
	}

	if len(m.Kubeconfig) == 0 {
		return errors.New("no kubeconfig defined")
	}

	if len(m.SeedClusters) == 0 {
		return errors.New("no seed clusters defined")
	}

	for key, dc := range m.Datacenters {
		if err := dc.Validate(m.SeedClusters); err != nil {
			return fmt.Errorf("datacenter %s is invalid: %v", key, err)
		}
	}

	if err := m.Secrets.Validate(); err != nil {
		return fmt.Errorf("secrets configuration is invalid: %v", err)
	}

	if err := m.Monitoring.Validate(); err != nil {
		return fmt.Errorf("monitoring configuration is invalid: %v", err)
	}

	if err := m.Logging.Validate(); err != nil {
		return fmt.Errorf("logging configuration is invalid: %v", err)
	}

	if err := m.Authentication.Validate(); err != nil {
		return fmt.Errorf("authentication configuration is invalid: %v", err)
	}

	if err := m.Settings.Validate(); err != nil {
		return fmt.Errorf("settings are invalid: %v", err)
	}

	return nil
}

type kubermaticDatacenters struct {
	Datacenters map[string]provider.DatacenterMeta `yaml:"datacenters"`
}

func (m *Manifest) KubermaticDatacenters() *kubermaticDatacenters {
	spec := &kubermaticDatacenters{
		Datacenters: make(map[string]provider.DatacenterMeta),
	}

	for _, name := range m.SeedClusters {
		spec.Datacenters[name] = provider.DatacenterMeta{
			IsSeed: true,
			Spec: provider.DatacenterSpec{
				BringYourOwn: &provider.BringYourOwnSpec{},
			},
		}
	}

	for key, dc := range m.Datacenters {
		spec.Datacenters[key] = dc.KubermaticMeta()
	}

	return spec
}

type SecretsManifest struct {
	DockerAuth string `yaml:"dockerAuth"`
}

func (m *SecretsManifest) Validate() error {
	if len(m.DockerAuth) == 0 {
		return errors.New("no docker authentication specified")
	}

	var tmp interface{}

	if err := json.Unmarshal([]byte(m.DockerAuth), &tmp); err != nil {
		return fmt.Errorf("docker authentication is not valid JSON: %v", err)
	}

	return nil
}

type MonitoringManifest struct {
	Enabled bool `yaml:"enabled"`
}

func (m *MonitoringManifest) Validate() error {
	return nil
}

type LoggingManifest struct {
	Enabled       bool `yaml:"enabled"`
	RetentionDays int  `yaml:"retentionDays"`
}

func (m *LoggingManifest) Validate() error {
	if m.Enabled && m.RetentionDays <= 0 {
		return errors.New("retentionDays must be greater than zero")
	}

	return nil
}

type AuthenticationManifest struct {
	Google GoogleAuthenticationManifest `yaml:"google"`
}

func (m *AuthenticationManifest) Validate() error {
	if err := m.Google.Validate(); err != nil {
		return fmt.Errorf("invalid Google OAuth configuration: %v", err)
	}

	return nil
}

type GoogleAuthenticationManifest struct {
	ClientID  string `yaml:"clientID"`
	SecretKey string `yaml:"secretKey"`
}

func (m *GoogleAuthenticationManifest) Validate() error {
	if len(m.ClientID) == 0 {
		return errors.New("no client ID specified")
	}

	if len(m.SecretKey) == 0 {
		return errors.New("no secret key specified")
	}

	return nil
}

type SettingsManifest struct {
	URL string `yaml:"url"`
}

func (m *SettingsManifest) Validate() error {
	if len(m.URL) == 0 {
		return errors.New("no base URL specified")
	}

	return nil
}
