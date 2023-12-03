# 此文件用于构建每一个Go服务
# 用法如下：
# docker build --build-arg SVC=user --build-arg MICRO_SVC_ENV=dev . -t leigg/go-svc-user:1.0.0
# 建议直接使用脚本 build_image.sh

FROM golang:1.20-alpine AS builder

ARG SVC
#ARG MICRO_SVC_ENV
WORKDIR /go/cache
COPY go.mod .
COPY go.sum .
RUN GOPROXY=https://goproxy.cn,direct go mod download

WORKDIR /build

COPY . .

# 关闭cgo的原因：使用了多阶段构建，go程序的编译环境和运行环境不同，不关就无法运行go程序
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 cd service/$SVC && go build -tags=k8s -o main -ldflags "-w -extldflags -static"

#FROM scratch as prod
FROM alpine as prod
# 通过 http://www.asznl.com/post/48 了解docker基础镜像：scratc、busybox、alpine

# 参数需要再次指定才会有效
ARG SVC
ARG MICRO_SVC_ENV

# alpine设置时区
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories &&  \
    apk add -U tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && apk del tzdata && date


# 创建程序读取的目录结构
RUN mkdir -p deploy/$MICRO_SVC_ENV && mkdir -p service/$SVC/deploy/$MICRO_SVC_ENV

# 拷贝全局配置
COPY --from=builder /build/deploy/$MICRO_SVC_ENV deploy/$MICRO_SVC_ENV
# 拷贝服务专有配置
COPY --from=builder /build/service/$SVC/deploy/$MICRO_SVC_ENV service/$SVC/deploy/$MICRO_SVC_ENV
# 拷贝二进制
COPY --from=builder /build/service/$SVC/main .

ENTRYPOINT ["/main"]
