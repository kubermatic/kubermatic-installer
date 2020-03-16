package helm

type releaseStatus string

const (
	releaseCheckFailed releaseStatus = ""

	// these constants mirror the Helm status from
	// `helm status --help`

	releaseUnknown         releaseStatus = "unknown"
	releaseDeployed        releaseStatus = "deployed"
	releaseDeleted         releaseStatus = "uninstalled"
	releaseSuperseded      releaseStatus = "superseded"
	releaseFailed          releaseStatus = "failed"
	releaseDeleting        releaseStatus = "uninstalling"
	releasePendingInstall  releaseStatus = "pending-install"
	releasePendingUpgrade  releaseStatus = "pending-upgrade"
	releasePendingRollback releaseStatus = "pending-rollback"
)
