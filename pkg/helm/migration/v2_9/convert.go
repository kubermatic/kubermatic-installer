package v2_9

import (
	"encoding/json"
	"fmt"

	"github.com/kubermatic/kubermatic-installer/pkg/helm/migration/util"
	"github.com/kubermatic/kubermatic-installer/pkg/yamled"
	"github.com/sirupsen/logrus"
)

const (
	alertmanagerVersion     = "v0.16.0"
	certManagerVersion      = "v0.6.0"
	curatorVersion          = "5.6.0-1"
	dexVersion              = "v2.12.0"
	elasticsearchVersion    = "6.5.1"
	grafanaVersion          = "5.4.3"
	kibanaVersion           = "6.5.1"
	kubermaticAddonsVersion = "v0.1.16"
	kubermaticAPIVersion    = "v2.9.1"
	kubermaticUIVersion     = "v1.1.0"
	kubeStateMetricsVersion = "v1.5.0"
	minioVersion            = "RELEASE.2019-01-16T21-44-08Z"
	nginxVersion            = "0.22.0"
	nodeExporterVersion     = "v0.17.0"
	prometheusVersion       = "v2.7.1"
)

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

	if err := c.updateKubermaticUIImage(doc); err != nil {
		return fmt.Errorf("failed to update Kubermatic UI image: %v", err)
	}

	if err := c.updateKubermaticUIConfig(doc); err != nil {
		return fmt.Errorf("failed to update Kubermatic UI configuration : %v", err)
	}

	if err := c.updateCertManager(doc); err != nil {
		return fmt.Errorf("failed to update cert-manager: %v", err)
	}

	if err := c.updateNginx(doc); err != nil {
		return fmt.Errorf("failed to update nginx-ingress: %v", err)
	}

	if err := c.updateDex(doc); err != nil {
		return fmt.Errorf("failed to update dex: %v", err)
	}

	if err := c.updateMinio(doc); err != nil {
		return fmt.Errorf("failed to update minio: %v", err)
	}

	if err := c.updateAlertmanager(doc); err != nil {
		return fmt.Errorf("failed to update alertmanager: %v", err)
	}

	if err := c.updateGrafana(doc); err != nil {
		return fmt.Errorf("failed to update grafana: %v", err)
	}

	if err := c.updateKubeStateMetrics(doc); err != nil {
		return fmt.Errorf("failed to update kube-state-metrics: %v", err)
	}

	if err := c.updateNodeExporter(doc); err != nil {
		return fmt.Errorf("failed to update node-exporter: %v", err)
	}

	if err := c.updatePrometheus(doc); err != nil {
		return fmt.Errorf("failed to update prometheus: %v", err)
	}

	if err := c.updateElasticsearch(doc); err != nil {
		return fmt.Errorf("failed to update elasticsearch: %v", err)
	}

	if err := c.updateKibana(doc); err != nil {
		return fmt.Errorf("failed to update kibana: %v", err)
	}

	if err := c.updateFluentbit(doc); err != nil {
		return fmt.Errorf("failed to update fluentbit: %v", err)
	}

	if err := c.removeMetricsServerAddon(doc); err != nil {
		return fmt.Errorf("failed to remove metrics-server addon: %v", err)
	}

	if err := c.removeS3Exporter(doc); err != nil {
		return fmt.Errorf("failed to remove S3 exporter: %v", err)
	}

	return nil
}

func (c *converter) updateKubermaticController(doc *yamled.Document) error {
	changes := updateDockerImage(doc, yamled.Path{"kubermatic", "controller"}, kubermaticAPIVersion) ||
		updateDockerImage(doc, yamled.Path{"kubermatic", "api"}, kubermaticAPIVersion) ||
		updateDockerImage(doc, yamled.Path{"kubermatic", "rbac"}, kubermaticAPIVersion) ||
		updateDockerImage(doc, yamled.Path{"kubermatic", "controller", "addons"}, kubermaticAddonsVersion)

	if changes {
		c.logger.Info("Updated Kubermatic versions.")
	}

	return nil
}

func (c *converter) updateKubermaticUIImage(doc *yamled.Document) error {
	if updateDockerImage(doc, yamled.Path{"kubermatic", "ui"}, kubermaticUIVersion) {
		c.logger.Info("Updated Kubermatic UI version.")
	}

	path := yamled.Path{"kubermatic", "ui", "image", "repository"}
	image, _ := doc.GetString(path)

	if image == "kubermatic/ui-v2" {
		doc.Set(path, "quay.io/kubermatic/ui-v2")
		c.logger.Info("Updated Kubermatic UI Docker image.")
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

	if _, exists := cfg["share_kubeconfig"]; !exists {
		cfg["share_kubeconfig"] = false
	}

	if _, exists := cfg["show_terms_of_service"]; !exists {
		cfg["show_terms_of_service"] = false
	}

	if _, exists := cfg["cleanup_cluster"]; !exists {
		cfg["cleanup_cluster"] = false
	}

	marshalled, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to re-encode config JSON: %v", err)
	}

	doc.Set(path, string(marshalled))
	c.logger.Info("Updated Kubermatic UI configuration flags.")

	return nil
}

func (c *converter) updateCertManager(doc *yamled.Document) error {
	if updateDockerImage(doc, yamled.Path{"certManager"}, certManagerVersion) {
		c.logger.Info("Updated cert-manager version.")
	}

	return nil
}

func (c *converter) updateNginx(doc *yamled.Document) error {
	path := yamled.Path{"nginx", "defaultBackend"}

	if doc.Has(path) {
		doc.Remove(path)
		c.logger.Info("Removed NGINX default backend.")
	}

	if updateDockerImage(doc, yamled.Path{"nginx"}, nginxVersion) {
		c.logger.Info("Updated NGINX ingress version.")
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

	return nil
}

func (c *converter) updateAlertmanager(doc *yamled.Document) error {
	path := yamled.Path{"alertmanager", "auth"}

	if doc.Has(path) {
		doc.Remove(path)
		c.logger.Info("Removed Alertmanager auth section.")
	}

	if updateDockerImage(doc, yamled.Path{"alertmanager"}, alertmanagerVersion) {
		c.logger.Info("Updated Alertmanager version.")
	}

	return nil
}

func (c *converter) updateGrafana(doc *yamled.Document) error {
	path := yamled.Path{"grafana", "host"}

	if doc.Has(path) {
		doc.Remove(path)
		c.logger.Info("Removed Grafana host.")
	}

	if updateDockerImage(doc, yamled.Path{"grafana"}, grafanaVersion) {
		c.logger.Info("Updated Grafana version.")
	}

	return nil
}

func (c *converter) updateKubeStateMetrics(doc *yamled.Document) error {
	if updateDockerImage(doc, yamled.Path{"kubeStateMetrics"}, kubeStateMetricsVersion) {
		c.logger.Info("Updated kube-state-metrics version.")
	}

	return nil
}

func (c *converter) updateNodeExporter(doc *yamled.Document) error {
	if updateDockerImage(doc, yamled.Path{"nodeExporter"}, nodeExporterVersion) {
		c.logger.Info("Updated node-exporter version.")
	}

	return nil
}

func (c *converter) updatePrometheus(doc *yamled.Document) error {
	path := yamled.Path{"prometheus", "auth"}

	if doc.Has(path) {
		doc.Remove(path)
		c.logger.Info("Removed Prometheus auth section.")
	}

	path = yamled.Path{"kubermatic", "ruleFiles"}

	rules, ok := doc.GetArray(path)
	if !ok {
		return nil
	}

	newRules := make([]string, 0)

	for _, addon := range rules {
		if a, ok := addon.(string); ok {
			if a == "/etc/prometheus/rules/*.yaml" {
				newRules = append(newRules, "/etc/prometheus/rules/general-*.yaml")
				newRules = append(newRules, "/etc/prometheus/rules/kubermatic-*.yaml")
				newRules = append(newRules, "/etc/prometheus/rules/managed-*.yaml")
			} else {
				newRules = append(newRules, a)
			}
		}
	}

	doc.Set(path, newRules)
	c.logger.Info("Update Prometheus rule files list.")

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

	path := yamled.Path{"logging", "elasticsearch", "image", "repository"}
	if repo, _ := doc.GetString(path); repo == "quay.io/pires/docker-elasticsearch-kubernetes" {
		doc.Set(path, "docker.elastic.co/elasticsearch/elasticsearch")
		c.logger.Info("Updated Elasticsearch Docker repository.")
	}

	path = yamled.Path{"logging", "elasticsearch", "curator", "image", "repository"}
	if repo, _ := doc.GetString(path); repo == "quay.io/pires/docker-elasticsearch-curator" {
		doc.Set(path, "quay.io/kubermatic/elasticsearch-curator")
		c.logger.Info("Updated Curator Docker repository.")
	}

	return nil
}

func (c *converter) updateKibana(doc *yamled.Document) error {
	doc.Remove(yamled.Path{"logging", "kibana", "auth"})
	doc.Remove(yamled.Path{"logging", "kibana", "host"})

	if updateDockerImage(doc, yamled.Path{"logging", "kibana"}, kibanaVersion) {
		c.logger.Info("Updated Kibana version.")
	}

	path := yamled.Path{"logging", "kibana", "image", "repository"}
	if repo, _ := doc.GetString(path); repo == "docker.elastic.co/kibana/kibana-oss" {
		doc.Set(path, "docker.elastic.co/kibana/kibana")
		c.logger.Info("Updated Kibana Docker repository.")
	}

	return nil
}

func (c *converter) updateFluentbit(doc *yamled.Document) error {
	doc.Remove(yamled.Path{"logging", "fluentd"})

	return nil
}

func (c *converter) removeMetricsServerAddon(doc *yamled.Document) error {
	path := yamled.Path{"kubermatic", "controller", "addons", "defaultAddons"}

	addons, ok := doc.GetArray(path)
	if !ok {
		return nil
	}

	newAddons := make([]string, 0)

	for _, addon := range addons {
		if a, _ := addon.(string); a != "metrics-server" {
			newAddons = append(newAddons, a)
		}
	}

	doc.Set(path, newAddons)

	return nil
}

func (c *converter) removeS3Exporter(doc *yamled.Document) error {
	doc.Remove(yamled.Path{"kubermatic", "s3_exporter"})

	return nil
}

func updateDockerImage(doc *yamled.Document, path yamled.Path, version string) bool {
	return util.UpdateVersion(doc, append(path, "image", "tag"), version)
}
