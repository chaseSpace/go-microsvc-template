package ip

import "net"

func GetLocalPrivateIPs(v4 bool, netCard string) (ips []net.IP, err error) {
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		// 过滤掉Loopback和无效接口
		//println(iface.Name)
		if netCard != "" && iface.Name != netCard {
			continue
		}
		// 过滤掉Loopback和无效接口
		if iface.Flags&net.FlagLoopback == 0 && iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagMulticast != 0 {
			addrs, err := iface.Addrs()
			if err != nil {
				return nil, err
			}
			for _, addr := range addrs {
				ipnet, ok := addr.(*net.IPNet)
				if ok && ipnet.IP.IsPrivate() {
					if (v4 && ipnet.IP.To4() != nil) || (!v4 && ipnet.IP.To4() == nil) {
						ips = append(ips, ipnet.IP)
					}
				}
			}
		}
	}

	return ips, nil
}
