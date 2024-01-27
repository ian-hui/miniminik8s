package runtime

import (
	"minik8s/minik8sTypes"
	"minik8s/pkg/apis"
	"testing"
)

var testPod = apis.Pod{
	Kind: "Pod",
	ObjectMeta: apis.ObjectMeta{
		Name:      "testPod",
		Namespace: "testNamespace",
		UID:       "3e057e86-7db7-4ad5-a186-101da380ebb6",
		Labels: map[string]string{
			"app": "test",
		},
	},
	Spec: apis.PodSpec{
		Volumes: []apis.HostVolume{
			{
				Name: "testVolume",
				Type: "HostPath",
				Path: "/Users/ianhui/code/golang/minik8s/minik8sTest",
			}},
		Containers: []apis.Container{
			{
				Name:            "testContainer-1",
				Image:           "docker.io/library/redis",
				ImagePullPolicy: minik8sTypes.IfNotPresent,
				//redis的mount
				VolumeMounts: []apis.VolumeMount{{Name: "testVolume", MountPath: "/data"}},
				Ports:        []apis.ContainerPort{{Name: "redis", HostPort: 6379, ContainerPort: "6379", Protocol: "tcp"}},
			},
			{
				Name:            "testContainer-2",
				Image:           "docker.io/library/nginx",
				ImagePullPolicy: minik8sTypes.IfNotPresent,
				Ports:           []apis.ContainerPort{{Name: "nginx", HostPort: 8080, ContainerPort: "80", Protocol: "tcp"}},
			},
		},
	},
}

func TestCreatePod(t *testing.T) {
	// 创建一个runtimeManager

	r := NewRuntimeManager()
	// err := r.DeletePod(&testPod)
	// if err != nil {
	// 	t.Error(err)
	// }

	// 创建pod
	s, err := r.createPod(&testPod)
	if err != nil {
		t.Error(err)
	}
	t.Log(s)

}

func TestDeletePod(t *testing.T) {
	// 创建一个runtimeManager
	r := NewRuntimeManager()
	err := r.killPod(&testPod)
	if err != nil {
		t.Error(err)
	}

}
