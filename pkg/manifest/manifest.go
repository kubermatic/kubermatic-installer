package manifest

import (
	"encoding/json"
	"errors"
	"fmt"
)

// VERSION is the manifest version, used to distinguish
// different versions during the installer lifetime.
const VERSION = "1"

// Manifest is a description of all options a customer
// can tweak for their Kubermatic installation. It
// contains options grouped by use-case which are then
// used to customize the various Helm charts during
// installation.
type Manifest struct {
	Version        string                        `yaml:"version"`
	Kubeconfig     string                        `yaml:"kubeconfig"`
	CloudProvider  CloudProvider                 `yaml:"cloudProvider"`
	Secrets        SecretsManifest               `yaml:"secrets"`
	SeedClusters   []string                      `yaml:"seedClusters"`
	Datacenters    map[string]DatacenterManifest `yaml:"datacenters"`
	Monitoring     MonitoringManifest            `yaml:"monitoring"`
	Logging        LoggingManifest               `yaml:"logging"`
	Authentication AuthenticationManifest        `yaml:"authentication"`
	Settings       SettingsManifest              `yaml:"settings"`

	// values determined during installation which at
	// some point might be configured explicitely in
	// the manifest
	MinioStorageClass string `yaml:"-"`
}

// Validate checks the manifest for semantical correctness
// and aborts at the first error.
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

// SupportsLoadBalancers determines whether the Kubernetes
// cluster is running at a provider that has support for
// creating load balancers. Right now, as we only support
// Google/Amazon, this always returns true, but in the
// future this needs to make a decision based on a user-
// selected cloud provider.
func (m *Manifest) SupportsLoadBalancers() bool {
	return true
}

// ServiceDomain returns the full domain for one of the
// services provided by Kubermatic, e.g. Prometheus or
// Grafana.
func (m *Manifest) ServiceDomain(service string) string {
	return fmt.Sprintf("%s.%s", service, m.Settings.BaseDomain)
}

// BaseURL returns the full URL, including protocol, for
// the Kubermatic dashboard. The URL includes a trailing
// slash.
func (m *Manifest) BaseURL() string {
	return fmt.Sprintf("https://%s/", m.Settings.BaseDomain)
}

// KubermaticDatacenters describes the datacenters that
// Kubermatic allows to create clusters in and where seed
// installations are running. It mimics the structure of
// Kubermatic's pkg/provider package.
type KubermaticDatacenters struct {
	Datacenters map[string]DatacenterMeta `yaml:"datacenters"`
}

// KubermaticDatacenters transforms the manifest's version
// of the datacenters into the structure that Kubermatic
// expects, so it can easily be marshalled into YAML.
func (m *Manifest) KubermaticDatacenters() *KubermaticDatacenters {
	spec := &KubermaticDatacenters{
		Datacenters: make(map[string]DatacenterMeta),
	}

	for _, name := range m.SeedClusters {
		spec.Datacenters[name] = DatacenterMeta{
			IsSeed: true,
			Spec: DatacenterSpec{
				BringYourOwn: &BringYourOwnSpec{},
			},
		}
	}

	for key, dc := range m.Datacenters {
		spec.Datacenters[key] = dc.KubermaticMeta()
	}

	return spec
}

// MasterDatacenterName returns the name of the master datacenter
// for the Kubermatic installation. It assumes that the Validate()
// function has been run and that there is at least one datacenter
// defined, because it returns the first datacenter's name.
func (m *Manifest) MasterDatacenterName() string {
	return m.SeedClusters[0]
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
	Google *GoogleAuthenticationManifest `yaml:"google"`
	GitHub *GitHubAuthenticationManifest `yaml:"github"`
}

func (m *AuthenticationManifest) Validate() error {
	if m.Google != nil && m.Google.ClientID != "" {
		if err := m.Google.Validate(); err != nil {
			return fmt.Errorf("invalid Google OAuth configuration: %v", err)
		}
	}

	if m.GitHub != nil && m.GitHub.ClientID != "" {
		if err := m.GitHub.Validate(); err != nil {
			return fmt.Errorf("invalid GitHub OAuth configuration: %v", err)
		}
	}

	if m.GitHub == nil && m.Google == nil {
		return errors.New("must configure at least one authentication provider")
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

type GitHubAuthenticationManifest struct {
	ClientID     string `yaml:"clientID"`
	SecretKey    string `yaml:"secretKey"`
	Organization string `yaml:"organization"`
}

func (m *GitHubAuthenticationManifest) Validate() error {
	if len(m.ClientID) == 0 {
		return errors.New("no client ID specified")
	}

	if len(m.SecretKey) == 0 {
		return errors.New("no secret key specified")
	}

	return nil
}

type SettingsManifest struct {
	BaseDomain string `yaml:"baseDomain"`
}

func (m *SettingsManifest) Validate() error {
	if len(m.BaseDomain) == 0 {
		return errors.New("no base domain specified")
	}

	return nil
}
