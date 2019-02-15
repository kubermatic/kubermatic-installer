package util

import (
	"strings"

	yaml "gopkg.in/yaml.v2"
)

func LoadYAML(content string) map[string]interface{} {
	result := make(map[string]interface{})
	yaml.Unmarshal([]byte(strings.TrimSpace(content)), &result)

	return result
}
