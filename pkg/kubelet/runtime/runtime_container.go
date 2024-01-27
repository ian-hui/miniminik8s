package runtime

import (
	"context"
	"fmt"
	"minik8s/minik8sTypes"
	"minik8s/pkg/apis"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// -----------------------------------------------------
// 这个文件主要处理的是pod中的非pause容器的操作
// -----------------------------------------------------

// 这里的startContainer 跟 k8s中的startContainer不一样
func (r *runtimeManager) startPodContainer(pod *apis.Pod, container apis.Container) error {
	filter := filters.NewArgs()
	filter.Add("label", minik8sTypes.KubernetesPodNameLabel+"="+pod.Name)
	filter.Add("label", minik8sTypes.KubernetesPodNamespaceLabel+"="+pod.Namespace)
	filter.Add("label", minik8sTypes.KubernetesPodUIDLabel+"="+string(pod.UID))
	filter.Add("label", minik8sTypes.Minik8sPodTypeLabel+"="+minik8sTypes.Minik8sGenericPodType) //普通容器的标签
	filter.Add("label", minik8sTypes.LabelsContainerName+"="+container.Name)                     //通过容器名字定位到一个容器
	res, err := r.containerManager.ListContainerWithOpts(context.TODO(), types.ContainerListOptions{
		All:     true,
		Filters: filter,
	})
	if err != nil {
		K8sLogger.Errorln("startContainer error: ", err)
		return err
	}
	for _, container := range res {
		err := r.containerManager.StartContainer(context.Background(), container.ID)
		if err != nil {
			K8sLogger.Errorln("startContainer error: ", err)
			return err
		}
	}
	return nil
}

func (r *runtimeManager) createPodContainer(pod *apis.Pod, container apis.Container, sandboxName string) error {
	//拉容器
	err := r.imagemanager.PullImage(context.Background(), container.ImagePullPolicy, container.Image)
	if err != nil {
		K8sLogger.Errorln("pullImage error: ", err)
		return err
	}
	//创建容器的配置
	config, hostConfig, err := r.generatePodContainerConfig(pod, container, sandboxName)
	if err != nil {
		K8sLogger.Errorln("createContainer error: ", err)
		return err
	}
	//创建容器的上下文
	ctx := context.Background()
	//创建容器
	_, err = r.containerManager.NewContainer(ctx, &config, &hostConfig, container.Name)
	if err != nil {
		K8sLogger.Errorln("createContainer error: ", err)
		return err
	}
	return nil
}

// 一个pod中的单个容器的配置（非pause容器）
func (r *runtimeManager) generatePodContainerConfig(pod *apis.Pod, container apis.Container, sandboxName string) (minik8sTypes.Config, minik8sTypes.HostConfig, error) {
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
	hostcfg := minik8sTypes.HostConfig{
		Binds: bindings,
		// 容器的ns加入到pause容器的ns中，仔细阅读 https://k8s.iswbm.com/c02/p02_learn-kubernetes-pod-via-pause-container.html
		NetworkMode:      minik8sTypes.NsModeContainerPrefix + sandboxName,
		IpcMode:          minik8sTypes.NsModeContainerPrefix + sandboxName,
		PidMode:          minik8sTypes.NsModeContainerPrefix + sandboxName,
		CPUResourceLimit: int64(container.Resources.Limits.Cpu),
		MemoryLimit:      int64(container.Resources.Limits.Memory),
	}
	return config, hostcfg, nil
}

func (r *runtimeManager) bindingVolumeInPod(pod *apis.Pod, container *apis.Container) ([]string, error) {
	//把pod级别的volume挂载到容器中
	//先获取pod级别的volume

	//把pod中的volume找出来，类似
	/*
			 	volumes:
		  		- name: my-vol
		    	  hostPath:
		      		path: /path/on/host
	*/
	volumeMap := map[string]*apis.HostVolume{}
	for _, volume := range pod.Spec.Volumes {
		if volume.Name != "" {
			volumeMap[volume.Name] = &volume
		}
	}
	// 代码检查pod中容器的 volumeMounts 配置。如果发现容器希望将 my-vol 卷挂载在 /path/in/container。则生成挂载字符串：由于 my-vol 存在于 volumeMap 中，
	// 代码将创建一个挂载字符串："/path/on/host:/path/in/container"。这个字符串表示宿主机的 /path/on/host 目录将挂载到容器的 /path/in/container 目录。
	var bind []string
	for _, containerVol := range container.VolumeMounts {
		// 如果容器中的volumemount的名字和pod中的volumes的名字一样，那么就把宿主机的路径和容器的路径绑定起来
		if value, ok := volumeMap[containerVol.Name]; !ok {
			K8sLogger.Errorln("volume %s not found in pod", containerVol.Name)
			return nil, fmt.Errorf("volume %s not found in pod", containerVol.Name)
		} else {
			//把宿主机的路径和容器的路径绑定起来
			// "/path/on/host:/path/in/container"
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
	res, err := r.containerManager.ListContainerWithOpts(context.TODO(), types.ContainerListOptions{
		All:     true,
		Filters: filter,
	})

	if err != nil {
		return "", err
	}
	if len(res) == 0 {
		return "", fmt.Errorf("container not found")
	}
	// 遍历所有的容器，然后删除
	for _, container := range res {
		err2 := r.containerManager.RemoveContainer(context.TODO(), container.ID)
		if err2 != nil {
			K8sLogger.Errorln("removeContainer error: ", err2)
			return "", err2
		}
	}

	return "", nil
}
