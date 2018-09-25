package installer

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/icza/dyno"
	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
	yaml "gopkg.in/yaml.v2"
)

type KubermaticValues map[string]interface{}

func LoadValuesFromFile(filename string) (KubermaticValues, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %v", filename, err)
	}

	var parsed KubermaticValues

	err = yaml.Unmarshal(content, &parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s as YAML: %v", filename, err)
	}

	return parsed, nil
}

func (v *KubermaticValues) ApplyManifest(m *manifest.Manifest) {
	v.SetKubeconfig(m.Kubeconfig)
	v.SetDatacenters(m.KubermaticDatacenters())

	domains := []string{m.Settings.URL}

	if m.Monitoring.Enabled {
		domains = append(
			domains,
			fmt.Sprintf("prometheus.%s", m.Settings.URL),
			fmt.Sprintf("grafana.%s", m.Settings.URL),
			fmt.Sprintf("alertmanager.%s", m.Settings.URL),
		)
	}

	if m.Logging.Enabled {
		domains = append(domains, fmt.Sprintf("kibana.%s", m.Settings.URL))
	}

	v.set("kubermatic.domain", m.Settings.URL)
	v.set("certificates.domains", domains)
}

func (v *KubermaticValues) SetKubeconfig(kubeconfig string) {
	v.set("kubermatic.kubeconfig", base64.StdEncoding.EncodeToString([]byte(kubeconfig)))
}

func (v *KubermaticValues) SetDatacenters(dcs *manifest.KubermaticDatacenters) {
	encoded, _ := yaml.Marshal(dcs)

	v.set("kubermatic.datacenters", base64.StdEncoding.EncodeToString(encoded))
}

func (v *KubermaticValues) set(path string, value interface{}) {
	elements := make([]interface{}, 0)

	for _, e := range strings.Split(path, ".") {
		elements = append(elements, e)
	}

	dyno.Set(v, value, elements...)
}
