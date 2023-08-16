## Go微服务模板

一个简洁、清爽的微服务项目架构，从变量命名到不同职责的（多层）目录结构定义。

计划支持以下模式或特性：

- 使用单仓库多服务模式
- 使用grpc+protobuf作为内部rpc通讯协议
- 使用grpc-gateway作为http网关
- 使用consul作为服务注册发现组件，支持扩展
- 使用gorm作为orm组件，支持扩展
- 使用redis作为cache组件，支持扩展
- 支持本地启动**多个**微服务

其他有用的特性：
- shell脚本支持mac环境（默认linux）


### 工具下载

#### 下载protoc
linux版本已经包含在本仓库的`tool/`目录下，无需再下载，如果是本地mac开发环境，则需要单独下载。

下载protoc（二进制文件放到项目的`tool_mac/protoc_v24/`目录下，目录需要新建）:

https://github.com/protocolbuffers/protobuf/releases/tag/v24.0

>windows开发环境暂不支持。

#### 下载protoc插件
```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

下载pb代码引用的库(否则生成的pb文件会报红)：
```shell
go get google.golang.org/grpc@v1.57.0
```