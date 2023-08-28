## Go微服务模板

一个简洁、清爽的微服务项目架构，从变量命名到不同职责的（多层）目录结构定义。

> **完成进度：80%**

计划支持以下模式或特性：

- ✅ 使用单仓库多服务模式
- ✅ 使用grpc+protobuf作为内部rpc通讯协议
- ✅ 使用grpc-gateway插件生成grpc服务的http反向代理
- ✅ 使用consul作为服务注册发现组件，支持扩展
    - 包含健康检查、超时重试与熔断功能
    - 包含服务之间通信流量的负载均衡
    - 包含服务之间通信的认证与授权
- ✅ 使用gorm作为orm组件，支持扩展
- ✅ 使用redis作为cache组件，支持扩展
- ✅ 支持本地启动**多个**微服务
    - 支持本地无注册中心启动多个微服务

其他有用的特性：

- ✅ shell脚本支持mac环境（默认linux）
- ✅ 定义微服务ERROR类型，以便跨服务传递error（已实现对应GRPC拦截器）

运行通过的示例：
- ✅ 单服务GRPC接口测试用例（[ext_api_test](./test/user/ext_api_test.go)）
- ✅ 跨服务GRPC调用测试用例（[admin-ext_api_test](./test/admin/ext_api_test.go)）

### Preview

🍡 一瞥 🍡

```go
// service/user/main.go
package main

import (
  "google.golang.org/grpc"
  "microsvc/deploy"
  "microsvc/infra"
  "microsvc/infra/sd"
  "microsvc/infra/svccli"
  "microsvc/infra/xgrpc"
  "microsvc/pkg"
  "microsvc/pkg/xlog"
  "microsvc/protocol/svc/user"
  deploy2 "microsvc/service/user/deploy"
  "microsvc/service/user/handler"
  "microsvc/util/graceful"
)

func main() {
  graceful.SetupSignal()
  defer graceful.OnExit()

  // 初始化config
  deploy.Init("user", deploy2.UserConf)
  // 初始化服务用到的基础组件（封装于pkg目录下），如log, kafka等
  pkg.Init(
    xlog.Init,
    // 假如我要新增kafka等组件，也是新增 pkg/xkafka目录，然后实现其init函数并添加在这里
  )

  // 初始化几乎每个服务都需要的infra组件，must参数指定是否必须初始化成功，若must=true且err非空则panic
  infra.MustSetup(
    //cache.InitRedis(true),
    //orm.InitGorm(true),
    sd.Init(true),
    svccli.Init(true),
  )

  x := xgrpc.New() // New一个封装好的grpc对象
  x.Apply(func(s *grpc.Server) {
    // 注册外部和内部的rpc接口对象
    user.RegisterUserExtServer(s, new(handler.UserExtCtrl))
    user.RegisterUserIntServer(s, new(handler.UserIntCtrl))
  })
  // 仅开发环境需要启动HTTP端口来代理gRPC服务
  if deploy.XConf.IsDevEnv() {
    x.SetHTTPExtRegister(user.RegisterUserExtHandler)
  }

  x.Start(deploy.XConf)
  // GRPC服务启动后 再注册服务
  sd.Register(deploy.XConf)

  graceful.Run()
}
```

### 1. 目录结构释义

```
├── consts  # 公共常量（不含单个svc独享的常量）
├── enums   # 所有枚举（含svc独享的枚举，enums数量一般小于consts，且大部分需要跨服务使用）
├── deploy  # 部署需要的公共配置文件，如db配置
│   ├── beta
│   ├── dev
│   └── prod
├── infra   # 基础设施（的初始化或内部逻辑），不含业务代码
│   ├── cache
│   ├── orm
│   ├── svccli
│   ├── svcregistar
│   ├── util
│   └── xgrpc
├── pkg     # 项目封装的常用包，比如err,time等，不含业务代码
│   └── xerr
├── proto   # proto文件
│   ├── include    # 可能引用的第三方proto文件，比如Google发布的proto类型
│   │   └── google
│   ├── model      # 项目内的表结构对应的struct定义，以服务划分目录
│   │   ├── admin
│   │   └── user
│   └── svc        # 各微服务使用的proto文件
│       ├── admin
│       ├── assets
│       └── user
│           ├── user.ext.proto    # user服务的外部接口组，仅允许外部调用，需要鉴权
│           └── user.int.proto    # ...内部接口组，仅允许内部调用，可不鉴权
├── protocol  # 生成的go文件
│   └── svc
│       ├── admin
│       ├── assets
│       └── user
├── service   # 微服务目录，存放业务代码
│   ├── admin
│   ├── gateway
│   └── user
│       └── deploy   # 每个微服务都有的目录，存放各自使用的专属配置目录（不含公共db配置，所以内容更少）
│           ├── beta
│           ├── dev
│           └── prod
├── tool   # 项目使用的外部工具，主要是二进制文件，如protoc等
│   └── protoc_v24   # 更改工具时，建议目录名包含版本
├── tool_mac # mac环境使用的外部工具
│   └── protoc_v24
└── bizcomm  # 存放可共用的业务逻辑
│  
└── util  # 存放可共用的其他逻辑
```

### 2. 如何使用

```shell
git clone https://github.com/chaseSpace/go-microsvc-template.git
cd go-microsvc-template/
go mod tidy
```

### 3. 工具下载（更新）

#### 下载protoc

linux、mac版本都已经包含在本仓库的`tool/`,`tool_mac/`目录下，无需再下载，已下载的是protoc
v24版本，其余插件也是编写本文档时的最新版本（下载时间2023年8月17日）。

如需更换版本，可点击下方链接自行下载：

https://github.com/protocolbuffers/protobuf/releases

> windows环境暂未支持，请自行配置环境。  
> 本模板配套的是shell脚本，在windows环境运行可能有问题，（但仍然建议使用类unix环境进行开发，以减少不必要的工作和麻烦）。

#### 下载protoc插件

本仓库的`tool/`,`tool_mac/`都已经包含这些插件，这里只是演示如何下载，以便你了解如何更新插件版本。

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.16
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.16

# 检查是否下载成功
[root@localhost go-microsvc-template]# ls $GOPATH/bin/
protoc         protoc-gen-go-grpc     protoc-gen-grpc-gateway  protoc-gen-swagger
protoc-gen-go  protoc-gen-openapiv2   

# 下载后需要复制到仓库下的tool目录（以及tool_mac），其他人拉取代码后，无需再下载
cp $GOPATH/bin/* tool/protoc_v24
```

若要更改版本，建议同时修改`tool/proto_v24/`目录名称，并同步修改`build_pb.sh`脚本中对该目录的引用部分，以便更新版本后脚本能够正常运行。

### 4. 本地启动微服务的原理

理论上来说，调用微服务是走注册中心的，要想在本地启动多个微服务且能正常互相调用，又不想在本地部署一个类似etcd/consul/zookeeper
的注册中心，最简单的办法是：

```
实现一个简单的单进程注册中心，当启动一个微服务且env=dev时，内部组件会检测本地是否有注册中心服务运行，若有则直接调用其接口进行注册；
若没有则会启动一个注册中心服务，供其他服务使用。

> 本地的注册中心使用一个可配置的固定端口。
```

注意：本地启动的微服务仍然连接的是**beta环境的数据库**。

### 其他建议

- `protocol/`是存放生成协议代码的目录，在实际项目开发中可以加入`.gitignore`文件，以避免在PR review时产生困扰；

#### 资源链接

- [Consul 官网介绍](https://developer.hashicorp.com/consul/docs/intro)
- [Consul 服务发现原理](https://developer.hashicorp.com/consul/docs/concepts/service-discovery)