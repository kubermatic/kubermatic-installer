package server

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kubermatic/kubermatic-installer/pkg/installer"
	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
	"github.com/sirupsen/logrus"
)

type installManager struct {
	logger *logrus.Logger
	logs   map[string]chan logItem
}

func newInstallManager(logger *logrus.Logger) *installManager {
	return &installManager{
		logger: logger,
		logs:   make(map[string]chan logItem),
	}
}

type valuesLogItem struct {
	Type   string `json:"type"`
	Values string `json:"values"`
}

func (i *installManager) Start(m manifest.Manifest) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", errors.New("failed to create UUID")
	}

	installID := id.String()

	// Create a buffer for log items, such as messages and the
	// values.yaml as the final item; make this channel large
	// enough to hold all conceivable items (usual installations
	// log no more than maybe 50 items), so that the installer
	// routine can start to run away and do its thing without
	// having to fear that nobody will read its log messages
	// and it would eventually block. It's better to have a
	// buffer that's too large and never gets read (i.e.
	// because the user never opened the websocket to read it
	// or closed it prematurely) than having an installer
	// routine that aborts in the middle and leaves the cluster
	// in a half-finished state.
	logItems := make(chan logItem, 200)
	i.logs[installID] = logItems

	// setup a multiplexing logger than logs both to the
	// buffer channel and to a CLI-bound logrus logger
	logger := newLogger(i.logger.WithField("proc", installID), logItems)

	// begin the actual installation
	go func() {
		options := installer.InstallerOptions{
			KeepFiles:   true,
			HelmTimeout: 600,
			ValuesFile:  "",
		}

		values, err := installer.NewInstaller(&m, logger).Run(options)
		if err != nil {
			logger.Errorf("Installation failed: %v", err)
		}

		// send out the values.yaml to allow the user to download it
		item := valuesLogItem{
			Type:   "values",
			Values: string(values.YAML()),
		}

		encoded, _ := json.Marshal(item)
		logItems <- encoded

		// do not let the reading websocket read and hang forever
		close(logItems)
	}()

	return installID, err
}

func (i *installManager) Logs(id string) (<-chan logItem, error) {
	channel, ok := i.logs[id]
	if !ok {
		return nil, fmt.Errorf("no installation process '%s' found", id)
	}

	return channel, nil
}
