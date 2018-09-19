package shared

type Manifest struct {
	Version      string
	AdvancedMode bool

	Secrets        SecretsManifest
	Datacenters    []DatacenterManifest
	Monitoring     MonitoringManifest
	Logging        LoggingManifest
	Authentication AuthenticationManifest
	Settings       SettingsManifest
}

func (m *Manifest) Validate() error {
	if err := m.Secrets.Validate(); err != nil {
		return err
	}

	for _, dc := range m.Datacenters {
		if err := dc.Validate(); err != nil {
			return err
		}
	}

	if err := m.Monitoring.Validate(); err != nil {
		return err
	}

	if err := m.Logging.Validate(); err != nil {
		return err
	}

	if err := m.Authentication.Validate(); err != nil {
		return err
	}

	if err := m.Settings.Validate(); err != nil {
		return err
	}

	return nil
}

type SecretsManifest struct {
	DockerHub string
	Quay      string
}

func (m *SecretsManifest) Validate() error {
	return nil
}

type DatacenterManifest struct {
	Location      string
	Country       string
	Region        string
	AMI           string
	ZoneCharacter string
}

func (m *DatacenterManifest) Validate() error {
	return nil
}

type MonitoringManifest struct {
	Enabled bool
}

func (m *MonitoringManifest) Validate() error {
	return nil
}

type LoggingManifest struct {
	Enabled       bool
	RetentionDays int
}

func (m *LoggingManifest) Validate() error {
	return nil
}

type AuthenticationManifest struct {
	Google GoogleAuthenticationManifest
}

func (m *AuthenticationManifest) Validate() error {
	return nil
}

type GoogleAuthenticationManifest struct {
	ClientID  string
	SecretKey string
}

func (m *GoogleAuthenticationManifest) Validate() error {
	return nil
}

type SettingsManifest struct {
	URL string
}

func (m *SettingsManifest) Validate() error {
	return nil
}
