package apis

import "minik8s/minik8sTypes"

type Container struct {
	Name       string
	Image      string
	Command    []string
	Args       []string
	WorkingDir string
	Ports      []ContainerPort
	// EnvFrom                  []EnvFromSource//configmap用的，我们不需要
	Env       []EnvVar //环境变量
	Resources ResourceRequirements
	// ResizePolicy             []ContainerResizePolicy //用于容器资源的动态变化，比如cpu，内存等
	// RestartPolicy            *ContainerRestartPolicy //默认always，不需要动了
	VolumeMounts []VolumeMount
	// VolumeDevices            []VolumeDevice //用于挂载设备，比如硬盘，这种一般是statefulset用的，因为需要把硬盘和容器绑定在一起
	LivenessProbe  *Probe
	ReadinessProbe *Probe
	StartupProbe   *Probe
	Lifecycle      *Lifecycle
	// TerminationMessagePath   string 终止信息默认写入容器的stdout，不需要动
	// TerminationMessagePolicy string
	ImagePullPolicy minik8sTypes.ImagePullPolicyType
}

type ContainerPort struct {
	Name          string
	HostPort      int32
	ContainerPort string
	Protocol      string
	HostIP        string
}

type EnvVar struct {
	Name  string
	Value string
	//ValueFrom *EnvVarSource //configmap才需要，我们简化就直接用value
}

type ResourceRequirements struct {
	Limits   ResourceList
	Requests ResourceList
}

type ResourceList map[string]string

type VolumeMount struct {
	Name      string
	MountPath string
	ReadOnly  bool
}

type Probe struct {
	Handler             Handler
	InitialDelaySeconds int32 //容器启动后多久开始探测
	TimeoutSeconds      int32 //探测超时时间，如果超过这个时间，就认为失败了
	PeriodSeconds       int32 //探测周期
	SuccessThreshold    int32 //成功门限，如果连续成功次数达到这个值，就认为成功了
	FailureThreshold    int32 //失败门限，如果连续失败次数达到这个值，就认为失败了
}

type Handler struct {
	HttpGet *HttpGetAction
	// Exec    *ExecAction
	// TcpSocket *TcpSocketAction
	// Grpc      *GrpcAction
}

type HttpGetAction struct {
	Path   string
	Port   int32
	Scheme string
	Host   string
}

type Lifecycle struct {
	PostStart *Handler
	PreStop   *Handler
}
