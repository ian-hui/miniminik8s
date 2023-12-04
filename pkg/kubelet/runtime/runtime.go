package runtime

import (
	"minik8s/minik8sTypes"
	"minik8s/pkg/apis"
	dockerclient "minik8s/pkg/kubelet/dockerClient"
	containermanager "minik8s/pkg/kubelet/runtime/containerManager"
)

type RuntimeManager interface {
	createPod(pod *apis.Pod) (string, error)
	generateSandBoxConfig(pod *apis.Pod) (minik8sTypes.Config, error)
	createPodSandbox(pod *apis.Pod) (string, error)
	removePodContainer(*apis.Pod, *apis.Container) (string, error)
	startPodContainer(*apis.Pod, apis.Container) error
	createPodContainer(*apis.Pod, apis.Container, string) error
	generatePodContainerConfig(*apis.Pod, apis.Container, string) (minik8sTypes.Config, minik8sTypes.HostConfig, error)
	// getPodSandbox(pod *apis.Pod) (*apis.PodSandbox, error)
	// getPodSandboxes() ([]*apis.PodSandbox, error)
	// getPodSandboxStatus(pod *apis.Pod) (*apis.PodSandboxStatus, error)
}

type runtimeManager struct {
	ContainerManager *containermanager.ContainerManager
}

func NewRuntimeManager() (r RuntimeManager) {
	cm := containermanager.NewContainerManager(dockerclient.GetDockerClient())
	runtimeMnanger := &runtimeManager{
		ContainerManager: cm,
	}
	r = runtimeMnanger
	return r
}

func (r *runtimeManager) createPod(pod *apis.Pod) (string, error) {
	s, err := r.createPodSandbox(pod)
	if err != nil {
		K8sLogger.Errorln("createPodSandbox error: ", err)
		return "", err
	}
	for _, container := range pod.Spec.Containers {
		// 创建容器
		err := r.createPodContainer(pod, container, s)
		if err != nil {
			K8sLogger.Errorln("createPodContainer error: ", err)
			return "", err
		}
		// 启动容器
		err = r.startPodContainer(pod, container)
		if err != nil {
			K8sLogger.Errorln("startPodContainer error: ", err)
			return "", err
		}
	}
	return s, nil
}
