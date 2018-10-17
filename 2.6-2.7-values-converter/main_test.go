package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestRemoveRoot(t *testing.T) {
	input := `
a:
  b:
    c: foo
    d: bar
  e:
    f: lol
g: wut`
	expectedOutput := `a:
  b:
    c: foo
    d: bar
  e:
    f: lol
`

	var values yaml.MapSlice

	err := yaml.Unmarshal([]byte(input), &values)
	assert.NoError(t, err)

	removedEntry, err := removeEntry(&values, []string{"g"})
	assert.NoError(t, err)
	assert.NotNil(t, removedEntry)

	data, err := yaml.Marshal(values)
	assert.NoError(t, err)

	assert.Equal(t, expectedOutput, string(data))
}

func TestRemoveEntry(t *testing.T) {
	input := `
a:
  b:
    c: foo
    d: bar
  e:
    f: lol
g: wut`
	expectedOutput := `a:
  b:
    d: bar
  e:
    f: lol
g: wut
`

	var values yaml.MapSlice

	err := yaml.Unmarshal([]byte(input), &values)
	assert.NoError(t, err)

	removedEntry, err := removeEntry(&values, []string{"a", "b", "c"})
	assert.NoError(t, err)
	assert.NotNil(t, removedEntry)

	data, err := yaml.Marshal(values)
	assert.NoError(t, err)

	assert.Equal(t, expectedOutput, string(data))
}
