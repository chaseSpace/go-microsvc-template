package integrated_test

import (
	"microsvc/infra/sd"
	"microsvc/infra/svccli"
	"testing"
)

func TestSvcCliNormal(t *testing.T) {
	sd.Init(true)
	svccli.User()
}
