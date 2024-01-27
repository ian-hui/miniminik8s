package runtime

import (
	"fmt"
	"minik8s/minik8sTypes"
	"minik8s/pkg/apis"
	dockerclient "minik8s/pkg/kubelet/dockerClient"
	containermanager "minik8s/pkg/kubelet/runtime/containerManager"
	imagemanager "minik8s/pkg/kubelet/runtime/imageManager"
	"sync"
)

type RuntimeManager interface {
	createPod(pod *apis.Pod) (string, error)
	generateSandBoxConfig(pod *apis.Pod) (minik8sTypes.Config, minik8sTypes.HostConfig, error)
	createPodSandbox(pod *apis.Pod) (string, error)
	removePodContainer(*apis.Pod, *apis.Container) (string, error)
	startPodContainer(*apis.Pod, apis.Container) error
	createPodContainer(*apis.Pod, apis.Container, string) error
	generatePodContainerConfig(*apis.Pod, apis.Container, string) (minik8sTypes.Config, minik8sTypes.HostConfig, error)
	killPod(pod *apis.Pod) error
	// getPodSandbox(pod *apis.Pod) (*apis.PodSandbox, error)
	// getPodSandboxes() ([]*apis.PodSandbox, error)
	// getPodSandboxStatus(pod *apis.Pod) (*apis.PodSandboxStatus, error)
}

type runtimeManager struct {
	containerManager *containermanager.ContainerManager
	imagemanager     *imagemanager.ImageManager
}

func NewRuntimeManager() (r RuntimeManager) {
	cm := containermanager.NewContainerManager(dockerclient.GetDockerClient())
	im := imagemanager.NewImageManager(dockerclient.GetDockerClient())
	runtimeMnanger := &runtimeManager{
		containerManager: cm,
		imagemanager:     im,
	}
	r = runtimeMnanger
	return
}

// 创建pod
func (r *runtimeManager) createPod(pod *apis.Pod) (string, error) {
	s, err := r.createPodSandbox(pod)
	if err != nil {
		K8sLogger.Errorln("createPodSandbox error: ", err)
		return "", err
	}
	// 依次创建pod中所有的容器
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

func (r *runtimeManager) killPod(pod *apis.Pod) error {
	// 删除pod中所有的容器
	wg := sync.WaitGroup{}
	wg.Add(len(pod.Spec.Containers))
	errChan := make(chan error, len(pod.Spec.Containers)) // 创建一个错误通道
	// 依次删除pod中所有的容器
	for _, container := range pod.Spec.Containers {
		go func(container apis.Container) {
			defer wg.Done()
			// 删除容器
			fmt.Println("container: ", container)
			s, err := r.removePodContainer(pod, &container)
			if err != nil {
				K8sLogger.Errorln("removePodContainer error: ", err)
				errChan <- fmt.Errorf("removePodContainer error: %v", err)
			}
			K8sLogger.Infoln("removePodContainer success: ", s)
		}(container)
	}
	wg.Wait()
	close(errChan) // 关闭通道

	//如果err通道只有一个错误
	if len(errChan) == 1 {
		return fmt.Errorf("removePodContainer one container error: %v", <-errChan)
	}
	//如果err通道有多个错误
	for err := range errChan {
		return fmt.Errorf("removePodContainer multi container error: %v", err)
	}

	// 删除pod沙箱容器
	err := r.removePodSandbox(pod)
	if err != nil {
		K8sLogger.Errorln("removePodSandbox error: ", err)
		return err
	}
	return nil
}
