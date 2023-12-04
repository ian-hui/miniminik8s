package imagemanager

import (
	"context"
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

func TestPullImage(t *testing.T) {
	im := NewImageManager(dockerclient.GetDockerClient())
	err := im.PullImage(context.TODO(), minik8sTypes.IfNotPresent, minik8sTypes.Minik8sPauseImage)
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveImage(t *testing.T) {
	im := NewImageManager(dockerclient.GetDockerClient())
	err := im.RemoveImage(context.Background(), minik8sTypes.Minik8sPauseImage)
	if err != nil {
		t.Error(err)
	}
}
