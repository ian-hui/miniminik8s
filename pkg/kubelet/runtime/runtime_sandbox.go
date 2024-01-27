package runtime

import (
	"context"
	"minik8s/logger"
	"minik8s/minik8sTypes"
	"minik8s/pkg/apis"
	containerManager "minik8s/pkg/kubelet/runtime/containerManager"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/go-connections/nat"
)

// -----------------------------------------------------
// 这个文件主要处理的是pod中的pause容器（也就是k8s官方所说的沙箱容器（sandbox））的操作
// -----------------------------------------------------

var (
	K8sLogger = logger.K8sLogger
)

// 目的是生成一个pause容器的配置
// 参照pkg/kubelet/kuberuntime/kuberuntime_sandbox.go
func (r *runtimeManager) generateSandBoxConfig(pod *apis.Pod) (minik8sTypes.Config, minik8sTypes.HostConfig, error) {
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
		err := containerManager.MakeContainerMapper(&c, &sandboxExposePorts)
		if err != nil {
			K8sLogger.Errorln("GenerateSandBoxConfig error: ", err)
			return minik8sTypes.Config{}, minik8sTypes.HostConfig{}, err
		}
	}
	//组合成docker的config
	config := minik8sTypes.Config{
		Image:           minik8sTypes.Minik8sPauseImage,
		Labels:          podSandboxConfig.Labels,
		ExposedPorts:    sandboxExposePorts,
		ImagePullPolicy: minik8sTypes.IfNotPresent,
	}
	//组合成docker的hostconfig
	hostConfig := minik8sTypes.HostConfig{
		IpcMode: minik8sTypes.IpcModeShareable,
	}
	return config, hostConfig, nil
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
	labels[minik8sTypes.Minik8sPodTypeLabel] = minik8sTypes.Minik8sPausePodType //这个标签是为了区分pause容器和其他容器
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

// 创建一个沙箱返回一个pause容器id
// 参照pkg/kubelet/kuberuntime/kuberuntime_sandbox.go
func (r *runtimeManager) createPodSandbox(pod *apis.Pod) (SandboxContainerName string, err error) {
	//创建一个容器管理器对象
	cm := r.containerManager
	ctx := context.Background()

	//生成沙箱配置
	config, hostcfg, err := r.generateSandBoxConfig(pod)
	if err != nil {
		K8sLogger.Errorln("CreateSandbox error: ", err)
		return "", err
	}
	//创建一个容器的配置对象
	SandboxContainerName = pod.Name + pod.UID
	//拉取pause镜像
	err = r.imagemanager.PullImage(ctx, minik8sTypes.IfNotPresent, minik8sTypes.Minik8sPauseImage)
	if err != nil {
		K8sLogger.Errorln("CreateSandbox error: ", err)
		return "", err
	}
	//创建一个容器
	ID, err := cm.NewContainer(ctx, &config, &hostcfg, SandboxContainerName)
	if err != nil {
		K8sLogger.Errorln("CreateSandbox error: ", err)
		return "", err
	}
	//启动容器
	err = cm.StartContainer(ctx, ID)
	if err != nil {
		K8sLogger.Errorln("CreateSandbox error: ", err)
		return "", err
	}
	return
}

// 删除pod中的sandbox
func (r *runtimeManager) removePodSandbox(pod *apis.Pod) error {
	filter := filters.NewArgs()
	// 在filter中添加标签
	// 四个标签：PodName、PodNamespace、PodUID、IfPause
	filter.Add("label", minik8sTypes.KubernetesPodNameLabel+"="+pod.Name)
	filter.Add("label", minik8sTypes.KubernetesPodNamespaceLabel+"="+pod.Namespace)
	filter.Add("label", minik8sTypes.KubernetesPodUIDLabel+"="+string(pod.UID))
	filter.Add("label", minik8sTypes.Minik8sPodTypeLabel+"="+minik8sTypes.Minik8sPausePodType)
	c, err := r.containerManager.ListContainerWithOpts(context.TODO(), types.ContainerListOptions{
		All:     true,
		Filters: filter,
	})
	if err != nil {
		K8sLogger.Errorln("StopPodSandbox error: ", err)
		return err
	}
	for _, container := range c {
		err := r.containerManager.RemoveContainer(context.Background(), container.ID)
		if err != nil {
			K8sLogger.Errorln("StopPodSandbox error: ", err)
			return err
		}
	}
	return nil
}
