package helm

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/icza/dyno"
	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
	yaml "gopkg.in/yaml.v2"
)

func base64encode(s string) string {
	return base64.StdEncoding.EncodeToString(bytes.Trim([]byte(s), "\n"))
}

func generateSecret() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

// Values represents the values.yaml for Helm plus some
// temporary runtime data.
type Values struct {
	data map[string]interface{}

	domains map[string]string
	secrets map[string]string
	baseURL string
}

func newValues(data map[string]interface{}) *Values {
	return &Values{
		data:    data,
		domains: make(map[string]string),
		secrets: make(map[string]string),
	}
}

func NewValues() *Values {
	return newValues(make(map[string]interface{}))
}

func LoadValuesFromFile(filename string) (*Values, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %v", filename, err)
	}

	parsed := make(map[string]interface{})

	err = yaml.Unmarshal(content, &parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s as YAML: %v", filename, err)
	}

	return newValues(parsed), nil
}

func (v *Values) ApplyManifest(m *manifest.Manifest) error {
	v.setKubeconfig(m.Kubeconfig)
	v.setDatacenters(m.KubermaticDatacenters())

	// configure domains
	v.configureDomains(m)

	// configure Google authentication
	v.set("kubermatic.auth.tokenIssuer", fmt.Sprintf("%s/dex", v.baseURL))
	v.set("kubermatic.auth.clientID", "kubermatic")

	// disable LoadBalancer services on providers that do not support it
	if !m.SupportsLoadBalancers() {
		v.set("nginx.hostNetwork", true)
	}

	// configure controller
	v.set("kubermatic.controller.datacenterName", m.MasterDatacenterName())

	// configure Docker secrets
	v.configureDockerSecrets(m)

	if m.Monitoring.Enabled {
		// configure prometheus
		v.set("prometheus.externalLabels.seed_cluster", m.MasterDatacenterName())
		v.set("prometheus.host", v.domains["prometheus"])

		// configure grafana
		v.set("grafana.host", v.domains["grafana"])

		// configure alertmanager
		v.set("alertmanager.host", v.domains["alertmanager"])
	}

	if m.Logging.Enabled {
		v.set("logging.elasticsearch.curator.interval", m.Logging.RetentionDays)
	}

	// configure dex
	if err := v.configureDex(m); err != nil {
		return err
	}

	// configure IAP
	if err := v.configureIAP(m); err != nil {
		return err
	}

	// configure minio
	minioAccessKey, err := generateSecret()
	if err != nil {
		return err
	}

	minioSecretKey, err := generateSecret()
	if err != nil {
		return err
	}

	v.set("minio.credentials.accessKey", minioAccessKey)
	v.set("minio.credentials.secretKey", minioSecretKey)

	if m.MinioStorageClass != "" {
		v.set("minio.storageClass", m.MinioStorageClass)
	}

	return nil
}

func (v *Values) configureDomains(m *manifest.Manifest) {
	// configure domains
	v.domains[""] = m.Settings.BaseDomain
	v.baseURL = fmt.Sprintf("https://%s", v.domains[""])

	if m.Monitoring.Enabled {
		v.domains["prometheus"] = m.ServiceDomain("prometheus")
		v.domains["grafana"] = m.ServiceDomain("grafana")
		v.domains["alertmanager"] = m.ServiceDomain("alertmanager")
	}

	if m.Logging.Enabled {
		v.domains["kibana"] = m.ServiceDomain("kibana")
	}

	domains := make([]string, 0)
	for _, domain := range v.domains {
		domains = append(domains, domain)
	}

	v.set("kubermatic.domain", v.domains[""])
	v.set("certificates.domains", domains)
}

func (v *Values) configureDex(m *manifest.Manifest) error {
	secret, err := generateSecret()
	if err != nil {
		return err
	}
	v.secrets["kubermatic"] = secret

	dexClients := []DexClient{
		{
			ID:     "kubermatic",
			Name:   "Kubermatic",
			Secret: v.secrets["kubermatic"],
			RedirectURIs: []string{
				v.baseURL,
				fmt.Sprintf("%s/clusters", v.baseURL),
				fmt.Sprintf("%s/projects", v.baseURL),
			},
		},
	}

	if m.Monitoring.Enabled {
		for _, key := range []string{"prometheus", "grafana", "alertmanager"} {
			secret, err = generateSecret()
			if err != nil {
				return err
			}
			v.secrets[key] = secret

			dexClients = append(
				dexClients,
				DexClient{
					ID:           key,
					Name:         key,
					Secret:       secret,
					RedirectURIs: []string{fmt.Sprintf("https://%s/oauth/callback", v.domains[key])},
				},
			)
		}
	}

	if m.Logging.Enabled {
		secret, err = generateSecret()
		if err != nil {
			return err
		}
		v.secrets["kibana"] = secret

		dexClients = append(
			dexClients,
			DexClient{
				ID:           "kibana",
				Name:         "kibana",
				Secret:       secret,
				RedirectURIs: []string{fmt.Sprintf("https://%s/oauth/callback", v.domains["kibana"])},
			},
		)
	}

	connectors := []DexConnector{}

	if m.Authentication.Google.ClientID != "" {
		connectors = append(connectors, NewGoogleDexConnector(m.Authentication.Google.ClientID, m.Authentication.Google.SecretKey, v.baseURL))
	}

	if m.Authentication.GitHub.ClientID != "" {
		connectors = append(connectors, NewGitHubDexConnector(m.Authentication.GitHub.ClientID, m.Authentication.GitHub.SecretKey, v.baseURL, m.Authentication.GitHub.Organization))
	}

	v.set("dex.connectors", connectors)
	v.set("dex.clients", dexClients)
	v.set("dex.ingress.host", v.domains[""])

	return nil
}

func (v *Values) configureIAP(m *manifest.Manifest) error {
	deployments := make(map[string]IAPDeployment)

	if m.Monitoring.Enabled {
		keys := make(map[string]string)

		for _, key := range []string{"prometheus", "grafana", "alertmanager"} {
			secret, err := generateSecret()
			if err != nil {
				return err
			}
			keys[key] = secret
		}

		resources := []IAPResource{NewNullIAPResource()}

		deployments["grafana"] = IAPDeployment{
			Name:            "grafana",
			ClientID:        "grafana",
			ClientSecret:    v.secrets["grafana"],
			EncryptionKey:   keys["grafana"],
			UpstreamService: "grafana.monitoring.svc.cluster.local",
			UpstreamPort:    3000,
			Ingress: IAPDeploymentIngress{
				Host: v.domains["grafana"],
			},
			Config: IAPDeploymentConfig{
				"enable-authorization-header": false,
				"scopes":                      []string{"groups"},
				"resources":                   resources,
			},
		}

		deployments["prometheus"] = IAPDeployment{
			Name:            "prometheus",
			ClientID:        "prometheus",
			ClientSecret:    v.secrets["prometheus"],
			EncryptionKey:   keys["prometheus"],
			UpstreamService: "prometheus-kubermatic.monitoring.svc.cluster.local",
			UpstreamPort:    9090,
			Ingress: IAPDeploymentIngress{
				Host: v.domains["prometheus"],
			},
			Config: IAPDeploymentConfig{
				"scopes":    []string{"groups"},
				"resources": resources,
			},
		}

		deployments["alertmanager"] = IAPDeployment{
			Name:            "alertmanager",
			ClientID:        "alertmanager",
			ClientSecret:    v.secrets["alertmanager"],
			EncryptionKey:   keys["alertmanager"],
			UpstreamService: "alertmanager-kubermatic.monitoring.svc.cluster.local",
			UpstreamPort:    9093,
			Ingress: IAPDeploymentIngress{
				Host: v.domains["alertmanager"],
			},
			Config: IAPDeploymentConfig{
				"scopes":    []string{"groups"},
				"resources": resources,
			},
		}
	}

	if m.Logging.Enabled {
		key, err := generateSecret()
		if err != nil {
			return err
		}

		deployments["kibana"] = IAPDeployment{
			Name:            "kibana",
			ClientID:        "kibana",
			ClientSecret:    v.secrets["kibana"],
			EncryptionKey:   key,
			UpstreamService: "kibana-logging.logging.svc.cluster.local",
			UpstreamPort:    5601,
			Ingress: IAPDeploymentIngress{
				Host: v.domains["kibana"],
			},
			Config: IAPDeploymentConfig{
				"scopes":    []string{"groups"},
				"resources": []IAPResource{NewNullIAPResource()},
			},
		}
	}

	if len(deployments) > 0 {
		v.set("iap.deployments", deployments)
		v.set("iap.discovery_url", fmt.Sprintf("%s/dex/.well-known/openid-configuration", v.baseURL))
		v.set("iap.port", 3000)
	}

	return nil
}

func (v *Values) configureDockerSecrets(m *manifest.Manifest) {
	v.set("kubermatic.imagePullSecretData", base64encode(m.Secrets.DockerAuth))
}

func (v *Values) setKubeconfig(kubeconfig string) {
	v.set("kubermatic.kubeconfig", base64encode(kubeconfig))
}

func (v *Values) setDatacenters(dcs *manifest.KubermaticDatacenters) {
	encoded, _ := yaml.Marshal(dcs)

	v.set("kubermatic.datacenters", base64encode(string(encoded)))
}

func (v *Values) YAML() []byte {
	encoded, _ := yaml.Marshal(v.data)

	return encoded
}

func (v *Values) set(path string, value interface{}) {
	elements := make([]interface{}, 0)

	for _, e := range strings.Split(path, ".") {
		elements = append(elements, e)
	}

	dyno.Set(v.data, value, elements...)
}
