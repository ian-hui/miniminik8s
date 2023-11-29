package containermanager

import (
	"minik8s/minik8sTypes"
	"minik8s/pkg/apis"

	"github.com/docker/go-connections/nat"
)

func MakeContainerMapper(container *apis.Container, portsetInPod *nat.PortSet) error {
	// PortBindingsInPod := nat.PortMap{} 我们不需要portmap，portmap是直接映射到宿主机的端口 我们用service来做端口映射

	for _, containerNetwork := range container.Ports {
		//有些信息是空的，比如protocol
		if containerNetwork.Protocol == "" {
			containerNetwork.Protocol = minik8sTypes.Minik8sTcpProtocol
		}
		if containerNetwork.HostIP == "" {
			containerNetwork.HostIP = minik8sTypes.Minik8sLocalHostIp
		}
		//这里的端口号是容器内部的端口号
		p, err := nat.NewPort(containerNetwork.Protocol, containerNetwork.ContainerPort)
		if err != nil {
			K8sLogger.Errorln("MakeContainerMapper error: ", err)
			return err
		}
		(*portsetInPod)[p] = struct{}{}
	}
	return nil
}
