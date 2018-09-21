package tasks

import (
	"fmt"

	"github.com/kubermatic/kubermatic-installer/pkg/shared"
)

type ValidateManifestTask struct {
	shared.BaseTask
}

func (t *ValidateManifestTask) Execute(ctx *shared.Context) error {
	if err := ctx.Manifest.Validate(); err != nil {
		return fmt.Errorf("manifest is invalid: %v", err)
	}

	return nil
}
