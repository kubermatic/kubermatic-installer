package kubermatic

import (
	"fmt"

	"github.com/kubermatic/kubermatic-installer/pkg/installer/state"
	"github.com/kubermatic/kubermatic-installer/pkg/installer/task"
	"github.com/kubermatic/kubermatic-installer/pkg/manifest"
)

const (
	NginxIngressControllerChartName   = "nginx-ingress-controller"
	NginxIngressControllerReleaseName = NginxIngressControllerChartName
	NginxIngressControllerNamespace   = NginxIngressControllerChartName

	CertManagerChartName   = "cert-manager"
	CertManagerReleaseName = CertManagerChartName
	CertManagerNamespace   = CertManagerChartName

	DexChartName   = "oauth"
	DexReleaseName = DexChartName
	DexNamespace   = DexChartName

	KubermaticOperatorChartName   = "kubermatic-operator"
	KubermaticOperatorReleaseName = KubermaticOperatorChartName
	KubermaticOperatorNamespace   = "kubermatic"

	StorageClassName = "kubermatic-fast"
)

func DeploymentTasks(installerState *state.InstallerState, clusterState *state.ClusterState) ([]task.Task, error) {
	tasks := []task.Task{}

	tasks, err := planStorageClass(tasks, installerState, clusterState)
	if err != nil {
		return tasks, fmt.Errorf("failed to plan StorageClass: %v", err)
	}

	tasks, err = planCertManager(tasks, installerState, clusterState)
	if err != nil {
		return tasks, fmt.Errorf("failed to plan cert-manager: %v", err)
	}

	tasks, err = planNginxIngressController(tasks, installerState, clusterState)
	if err != nil {
		return tasks, fmt.Errorf("failed to plan nginx-ingress-controller: %v", err)
	}

	tasks, err = planDex(tasks, installerState, clusterState)
	if err != nil {
		return tasks, fmt.Errorf("failed to plan oauth: %v", err)
	}

	// tasks, err = planKubermaticOperator(tasks, installerState, clusterState)
	// if err != nil {
	// 	return tasks, fmt.Errorf("failed to plan kubermatic-operator: %v", err)
	// }

	return tasks, nil
}

func planStorageClass(tasks []task.Task, installerState *state.InstallerState, clusterState *state.ClusterState) ([]task.Task, error) {
	// The check for existence is not because the EnsureStorageClassTask cannot handle existing
	// StorageClasses, but to prevent us from attempting to always create a StorageClass,
	// which would always fail for custom environments.
	if !clusterState.HasStorageClass(StorageClassName) {
		sc := storageClassForProvider(StorageClassName, manifest.ProviderGKE)
		if sc == nil {
			return tasks, fmt.Errorf("cannot automatically create StorageClass '%s' for this cloud provider, please create it manually.", StorageClassName)
		}

		tasks = append(tasks, &task.EnsureStorageClassTask{
			StorageClass: sc,
		})
	}

	return tasks, nil
}

func planCertManager(tasks []task.Task, installerState *state.InstallerState, clusterState *state.ClusterState) ([]task.Task, error) {
	tasks = append(tasks, &task.EnsureHelmReleaseTask{
		ChartName:   CertManagerChartName,
		ReleaseName: CertManagerReleaseName,
		Namespace:   CertManagerNamespace,
	})

	return tasks, nil
}

func planNginxIngressController(tasks []task.Task, installerState *state.InstallerState, clusterState *state.ClusterState) ([]task.Task, error) {
	tasks = append(tasks, &task.EnsureHelmReleaseTask{
		ChartName:   NginxIngressControllerChartName,
		ReleaseName: NginxIngressControllerReleaseName,
		Namespace:   NginxIngressControllerNamespace,
	})

	return tasks, nil
}

func planDex(tasks []task.Task, installerState *state.InstallerState, clusterState *state.ClusterState) ([]task.Task, error) {
	tasks = append(tasks, &task.EnsureHelmReleaseTask{
		ChartName:   DexChartName,
		ReleaseName: DexReleaseName,
		Namespace:   DexNamespace,
	})

	return tasks, nil
}

func planKubermaticOperator(tasks []task.Task, installerState *state.InstallerState, clusterState *state.ClusterState) ([]task.Task, error) {
	tasks = append(tasks, &task.EnsureHelmReleaseTask{
		ChartName:   KubermaticOperatorChartName,
		ReleaseName: KubermaticOperatorReleaseName,
		Namespace:   KubermaticOperatorNamespace,
	})

	return tasks, nil
}
