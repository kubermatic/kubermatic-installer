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

func TestUpdatePrometheusConfig(t *testing.T) {
	input := `kubermatic:
  b:
    c: foo
    d: bar
  e:
    f: lol
prometheus:
  auth: 'Y3VyaW9zaXR5IGtpbGxlZCB0aGUgY2F0Cg=='
  version: 'v2.2.1'
  host: ""
`
	expectedOutput := `kubermatic:
  b:
    c: foo
    d: bar
  e:
    f: lol
prometheus:
  auth: Y3VyaW9zaXR5IGtpbGxlZCB0aGUgY2F0Cg==
  version: v2.2.1
  host: ""
  storageSize: 100Gi
  externalLabels:
    region: default
  containers:
    prometheus:
      resources:
        limits:
          cpu: 1
          memory: 2Gi
        requests:
          cpu: 100m
          memory: 512Mi
    reloader:
      resources:
        limits:
          cpu: 100m
          memory: 64Mi
        requests:
          cpu: 25m
          memory: 16Mi
`

	var values yaml.MapSlice

	err := yaml.Unmarshal([]byte(input), &values)
	assert.NoError(t, err)

	err = updatePrometheusConfig(&values)
	assert.NoError(t, err)

	data, err := yaml.Marshal(values)
	assert.NoError(t, err)

	assert.Equal(t, expectedOutput, string(data))
}
