package task

import (
	"github.com/sirupsen/logrus"
)

// InfoTask is used to just display some status information, but not actually do anything.
type InfoTask struct {
	Message string
}

func (t *InfoTask) Run(_ *Config, _ *State, _ *Clients, log logrus.FieldLogger, _ bool) error {
	log.Info(t.Message)
	return nil
}
