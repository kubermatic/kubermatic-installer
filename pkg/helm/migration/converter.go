package migration

import yaml "gopkg.in/yaml.v2"

type converter interface {
	Convert(v *yaml.MapSlice, isMaster bool) error
}
