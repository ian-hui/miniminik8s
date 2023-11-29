package runtime

import (
	"context"
	"minik8s/logger"
	"minik8s/minik8sTypes"
	"minik8s/pkg/apis"
	dockerclient "minik8s/pkg/kubelet/dockerClient"
	containermanager "minik8s/pkg/kubelet/runtime/containerManager"

	"github.com/docker/go-connections/nat"
)

var (
	K8sLogger = logger.K8sLogger
)

func NewRuntimeManager() (r RuntimeManager) {
	cm := containermanager.NewContainerManager(dockerclient.GetDockerClient())
	runtimeMnanger := runtimeManager{
		ContainerManager: cm,
	}
	r = runtimeMnanger
	return r
}

// 参照pkg/kubelet/kuberuntime/kuberuntime_sandbox.go
func (r runtimeManager) GenerateSandBoxConfig(pod *apis.Pod) (minik8sTypes.Config, error) {
	//标签修改
	podUID := string(pod.UID)
	podSandboxConfig := &apis.PodSandboxConfig{
		Metadata: &apis.PodSandboxMetadata{
			Name: pod.Name,
			Uid:  podUID,
		},
		Labels:      newPodLabels(pod),
		Annotations: newPodAnnotations(pod),
	}
	//todo：dns配置
	//hostname和domainname配置
	//pod映射配置
	sandboxExposePorts := nat.PortSet{}
	for _, c := range pod.Spec.Containers {
		err := containermanager.MakeContainerMapper(&c, &sandboxExposePorts)
		if err != nil {
			K8sLogger.Errorln("GenerateSandBoxConfig error: ", err)
			return minik8sTypes.Config{}, err
		}
	}
	//组合成docker的config
	config := minik8sTypes.Config{
		Image:        minik8sTypes.Minik8sPauseImage,
		Labels:       podSandboxConfig.Labels,
		ExposedPorts: sandboxExposePorts,
	}
	return config, nil
}

// 添加pod的标签
func newPodLabels(pod *apis.Pod) map[string]string {
	labels := map[string]string{}

	// Get labels from v1.Pod
	for k, v := range pod.Labels {
		labels[k] = v
	}

	labels[minik8sTypes.KubernetesPodNameLabel] = pod.Name
	labels[minik8sTypes.KubernetesPodNamespaceLabel] = pod.Namespace
	labels[minik8sTypes.KubernetesPodUIDLabel] = string(pod.UID)

	return labels
}

func newPodAnnotations(pod *apis.Pod) map[string]string {
	annotations := map[string]string{}
	// Get annotations from v1.Pod
	for k, v := range pod.Annotations {
		annotations[k] = v
	}
	return annotations
}

func (r runtimeManager) CreatePodSandbox(pod *apis.Pod) (string, error) {
	//创建一个容器管理器对象
	cm := r.ContainerManager
	//生成沙箱配置
	config, err := r.GenerateSandBoxConfig(pod)
	if err != nil {
		K8sLogger.Errorln("CreateSandbox error: ", err)
		return "", err
	}
	//创建一个容器的配置对象
	hostcfg := &minik8sTypes.HostConfig{}
	containerName := pod.Name
	//创建一个容器
	ctx := context.Background()
	ID, err := cm.NewContainer(ctx, &config, hostcfg, containerName)
	if err != nil {
		K8sLogger.Errorln("CreateSandbox error: ", err)
		return "", err
	}
	return ID, nil

}
