package consul

import capi "github.com/hashicorp/consul/api"

func newHealthCheck(uniqueName, tcpAddr string) *capi.AgentServiceCheck {
	return &capi.AgentServiceCheck{
		CheckID:  uniqueName,
		Interval: "6s",
		Timeout:  "3s",
		TCP:      tcpAddr,
		//Status:   capi.HealthPassing,
	}
}
