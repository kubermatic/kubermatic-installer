package v2_10

import (
	"encoding/json"
	"fmt"

	"github.com/kubermatic/kubermatic-installer/pkg/helm/migration/util"
	"github.com/kubermatic/kubermatic-installer/pkg/yamled"
	"github.com/sirupsen/logrus"
)

const (
	alertmanagerVersion               = "v0.16.2"
	busyboxVersion                    = "1.30.1"
	certManagerVersion                = "v0.7.0"
	curatorVersion                    = "5.6.0-1"
	dexVersion                        = "v2.15.0"
	elasticsearchVersion              = "6.6.2"
	fluentbitVersion                  = "1.0.6"
	grafanaVersion                    = "6.1.3"
	kibanaVersion                     = "6.6.2"
	kubermaticAddonsVersion           = "v0.2.8"
	kubermaticAPIVersion              = "v2.10.0"
	kubermaticMasterControllerVersion = "v2.10.0"
	kubermaticUIVersion               = "v1.2.0"
	minioVersion                      = "RELEASE.2019-04-09T01-22-30Z"
	nginxVersion                      = "0.24.1"
	nodePortProxyVersion              = "v2.0.0"
	prometheusVersion                 = "v2.8.1"
	veleroVersion                     = "v0.11.0"
	vpaVersion                        = "0.5.0"
)

type toleration struct {
	Key      string `yaml:"foo,omitempty"`
	Operator string `yaml:"operator,omitempty"`
	Value    string `yaml:"value,omitempty"`
	Effect   string `yaml:"effect,omitempty"`
}

type converter struct {
	logger logrus.FieldLogger
}

func NewConverter(logger logrus.FieldLogger) *converter {
	return &converter{
		logger: logger,
	}
}

func (c *converter) Convert(doc *yamled.Document, isMaster bool) error {
	if err := c.updateKubermaticController(doc); err != nil {
		return fmt.Errorf("failed to update Kubermatic controller: %v", err)
	}

	if err := c.updateKubermaticMasterController(doc); err != nil {
		return fmt.Errorf("failed to update Kubermatic master controller: %v", err)
	}

	if err := c.updateKubermaticUIImage(doc); err != nil {
		return fmt.Errorf("failed to update Kubermatic UI image: %v", err)
	}

	if err := c.updateKubermaticUIConfig(doc); err != nil {
		return fmt.Errorf("failed to update Kubermatic UI configuration: %v", err)
	}

	if err := c.updateVPA(doc); err != nil {
		return fmt.Errorf("failed to update VPA: %v", err)
	}

	if err := c.updateCertManager(doc); err != nil {
		return fmt.Errorf("failed to update cert-manager: %v", err)
	}

	if err := c.updateNginx(doc); err != nil {
		return fmt.Errorf("failed to update nginx-ingress: %v", err)
	}

	if err := c.updateNodePortProxy(doc); err != nil {
		return fmt.Errorf("failed to update node port proxy: %v", err)
	}

	if err := c.updateDex(doc); err != nil {
		return fmt.Errorf("failed to update Dex: %v", err)
	}

	if err := c.updateMinio(doc); err != nil {
		return fmt.Errorf("failed to update Minio: %v", err)
	}

	if err := c.updateAlertmanager(doc); err != nil {
		return fmt.Errorf("failed to update Alertmanager: %v", err)
	}

	if err := c.updateGrafana(doc); err != nil {
		return fmt.Errorf("failed to update Grafana: %v", err)
	}

	if err := c.updatePrometheus(doc); err != nil {
		return fmt.Errorf("failed to update Prometheus: %v", err)
	}

	if err := c.updateElasticsearch(doc); err != nil {
		return fmt.Errorf("failed to update Elasticsearch: %v", err)
	}

	if err := c.updateKibana(doc); err != nil {
		return fmt.Errorf("failed to update Kibana: %v", err)
	}

	if err := c.updateFluentbit(doc); err != nil {
		return fmt.Errorf("failed to update fluentbit: %v", err)
	}

	if err := c.updateVelero(doc); err != nil {
		return fmt.Errorf("failed to update Velero: %v", err)
	}

	return nil
}

func (c *converter) updateKubermaticController(doc *yamled.Document) error {
	path := yamled.Path{"kubermatic", "deployVPA"}

	if doc.Has(path) {
		c.logger.Info("Removing kubermatic.redundant deployVPA flag in favor of kubermatic.controller.featureFlags...")
		doc.Remove(path)
	}

	addonNode, exists := doc.Get(yamled.Path{"kubermatic", "controller", "addons"})
	if exists {
		c.logger.Infof("Moving kubermatic.defaultAddons to new structure...")
		doc.Set(yamled.Path{"kubermatic", "controller", "addons", "kubernetes"}, addonNode)

		addonPath := yamled.Path{"kubermatic", "controller", "addons", "kubernetes", "defaultAddons"}
		addons, ok := doc.GetArray(addonPath)
		if ok {
			hasNodeExporter := false

			for _, addon := range addons {
				if addon.(string) == "node-exporter" {
					hasNodeExporter = true
					break
				}
			}

			if !hasNodeExporter {
				c.logger.Infof("Adding new node-exporter default addon for Kubernetes clusters...")
				addons = append(addons, "node-exporter")
				doc.Set(addonPath, addons)
			}
		}

		if updateDockerImage(doc, yamled.Path{"kubermatic", "controller", "addons", "kubernetes"}, kubermaticAddonsVersion) {
			c.logger.Infof("Updated docker image for Kubermatic addons for Kubernetes clusters to %s.", kubermaticAddonsVersion)
		}
	}

	updated := updateDockerImage(doc, yamled.Path{"kubermatic", "controller"}, kubermaticAPIVersion) ||
		updateDockerImage(doc, yamled.Path{"kubermatic", "api"}, kubermaticAPIVersion)

	if updated {
		c.logger.Infof("Updated Kubermatic API version to %s.", kubermaticAPIVersion)
	}

	return nil
}

func (c *converter) updateKubermaticUIImage(doc *yamled.Document) error {
	if updateDockerImage(doc, yamled.Path{"kubermatic", "ui"}, kubermaticUIVersion) {
		c.logger.Info("Updated Kubermatic UI version.")
	}

	return nil
}

func (c *converter) updateKubermaticUIConfig(doc *yamled.Document) error {
	path := yamled.Path{"kubermatic", "ui", "config"}

	config, ok := doc.GetString(path)
	if !ok {
		return nil
	}

	cfg := make(map[string]interface{})
	err := json.Unmarshal([]byte(config), &cfg)
	if err != nil {
		return fmt.Errorf("failed to decode config JSON: %v", err)
	}

	if _, exists := cfg["default_node_count"]; !exists {
		c.logger.Info("Adding new Kubermatic API config flag default_node_count=3.")
		cfg["default_node_count"] = 3
	}

	if _, exists := cfg["custom_links"]; !exists {
		c.logger.Info("Adding new Kubermatic API config flag custom_links=[].")
		cfg["custom_links"] = []string{}
	}

	marshalled, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to re-encode config JSON: %v", err)
	}

	doc.Set(path, string(marshalled))
	c.logger.Info("Updated Kubermatic UI configuration flags.")

	return nil
}

func (c *converter) updateKubermaticMasterController(doc *yamled.Document) error {
	config, exists := doc.Get(yamled.Path{"kubermatic", "rbac"})
	if exists {
		c.logger.Infof("Moving kubermatic.rbac to kubermatic.masterController...")
		doc.Set(yamled.Path{"kubermatic", "masterController"}, config)
	}

	if updateDockerImage(doc, yamled.Path{"kubermatic", "masterController"}, kubermaticMasterControllerVersion) {
		c.logger.Info("Updated Kubermatic master controller version.")
	}

	return nil
}

func (c *converter) updateVPA(doc *yamled.Document) error {
	updated := updateDockerImage(doc, yamled.Path{"kubermatic", "vpa", "updater"}, vpaVersion) ||
		updateDockerImage(doc, yamled.Path{"kubermatic", "vpa", "recommender"}, vpaVersion) ||
		updateDockerImage(doc, yamled.Path{"kubermatic", "vpa", "admissioncontroller"}, vpaVersion)

	if updated {
		c.logger.Infof("Updated VPA version to %s.", vpaVersion)
	}

	return nil
}

func (c *converter) updateCertManager(doc *yamled.Document) error {
	config, exists := doc.Get(yamled.Path{"certManager", "image"})
	if exists {
		c.logger.Infof("Moving certManager.image to certManager.controller.image...")
		doc.Set(yamled.Path{"certManager", "controller", "image"}, config)
	}

	config, exists = doc.Get(yamled.Path{"certManager", "webhookImage"})
	if exists {
		c.logger.Infof("Moving certManager.webhookImage to certManager.webhook.image...")
		doc.Set(yamled.Path{"certManager", "webhook", "image"}, config)
	}

	config, exists = doc.Get(yamled.Path{"certManager", "caSyncImage"})
	if exists {
		c.logger.Infof("Moving certManager.caSyncImage to certManager.cainjector.image...")
		doc.Set(yamled.Path{"certManager", "cainjector", "image"}, config)
	}

	doc.Set(yamled.Path{"certManager", "cainjector", "image", "repository"}, "quay.io/jetstack/cert-manager-cainjector")

	updated := updateDockerImage(doc, yamled.Path{"certManager", "controller"}, certManagerVersion) ||
		updateDockerImage(doc, yamled.Path{"certManager", "webhook"}, certManagerVersion) ||
		updateDockerImage(doc, yamled.Path{"certManager", "cainjector"}, certManagerVersion)

	if updated {
		c.logger.Infof("Updated cert-manager version to %s.", certManagerVersion)
	}

	return nil
}

func (c *converter) updateNginx(doc *yamled.Document) error {
	path := yamled.Path{"nginx", "prometheus"}

	if doc.Has(path) {
		doc.Remove(path)
		c.logger.Info("Removed NGINX Prometheus configuration, port is now always set to 10254.")
	}

	if updateDockerImage(doc, yamled.Path{"nginx"}, nginxVersion) {
		c.logger.Info("Updated NGINX ingress version.")
	}

	if ignored, _ := doc.GetBool(yamled.Path{"nginx", "ignoreMasterTaint"}); ignored {
		doc.Remove(yamled.Path{"nginx", "ignoreMasterTaint"})
		doc.Set(yamled.Path{"nginx", "tolerations"}, []toleration{
			{
				Key:      "only_critical",
				Operator: "Equal",
				Value:    "true",
				Effect:   "NoSchedule",
			},
			{
				Key:      "dedicated",
				Operator: "Equal",
				Value:    "master",
				Effect:   "NoSchedule",
			},
			{
				Key:    "node-role.kubernetes.io/master",
				Effect: "NoSchedule",
			},
		})

		c.logger.Info("Replaced nginx.ignoreMasterTaint flag with explicit tolerations.")
	}

	return nil
}

func (c *converter) updateNodePortProxy(doc *yamled.Document) error {
	if updateDockerImage(doc, yamled.Path{"nodePortPoxy"}, nodePortProxyVersion) {
		c.logger.Info("Updated Node Port Proxy version.")
	}

	return nil
}

func (c *converter) updateDex(doc *yamled.Document) error {
	if updateDockerImage(doc, yamled.Path{"dex"}, dexVersion) {
		c.logger.Info("Updated Dex version.")
	}

	return nil
}

func (c *converter) updateMinio(doc *yamled.Document) error {
	path := yamled.Path{"minio", "image", "tag"}

	version, exists := doc.GetString(path)
	if exists && version < minioVersion {
		doc.Set(path, minioVersion)
		c.logger.Info("Updated Minio version.")
	}

	backupFlag, exists := doc.GetBool(yamled.Path{"minio", "backups"})
	if exists {
		doc.Remove(yamled.Path{"minio", "backups"})
		doc.Set(yamled.Path{"minio", "backup"}, map[string]interface{}{
			"enabled": backupFlag,
		})

		c.logger.Info("Renamed minio.backups to minio.backup.enabled.")
	}

	return nil
}

func (c *converter) updateAlertmanager(doc *yamled.Document) error {
	if updateDockerImage(doc, yamled.Path{"alertmanager"}, alertmanagerVersion) {
		c.logger.Info("Updated Alertmanager version.")
	}

	return nil
}

func (c *converter) updateGrafana(doc *yamled.Document) error {
	if updateDockerImage(doc, yamled.Path{"grafana"}, grafanaVersion) {
		c.logger.Info("Updated Grafana version.")
	}

	return nil
}

func (c *converter) updatePrometheus(doc *yamled.Document) error {
	backupFlag, exists := doc.GetBool(yamled.Path{"prometheus", "backups"})
	if exists {
		doc.Remove(yamled.Path{"prometheus", "backups"})
		doc.Set(yamled.Path{"prometheus", "backup"}, map[string]interface{}{
			"enabled": backupFlag,
		})

		c.logger.Info("Renamed prometheus.backups to prometheus.backup.enabled.")
	}

	if updateDockerImage(doc, yamled.Path{"prometheus"}, prometheusVersion) {
		c.logger.Info("Updated Prometheus version.")
	}

	return nil
}

func (c *converter) updateElasticsearch(doc *yamled.Document) error {
	doc.Remove(yamled.Path{"logging", "elasticsearch", "optimizations"})

	if updateDockerImage(doc, yamled.Path{"logging", "elasticsearch"}, elasticsearchVersion) {
		c.logger.Info("Updated Elasticsearch version.")
	}

	if updateDockerImage(doc, yamled.Path{"logging", "elasticsearch", "curator"}, curatorVersion) {
		c.logger.Info("Updated Curator version.")
	}

	if updateDockerImage(doc, yamled.Path{"logging", "elasticsearch", "init"}, busyboxVersion) {
		c.logger.Info("Updated Busybox version.")
	}

	path := yamled.Path{"logging", "elasticsearch", "image", "repository"}
	if repo, _ := doc.GetString(path); repo == "docker.elastic.co/elasticsearch/elasticsearch" {
		doc.Set(path, "docker.elastic.co/elasticsearch/elasticsearch-oss")
		c.logger.Info("Switched to Open-Source Elasticsearch Docker repository.")
	}

	return nil
}

func (c *converter) updateKibana(doc *yamled.Document) error {
	if updateDockerImage(doc, yamled.Path{"logging", "kibana"}, kibanaVersion) {
		c.logger.Info("Updated Kibana version.")
	}

	path := yamled.Path{"logging", "kibana", "image", "repository"}
	if repo, _ := doc.GetString(path); repo == "docker.elastic.co/kibana/kibana" {
		doc.Set(path, "docker.elastic.co/kibana/kibana-oss")
		c.logger.Info("Switched to Open-Source Kibana Docker repository.")
	}

	return nil
}

func (c *converter) updateFluentbit(doc *yamled.Document) error {
	if updateDockerImage(doc, yamled.Path{"logging", "fluentbit"}, fluentbitVersion) {
		c.logger.Info("Updated fluentbit version.")
	}

	return nil
}

func (c *converter) updateVelero(doc *yamled.Document) error {
	config, exists := doc.Get(yamled.Path{"ark"})
	if !exists {
		return nil
	}

	doc.Remove(yamled.Path{"ark"})
	doc.Set(yamled.Path{"velero"}, config)
	c.logger.Info("Copied Ark configuration to Velero configuration.")

	if updateDockerImage(doc, yamled.Path{"velero"}, veleroVersion) {
		c.logger.Info("Updated Velero version.")
	}

	doc.Set(yamled.Path{"velero", "image", "repository"}, "gcr.io/heptio-images/velero")

	resticFlag, exists := doc.GetBool(yamled.Path{"velero", "restic"})
	if exists {
		doc.Remove(yamled.Path{"velero", "restic"})
		doc.Set(yamled.Path{"velero", "restic", "deploy"}, resticFlag)

		c.logger.Info("Renamed ark.restic flag to velero.restic.deploy.")
	}

	return nil
}

func updateDockerImage(doc *yamled.Document, path yamled.Path, version string) bool {
	return util.UpdateVersion(doc, append(path, "image", "tag"), version)
}
