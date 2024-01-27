package status

import (
	"minik8s/pkg/apis"
	"time"
)

type StatusManager interface {
	Start()
	GetPodStatus(uid string) (apis.PodStatus, bool)
	// 设置 pod 的状态并会触发一个状态同步操作
	SetPodStatus(pod *apis.Pod, status apis.PodStatus)
	// 设置 pod .status.containerStatuses 中 container 是否为 ready 状态并触发状态同步操作
	SetContainerReadiness(podUID string, containerID string, ready bool)
	SetContainerStartup(podUID string, containerID string, started bool)
	// 将 pod .status.containerStatuses 和 .status.initContainerStatuses 中 container 的 state 置为 Terminated 状态并触发状态同步操作
	TerminalPod(pod *apis.Pod)
	// 从 statusManager 缓存 podStatuses 中删除对应的 pod
	RemoveOrphanedStatuses(podUIDs map[string]bool)
}

const syncPeriod = 10 * time.Second

type statusManager struct {
}

func NewStatusManager() *statusManager {
	return &statusManager{
		// podStatuses: make(map[string]apis.PodStatus),
	}
}
