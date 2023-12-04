package apis

import (
	"minik8s/minik8sTypes"
	"time"

	"github.com/docker/docker/api/types"
)

type PodSandboxConfig struct {
	Metadata     *PodSandboxMetadata
	LogDirectory string
	Labels       map[string]string
	Annotations  map[string]string
}

type PodSandboxMetadata struct {
	Name      string
	Uid       string
	Namespace string
}

type Pod struct {
	ObjectMeta
	Spec PodSpec
	Kind string
	PodStatus
}

type ObjectMeta struct {
	Name        string
	UID         string
	Namespace   string
	Labels      map[string]string
	Annotations map[string]string
}

type PodSpec struct {
	Volumes        []HostVolume
	Containers     []Container
	RestartPolicy  minik8sTypes.RestartPolicy
	InitContainers []Container
}

type HostVolume struct {
	Name string
	Type string
	Path string
}

type PodStatus struct {
	// IP address allocated to the pod. Routable at least within the cluster. Empty if not yet allocated.
	PodIP string `json:"podIP" yaml:"podIP"` //在cni设置完毕（add）后，这个值将会被设置

	Phase PodPhase

	// 容器的状态数组
	ContainerStatuses []types.ContainerState `json:"containerStatuses" yaml:"containerStatuses"`

	// 最新的更新时间
	// UpdateTime string `json:"lastUpdateTime" yaml:"lastUpdateTime"`
	UpdateTime time.Time `json:"lastUpdateTime" yaml:"lastUpdateTime"`

	// Pod的容器的资源使用情况
	CpuPercent float64 `json:"cpuPercent" yaml:"cpuPercent"`
	MemPercent float64 `json:"memPercent" yaml:"memPercent"`
}

// 直接抄过来
type PodPhase string

// These are the valid statuses of pods.
const (
	// PodPending means the pod has been accepted by the system, but one or more of the containers
	// has not been started. This includes time before being bound to a node, as well as time spent
	// pulling images onto the host.
	PodPending PodPhase = "Pending"
	// PodRunning means the pod has been bound to a node and all of the containers have been started.
	// At least one container is still running or is in the process of being restarted.
	PodRunning PodPhase = "Running"
	// PodSucceeded means that all containers in the pod have voluntarily terminated
	// with a container exit code of 0, and the system is not going to restart any of these containers.
	PodSucceeded PodPhase = "Succeeded"
	// PodFailed means that all containers in the pod have terminated, and at least one container has
	// terminated in a failure (exited with a non-zero exit code or was stopped by the system).
	PodFailed PodPhase = "Failed"
	// PodUnknown means that for some reason the state of the pod could not be obtained, typically due
	// to an error in communicating with the host of the pod.
	// Deprecated: It isn't being set since 2015 (74da3b14b0c0f658b3bb8d2def5094686d0e9095)
	PodUnknown PodPhase = "Unknown"
)
