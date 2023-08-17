## Go微服务模板

一个简洁、清爽的微服务项目架构，从变量命名到不同职责的（多层）目录结构定义。

>**完成进度：30%**

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

linux、mac版本都已经包含在本仓库的`tool/`,`tool_mac/`目录下，无需再下载，已下载的是protoc v24版本，其余插件也是编写本文档时的最新版本（下载时间2023年8月17日）。

如需更换版本，可点击下方链接自行下载：

https://github.com/protocolbuffers/protobuf/releases

> windows环境暂未支持，请自行配置环境。  
> 本模板配套的是shell脚本，在windows环境运行可能有问题，（但仍然建议使用类unix环境进行开发，以减少不必要的工作和麻烦）。

#### 下载protoc插件

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

下载pb代码引用的库(否则生成的pb文件会报红)：

```shell
go get google.golang.org/grpc@v1.57.0
```