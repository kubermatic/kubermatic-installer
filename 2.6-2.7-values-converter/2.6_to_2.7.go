package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func convert_26_to_27(v *yaml.MapSlice) error {
	if err := moveAddonsFromKubermaticToController(v); err != nil {
		return fmt.Errorf("moving the addons section: %s", err)
	}
	logrus.Info("Moved 'kubermatic->addons->defaultAddons' section to 'kubermatic->controller->addons->defaultAddons'.")

	if err := updateImageTags(v); err != nil {
		return fmt.Errorf("updating image versions: %s", err)
	}
	logrus.Info("Updated image versions.")

	if err := addS3ExporterSection(v); err != nil {
		return fmt.Errorf("adding S3 exporter config section: %s", err)
	}
	logrus.Info("Added 'kubermatic->s3_exporter' section.")

	if _, err := removeEntry(v, []string{"kubeStateMetrics", "rbacProxy"}); err != nil {
		return fmt.Errorf(`failed to remove 'kubeStateMetrics->rbacProxy': %s`, err)
	}
	logrus.Info("Removed 'kubeStateMetrics->rbacProxy' section.")

	if err := updatePrometheusConfig(v); err != nil {
		return fmt.Errorf("updating prometheus config: %s", err)
	}
	logrus.Info("Added new settings to the 'prometheus' section.")

	if _, err := removeEntry(v, []string{"prometheusOperator"}); err != nil {
		return fmt.Errorf(`failed to remove 'prometheusOperator': %s`, err)
	}
	logrus.Info("Removed 'prometheusOperator' section.")

	return nil
}

func moveAddonsFromKubermaticToController(v *yaml.MapSlice) error {
	// remove the addons list from under 'kubermatic'
	addons, err := removeEntry(v, []string{"kubermatic", "addons"})
	if err != nil {
		return fmt.Errorf(`failed to remove 'kubermatic->addons': %s`, err)
	}
	addonsMapSlice := addons.Value.(yaml.MapSlice)

	defaultAddons := getEntry(&addonsMapSlice, "defaultAddons")
	if defaultAddons == nil {
		return fmt.Errorf(`section 'kubermatic->addons->defaultAddons' not found`)
	}

	// add under 'kubermatic->controller->addons'
	controllerAddons, err := getPath(v, []string{"kubermatic", "controller", "addons"})
	if err != nil {
		return fmt.Errorf(`failed to get 'kubermatic->controller->addons': %s`, err)
	}
	controllerAddonsMapSlice := controllerAddons.Value.(yaml.MapSlice)

	controllerAddonsMapSlice = append(controllerAddonsMapSlice, *defaultAddons)
	controllerAddons.Value = controllerAddonsMapSlice

	return nil
}

func updateImageTags(v *yaml.MapSlice) error {
	kubermaticVersion := "v2.7.7"
	uiVersion := "v0.38.0"
	addonsVersion := "v0.1.11"
	nginxVersion := "0.18.0"
	alertManagerVersion := "v0.15.0"
	kubeStateMetricsRepo := "k8s.gcr.io/addon-resizer"
	kubeStateMetricsVersion := "1.7"

	if err := modifyEntry(v, []string{"kubermatic", "controller", "image", "tag"}, kubermaticVersion); err != nil {
		return fmt.Errorf("Failed to set 'kubermatic->controller->image->tag': %s", err)
	}

	if err := modifyEntry(v, []string{"kubermatic", "controller", "addons", "image", "tag"}, addonsVersion); err != nil {
		return fmt.Errorf("Failed to set 'kubermatic->controller->addons->image->tag': %s", err)
	}

	if err := modifyEntry(v, []string{"kubermatic", "api", "image", "tag"}, kubermaticVersion); err != nil {
		return fmt.Errorf("Failed to set 'kubermatic->api->image->tag': %s", err)
	}

	if err := modifyEntry(v, []string{"kubermatic", "ui", "image", "tag"}, uiVersion); err != nil {
		return fmt.Errorf("Failed to set 'kubermatic->ui->image->tag': %s", err)
	}

	if err := modifyEntry(v, []string{"nginx", "image", "tag"}, nginxVersion); err != nil {
		return fmt.Errorf("Failed to set 'nginx->image->tag': %s", err)
	}

	if err := modifyEntry(v, []string{"alertmanager", "version"}, alertManagerVersion); err != nil {
		return fmt.Errorf("Failed to set 'alertmanager->version': %s", err)
	}

	if err := modifyEntry(v, []string{"kubeStateMetrics", "resizer", "image", "repository"}, kubeStateMetricsRepo); err != nil {
		return fmt.Errorf("Failed to set 'kubeStateMetrics->resizer->image->repository': %s", err)
	}
	if err := modifyEntry(v, []string{"kubeStateMetrics", "resizer", "image", "tag"}, kubeStateMetricsVersion); err != nil {
		return fmt.Errorf("Failed to set 'kubeStateMetrics->resizer->image->tag': %s", err)
	}

	return nil
}

func addS3ExporterSection(v *yaml.MapSlice) error {
	kubermatic := getEntry(v, "kubermatic")
	if kubermatic == nil {
		return fmt.Errorf(`section 'kubermatic' not found`)
	}

	newEntry := yaml.MapItem{
		Key: "s3_exporter",
		Value: yaml.MapSlice{
			yaml.MapItem{
				Key: "image",
				Value: yaml.MapSlice{
					yaml.MapItem{
						Key:   "repository",
						Value: "quay.io/kubermatic/s3-exporter",
					},
					yaml.MapItem{
						Key:   "tag",
						Value: "v0.2",
					},
				},
			},
			yaml.MapItem{
				Key:   "endpoint",
				Value: "http://minio.minio.svc.cluster.local:9000",
			},
			yaml.MapItem{
				Key:   "bucket",
				Value: "kubermatic-etcd-backups",
			},
		},
	}

	val := kubermatic.Value.(yaml.MapSlice)
	val = append(val, newEntry)
	kubermatic.Value = val

	return nil
}

func updatePrometheusConfig(v *yaml.MapSlice) error {
	prometheus := getEntry(v, "prometheus")
	if prometheus == nil {
		return fmt.Errorf(`section 'prometheus' not found`)
	}

	newEntries := []yaml.MapItem{
		yaml.MapItem{
			Key:   "storageSize",
			Value: "100Gi",
		},
		yaml.MapItem{
			Key: "externalLabels",
			Value: yaml.MapSlice{
				yaml.MapItem{
					Key:   "region",
					Value: "default",
				},
			},
		},
		yaml.MapItem{
			Key: "containers",
			Value: yaml.MapSlice{
				yaml.MapItem{
					Key: "prometheus",
					Value: yaml.MapSlice{
						yaml.MapItem{
							Key: "resources",
							Value: yaml.MapSlice{
								yaml.MapItem{
									Key: "limits",
									Value: yaml.MapSlice{
										yaml.MapItem{
											Key:   "cpu",
											Value: 1,
										},
										yaml.MapItem{
											Key:   "memory",
											Value: "2Gi",
										},
									},
								},
								yaml.MapItem{
									Key: "requests",
									Value: yaml.MapSlice{
										yaml.MapItem{
											Key:   "cpu",
											Value: "100m",
										},
										yaml.MapItem{
											Key:   "memory",
											Value: "512Mi",
										},
									},
								},
							},
						},
					},
				},
				yaml.MapItem{
					Key: "reloader",
					Value: yaml.MapSlice{
						yaml.MapItem{
							Key: "resources",
							Value: yaml.MapSlice{
								yaml.MapItem{
									Key: "limits",
									Value: yaml.MapSlice{
										yaml.MapItem{
											Key:   "cpu",
											Value: "100m",
										},
										yaml.MapItem{
											Key:   "memory",
											Value: "64Mi",
										},
									},
								},
								yaml.MapItem{
									Key: "requests",
									Value: yaml.MapSlice{
										yaml.MapItem{
											Key:   "cpu",
											Value: "25m",
										},
										yaml.MapItem{
											Key:   "memory",
											Value: "16Mi",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	val := prometheus.Value.(yaml.MapSlice)
	val = append(val, newEntries...)
	prometheus.Value = val

	return nil
}
