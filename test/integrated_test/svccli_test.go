package integrated_test

import (
	"microsvc/infra/sd"
	"microsvc/infra/svccli/rpc"
	"testing"
)

func TestSvcCliNormal(t *testing.T) {
	sd.Init(true)
	rpc.User()
}
