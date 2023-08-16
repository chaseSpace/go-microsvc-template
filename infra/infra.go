package infra

type initFunc func(onEnd func(must bool, err error))

func MustSetup(initFn ...initFunc) {
	for _, fn := range initFn {
		fn(func(must bool, err error) {
			if must && err != nil {
				panic(err)
			}
			// TODO LOG
		})
	}
}
