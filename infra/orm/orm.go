package orm

func InitGorm(must bool) func(func(must bool, err error)) {
	return func(onEnd func(must bool, err error)) {
		// todo
		var err error
		onEnd(must, err)
	}
}
