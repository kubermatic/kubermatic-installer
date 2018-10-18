package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestSetIsMaster(t *testing.T) {
	input := `kubermatic:
  b:
    c: foo
    d: bar
  e:
    f: lol
prometheus: wut
`
	expectedOutput := `kubermatic:
  isMaster: true
  b:
    c: foo
    d: bar
  e:
    f: lol
prometheus: wut
`

	var values yaml.MapSlice

	err := yaml.Unmarshal([]byte(input), &values)
	assert.NoError(t, err)

	err = setIsMaster(&values, true)
	assert.NoError(t, err)

	data, err := yaml.Marshal(values)
	assert.NoError(t, err)

	assert.Equal(t, expectedOutput, string(data))
}

func TestAddRBACController(t *testing.T) {
	input := `
kubermatic:
   api:
    replicas: 2
    image:
      repository: "kubermatic/api"
      tag: "v2.8.0-rc.4"
      pullPolicy: "IfNotPresent"
`
	expectedOutput := `kubermatic:
  api:
    replicas: 2
    image:
      repository: kubermatic/api
      tag: v2.8.0-rc.4
      pullPolicy: IfNotPresent
  rbac:
    replicas: 1
    image:
      repository: kubermatic/api
      tag: v2.8.0-rc.4
      pullPolicy: IfNotPresent
`

	var values yaml.MapSlice

	err := yaml.Unmarshal([]byte(input), &values)
	assert.NoError(t, err)

	err = addRBACController(&values)
	assert.NoError(t, err)

	data, err := yaml.Marshal(values)
	assert.NoError(t, err)

	assert.Equal(t, expectedOutput, string(data))
}

func TestMergeDockerAuthJSON(t *testing.T) {
	input := `
kubermatic:
  docker:
    secret: ewogICJhdXRocyI6IHsKICAgICJodHRwczovL2luZGV4LmRvY2tlci5pby92MS8iOiB7CiAgICAgICJhdXRoIjogImZvbyIsCiAgICAgICJlbWFpbCI6ICJ1c2VyMUBleGFtcGxlLmNvbSIKICAgIH0KICB9Cn0K
  quay:
    secret: ewogICJhdXRocyI6IHsKICAgICJxdWF5LmlvIjogewogICAgICAiYXV0aCI6ICJiYXIiLAogICAgICAiZW1haWwiOiAidXNlcjJAZXhhbXBsZS5jb20iCiAgICB9CiAgfQp9Cg==
`
	expectedOutput := `kubermatic:
  imagePullSecretData: ewogICJhdXRocyI6IHsKICAgICJodHRwczovL2luZGV4LmRvY2tlci5pby92MS8iOiB7CiAgICAgICJhdXRoIjogImZvbyIsCiAgICAgICJlbWFpbCI6ICJ1c2VyMUBleGFtcGxlLmNvbSIKICAgIH0sCiAgICAicXVheS5pbyI6IHsKICAgICAgImF1dGgiOiAiYmFyIiwKICAgICAgImVtYWlsIjogInVzZXIyQGV4YW1wbGUuY29tIgogICAgfQogIH0KfQ==
`

	var values yaml.MapSlice

	err := yaml.Unmarshal([]byte(input), &values)
	assert.NoError(t, err)

	err = mergeDockerAuthData(&values)
	assert.NoError(t, err)

	data, err := yaml.Marshal(values)
	assert.NoError(t, err)

	assert.Equal(t, expectedOutput, string(data))
}
