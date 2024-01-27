package minik8sTypes

import "github.com/docker/go-connections/nat"

type Config struct {
	//****************************************************//
	//********** docker standard config ***********************//
	// Hostname        string              // Hostname 不需要，直接使用默认生成的
	// Domainname      string              // Domainname 不需要，直接使用默认生成的
	// User            string              // User that will run the command(s) inside the container, also support user:group 不需要，一般是用户yaml文件中指定，但是我们这里不做这部分
	// AttachStdin     bool                // Attach the standard input, makes possible user interaction
	// AttachStdout    bool                // Attach the standard output 默认是开的，与日志相关
	// AttachStderr    bool                // Attach the standard error	默认是开的，与日志相关
	// ExposedPorts    nat.PortSet         `json:",omitempty"` // List of exposed ports 声明容器暴露的端口，这个是必须的
	// Tty             bool                // Attach standard streams to a tty, including stdin if it is not closed.
	// OpenStdin       bool                // Open stdin
	// StdinOnce       bool                // If true, close stdin after the 1 attached client disconnects.
	// Env             []string            // List of environment variable to set in the container
	// Cmd             []string   // Command to run when starting the container
	// Healthcheck     *HealthConfig       `json:",omitempty"` // Healthcheck describes how to check the container is healthy
	// ArgsEscaped     bool                `json:",omitempty"` // True if command is already escaped (Windows specific)
	// Image           string              // Name of the image as it was passed by the operator (e.g. could be symbolic)
	// Volumes         map[string]struct{} // List of volumes (mounts) used for the container
	// WorkingDir      string              // Current directory (PWD) in the command will be launched
	// Entrypoint      []string  // Entrypoint to run when starting the container
	// NetworkDisabled bool                `json:",omitempty"` // Is network disabled
	// MacAddress      string              `json:",omitempty"` // Mac Address of the container
	// OnBuild         []string            // ONBUILD metadata that were defined on the image Dockerfile
	// Labels          map[string]string   // List of labels set to this container
	// StopSignal      string              `json:",omitempty"` // Signal to stop a container
	// StopTimeout     *int                `json:",omitempty"` // Timeout (in seconds) to stop a container
	// Shell           strslice.StrSlice   `json:",omitempty"` // Shell for shell-form of RUN, CMD, ENTRYPOINT

	//****************************************************//
	//我们不需要使用太多，所以只保留了一些常用的
	//********** docker custom config ***********************//
	Tty             bool                // 是否需要Tty终端 Attach standard streams to a tty, including stdin if it is not closed.
	Env             []string            // 环境变量 List of environment variable to set in the container
	Cmd             []string            // 启动子容器的时候执行的命令 Command to run when starting the container
	Entrypoint      []string            // Entrypoint to run when starting the container
	Image           string              // Name of the image as it was passed by the operator (e.g. could be symbolic)
	ImagePullPolicy ImagePullPolicyType `default:"IfNotPresent"` // 拉取镜像的策略, Always, Never, IfNotPresent
	Volumes         map[string]struct{} // List of volumes (mounts) used for the container
	Labels          map[string]string   // List of labels set to this container
	ExposedPorts    nat.PortSet         // List of exposed ports // List of exposed ports 声明容器暴露的端口，这个是必须的
}

type ImagePullPolicyType string

const (
	Always       ImagePullPolicyType = "Always"
	Never        ImagePullPolicyType = "Never"
	IfNotPresent ImagePullPolicyType = "IfNotPresent"
)

type HostConfig struct {
	//定义：HostConfig 包含了容器与宿主机交互的配置，如端口绑定、卷挂载、资源限制等。这些配置影响容器如何与宿主机的系统和资源交互。
	//一般参数都是宿主机操作容器需要的参数，比如端口映射，宿主机的端口映射到容器的端口，这样宿主机就可以通过宿主机的端口访问容器的端口了。
	//****************************************************//
	//********** docker standard config ***********************//
	// // Applicable to all platforms
	// Binds           []string          // List of volume bindings for this container
	// ContainerIDFile string            // File (path) where the containerId is written
	// LogConfig       LogConfig         // Configuration of the logs for this container
	// NetworkMode     NetworkMode       // Network mode to use for the container
	// PortBindings    nat.PortMap       // Port mapping between the exposed port (container) and the host 配置端口映射，这通常与exposedPorts配合使用
	// RestartPolicy   RestartPolicy     // Restart policy to be used for the container
	// AutoRemove      bool              // Automatically remove container when it exits
	// VolumeDriver    string            // Name of the volume driver used to mount volumes
	// VolumesFrom     []string          // List of volumes to take from other container
	// ConsoleSize     [2]uint           // Initial console size (height,width)
	// Annotations     map[string]string `json:",omitempty"` // Arbitrary non-identifying metadata attached to container and provided to the runtime

	// // Applicable to UNIX platforms
	// CapAdd          strslice.StrSlice // List of kernel capabilities to add to the container
	// CapDrop         strslice.StrSlice // List of kernel capabilities to remove from the container
	// CgroupnsMode    CgroupnsMode      // Cgroup namespace mode to use for the container
	// DNS             []string          `json:"Dns"`        // List of DNS server to lookup
	// DNSOptions      []string          `json:"DnsOptions"` // List of DNSOption to look for
	// DNSSearch       []string          `json:"DnsSearch"`  // List of DNSSearch to look for
	// ExtraHosts      []string          // List of extra hosts
	// GroupAdd        []string          // List of additional groups that the container process will run as
	// IpcMode         IpcMode           // IPC namespace to use for the container
	// Cgroup          CgroupSpec        // Cgroup to use for the container
	// Links           []string          // List of links (in the name:alias form)
	// OomScoreAdj     int               // Container preference for OOM-killing
	// PidMode         PidMode           // PID namespace to use for the container
	// Privileged      bool              // Is the container in privileged mode
	// PublishAllPorts bool              // Should docker publish all exposed port for the container
	// ReadonlyRootfs  bool              // Is the container root filesystem in read-only
	// SecurityOpt     []string          // List of string values to customize labels for MLS systems, such as SELinux.
	// StorageOpt      map[string]string `json:",omitempty"` // Storage driver options per container.
	// Tmpfs           map[string]string `json:",omitempty"` // List of tmpfs (mounts) used for the container
	// UTSMode         UTSMode           // UTS namespace to use for the container
	// UsernsMode      UsernsMode        // The user namespace to use for the container
	// ShmSize         int64             // Total shm memory usage
	// Sysctls         map[string]string `json:",omitempty"` // List of Namespaced sysctls used for the container
	// Runtime         string            `json:",omitempty"` // Runtime to use with this container

	// // Applicable to Windows
	// Isolation Isolation // Isolation technology of the container (e.g. default, hyperv)

	// // Contains container's resources (cgroups, ulimits)
	// Resources

	// // Mounts specs used by the container
	// Mounts []mount.Mount `json:",omitempty"`

	// // MaskedPaths is the list of paths to be masked inside the container (this overrides the default set of paths)
	// MaskedPaths []string

	// // ReadonlyPaths is the list of paths to be set as read-only inside the container (this overrides the default set of paths)
	// ReadonlyPaths []string

	// // Run a custom init inside the container, if null, use the daemon's configured settings
	// Init *bool `json:",omitempty"`

	//****************************************************//
	//我们不需要使用太多，所以只保留了一些常用的
	//********** docker custom config ***********************//
	VolumesFrom      []string    // List of volumes to take from other containers
	Links            []string    // List of links (in the name:alias form)
	NetworkMode      string      // [网络模式] Network mode to use for the container
	PidMode          string      // [PidMode] PID namespace to use for the container
	IpcMode          string      // [IPC Mode ]IPC namespace to use for the container(设置这三个可以让容器共享网络、PID、IPC的ns)
	Binds            []string    // List of volume bindings for this container
	PortBindings     nat.PortMap // List of port bindings for this container // Port mapping between the exposed port (container) and the host 配置端口映射，这通常与exposedPorts配合使用
	CPUResourceLimit int64       // CPU资源限制 单位是10的负9次方核
	MemoryLimit      int64       // 内存资源限制 单位是字节
}

type RunningSystem string

const (
	LabelsContainerName  string        = "containerName"
	RunningSystemDocker  RunningSystem = "docker"
	RunningSystemMinik8s RunningSystem = "minik8s"
	IsTrue               string        = "_true"
)

const (
	IpcModeShareable      = "shareable"
	NsModeContainerPrefix = "container:"
)
