package migration

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/kubermatic/kubermatic-installer/pkg/helm/migration/v2_8"
	"github.com/kubermatic/kubermatic-installer/pkg/helm/migration/v2_9"
	"github.com/kubermatic/kubermatic-installer/pkg/yamled"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type conversion struct {
	from      string
	to        string
	converter converter
}

func Migrate(values *yaml.MapSlice, isMaster bool, from string, to string, logger logrus.FieldLogger) error {
	conversions, err := getConversions(from, to, logger)
	if err != nil {
		return fmt.Errorf("failed to determine required migration steps: %v", err)
	}

	if len(conversions) == 0 {
		return fmt.Errorf("migration would be a no-op")
	}

	document, err := yamled.NewFromMapSlice(values)
	if err != nil {
		return fmt.Errorf("failed to prepare YAML document for editing: %v", err)
	}

	for _, conversion := range conversions {
		logger.Infof("Converting from %s to %s...", conversion.from, conversion.to)

		err := conversion.converter.Convert(document, isMaster)
		if err != nil {
			return err
		}
	}

	return nil
}

func getConversions(from string, to string, logger logrus.FieldLogger) ([]conversion, error) {
	converters := make([]conversion, 0)

	for from != to {
		v, err := semver.NewVersion(from)
		if err != nil {
			return converters, fmt.Errorf("could not parse version: %v", err)
		}

		from = fmt.Sprintf("%d.%d", v.Major(), v.Minor())
		next := ""

		var converter converter

		switch from {
		case "2.7":
			converter = v2_8.NewConverter(logger)
			next = "2.8"

		case "2.8":
			converter = v2_9.NewConverter(logger)
			next = "2.9"

		default:
			return converters, fmt.Errorf("unrecognized version %s", from)
		}

		converters = append(converters, conversion{
			from:      from,
			to:        next,
			converter: converter,
		})
		from = next
	}

	return converters, nil
}
