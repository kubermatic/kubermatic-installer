package state

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kubermatic/kubermatic-installer/pkg/client/helm"

	storagev1 "k8s.io/api/storage/v1"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ClusterState struct {
	KubernetesVersion string
	StorageClasses    []storagev1.StorageClass
	HelmReleases      []helm.Release
}

func NewClusterState(ctx context.Context, kubeClient ctrlruntimeclient.Client, helmClient helm.Client) (*ClusterState, error) {
	classes := storagev1.StorageClassList{}
	if err := kubeClient.List(ctx, &classes); err != nil {
		return nil, fmt.Errorf("failed to determine storage classes: %v", err)
	}

	releases, err := helmClient.ListReleases("") // "" = all namespaces
	if err != nil {
		return nil, fmt.Errorf("failed to list Helm releases: %v", err)
	}

	clusterState := &ClusterState{
		KubernetesVersion: "TODO",
		StorageClasses:    classes.Items,
		HelmReleases:      releases,
	}

	return clusterState, nil
}

func (s *ClusterState) Clone() ClusterState {
	result := ClusterState{
		KubernetesVersion: s.KubernetesVersion,
		StorageClasses:    []storagev1.StorageClass{},
		HelmReleases:      []helm.Release{},
	}

	for _, sc := range s.StorageClasses {
		copy := sc.DeepCopy()
		result.StorageClasses = append(result.StorageClasses, *copy)
	}

	for _, release := range s.HelmReleases {
		result.HelmReleases = append(result.HelmReleases, release.Clone())
	}

	return result
}

func (s *ClusterState) ReleasesByName(name string, namespace string) []helm.Release {
	result := []helm.Release{}

	for _, r := range s.HelmReleases {
		if r.Name == name && (namespace == "" || r.Namespace == namespace) {
			result = append(result, r)
		}
	}

	return result
}

func (s *ClusterState) ReleasesByChart(chart string, namespace string) []helm.Release {
	result := []helm.Release{}

	for _, r := range s.HelmReleases {
		if r.Chart == chart && (namespace == "" || r.Namespace == namespace) {
			result = append(result, r)
		}
	}

	return result
}

func (s *ClusterState) HasChart(chart string, namespace string) bool {
	return len(s.ReleasesByChart(chart, namespace)) > 0
}

func (s *ClusterState) HasRelease(name string, namespace string) bool {
	return len(s.ReleasesByName(name, namespace)) > 0
}

func (s *ClusterState) Release(name string, namespace string) *helm.Release {
	releases := s.ReleasesByName(name, namespace)
	if len(releases) == 0 {
		return nil
	}

	return &releases[0]
}

func (s *ClusterState) HasStorageClass(name string) bool {
	for _, s := range s.StorageClasses {
		if s.Name == name {
			return true
		}
	}

	return false
}

func (s *ClusterState) UpdateRelease(name string, namespace string, chart *helm.Chart) {
	for idx, r := range s.HelmReleases {
		if r.Name == name && (namespace == "" || r.Namespace == namespace) {
			r.AppVersion = chart.AppVersion
			r.Version = chart.Version

			rev, err := strconv.Atoi(r.Revision)
			if err == nil {
				r.Revision = strconv.Itoa(rev + 1)
			} else {
				r.Revision = "1"
			}

			s.HelmReleases[idx] = r
		}
	}
}
