## 在K8s上部署此项目

因为K8s原生支持使用DNS作为服务（Pod）的寻址方式，所以如果使用K8s部署此项目，我们无需再部署类似Consul/etcd等组件用于服务注册/发现。

本文档介绍如何使用Helm在K8s上部署此项目。

首先，本项目已经使用`helm create CHART_NAME`创建了一个Helm chart，位于`根/deploy/service-chart`。这是一个通用的微服务chart，
在部署时会引用每个微服务目录中的`values.yaml`中的配置（例如`service/user/helmValues.yaml`），并使用`templates`
目录下的模板文件生成K8s资源。

也就是说，所有微服务共用于一套helm chart模板，各自管理自己的`values.yaml`即可。

> 你可以查看 [Helm手记](https://github.com/chaseSpace/k8s-tutorial-cn/blob/main/doc_helm.md) 来快速入门Helm。


部署步骤如下：

- 构建服务镜像并推送至仓库
- 部署服务的helm chart

下面的步骤以部署dev环境的user和gateway服务为例：

```shell
# 在shell中进入项目根目录，执行构建脚本
$ sh build_image.sh user dev

# 部署helm chart
$ helm install go-svc-user ./deploy/go-svc-chart -f ./service/user/helmValues.yaml

# 查看helm部署
$ helm ls                                                                         
NAME          	NAMESPACE	REVISION	UPDATED                                	STATUS  	CHART             	APP VERSION
go-svc-gateway	default  	1       	2023-12-05 22:49:55.862068237 +0800 CST	deployed	go-svc-chart-0.1.0	1.16.0     
go-svc-user   	default  	1       	2023-12-05 22:56:08.625016541 +0800 CST	deployed	go-svc-chart-0.1.0	1.16.0

# 查看chart创建的k8s资源
$ helm status go-svc-user --show-resources

# 以同样的步骤部署gateway
sh build_image.sh gayeway dev
helm install go-svc-gayeway ./deploy/go-svc-chart -f ./service/gayeway/helmValues.yaml
```

访问服务：

```shell
$ kubectl get svc                                                    
NAME             TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
go-svc-gateway   ClusterIP   20.1.235.201   <none>        8000/TCP   103m
go-svc-user      ClusterIP   20.1.171.170   <none>        3000/TCP   97m
kubernetes       ClusterIP   20.1.0.1       <none>        443/TCP    16d

# 通过service ip访问gateway
$ curl 20.1.235.201:8000/ping                                 
pong
# 通过service ip访问user的注册接口（gateway转发）
$ curl -X POST  http://20.1.235.201:8000/forward/svc.user.UserExt/SignUp -d {}
{"code":400,"msg":"无效昵称或超出长度","data":null}

# 若要通过ingress访问
# 先获取ingress控制器的svc ip和端口，这里是20.1.166.146:30189
$ kubectl get svc -ningress-nginx                                             
NAME                                 TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             LoadBalancer   20.1.166.146   <pending>     80:30189/TCP,443:30415/TCP   10d
ingress-nginx-controller-admission   ClusterIP      20.1.75.220    <none>        443/TCP                      10d

# 然后访问节点的30189端口会自动转发至nginx控制器的80端口
$ curl 127.0.0.1:30189/ping                                 
pong
$ curl 127.0.0.1:30189/forward/svc.user.UserExt/SignUp -d {}
{"code":400,"msg":"无效昵称或超出长度","data":null}
```

此后，若要更新服务，按下面的步骤进行：

- 修改代码
 - 推送项目
- 使用Helm更新发布：`helm upgrade go-svc-$SVC ./deploy/go-svc-chart -f $helmValues.yaml --set image-tag=x.x.x --description <更新说明>`

