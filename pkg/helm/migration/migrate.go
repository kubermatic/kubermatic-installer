package migration

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/kubermatic/kubermatic-installer/pkg/helm/migration/v2_8"
	"github.com/kubermatic/kubermatic-installer/pkg/helm/migration/v2_9"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

func Migrate(values *yaml.MapSlice, isMaster bool, from string, to string, logger logrus.FieldLogger) error {
	converters, err := getConvertors(from, to, logger)
	if err != nil {
		return fmt.Errorf("failed to determine required migration steps: %v", err)
	}

	fmt.Println(converters)

	return nil
}

func getConvertors(from string, to string, logger logrus.FieldLogger) ([]converter, error) {
	converters := make([]converter, 0)

	for from != to {
		v, err := semver.NewVersion(from)
		if err != nil {
			return converters, fmt.Errorf("could not parse version: %v", err)
		}

		switch fmt.Sprintf("%d.%d", v.Major(), v.Minor()) {
		case "2.7":
			converters = append(converters, v2_8.NewConverter(logger))
			from = "2.8"

		case "2.8":
			converters = append(converters, v2_9.NewConverter(logger))
			from = "2.9"

		default:
			return converters, fmt.Errorf("unrecognized version %s", from)
		}
	}

	return converters, nil
}
