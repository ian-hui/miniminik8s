package runtime

import (
	"minik8s/minik8sTypes"
	"minik8s/pkg/apis"
	containermanager "minik8s/pkg/kubelet/runtime/containerManager"
)

type RuntimeManager interface {
	GenerateSandBoxConfig(pod *apis.Pod) (minik8sTypes.Config, error)
	CreatePodSandbox(pod *apis.Pod) (string, error)
	// DeletePodSandbox(pod *apis.Pod) error
	// getPodSandbox(pod *apis.Pod) (*apis.PodSandbox, error)
	// getPodSandboxes() ([]*apis.PodSandbox, error)
	// getPodSandboxStatus(pod *apis.Pod) (*apis.PodSandboxStatus, error)
}

type runtimeManager struct {
	ContainerManager *containermanager.ContainerManager
}
