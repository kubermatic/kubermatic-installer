package installer

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/icza/dyno"
	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
	yaml "gopkg.in/yaml.v2"
)

func generateSecret() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

type KubermaticValues struct {
	data map[string]interface{}
}

func LoadValuesFromFile(filename string) (KubermaticValues, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return KubermaticValues{}, fmt.Errorf("failed to read %s: %v", filename, err)
	}

	parsed := make(map[string]interface{})

	err = yaml.Unmarshal(content, &parsed)
	if err != nil {
		return KubermaticValues{}, fmt.Errorf("failed to parse %s as YAML: %v", filename, err)
	}

	return KubermaticValues{parsed}, nil
}

func (v *KubermaticValues) ApplyManifest(m *manifest.Manifest) error {
	v.setKubeconfig(m.Kubeconfig)
	v.setDatacenters(m.KubermaticDatacenters())

	// configure domains
	baseDomain := m.Settings.BaseDomain
	baseURL := fmt.Sprintf("https://%s", baseDomain)
	domains := []string{baseDomain}

	if m.Monitoring.Enabled {
		domains = append(
			domains,
			fmt.Sprintf("prometheus.%s", baseDomain),
			fmt.Sprintf("grafana.%s", baseDomain),
			fmt.Sprintf("alertmanager.%s", baseDomain),
		)
	}

	//	if m.Logging.Enabled {
	//		domains = append(domains, fmt.Sprintf("kibana.%s", baseDomain))
	//	}

	v.set("kubermatic.domain", baseDomain)
	v.set("certificates.domains", domains)

	// configure Google authentication
	v.set("kubermatic.auth.tokenIssuer", fmt.Sprintf("%s/dex", baseURL))
	v.set("kubermatic.auth.clientID", "kubermatic")

	// disable LoadBalancer services on providers that do not support it
	if !m.SupportsLoadBalancers() {
		v.set("nginx.hostNetwork", true)
	}

	// configure controller
	v.set("kubermatic.controller.datacenterName", m.MasterDatacenterName())

	// configure Docker secrets
	v.configureDockerSecrets(m)

	// configure prometheus
	v.set("prometheus.externalLabels.seed_cluster", m.MasterDatacenterName())
	v.set("prometheus.host", fmt.Sprintf("prometheus.%s", baseDomain))

	// configure grafana
	v.set("grafana.host", fmt.Sprintf("grafana.%s", baseDomain))

	// configure alertmanager
	v.set("alertmanager.host", fmt.Sprintf("alertmanager.%s", baseDomain))

	// configure dex
	secrets := make(map[string]string)

	secret, err := generateSecret()
	if err != nil {
		return err
	}
	secrets["kubermatic"] = secret

	connector := NewGoogleDexConnector(m.Authentication.Google.ClientID, m.Authentication.Google.SecretKey, baseURL)
	dexClients := []DexClient{
		{
			ID:     "kubermatic",
			Name:   "Kubermatic",
			Secret: secrets["kubermatic"],
			RedirectURIs: []string{
				baseURL,
				fmt.Sprintf("%s/clusters", baseURL),
				fmt.Sprintf("%s/projects", baseURL),
			},
		},
	}

	if m.Monitoring.Enabled {
		for _, key := range []string{"prometheus", "grafana", "alertmanager"} {
			secret, err = generateSecret()
			if err != nil {
				return err
			}
			secrets[key] = secret

			dexClients = append(
				dexClients,
				DexClient{
					ID:           key,
					Name:         key,
					Secret:       secret,
					RedirectURIs: []string{fmt.Sprintf("https://%s.%s/oauth/callback", key, baseDomain)},
				},
			)
		}
	}

	v.set("dex.connectors", []DexConnector{connector})
	v.set("dex.clients", dexClients)
	v.set("dex.ingress.host", baseDomain)

	// configure IAP
	deployments := make(map[string]IAPDeployment)

	if m.Monitoring.Enabled {
		keys := make(map[string]string)

		for _, key := range []string{"prometheus", "grafana", "alertmanager"} {
			secret, err = generateSecret()
			if err != nil {
				return err
			}
			keys[key] = secret
		}

		deployments["grafana"] = IAPDeployment{
			Name:            "grafana",
			ClientID:        "grafana",
			ClientSecret:    secrets["grafana"],
			EncryptionKey:   keys["grafana"],
			UpstreamService: "grafana.monitoring.svc.cluster.local",
			UpstreamPort:    3000,
			Ingress: IAPDeploymentIngress{
				Host: fmt.Sprintf("grafana.%s", baseDomain),
			},
			Config: IAPDeploymentConfig{
				"enable-authorization-header": false,
				"scopes":                      []string{"groups"},
				"resources":                   NewNullIAPResource(),
			},
		}

		deployments["prometheus"] = IAPDeployment{
			Name:            "prometheus",
			ClientID:        "prometheus",
			ClientSecret:    secrets["prometheus"],
			EncryptionKey:   keys["prometheus"],
			UpstreamService: "prometheus-kubermatic.monitoring.svc.cluster.local",
			UpstreamPort:    9090,
			Ingress: IAPDeploymentIngress{
				Host: fmt.Sprintf("prometheus.%s", baseDomain),
			},
			Config: IAPDeploymentConfig{
				"scopes":    []string{"groups"},
				"resources": NewNullIAPResource(),
			},
		}

		deployments["alertmanager"] = IAPDeployment{
			Name:            "alertmanager",
			ClientID:        "alertmanager",
			ClientSecret:    secrets["alertmanager"],
			EncryptionKey:   keys["alertmanager"],
			UpstreamService: "alertmanager-kubermatic.monitoring.svc.cluster.local",
			UpstreamPort:    9093,
			Ingress: IAPDeploymentIngress{
				Host: fmt.Sprintf("alertmanager.%s", baseDomain),
			},
			Config: IAPDeploymentConfig{
				"scopes":    []string{"groups"},
				"resources": NewNullIAPResource(),
			},
		}
	}

	v.set("iap.deployments", deployments)
	v.set("iap.discovery_url", fmt.Sprintf("%s/dex/.well-known/openid-configuration", baseURL))
	v.set("iap.port", 3000)

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

	return nil
}

func (v *KubermaticValues) configureDockerSecrets(m *manifest.Manifest) {
	type dockerAuth struct {
		Auth  string `json:"auth"`
		EMail string `json:"email"`
	}

	type dockerConfig struct {
		Auths map[string]dockerAuth `json:"auths"`
	}

	cfg := dockerConfig{}
	json.Unmarshal([]byte(m.Secrets.DockerAuth), &cfg)

	secrets := map[string]dockerConfig{
		// the new kubermatic 2.8+ way
		"kubermatic.imagePullSecretData": cfg,
	}

	// go through the provided JSON and find the credentials
	// for docker.io and quay.io to split them into the two
	// seperate secrets that Kubermatic pre-2.8 require
	for registry, auth := range cfg.Auths {
		subcfg := dockerConfig{
			Auths: make(map[string]dockerAuth),
		}

		if strings.Contains(registry, "quay.io") {
			subcfg.Auths["quay.io"] = auth
			secrets["kubermatic.quay.secret"] = subcfg
		} else if strings.Contains(registry, "docker.io") {
			subcfg.Auths["https://index.docker.io/v1/"] = auth
			secrets["kubermatic.docker.secret"] = subcfg
		}
	}

	for path, val := range secrets {
		blob, _ := json.Marshal(val)
		v.set(path, base64.StdEncoding.EncodeToString(blob))
	}
}

func (v *KubermaticValues) setKubeconfig(kubeconfig string) {
	v.set("kubermatic.kubeconfig", base64.StdEncoding.EncodeToString([]byte(kubeconfig)))
}

func (v *KubermaticValues) setDatacenters(dcs *manifest.KubermaticDatacenters) {
	encoded, _ := yaml.Marshal(dcs)

	v.set("kubermatic.datacenters", base64.StdEncoding.EncodeToString(encoded))
}

func (v *KubermaticValues) YAML() []byte {
	encoded, _ := yaml.Marshal(v.data)

	return encoded
}

func (v *KubermaticValues) set(path string, value interface{}) {
	elements := make([]interface{}, 0)

	for _, e := range strings.Split(path, ".") {
		elements = append(elements, e)
	}

	dyno.Set(v.data, value, elements...)
}
