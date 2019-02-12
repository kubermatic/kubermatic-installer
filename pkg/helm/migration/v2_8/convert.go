package v2_8

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/kubermatic/kubermatic-installer/pkg/helm/migration/util"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type converter struct {
	logger logrus.FieldLogger
}

func NewConverter(logger logrus.FieldLogger) *converter {
	return &converter{
		logger: logger,
	}
}

func (c *converter) Convert(v *yaml.MapSlice, isMaster bool) error {
	if err := c.setIsMaster(v, isMaster); err != nil {
		return fmt.Errorf("setting isMaster: %s", err)
	}
	c.logger.Info("Added 'kubermatic.isMaster' entry.")

	if err := c.mergeDockerAuthData(v); err != nil {
		return fmt.Errorf("merging Docker auth data: %s", err)
	}
	c.logger.Info("Merged Docker auth data.")

	if err := c.updateCertManagerSettings(v); err != nil {
		return fmt.Errorf("removing old cert-manager settings: %s", err)
	}
	c.logger.Info("Removed old cert-manager settings.")

	return nil
}

func (c *converter) setIsMaster(v *yaml.MapSlice, isMaster bool) error {
	kubermatic := util.GetEntry(v, "kubermatic")
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

func (c *converter) mergeDockerAuthData(v *yaml.MapSlice) error {
	dockerData, err := c.mergeDockerAuthDataGetDocker(v)
	if err != nil {
		return fmt.Errorf("extracting 'kubermatic.docker.secret': %s", err)
	}

	quayData, err := c.mergeDockerAuthDataGetQuay(v)
	if err != nil {
		return fmt.Errorf("extracting 'kubermatic.quay.secret': %s", err)
	}

	mergedJSONData, err := c.mergeDockerAuthMergeJSONs(dockerData, quayData)
	if err != nil {
		return fmt.Errorf("merging auth JSONs: %s", err)
	}

	newEntry := yaml.MapItem{
		Key:   "imagePullSecretData",
		Value: base64.StdEncoding.EncodeToString(mergedJSONData),
	}

	kubermatic := util.GetEntry(v, "kubermatic")
	if kubermatic == nil {
		return fmt.Errorf("section 'kubermatic' not found")
	}

	val := kubermatic.Value.(yaml.MapSlice)
	val = append(yaml.MapSlice{newEntry}, val...)
	kubermatic.Value = val

	return nil
}

func (c *converter) mergeDockerAuthDataGetDocker(v *yaml.MapSlice) ([]byte, error) {
	dockerSection, err := util.RemoveEntry(v, []string{"kubermatic", "docker"})
	if err != nil {
		return nil, fmt.Errorf("removing 'kubermatic.docker': %s", err)
	}

	dockerSlice := dockerSection.Value.(yaml.MapSlice)
	secretEntry := util.GetEntry(&dockerSlice, "secret")
	if secretEntry == nil {
		return nil, fmt.Errorf("section 'kubermatic.docker.secret' not found")
	}

	return base64.StdEncoding.DecodeString(secretEntry.Value.(string))
}

func (c *converter) mergeDockerAuthDataGetQuay(v *yaml.MapSlice) ([]byte, error) {
	quaySection, err := util.RemoveEntry(v, []string{"kubermatic", "quay"})
	if err != nil {
		return nil, fmt.Errorf("removing 'kubermatic.quay': %s", err)
	}

	quaySlice := quaySection.Value.(yaml.MapSlice)
	secretEntry := util.GetEntry(&quaySlice, "secret")
	if secretEntry == nil {
		return nil, fmt.Errorf("section 'kubermatic.quay.secret' not found")
	}

	return base64.StdEncoding.DecodeString(secretEntry.Value.(string))
}

func (c *converter) mergeDockerAuthMergeJSONs(input ...[]byte) ([]byte, error) {
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

func (c *converter) updateCertManagerSettings(v *yaml.MapSlice) error {
	if _, err := util.RemoveEntry(v, []string{"replicaCount"}); err != nil {
		logrus.Errorf(`cannot remove section 'replicaCount': %s`, err)
	}

	if _, err := util.RemoveEntry(v, []string{"image"}); err != nil {
		logrus.Errorf(`cannot remove section 'image': %s`, err)
	}

	if _, err := util.RemoveEntry(v, []string{"createCustomResource"}); err != nil {
		logrus.Errorf(`cannot remove section 'createCustomResource': %s`, err)
	}

	if _, err := util.RemoveEntry(v, []string{"rbac"}); err != nil {
		logrus.Errorf(`cannot remove section 'rbac': %s`, err)
	}

	if _, err := util.RemoveEntry(v, []string{"resources"}); err != nil {
		logrus.Errorf(`cannot remove section 'resources': %s`, err)
	}

	return nil
}
