package xrand

func init() {
	// Begin with go 1.20, `Seed` initialize operation is not need again.
	//rand.Seed(time.Now().UnixNano())
}
