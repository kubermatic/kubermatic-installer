package v2_9

import (
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
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
	return nil
}
