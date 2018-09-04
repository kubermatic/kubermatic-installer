package tasks

import (
	"fmt"
	"github.com/kubermatic/kubermatic-installer/pkg/shared"
)

const TASK_VALIDATE_VERSION = "validate-version"

func ValidateVersion(ctx *shared.Context) error {
	if ctx.Manifest.Version != shared.INSTALLER_VERSION {
		return fmt.Errorf("version mismatch in manifest, expected %s got %s", shared.INSTALLER_VERSION, ctx.Manifest.Version)
	}

	return nil
}
