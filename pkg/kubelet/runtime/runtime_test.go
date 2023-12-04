package runtime

import (
	"minik8s/pkg/apis"
	"minik8s/pkg/uuid"
	"testing"
)

var testPod = apis.Pod{
	Kind: "Pod",
	ObjectMeta: apis.ObjectMeta{
		Name:      "testPod",
		Namespace: "testNamespace",
		UID:       uuid.NewUID(),
		Labels: map[string]string{
			"app": "test",
		},
	},
	Spec: apis.PodSpec{
		Containers: []apis.Container{
			{
				Name:  "testContainer-1",
				Image: "docker.io/library/redis",
			},
			{
				Name:  "testContainer-2",
				Image: "docker.io/library/nginx",
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
