package svcregistar

func Init(must bool) func(func(must bool, err error)) {
	// 在这里 决定使用 etcd/consul
	return func(onEnd func(must bool, err error)) {
		// todo
		var err error
		onEnd(must, err)
	}
}
