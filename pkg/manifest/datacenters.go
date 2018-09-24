package manifest

import (
	"errors"
	"fmt"

	"k8s.io/apimachinery/pkg/util/sets"
)

var (
	AllOperatingSystems = sets.NewString(string(OperatingSystemCoreos), string(OperatingSystemCentOS), string(OperatingSystemUbuntu))
)

type DatacenterManifest struct {
	Location string                 `yaml:"location"`
	Country  string                 `yaml:"country"`
	Seed     string                 `yaml:"seed"`
	Spec     DatacenterSpecManifest `yaml:"spec"`
}

func (m *DatacenterManifest) KubermaticMeta() DatacenterMeta {
	return DatacenterMeta{
		Location: m.Location,
		Country:  m.Country,
		Seed:     m.Seed,
		Spec:     m.Spec.KubermaticSpec(),
	}
}

func (m *DatacenterManifest) Validate(seeds []string) error {
	if len(m.Location) == 0 {
		return errors.New("no location specified")
	}

	if len(m.Country) != 2 {
		return errors.New("no or invalid country specified")
	}

	if seeds != nil {
		validSeed := false

		for _, seed := range seeds {
			if seed == m.Seed {
				validSeed = true
				break
			}
		}

		if !validSeed {
			return fmt.Errorf("invalid seed cluster '%s' specified", m.Seed)
		}
	}

	if err := m.Spec.Validate(); err != nil {
		return fmt.Errorf("invalid cloud spec: %v", err)
	}

	return nil
}

type DatacenterSpecManifest struct {
	AWS          *DatacenterAWSManifest          `yaml:"aws"`
	DigitalOcean *DatacenterDigitalOceanManifest `yaml:"digitalocean"`
	Hetzner      *DatacenterHetznerManifest      `yaml:"hetzner"`
	Azure        *DatacenterAzureManifest        `yaml:"azure"`
	VSphere      *DatacenterVSphereManifest      `yaml:"vsphere"`
	OpenStack    *DatacenterOpenStackManifest    `yaml:"openstack"`
	BringYourOwn *DatacenterBringYourOwnManifest `yaml:"bringyourown"`
}

func (m *DatacenterSpecManifest) Validate() error {
	type validatable interface {
		Validate() error
	}

	specs := make([]validatable, 0)

	if m.AWS != nil {
		specs = append(specs, m.AWS)
	}

	if m.DigitalOcean != nil {
		specs = append(specs, m.DigitalOcean)
	}

	if m.Hetzner != nil {
		specs = append(specs, m.Hetzner)
	}

	if m.Azure != nil {
		specs = append(specs, m.Azure)
	}

	if m.VSphere != nil {
		specs = append(specs, m.VSphere)
	}

	if m.OpenStack != nil {
		specs = append(specs, m.OpenStack)
	}

	if m.BringYourOwn != nil {
		specs = append(specs, m.BringYourOwn)
	}

	if len(specs) == 0 {
		return errors.New("no spec configured")
	}

	if len(specs) > 1 {
		return errors.New("more than one spec configured")
	}

	if err := specs[0].Validate(); err != nil {
		return fmt.Errorf("invalid spec: %v", err)
	}

	return nil
}

func (m *DatacenterSpecManifest) KubermaticSpec() DatacenterSpec {
	spec := DatacenterSpec{}

	if m.AWS != nil {
		spec.AWS = m.AWS.KubermaticSpec()
	}

	if m.DigitalOcean != nil {
		spec.Digitalocean = m.DigitalOcean.KubermaticSpec()
	}

	if m.Hetzner != nil {
		spec.Hetzner = m.Hetzner.KubermaticSpec()
	}

	if m.Azure != nil {
		spec.Azure = m.Azure.KubermaticSpec()
	}

	if m.VSphere != nil {
		spec.VSphere = m.VSphere.KubermaticSpec()
	}

	if m.OpenStack != nil {
		spec.Openstack = m.OpenStack.KubermaticSpec()
	}

	if m.BringYourOwn != nil {
		spec.BringYourOwn = m.BringYourOwn.KubermaticSpec()
	}

	return spec
}

type DatacenterAWSManifest struct {
	Region        string `yaml:"region"`
	AMI           string `yaml:"ami"`
	ZoneCharacter string `yaml:"zoneCharacter"`
}

func (m *DatacenterAWSManifest) Validate() error {
	if len(m.Region) == 0 {
		return errors.New("no region specified")
	}

	if len(m.AMI) == 0 {
		return errors.New("no AMI specified")
	}

	if len(m.ZoneCharacter) == 0 {
		return errors.New("no zone character specified")
	}

	return nil
}

func (m *DatacenterAWSManifest) KubermaticSpec() *AWSSpec {
	return &AWSSpec{
		Region:        m.Region,
		AMI:           m.AMI,
		ZoneCharacter: m.ZoneCharacter,
	}
}

type DatacenterDigitalOceanManifest struct {
	Region string `yaml:"region"`
}

func (m *DatacenterDigitalOceanManifest) Validate() error {
	if len(m.Region) == 0 {
		return errors.New("no region specified")
	}

	return nil
}

func (m *DatacenterDigitalOceanManifest) KubermaticSpec() *DigitaloceanSpec {
	return &DigitaloceanSpec{
		Region: m.Region,
	}
}

type DatacenterHetznerManifest struct {
	Datacenter string `yaml:"datacenter"`
	Location   string `yaml:"location"`
}

func (m *DatacenterHetznerManifest) Validate() error {
	if len(m.Datacenter) == 0 {
		return errors.New("no datacenter specified")
	}

	return nil
}

func (m *DatacenterHetznerManifest) KubermaticSpec() *HetznerSpec {
	return &HetznerSpec{
		Datacenter: m.Datacenter,
		Location:   m.Location,
	}
}

type DatacenterAzureManifest struct {
	Location string `yaml:"location"`
}

func (m *DatacenterAzureManifest) Validate() error {
	if len(m.Location) == 0 {
		return errors.New("no location specified")
	}
	return nil
}

func (m *DatacenterAzureManifest) KubermaticSpec() *AzureSpec {
	return &AzureSpec{
		Location: m.Location,
	}
}

type DatacenterVSphereManifest struct {
	Endpoint      string `yaml:"endpoint"`
	AllowInsecure bool   `yaml:"allowInsecure"`

	Datastore  string    `yaml:"datastore"`
	Datacenter string    `yaml:"datacenter"`
	Cluster    string    `yaml:"cluster"`
	RootPath   string    `yaml:"rootPath"`
	Templates  ImageList `yaml:"templates"`

	// Infra management user is an optional user that will be used only
	// for everything except the cloud provider functionality which will
	// still use the credentials passed in via the frontend/api
	InfraManagementUser *VSphereCredentials `yaml:"infraManagementUser"`
}

func (m *DatacenterVSphereManifest) Validate() error {
	for image := range m.Templates {
		if !AllOperatingSystems.Has(string(image)) {
			return fmt.Errorf("template for unknown operating system '%s' specified", image)
		}
	}

	return nil
}

func (m *DatacenterVSphereManifest) KubermaticSpec() *VSphereSpec {
	return &VSphereSpec{
		Endpoint:      m.Endpoint,
		AllowInsecure: m.AllowInsecure,
		Datastore:     m.Datastore,
		Datacenter:    m.Datacenter,
		Cluster:       m.Cluster,
		RootPath:      m.RootPath,
		Templates:     m.Templates,
		InfraManagementUser: &VSphereCredentials{
			Username: m.InfraManagementUser.Username,
			Password: m.InfraManagementUser.Password,
		},
	}
}

type DatacenterOpenStackManifest struct {
	AuthURL          string    `yaml:"authUrl"`
	AvailabilityZone string    `yaml:"availabilityZone"`
	Region           string    `yaml:"region"`
	IgnoreVolumeAZ   bool      `yaml:"ignoreVolumeAZ"`
	DNSServers       []string  `yaml:"dnsServers"`
	Images           ImageList `yaml:"images"`
}

func (m *DatacenterOpenStackManifest) Validate() error {
	for image := range m.Images {
		if !AllOperatingSystems.Has(string(image)) {
			return fmt.Errorf("image for unknown operating system '%s' specified", image)
		}
	}

	return nil
}

func (m *DatacenterOpenStackManifest) KubermaticSpec() *OpenstackSpec {
	return &OpenstackSpec{
		AuthURL:          m.AuthURL,
		AvailabilityZone: m.AvailabilityZone,
		Region:           m.Region,
		IgnoreVolumeAZ:   m.IgnoreVolumeAZ,
		DNSServers:       m.DNSServers,
		Images:           m.Images,
	}
}

type DatacenterBringYourOwnManifest struct {
}

func (m *DatacenterBringYourOwnManifest) Validate() error {
	return nil
}

func (m *DatacenterBringYourOwnManifest) KubermaticSpec() *BringYourOwnSpec {
	return &BringYourOwnSpec{}
}
