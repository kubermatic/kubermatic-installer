package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func convert_27_to_28(v *yaml.MapSlice, isMaster bool) error {
	if err := updateImageTags(v); err != nil {
		return fmt.Errorf("updating image versions: %s", err)
	}
	logrus.Info("Updated image versions.")

	if err := addRBACController(v); err != nil {
		return fmt.Errorf("addding rbac controller section: %s", err)
	}
	logrus.Info("Added RBAC controller section at 'kubermatic->rbac'.")

	if err := setIsMaster(v, isMaster); err != nil {
		return fmt.Errorf("setting isMaster: %s", err)
	}
	logrus.Info("Added 'kubermatic->isMaster' entry.")

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

func addRBACController(v *yaml.MapSlice) error {
	section := `
kubermatic:
  rbac:
    replicas: 1
    image:
      repository: kubermatic/api
      tag: v2.8.0-rc.4
      pullPolicy: IfNotPresent
`
	return mergeSection(v, section)
}

func updateImageTags(v *yaml.MapSlice) error {
	kubermaticVersion := "v2.8.0-rc.4"
	uiVersion := "v1.0.1"
	addonsVersion := "v0.1.12"

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

	return nil
}

func setIsMaster(v *yaml.MapSlice, isMaster bool) error {
	kubermatic := getEntry(v, "kubermatic")
	if kubermatic == nil {
		return fmt.Errorf(`section 'kubermatic' not found`)
	}

	isMasterEntry := yaml.MapItem{
		Key:   "isMaster",
		Value: isMaster,
	}

	val := kubermatic.Value.(yaml.MapSlice)
	val = append(yaml.MapSlice{isMasterEntry}, val...)
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
