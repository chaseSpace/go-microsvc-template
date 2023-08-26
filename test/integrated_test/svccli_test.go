package integrated_test

import (
	"microsvc/infra/svccli"
	"microsvc/infra/svcdiscovery"
	"testing"
)

func TestSvcCliNormal(t *testing.T) {
	svcdiscovery.Init(true)
	svccli.User()
}
