package dockerclient

import (
	"minik8s/logger"

	"github.com/docker/docker/client"
)

var (
	dclient   *client.Client
	K8sLogger = logger.K8sLogger
)

func NewDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		K8sLogger.Error("NewDockerClient error: ", err)
		return nil, err
	}
	dclient = cli
	return dclient, nil
}

func GetDockerClient() *client.Client {
	if dclient == nil {
		_, err := NewDockerClient()
		if err != nil {
			K8sLogger.Error("GetDockerClient error: ", err)
			return nil
		}
	}
	return dclient
}
