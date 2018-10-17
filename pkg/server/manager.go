package server

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kubermatic/kubermatic-installer/pkg/client/kubernetes"
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

type resultLogItem struct {
	Type              string               `json:"type"`
	HelmValues        string               `json:"helmValues"`
	NginxIngresses    []kubernetes.Ingress `json:"nginxIngresses"`
	NodeportIngresses []kubernetes.Ingress `json:"nodeportIngresses"`
}

type installPhaseBuilder func(installer.InstallerOptions, *manifest.Manifest, *logrus.Logger) installer.Installer

func (i *installManager) Start(m manifest.Manifest, builder installPhaseBuilder) (string, error) {
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

		result, err := builder(options, &m, logger).Run()
		if err != nil {
			logger.Errorf("Installation failed: %v", err)
		}

		// send out the install result to allow the user to download
		// the values.yaml and show helpful information about DNS settings
		item := resultLogItem{
			Type:              "result",
			HelmValues:        string(result.HelmValues.YAML()),
			NginxIngresses:    result.NginxIngresses,
			NodeportIngresses: result.NodeportIngresses,
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
