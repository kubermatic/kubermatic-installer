package command

import (
	"fmt"
	"io"
	"os"

	"github.com/kubermatic/kubermatic-installer/pkg/helm/migration"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

func MigrateValuesCommand(logger *logrus.Logger) cli.Command {
	return cli.Command{
		Name:      "migrate-values",
		Usage:     "Upgrades a Helm values.yaml to a newer Kubermatic version.",
		Action:    MigrateValuesAction(logger),
		ArgsUsage: "[YAML_FILE]",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "master",
				Usage: "whether the given values are for a Kubermatic master cluster",
			},
			cli.StringFlag{
				Name:  "from",
				Usage: "verstion to migrate from (e.g. '2.8')",
			},
			cli.StringFlag{
				Name:  "to",
				Usage: "verstion to migrate to (e.g. '2.9')",
			},
		},
	}
}

func MigrateValuesAction(logger *logrus.Logger) cli.ActionFunc {
	return handleErrors(logger, setupLogger(logger, func(ctx *cli.Context) error {
		input := ctx.Args().First()

		var source io.ReadCloser
		if len(input) > 0 {
			file, err := os.Open(input)
			if err != nil {
				return fmt.Errorf("could not open input file %s: %v", input, err)
			}

			source = file
		} else {
			source = os.Stdin
		}
		defer source.Close() // nolint: errcheck

		// We use a yaml.MapSlice because it preserves the order of the item during
		// decode-encode. This results in the output being identical to input except
		// comments and empty lines being stripped.
		var (
			values yaml.MapSlice
			err    error
		)

		if err := yaml.NewDecoder(source).Decode(&values); err != nil {
			return fmt.Errorf("failed to decode input YAML: %v", err)
		}

		isMaster := ctx.Bool("master")
		from := ctx.String("from")
		to := ctx.String("to")

		if values, err = migration.Migrate(values, isMaster, from, to, logger); err != nil {
			return fmt.Errorf("migration failed: %v", err)
		}

		if err := yaml.NewEncoder(os.Stdout).Encode(values); err != nil {
			return fmt.Errorf("failed to decode output YAML: %v", err)
		}

		return nil
	}))
}
