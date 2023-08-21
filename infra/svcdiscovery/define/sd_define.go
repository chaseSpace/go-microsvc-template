package define

type ServiceDiscovery interface {
	Register(serviceName string, address string, port int, metadata map[string]string) error
	Deregister(serviceName string) error
	Discover(serviceName string) ([]ServiceInstance, error)
}

// ServiceInstance 表示注册的单个实例
type ServiceInstance struct {
	ID       string
	Name     string
	Address  string
	Port     int
	Metadata map[string]string
}
