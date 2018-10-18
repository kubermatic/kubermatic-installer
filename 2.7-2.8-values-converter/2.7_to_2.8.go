package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func convert_27_to_28(v *yaml.MapSlice, isMaster bool) error {
	if err := setIsMaster(v, isMaster); err != nil {
		return fmt.Errorf("setting isMaster: %s", err)
	}
	logrus.Info("Added 'kubermatic->isMaster' entry.")

	if err := mergeDockerAuthData(v); err != nil {
		return fmt.Errorf("merging Docker auth data: %s", err)
	}
	logrus.Info("Merged Docker auth data.")

	if err := updateCertManagerSettings(v); err != nil {
		return fmt.Errorf("removing old cert manager settings: %s", err)
	}
	logrus.Info("Removed old cert manager settings.")

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

func updateCertManagerSettings(v *yaml.MapSlice) error {
	if _, err := removeEntry(v, []string{"replicaCount"}); err != nil {
		logrus.Errorf(`cannot remove section 'replicaCount': %s`, err)
	}

	if _, err := removeEntry(v, []string{"image"}); err != nil {
		logrus.Errorf(`cannot remove section 'image': %s`, err)
	}

	if _, err := removeEntry(v, []string{"createCustomResource"}); err != nil {
		logrus.Errorf(`cannot remove section 'createCustomResource': %s`, err)
	}

	if _, err := removeEntry(v, []string{"rbac"}); err != nil {
		logrus.Errorf(`cannot remove section 'rbac': %s`, err)
	}

	if _, err := removeEntry(v, []string{"resources"}); err != nil {
		logrus.Errorf(`cannot remove section 'resources': %s`, err)
	}

	return nil
}
