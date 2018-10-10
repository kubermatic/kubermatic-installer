package manifest

type CloudProvider string

const (
	ProviderEKS CloudProvider = "aws-eks"
	ProviderGKE CloudProvider = "google-gke"
	ProviderAKS CloudProvider = "azure-aks"
)

var AllProviders = []CloudProvider{ProviderAKS, ProviderEKS, ProviderGKE}
