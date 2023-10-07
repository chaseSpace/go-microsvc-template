# 示例集合（Examples）

## 1. 微服务接口调用
rpc包放在infra/下面：
```shell
infra/
├── svccli
│   ├── rpc   # 微服务的内部（internal）接口调用，不会鉴权，仅开启grpc tls
│   └── rpcext  # 微服务的外部接口调用，限制了仅允许gateway调用，会鉴权
```

### 1.1 内部接口调用
以admin调用user的`GetUser`为例：
```go
import "microsvc/infra/svccli/rpc"

func (a AdminCtrl) GetUser(ctx context.Context, req *admin.GetUserReq) (*admin.GetUserRsp, error) {
	// 调用user接口获取数据
	rsp, err := rpc.User().GetUser(ctx, &user.GetUserIntReq{
		Uids: req.Uids,
	})
	if err != nil {
		return nil, err
	}
	return &admin.GetUserRsp{Umap: rsp.Umap}, nil
}
```

### 1.2 外部接口调用
通常仅允许gateway进行以原生grpc方式调用，所以不会引用`rpcext`包。本项目的gateway已经实现了外部接口的透传，所以一般无需在代码中显式调用外部接口。

>通常不允许内部服务调用外部接口，但可以根据实际情况进行调整。