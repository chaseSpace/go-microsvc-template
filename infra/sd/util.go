package sd

import (
	"fmt"
	"microsvc/deploy"
	"microsvc/enums/svc"
)

// GetSvcTargetInK8s e.g. `go-svc-user` is `user` service's dns name.
func GetSvcTargetInK8s(svc svc.Svc) string {
	return fmt.Sprintf("go-svc-%s:%d", svc.Name(), deploy.XConf.GRPCPort)
}
