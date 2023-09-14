## Go微服务模板

一个简洁、清爽的微服务项目架构，从变量命名到不同职责的（多层）目录结构定义。

**目录**

<!-- TOC -->
  * [Preview](#preview)
  * [1. 启动&停止日志输出](#1-启动停止日志输出)
  * [2. 目录结构释义](#2-目录结构释义)
  * [3. 如何使用](#3-如何使用)
  * [4. 示例集合](#4-示例集合)
    * [4.1 美化proto unmarshal error response](#41-美化proto-unmarshal-error-response)
  * [5. 本地（dev）环境启动微服务的原理](#5-本地dev环境启动微服务的原理)
  * [6. 工具下载（更新）](#6-工具下载更新)
    * [6.1 下载protoc](#61-下载protoc)
    * [6.2 下载protoc插件](#62-下载protoc插件)
  * [7. 其他](#7-其他)
    * [计划](#计划)
    * [建议](#建议)
    * [使用的外部库](#使用的外部库)
    * [资源链接](#资源链接)
<!-- TOC -->

> **完成进度：90%**

计划支持以下模式或特性：

- ✅ 使用单仓库多服务模式
- ✅ 使用grpc+protobuf作为内部rpc通讯协议
- ✅ 统一API Gateway管理南北流量
    - ✅ 透明转发HTTP流量到后端服务，无编码开销
    - ✅ 能够动态转发流量至新增服务，无需重启（通过服务发现以及自定义gRPC编解码方式）
- ✅ 使用consul作为注册中心组件，支持扩展
    - ✅ 包含健康检查
    - ✅ 包含服务之间通信流量的负载均衡
    - ✅ 包含服务之间通信的加密、授权
- ✅ 使用gorm作为orm组件，支持扩展
- ✅ 使用redis作为cache组件，支持扩展
- RPC超时重试与熔断功能
- ✅ 支持本地启动**多个**微服务（不需要部署第三方注册）

其他有用的特性：

- ✅ shell脚本支持mac环境（默认linux）
- ✅ 定义微服务Error类型，以便跨服务传递error（在GRPC拦截器中解析），[查看代码](./pkg/xerr/err.go)
- ✅ 跨多个服务传递metadata示例（通过Context），搜索函数`TraceGRPC`
- ✅ gRPC Client 拦截器示例，包含`GRPCCallLog`, `ExtractGRPCErr`, `WithFailedClient`
- ✅ gRPC Server 拦截器示例，包含`RecoverGRPCRequest`, `ToCommonResponse`, `LogGRPCRequest`, `TraceGRPC`, `StandardizationGRPCErr`
- ✅ 美化proto解析错误response，[查看示例](#41-美化proto-unmarshal-error-response)
- gateway实现oauth2.0统一鉴权（admin服务单独简单token-自动续期鉴权）


运行通过的示例：

- ✅ **本地**单服务GRPC接口测试用例（[user-ext_api_test](./test/user/ext_api_test.go)）
- ✅ **本地**跨服务GRPC调用测试用例（[admin-ext_api_test](./test/admin/ext_api_test.go)）

目前已提供常见的微服务示例：
- admin: 管理后台
- user：用户模块（后续会实现基础的注册、登录功能）
- assets（TODO）：资产模块（后续会实现一个简单含流水、消费、进账的货币功能）
- review：审核模块（自行接入第三方）


本项目文档指引：
- [使用证书加密以及指定授权gRPC通信](./generate_cert_for_svc.md)

### Preview

🍡 一瞥 🍡

```go
// service/user/main.go
package main

import (
  "google.golang.org/grpc"
  "microsvc/deploy"
  "microsvc/enums"
  "microsvc/infra"
  "microsvc/infra/sd"
  "microsvc/infra/svccli"
  "microsvc/infra/xgrpc"
  _ "microsvc/infra/xgrpc/protobytes"
  "microsvc/pkg"
  "microsvc/pkg/xkafka"
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
  deploy.Init(enums.SvcUser, deploy2.UserConf)

  // 初始化服务用到的基础组件（封装于pkg目录下），如log, kafka等
  pkg.Setup(
    xlog.Init,
    xkafka.Init,
  )

  // 初始化几乎每个服务都需要的infra组件，must参数指定是否必须初始化成功，若must=true且err非空则panic
  infra.Setup(
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

  x.Start(deploy.XConf)
  // GRPC服务启动后 再注册服务
  sd.Register(deploy.XConf)

  graceful.Run()
}
```

### 1. 启动&停止日志输出

<details>
<summary>点击展开/折叠</summary>

```shell
************* init Share-Config OK *************
&deploy.XConfig{                                           
  Svc:   "admin",                                          
  Env:   "dev",                                            
  Mysql: map[string]*deploy.Mysql{                         
    "microsvc": &deploy.Mysql{                             
      DBname:   "microsvc",                                
      Host:     "0.0.0.0",                                 
      Port:     "3306",                                    
      User:     "root",                                    
      Password: "123",                                     
      Password: "123",
    },
  },
  SimpleSdHttpPort: 5000,
  gRPCPort: 0,
  httpPort: 0,
  svcConf:  nil,
}

************* init Svc-Config OK *************
&deploy.SvcConfig{
  CommConfig: deploy.CommConfig{
    Svc:      "admin",
    LogLevel: "debug",
  },
}

Congratulations! ^_^
Your service ["go-admin"] is serving gRPC on "localhost:60280"

{"LEVEL":"x-info","TS":"2023-08-29 15:44:41.625","CALLER":"sd/base.go:61","MSG":"sd: register svc success","reg_svc":"go-admin","addr":"127.0.0.1:60280","SERVICE":"go-admin"}

### 停止服务...

{"LEVEL":"x-warn","TS":"2023-08-29 15:44:43.163","CALLER":"graceful/base.go:46","MSG":"****** graceful ****** server ready to exit(signal)","SERVICE":"go-admin"}
{"LEVEL":"x-debug","TS":"2023-08-29 15:44:43.163","CALLER":"svccli/base.go:69","MSG":"svccli: resource released...","SERVICE":"go-admin"}
{"LEVEL":"x-debug","TS":"2023-08-29 15:44:43.164","CALLER":"sd/base.go:72","MSG":"sd: deregister success","svc":"go-admin","SERVICE":"go-admin"}
{"LEVEL":"x-debug","TS":"2023-08-29 15:44:43.165","CALLER":"cache/redis.go:77","MSG":"cache-redis: resource released...","SERVICE":"go-admin"}
{"LEVEL":"x-debug","TS":"2023-08-29 15:44:43.165","CALLER":"orm/mysql.go:85","MSG":"orm-mysql: resource released...","SERVICE":"go-admin"}
{"LEVEL":"x-info","TS":"2023-08-29 15:44:43.165","CALLER":"xgrpc/grpc.go:79","MSG":"xgrpc: gRPC server shutdown completed","SERVICE":"go-admin"}
{"LEVEL":"x-info","TS":"2023-08-29 15:44:43.165","CALLER":"graceful/base.go:30","MSG":"****** graceful ****** server exited","SERVICE":"go-admin"}
```

</details>

### 2. 目录结构释义

```
├── bizcomm # 业务公共代码
├── consts  # 公共常量（不含单个svc独享的常量）
├── enums   # 公共枚举（含svc独享的枚举，enums数量一般小于consts，且大部分需要跨服务使用）
├── deploy  # 部署需要的公共配置文件，如db配置
│   ├── beta
│   ├── dev
│       └── cert  # 证书目录，仅供模板演示，实际项目中不应和代码一起托管
│   └── prod
├── docs    # 项目各类文档，建议再划分子目录
│   └── sql   
├── infra   # 基础设施（的初始化或内部逻辑），不含业务代码
│   ├── cache   # 缓存基础代码
│   ├── orm     # ORM基础代码
│   ├── sd      # 服务注册发现基础代码
│   ├── svccli  # 服务client基础代码
│   └── xgrpc   # grpc基础代码
├── pkg     # 项目封装的常用包，比如err,time等，不含业务代码
│   └── xerr
│   └── xkafka
│   └── xlog
│   └── xtime
├── proto   # protobuf文件
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
├── protocol  # 生成的pb文件
│   └── svc
│       ├── admin
│       ├── assets
│       └── user
├── service   # 微服务目录，存放业务代码
│   ├── admin  # 示例服务：管理后台
│   ├── gateway  # 统一网关，转发所有流量到后端服务
│   └── user
│       └── deploy   # 每个微服务都有的目录，存放各自使用的专属配置目录（不含公共db配置，所以代码很少）
│       ├── cache    
│       ├── dao
│       ├── deploy
│       │   └── dev
│       └── handler
├── test
│   ├── admin
│   ├── gateway
│   ├── tbase
│   └── user
├── tool   # 项目使用的外部工具，主要是二进制文件，如protoc等
│   └── protoc_v24   # 更改工具时，建议目录名包含版本
├── tool_mac # mac环境使用的外部工具
│   └── protoc_v24
│  
└── util  # 存放可共用的其他逻辑
```

### 3. 如何使用

```shell
git clone https://github.com/chaseSpace/go-microsvc-template.git
cd go-microsvc-template/
go mod download

# 启动服务
go run service/user/main.go
go run service/admin/main.go
...
```

### 4. 示例集合

#### 4.1 美化proto unmarshal error response

**美化前**
```shell
POST http://localhost:8000/svc.admin.AdminExt/GetUser

HTTP/1.1 200 OK
Server: fasthttp
Date: Sun, 03 Sep 2023 02:55:58 GMT
Content-Type: application/json
Content-Length: 138

{
  "code": 500,
  "msg": "grpc: error unmarshalling request: json: cannot unmarshal number into Go struct field GetUserReq.uids of type []int64"
}
```

**美化后**
```shell
POST http://localhost:8000/svc.admin.AdminExt/GetUser

HTTP/1.1 200 OK
Server: fasthttp
Date: Sun, 03 Sep 2023 02:56:37 GMT
Content-Type: application/json
Content-Length: 153

{
  "code": 400,
  "msg": "ErrBadRequest ➜ /svc.admin.AdminExt/GetUser ➜ json: cannot unmarshal number into Go struct field GetUserReq.uids of type []int64"
}
```

### 5. 本地（dev）环境启动微服务的原理

理论上来说，调用微服务是走注册中心的，要想在本地启动多个微服务且能正常互相调用，又不想在本地部署一个类似etcd/consul/zookeeper
的注册中心，最简单的办法是：

实现一个简单的注册中心模块，然后**在开发环境**随服务启动。

- [~~网络协议之mDNS~~（由于windows支持不完善，不再采纳）](https://www.cnblogs.com/Alanf/p/8653223.html)
- [simple_sd实现](./xvendor/simple_sd)

注意：dev环境启动的微服务仍然连接的是**beta环境的数据库**。


### 6. 工具下载（更新）

#### 6.1 下载protoc

linux、mac版本都已经包含在本仓库的`tool/`,`tool_mac/`目录下，无需再下载，已下载的是protoc
v24版本，其余插件也是编写本文档时的最新版本（下载时间2023年8月17日）。

如需更换版本，可点击下方链接自行下载：

https://github.com/protocolbuffers/protobuf/releases

> windows环境暂未支持，请自行配置环境。  
> 本模板配套的是shell脚本，在windows环境运行可能有问题，（但仍然建议使用类unix环境进行开发，以减少不必要的工作和麻烦）。

#### 6.2 下载protoc插件

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

### 7. 其他

#### 计划
- 支持K8s部署

#### 建议

- `protocol/`是存放生成协议代码的目录，在实际项目开发中可以加入`.gitignore`文件，以避免在PR review时产生困扰；

#### 使用的外部库

- github.com/valyala/fasthttp v1.49.0
- github.com/hashicorp/consul/api v1.24.0
- github.com/k0kubun/pp v2.4.0+incompatible
- github.com/pkg/errors v0.9.1
- github.com/redis/go-redis/v9 v9.1.0
- github.com/spf13/viper v1.16.0
- go.uber.org/zap v1.21.0
- google.golang.org/genproto/googleapis/api v0.0.0-20230726155614-23370e0ffb3e
- google.golang.org/grpc v1.57.0
- google.golang.org/protobuf v1.31.0
- gorm.io/driver/mysql v1.5.1
- gorm.io/gorm v1.25.3
- github.com/samber/lo v1.38.1

#### 资源链接

- [Consul 官网介绍](https://developer.hashicorp.com/consul/docs/intro)
- [Consul 服务发现原理](https://developer.hashicorp.com/consul/docs/concepts/service-discovery)