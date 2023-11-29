package minik8sTypes

// minikube labels
const (
	KubernetesPodNameLabel      = "io.minik8s.pod.name"
	KubernetesPodNamespaceLabel = "io.minik8s.pod.namespace"
	KubernetesPodUIDLabel       = "io.minik8s.pod.uid"
)

type RestartPolicy string

// restart policy
const (
	Minik8sRestartPolicyAlways    = "Always"
	Minik8sRestartPolicyNever     = "Never"
	Minik8sRestartPolicyOnFailure = "OnFailure"
)

type PodType string

// podtype labels
const (
	Minik8sPodType        = "io.minik8s.pod.type"
	Minik8sPausePodType   = "pause"
	Minik8sPauseImage     = "k8s.gcr.io/pause:3.1"
	Minik8sGenericPodType = "generic"
)

// networks
const (
	Minik8sTcpProtocol = "TCP"
	Minik8sUdpProtocol = "UDP"
	Minik8sLocalHostIp = "127.0.0.1"
)
