package v2_8

import (
	"strings"
	"testing"

	"github.com/kubermatic/kubermatic-installer/pkg/yamled"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func assertYAML(t *testing.T, doc *yamled.Document, expected string) {
	data, err := yaml.Marshal(doc)
	assert.NoError(t, err)

	assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(string(data)))
}

func TestSetIsMaster(t *testing.T) {
	input := `
kubermatic:
  b:
    c: foo
    d: bar
  e:
    f: lol
prometheus: wut
`
	expectedOutput := `
kubermatic:
  b:
    c: foo
    d: bar
  e:
    f: lol
  isMaster: true
prometheus: wut
`

	doc, err := yamled.Load(strings.NewReader(input))
	assert.NoError(t, err)

	converter := NewConverter(nil)

	err = converter.setIsMaster(doc, true)
	assert.NoError(t, err)

	assertYAML(t, doc, expectedOutput)
}

func TestMergeDockerAuthJSON(t *testing.T) {
	input := `
kubermatic:
  docker:
    secret: ewogICJhdXRocyI6IHsKICAgICJodHRwczovL2luZGV4LmRvY2tlci5pby92MS8iOiB7CiAgICAgICJhdXRoIjogImZvbyIsCiAgICAgICJlbWFpbCI6ICJ1c2VyMUBleGFtcGxlLmNvbSIKICAgIH0KICB9Cn0K
  quay:
    secret: ewogICJhdXRocyI6IHsKICAgICJxdWF5LmlvIjogewogICAgICAiYXV0aCI6ICJiYXIiLAogICAgICAiZW1haWwiOiAidXNlcjJAZXhhbXBsZS5jb20iCiAgICB9CiAgfQp9Cg==
`
	expectedOutput := `
kubermatic:
  imagePullSecretData: ewogICJhdXRocyI6IHsKICAgICJodHRwczovL2luZGV4LmRvY2tlci5pby92MS8iOiB7CiAgICAgICJhdXRoIjogImZvbyIsCiAgICAgICJlbWFpbCI6ICJ1c2VyMUBleGFtcGxlLmNvbSIKICAgIH0sCiAgICAicXVheS5pbyI6IHsKICAgICAgImF1dGgiOiAiYmFyIiwKICAgICAgImVtYWlsIjogInVzZXIyQGV4YW1wbGUuY29tIgogICAgfQogIH0KfQ==
`

	doc, err := yamled.Load(strings.NewReader(input))
	assert.NoError(t, err)

	converter := NewConverter(nil)

	err = converter.mergeDockerAuthData(doc)
	assert.NoError(t, err)

	assertYAML(t, doc, expectedOutput)
}

func TestUpdateCertManagerSettings(t *testing.T) {
	input := `
kubermatic:
  foo: bar

# ====== cert-manager ======
# Default values for cert-manager.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.
replicaCount: 1

image:
  repository: quay.io/jetstack/cert-manager-controller
  tag: v0.2.3
  pullPolicy: Always

createCustomResource: true

rbac:
  enabled: true

resources: {}

`
	expectedOutput := `
kubermatic:
  foo: bar
`

	doc, err := yamled.Load(strings.NewReader(input))
	assert.NoError(t, err)

	converter := NewConverter(nil)

	err = converter.updateCertManagerSettings(doc)
	assert.NoError(t, err)

	assertYAML(t, doc, expectedOutput)
}
