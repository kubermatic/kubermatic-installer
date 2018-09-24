package manifest

// DatacenterMeta describes a Kubermatic datacenter.
type DatacenterMeta struct {
	Location         string         `yaml:"location"`
	Seed             string         `yaml:"seed"`
	Country          string         `yaml:"country"`
	Spec             DatacenterSpec `yaml:"spec"`
	Private          bool           `yaml:"private"`
	IsSeed           bool           `yaml:"is_seed"`
	SeedDNSOverwrite *string        `yaml:"seed_dns_overwrite,omitempty"`
}

// DatacenterSpec descriped kubermatic-installer datacernter
type DatacenterSpec struct {
	Digitalocean *DigitaloceanSpec `yaml:"digitalocean,omitempty"`
	BringYourOwn *BringYourOwnSpec `yaml:"bringyourown,omitempty"`
	AWS          *AWSSpec          `yaml:"aws,omitempty"`
	Azure        *AzureSpec        `yaml:"azure,omitempty"`
	Openstack    *OpenstackSpec    `yaml:"openstack,omitempty"`
	Hetzner      *HetznerSpec      `yaml:"hetzner,omitempty"`
	VSphere      *VSphereSpec      `yaml:"vsphere,omitempty"`
}

// AWSSpec describes a aws datacenter
type AWSSpec struct {
	Region        string `yaml:"region"`
	AMI           string `yaml:"ami"`
	ZoneCharacter string `yaml:"zone_character"`
}

// DigitaloceanSpec describes a DigitalOcean datacenter
type DigitaloceanSpec struct {
	Region string `yaml:"region"`
}

// BringYourOwnSpec describes a datacenter our of bring your own nodes
type BringYourOwnSpec struct{}

// AzureSpec describes an Azure cloud datacenter
type AzureSpec struct {
	Location string `yaml:"location"`
}

// OpenstackSpec describes a open stack datacenter
type OpenstackSpec struct {
	AuthURL          string `yaml:"auth_url"`
	AvailabilityZone string `yaml:"availability_zone"`
	Region           string `yaml:"region"`
	IgnoreVolumeAZ   bool   `yaml:"ignore_volume_az"`
	// Used for automatic network creation
	DNSServers []string  `yaml:"dns_servers"`
	Images     ImageList `yaml:"images"`
}

// HetznerSpec describes a Hetzner cloud datacenter
type HetznerSpec struct {
	Datacenter string `yaml:"datacenter"`
	Location   string `yaml:"location"`
}

// VSphereSpec describes a vsphere datacenter
type VSphereSpec struct {
	Endpoint      string `yaml:"endpoint"`
	AllowInsecure bool   `yaml:"allow_insecure"`

	Datastore  string    `yaml:"datastore"`
	Datacenter string    `yaml:"datacenter"`
	Cluster    string    `yaml:"cluster"`
	RootPath   string    `yaml:"root_path"`
	Templates  ImageList `yaml:"templates"`

	// Infra management user is an optional user that will be used only
	// for everything except the cloud provider functionality which will
	// still use the credentials passed in via the frontend/api
	InfraManagementUser *VSphereCredentials `yaml:"infra_management_user,omitempty"`
}

// VSphereCredentials describes the credentials used
// as the infra management user
type VSphereCredentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// ImageList defines a map of operating system and the image to use
type ImageList map[OperatingSystem]string

type OperatingSystem string

const (
	OperatingSystemCoreos OperatingSystem = "coreos"
	OperatingSystemUbuntu OperatingSystem = "ubuntu"
	OperatingSystemCentOS OperatingSystem = "centos"
)
