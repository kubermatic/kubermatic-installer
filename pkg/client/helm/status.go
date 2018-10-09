package helm

type releaseStatus int

const (
	releaseCheckFailed releaseStatus = -1

	// these constants mirror the Helm status from
	// https://github.com/helm/helm/blob/master/_proto/hapi/release/status.proto

	releaseUnknown        releaseStatus = 0
	releaseDeployed       releaseStatus = 1
	releaseDeleted        releaseStatus = 2
	releaseSuperseded     releaseStatus = 3
	releaseFailed         releaseStatus = 4
	releaseDeleting       releaseStatus = 5
	releasePendingInstall releaseStatus = 6
	releasePendingUpgrade releaseStatus = 7
)

func (s releaseStatus) String() string {
	switch s {
	case releaseCheckFailed:
		return "CheckFailed"
	case releaseUnknown:
		return "Unknown"
	case releaseDeployed:
		return "Deployed"
	case releaseDeleted:
		return "Deleted"
	case releaseSuperseded:
		return "Superseded"
	case releaseFailed:
		return "Failed"
	case releaseDeleting:
		return "Deleting"
	case releasePendingInstall:
		return "PendingInstall"
	case releasePendingUpgrade:
		return "PendingUpgrade"
	}

	return "???"
}
