package containermanager

import (
	"context"
	"encoding/json"
	"io"
	"minik8s/logger"
	"minik8s/minik8sTypes"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

/*
	这个文件是对docker容器的一些操作的封装
	比如创建容器、删除容器、获取容器状态等
*/

var (
	K8sLogger = logger.K8sLogger
)

type ContainerManagerInterface interface {
	NewContainer(ctx context.Context, config *minik8sTypes.Config, hostConfig *minik8sTypes.HostConfig, containerName string) (dockerID string, err error)
	StartContainer(ctx context.Context, dockerID string) error
	StopContainer(ctx context.Context, dockerID string) error
	RemoveContainer(ctx context.Context, dockerID string) error
	ListMinik8sContainer(ctx context.Context) ([]types.Container, error)
	ListALlContainer(ctx context.Context) ([]types.Container, error)
	InspectContainer(ctx context.Context, dockerID string) (types.ContainerJSON, error)
	GetContainerLogs(ctx context.Context, dockerID string) (io.ReadCloser, error)
	ContainerStats(ctx context.Context, dockerID string) (*types.StatsJSON, error)
	RestartContainer(ctx context.Context, dockerID string) error
}

type ContainerManager struct {
	client *client.Client
}

func NewContainerManager(c *client.Client) *ContainerManager {
	return &ContainerManager{
		client: c,
	}
}

func (cm *ContainerManager) NewContainer(ctx context.Context, config *minik8sTypes.Config, hostConfig *minik8sTypes.HostConfig, containerName string) (dockerID string, err error) {
	//为了和在宿主机上跑的docker分开来，我们给label中加上一个标识
	config.Labels[string(minik8sTypes.RunningSystemMinik8s)] = minik8sTypes.IsTrue
	//创建一个容器
	resp, err := cm.client.ContainerCreate(ctx,
		&container.Config{
			Tty:          config.Tty,
			Env:          config.Env,
			Cmd:          config.Cmd,
			Entrypoint:   config.Entrypoint,
			Image:        config.Image,
			ExposedPorts: config.ExposedPorts,
			Volumes:      config.Volumes,
			Labels:       config.Labels,
		}, &container.HostConfig{
			PortBindings: hostConfig.PortBindings,
			VolumesFrom:  hostConfig.VolumesFrom,
			Links:        hostConfig.Links,
			NetworkMode:  container.NetworkMode(hostConfig.NetworkMode),
			Binds:        hostConfig.Binds,
			PidMode:      container.PidMode(hostConfig.PidMode),
			IpcMode:      container.IpcMode(hostConfig.IpcMode),
			Resources: container.Resources{
				NanoCPUs: hostConfig.CPUResourceLimit,
				Memory:   hostConfig.MemoryLimit,
			},
		}, nil, nil, containerName)
	if err != nil {
		K8sLogger.Error("NewContainer error: ", err)
		return "", err
	}
	return resp.ID, nil
}

// 启动一个容器
func (cm *ContainerManager) StartContainer(ctx context.Context, dockerID string) error {
	err := cm.client.ContainerStart(ctx, dockerID, types.ContainerStartOptions{})
	if err != nil {
		K8sLogger.Error("StartContainer error: ", err)
		return err
	}
	return nil
}

// 停止一个容器
// stopOptions意思是停止容器的时候的一些选项，比如超时时间等
func (cm *ContainerManager) StopContainer(ctx context.Context, dockerID string) error {
	err := cm.client.ContainerStop(ctx, dockerID, container.StopOptions{Timeout: nil})
	if err != nil {
		K8sLogger.Error("StopContainer error: ", err)
		return err
	}
	return nil
}

// 删除一个容器
func (cm *ContainerManager) RemoveContainer(ctx context.Context, dockerID string) error {
	//先暂停
	err := cm.StopContainer(ctx, dockerID)
	if err != nil {
		K8sLogger.Error("RemoveContainer error: ", err)
		return err
	}
	err = cm.client.ContainerRemove(ctx, dockerID, types.ContainerRemoveOptions{Force: true})
	if err != nil {
		K8sLogger.Error("RemoveContainer error: ", err)
		return err
	}
	return nil
}

// 获取所有minik8s容器的信息
func (cm *ContainerManager) ListMinik8sContainer(ctx context.Context) ([]types.Container, error) {
	filter := filters.NewArgs()
	filter.Add("label", string(minik8sTypes.RunningSystemMinik8s)+"="+minik8sTypes.IsTrue)
	c, err := cm.client.ContainerList(ctx, types.ContainerListOptions{
		Filters: filter,
		All:     true,
	})
	if err != nil {
		K8sLogger.Error("ListContainer error: ", err)
		return nil, err
	}
	return c, nil
}

func (cm *ContainerManager) ListALlContainer(ctx context.Context) ([]types.Container, error) {
	c, err := cm.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		K8sLogger.Error("ListContainer error: ", err)
		return nil, err
	}
	return c, nil
}

// 获取一个容器的状态(是否运行、退出码等)
func (cm *ContainerManager) InspectContainer(ctx context.Context, dockerID string) (types.ContainerJSON, error) {
	cj, err := cm.client.ContainerInspect(ctx, dockerID)
	if err != nil {
		K8sLogger.Error("InspectContainer error: ", err)
		return types.ContainerJSON{}, err
	}
	return cj, nil
}

// 获取一个容器的日志
func (cm *ContainerManager) GetContainerLogs(ctx context.Context, dockerID string) (io.ReadCloser, error) {
	rc, err := cm.client.ContainerLogs(ctx, dockerID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		K8sLogger.Error("GetContainerLogs error: ", err)
		return nil, err
	}
	return rc, nil
}

// 获取容器状态（cpu、内存、网络等）
func (cm *ContainerManager) ContainerStats(ctx context.Context, dockerID string) (*types.StatsJSON, error) {
	rc, err := cm.client.ContainerStats(ctx, dockerID, false)
	if err != nil {
		K8sLogger.Error("GetContainerStats error: ", err)
		return nil, err
	}
	defer rc.Body.Close()

	decoder := json.NewDecoder(rc.Body)
	statsInfo := &types.StatsJSON{}
	err = decoder.Decode(statsInfo)
	if err != nil {
		return nil, err
	}
	return statsInfo, nil
}

// 重启一个容器
func (cm *ContainerManager) RestartContainer(ctx context.Context, dockerID string) error {
	err := cm.client.ContainerRestart(ctx, dockerID, container.StopOptions{Timeout: nil})
	if err != nil {
		K8sLogger.Error("RestartContainer error: ", err)
		return err
	}
	return nil
}
