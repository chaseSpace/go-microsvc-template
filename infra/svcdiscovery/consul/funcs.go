package consul

import capi "github.com/hashicorp/consul/api"

func newHealthCheck(uniqueName, tcpAddr string) *capi.AgentServiceCheck {
	return &capi.AgentServiceCheck{
		Name:     uniqueName,
		Interval: "10s",
		Timeout:  "3s",
		TCP:      tcpAddr,
	}
}
