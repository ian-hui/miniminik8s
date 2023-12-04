package runtime

import (
	"context"
	"fmt"
	"minik8s/minik8sTypes"
	"minik8s/pkg/apis"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// 这里的startContainer 跟 k8s中的startContainer不一样
func (r *runtimeManager) startPodContainer(pod *apis.Pod, container apis.Container) error {
	filter := filters.NewArgs()
	filter.Add("label", minik8sTypes.KubernetesPodNameLabel+"="+pod.Name)
	filter.Add("label", minik8sTypes.KubernetesPodNamespaceLabel+"="+pod.Namespace)
	filter.Add("label", minik8sTypes.KubernetesPodUIDLabel+"="+string(pod.UID))
	filter.Add("label", minik8sTypes.Minik8sPodTypeLabel+"="+minik8sTypes.Minik8sGenericPodType) //普通容器的标签
	filter.Add("label", minik8sTypes.LabelsContainerName+"="+container.Name)                     //通过容器名字定位到一个容器
	res, err := r.ContainerManager.ListContainerWithOpts(context.TODO(), types.ContainerListOptions{
		All:     true,
		Filters: filter,
	})
	if err != nil {
		K8sLogger.Errorln("startContainer error: ", err)
		return err
	}
	for _, container := range res {
		err := r.ContainerManager.StartContainer(context.Background(), container.ID)
		if err != nil {
			K8sLogger.Errorln("startContainer error: ", err)
			return err
		}
	}
	return nil
}

func (r *runtimeManager) createPodContainer(pod *apis.Pod, container apis.Container, sandboxID string) error {
	//创建容器的配置
	config, hostConfig, err := r.generatePodContainerConfig(pod, container, sandboxID)
	if err != nil {
		K8sLogger.Errorln("createContainer error: ", err)
		return err
	}
	//创建容器
	ctx := context.Background()
	_, err = r.ContainerManager.NewContainer(ctx, &config, &hostConfig, container.Name)
	if err != nil {
		K8sLogger.Errorln("createContainer error: ", err)
		return err
	}
	return nil
}

func (r *runtimeManager) generatePodContainerConfig(pod *apis.Pod, container apis.Container, sandboxID string) (minik8sTypes.Config, minik8sTypes.HostConfig, error) {
	labels := map[string]string{}
	for k, v := range pod.Labels {
		labels[k] = v
	}
	labels[minik8sTypes.KubernetesPodNameLabel] = pod.Name
	labels[minik8sTypes.KubernetesPodNamespaceLabel] = pod.Namespace
	labels[minik8sTypes.KubernetesPodUIDLabel] = string(pod.UID)
	labels[minik8sTypes.Minik8sPodTypeLabel] = minik8sTypes.Minik8sGenericPodType //普通容器的标签
	labels[minik8sTypes.LabelsContainerName] = container.Name
	//env
	var containerEnv []string
	for _, env := range container.Env {
		containerEnv = append(containerEnv, env.Name+"="+env.Value)
	}
	//把pod级别的volume挂载到容器中
	bindings, err := r.bindingVolumeInPod(pod, &container)
	if err != nil {
		K8sLogger.Errorln("generateContainerConfig error: ", err)
		return minik8sTypes.Config{}, minik8sTypes.HostConfig{}, err
	}
	//生成容器配置
	config := minik8sTypes.Config{
		Image:           container.Image,
		Env:             containerEnv,
		ImagePullPolicy: container.ImagePullPolicy,
		Labels:          labels,
		Tty:             true,
		// Tty: container,
	}
	//生成容器host配置
	hostcfg := &minik8sTypes.HostConfig{
		Binds:            bindings,
		CPUResourceLimit: int64(container.Resources.Limits.Cpu),
		MemoryLimit:      int64(container.Resources.Limits.Memory),
	}
	return config, *hostcfg, nil
}

func (r *runtimeManager) bindingVolumeInPod(pod *apis.Pod, container *apis.Container) ([]string, error) {
	//把pod级别的volume挂载到容器中
	//先获取pod级别的volume
	volumeMap := map[string]*apis.HostVolume{}
	for _, volume := range pod.Spec.Volumes {
		if volume.Name != "" {
			volumeMap[volume.Name] = &volume
		}
	}
	var bind []string
	for _, containerVol := range container.VolumeMounts {
		if value, ok := volumeMap[containerVol.Name]; !ok {
			K8sLogger.Errorln("volume %s not found in pod", containerVol.Name)
			return nil, fmt.Errorf("volume %s not found in pod", containerVol.Name)
		} else {
			bind = append(bind, value.Path+":"+containerVol.MountPath)
		}
	}
	return bind, nil
}

func (r *runtimeManager) removePodContainer(pod *apis.Pod, container *apis.Container) (string, error) {

	filter := filters.NewArgs()

	// 根据标签过滤器，过滤出来所有的容器
	filter.Add("label", minik8sTypes.KubernetesPodNameLabel+"="+pod.Name)
	filter.Add("label", minik8sTypes.KubernetesPodNamespaceLabel+"="+pod.Namespace)
	filter.Add("label", minik8sTypes.KubernetesPodUIDLabel+"="+string(pod.UID))
	filter.Add("label", minik8sTypes.Minik8sPodTypeLabel+"="+minik8sTypes.Minik8sGenericPodType) //普通容器的标签
	filter.Add("label", minik8sTypes.LabelsContainerName+"="+container.Name)                     //通过容器名字定位到一个容器
	// 根据容器的名字过滤器，过滤出来所有的容器
	res, err := r.ContainerManager.ListContainerWithOpts(context.TODO(), types.ContainerListOptions{
		All:     true,
		Filters: filter,
	})

	if err != nil {
		return "", err
	}

	// 遍历所有的容器，然后删除
	for _, container := range res {
		err2 := r.ContainerManager.RemoveContainer(context.TODO(), container.ID)
		if err2 != nil {
			K8sLogger.Errorln("removeContainer error: ", err2)
			return "", err2
		}
	}

	return "", nil
}
