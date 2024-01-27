package containermanager

import (
	"context"
	"fmt"
	"minik8s/minik8sTypes"
	dockerclient "minik8s/pkg/kubelet/dockerClient"
	"testing"
)

func TestMain(m *testing.M) {
	c, err := dockerclient.NewDockerClient()
	if err != nil {
		panic(err)
	}
	defer c.Close()
	m.Run()
}

func TestListALl(t *testing.T) {
	cm := NewContainerManager(dockerclient.GetDockerClient())
	ctx := context.Background()
	containers, err := cm.ListALlContainer(ctx)
	if err != nil {
		t.Error(err)
	}
	for _, container := range containers {
		fmt.Println(container.Image)
	}
}

func TestAddContainer(t *testing.T) {
	cm := NewContainerManager(dockerclient.GetDockerClient())
	ctx := context.Background()
	err := cm.StartContainer(ctx, "nginx:latest")
	if err != nil {
		t.Error(err)
	}
}

func TestCreateContainer(t *testing.T) {

	// 创建一个容器管理器对象
	cm := NewContainerManager(dockerclient.GetDockerClient())

	// 依次创建容器
	// 创建一个容器的配置对象
	option := &minik8sTypes.Config{
		Image:           "nginx:latest",
		Env:             []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
		ImagePullPolicy: minik8sTypes.IfNotPresent,
		Labels:          map[string]string{"test": "test"},
	}
	hostcfg := &minik8sTypes.HostConfig{}
	containerName := "test"
	ctx := context.Background()
	ID, err := cm.NewContainer(ctx, option, hostcfg, containerName)
	if err != nil {
		t.Error(err)
	}
	cm.StartContainer(ctx, ID)
	t.Logf("[Created] Container ID: %s", ID)

}

func TestRemoveContainer(t *testing.T) {
	cm := NewContainerManager(dockerclient.GetDockerClient())
	ctx := context.Background()
	err := cm.RemoveContainer(ctx, "test")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("remove success")
}

func TestInspectContainer(t *testing.T) {
	cm := NewContainerManager(dockerclient.GetDockerClient())
	ctx := context.Background()
	c, err := cm.ListALlContainer(ctx)
	if err != nil {
		t.Error(err)
	}
	for _, container := range c {
		thecontainer, err := cm.InspectContainer(ctx, container.ID)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(thecontainer.ID)
	}
}

func TestContainerStat(t *testing.T) {
	cm := NewContainerManager(dockerclient.GetDockerClient())
	ctx := context.Background()
	c, err := cm.ListALlContainer(ctx)
	if err != nil {
		t.Error(err)
	}
	for _, container := range c {
		stat, err := cm.ContainerStats(ctx, container.ID)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(stat.MemoryStats.Usage, stat.CPUStats)
	}
}

func TestGetContainerStatus(t *testing.T) {
	cm := NewContainerManager(dockerclient.GetDockerClient())
	ctx := context.Background()
	// c, err := cm.ListALlContainer(ctx)
	// if err != nil {
	// 	t.Error(err)
	// }
	// for _, container := range c {
	stat, err := cm.InspectContainer(ctx, "693b6913052c")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(stat.State)
	// }
}
