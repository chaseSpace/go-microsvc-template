## Go微服务模板

一个简洁、清爽的微服务项目架构，从变量命名到不同职责的（多层）目录结构定义。

> **完成进度：40%**

计划支持以下模式或特性：

- 使用单仓库多服务模式
- 使用grpc+protobuf作为内部rpc通讯协议
- 使用grpc-gateway插件生成grpc服务的http反向代理
- 使用consul作为服务注册发现组件，支持扩展
- 使用gorm作为orm组件，支持扩展
- 使用redis作为cache组件，支持扩展
- 支持本地启动**多个**微服务

其他有用的特性：

- shell脚本支持mac环境（默认linux）
- 定义了err类型，以便在跨服务传播

### 1. 目录结构释义

```
├── consts  # 常亮
├── enums   # 枚举
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

### 2. 工具下载

#### 下载protoc

linux、mac版本都已经包含在本仓库的`tool/`,`tool_mac/`目录下，无需再下载，已下载的是protoc v24版本，其余插件也是编写本文档时的最新版本（下载时间2023年8月17日）。

如需更换版本，可点击下方链接自行下载：

https://github.com/protocolbuffers/protobuf/releases

> windows环境暂未支持，请自行配置环境。  
> 本模板配套的是shell脚本，在windows环境运行可能有问题，（但仍然建议使用类unix环境进行开发，以减少不必要的工作和麻烦）。

#### 下载protoc插件

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.16
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.16
```

下载pb代码引用的库(否则生成的pb文件会报红)：

```shell
go get google.golang.org/grpc@v1.57.0
```