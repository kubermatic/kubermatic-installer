package v2_8

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/kubermatic/kubermatic-installer/pkg/yamled"
	"github.com/sirupsen/logrus"
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
	if err := c.setIsMaster(doc, isMaster); err != nil {
		return fmt.Errorf("setting isMaster: %s", err)
	}
	c.logger.Info("Added 'kubermatic.isMaster' entry.")

	if err := c.mergeDockerAuthData(doc); err != nil {
		return fmt.Errorf("merging Docker auth data: %s", err)
	}
	c.logger.Info("Merged Docker auth data.")

	if err := c.updateCertManagerSettings(doc); err != nil {
		return fmt.Errorf("removing old cert-manager settings: %s", err)
	}
	c.logger.Info("Removed old cert-manager settings.")

	return nil
}

func (c *converter) setIsMaster(doc *yamled.Document, isMaster bool) error {
	if !doc.Fill(yamled.Path{"kubermatic", "isMaster"}, isMaster) {
		return errors.New("failed to set isMaster flag")
	}

	return nil
}

func (c *converter) mergeDockerAuthData(doc *yamled.Document) error {
	if doc.Has(yamled.Path{"kubermatic", "imagePullSecretData"}) {
		return nil
	}

	dockerData, err := c.mergeDockerAuthDataGetDocker(doc)
	if err != nil {
		return fmt.Errorf("extracting 'kubermatic.docker.secret': %s", err)
	}

	quayData, err := c.mergeDockerAuthDataGetQuay(doc)
	if err != nil {
		return fmt.Errorf("extracting 'kubermatic.quay.secret': %s", err)
	}

	mergedJSONData, err := c.mergeDockerAuthMergeJSONs(dockerData, quayData)
	if err != nil {
		return fmt.Errorf("merging auth JSONs: %s", err)
	}

	doc.Set(
		yamled.Path{"kubermatic", "imagePullSecretData"},
		base64.StdEncoding.EncodeToString(mergedJSONData),
	)

	return nil
}

func (c *converter) mergeDockerAuthDataGetDocker(doc *yamled.Document) ([]byte, error) {
	secret, ok := doc.GetString(yamled.Path{"kubermatic", "docker", "secret"})

	doc.Remove(yamled.Path{"kubermatic", "docker"})

	if ok {
		return base64.StdEncoding.DecodeString(secret)
	}

	return nil, nil
}

func (c *converter) mergeDockerAuthDataGetQuay(doc *yamled.Document) ([]byte, error) {
	secret, ok := doc.GetString(yamled.Path{"kubermatic", "quay", "secret"})

	doc.Remove(yamled.Path{"kubermatic", "quay"})

	if ok {
		return base64.StdEncoding.DecodeString(secret)
	}

	return nil, nil
}

func (c *converter) mergeDockerAuthMergeJSONs(input ...[]byte) ([]byte, error) {
	type authDatum struct {
		Auth  string `json:"auth"`
		Email string `json:"email"`
	}

	type authData struct {
		Auths map[string]authDatum `json:"auths"`
	}

	mergedAuthData := authData{
		Auths: make(map[string]authDatum),
	}

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

func (c *converter) updateCertManagerSettings(doc *yamled.Document) error {
	// in 2.7 these keys were accidentally on the top-level of the values.yaml
	// before we moved them down into a `certManager` key

	doc.Remove(yamled.Path{"replicaCount"})
	doc.Remove(yamled.Path{"image"})
	doc.Remove(yamled.Path{"createCustomResource"})
	doc.Remove(yamled.Path{"rbac"})
	doc.Remove(yamled.Path{"resources"})

	return nil
}
