package migration

import "github.com/kubermatic/kubermatic-installer/pkg/yamled"

type converter interface {
	Convert(v *yamled.Document, isMaster bool) error
}
