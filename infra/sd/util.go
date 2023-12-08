package sd

import (
	"fmt"
	"microsvc/deploy"
	"microsvc/enums/svc"
)

// GetSvcTargetInK8s
// 注意: 这里的 .default.svc.cluster.local 是k8s中的default空间下svc的默认域名，你可以根据情况修改default为其他命名空间
// 修改域名的同时需要重新生成grpc通信使用的证书
func GetSvcTargetInK8s(svc svc.Svc) string {
	return fmt.Sprintf("go-svc-%s.default.svc.cluster.local:%d", svc.Name(), deploy.XConf.GRPCPort)
}
