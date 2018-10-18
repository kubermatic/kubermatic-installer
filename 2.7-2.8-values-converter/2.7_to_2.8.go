package main

import (
	"encoding/base64"
	"encoding/json"
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

	if err := mergeDockerAuthData(v); err != nil {
		return fmt.Errorf("merging Docker auth data: %s", err)
	}
	logrus.Info("Merged Docker auth data.")

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

func mergeDockerAuthData(v *yaml.MapSlice) error {
	dockerData, err := mergeDockerAuthDataGetDocker(v)
	if err != nil {
		return fmt.Errorf("extracting 'kubermatic->docker->secret': %s", err)
	}

	quayData, err := mergeDockerAuthDataGetQuay(v)
	if err != nil {
		return fmt.Errorf("extracting 'kubermatic->quay->secret': %s", err)
	}

	mergedJSONData, err := mergeDockerAuthMergeJSONs(dockerData, quayData)
	if err != nil {
		return fmt.Errorf("merging auth JSONs: %s", err)
	}

	newEntry := yaml.MapItem{
		Key:   "imagePullSecretData",
		Value: base64.StdEncoding.EncodeToString(mergedJSONData),
	}

	kubermatic := getEntry(v, "kubermatic")
	if kubermatic == nil {
		return fmt.Errorf("section 'kubermatic' not found")
	}

	val := kubermatic.Value.(yaml.MapSlice)
	val = append(yaml.MapSlice{newEntry}, val...)
	kubermatic.Value = val

	return nil
}

func mergeDockerAuthDataGetDocker(v *yaml.MapSlice) ([]byte, error) {
	dockerSection, err := removeEntry(v, []string{"kubermatic", "docker"})
	if err != nil {
		return nil, fmt.Errorf("removing 'kubermatic->docker': %s", err)
	}

	dockerSlice := dockerSection.Value.(yaml.MapSlice)
	secretEntry := getEntry(&dockerSlice, "secret")
	if secretEntry == nil {
		return nil, fmt.Errorf("section 'kubermatic->docker->secret' not found")
	}

	return base64.StdEncoding.DecodeString(secretEntry.Value.(string))
}

func mergeDockerAuthDataGetQuay(v *yaml.MapSlice) ([]byte, error) {
	quaySection, err := removeEntry(v, []string{"kubermatic", "quay"})
	if err != nil {
		return nil, fmt.Errorf("removing 'kubermatic->quay': %s", err)
	}

	quaySlice := quaySection.Value.(yaml.MapSlice)
	secretEntry := getEntry(&quaySlice, "secret")
	if secretEntry == nil {
		return nil, fmt.Errorf("section 'kubermatic->quay->secret' not found")
	}

	return base64.StdEncoding.DecodeString(secretEntry.Value.(string))
}

func mergeDockerAuthMergeJSONs(input ...[]byte) ([]byte, error) {
	type authDatum struct {
		Auth  string `json:"auth"`
		Email string `json:"email"`
	}

	type authData struct {
		Auths map[string]authDatum `json:"auths"`
	}

	mergedAuthData := authData{Auths: make(map[string]authDatum)}

	for _, in := range input {
		var inputAuthData authData
		err := json.Unmarshal(in, &inputAuthData)
		if err != nil {
			return nil, fmt.Errorf("unmarshalling docker auth data: %s", err)
		}

		for k, v := range inputAuthData.Auths {
			mergedAuthData.Auths[k] = v
		}
	}

	return json.MarshalIndent(mergedAuthData, "", "  ")
}
