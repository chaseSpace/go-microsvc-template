## 使用TLS微服务之间的通信安全

### 1. 介绍
微服务之间的通信如果不进行加密、授权可能导致以下问题和安全风险：

- 未授权访问：没有认证的微服务之间，任何服务都可以尝试连接到其他服务，这可能导致未经授权的访问。这意味着不受信任的服务可能会访问敏感数据或执行敏感操作。

- 数据泄漏：如果微服务之间的通信不受保护，敏感数据可能在传输过程中被窃听或截取。这可能导致数据泄漏和隐私问题。

- 数据篡改：未经认证的通信可能容易受到中间人攻击（Man-in-the-Middle, MitM），攻击者可以修改或篡改数据包，导致数据的完整性受到威胁。

- 服务冒充：没有认证的情况下，攻击者可能冒充合法的微服务，向其他服务发送虚假请求，从而欺骗其他服务执行不安全的操作或泄露敏感信息。

- 缺乏审计和跟踪：没有认证的情况下，很难跟踪和审计谁访问了哪些服务和数据。这会使问题排查和监视变得更加困难。

保障微服务之间的通信安全有多种方式，常见的如下：
1. 基于TOKEN，如JWT。令牌可以包含一些信息，如服务名称、版本、签发时间和有效期等；
2. 基于证书的加密与授权。为每个微服务生成不同的证书，服务之间的rpc开启TLS通信。此步骤不用编码，但需要管理证书；
3. OAuth 2.0。一个微服务可以充当 OAuth 2.0 服务提供者 (Authorization Server)，为其他微服务颁发访问令牌，并验证这些令牌，需要较多编码工作；
4. 单点登录（SSO）：相当于一个TOKEN在所有服务有效，编码不多，但不够安全；
5. 服务网格 (Service Mesh)：使用服务网格技术（如 Istio 或 Envoy）可以在微服务之间提供自动的身份认证和安全性。服务网格可以处理许多安全性和认证方面的细节，包括流量加密、认证、授权和监控；
   - 唯一缺点，增加架构复杂度；

**综上，笔者选择第二种，即证书。**

- 对于加密：TLS本身即是对流量的加密
- 对于授权：Server端可以对Client证书进行检查，只放行特定的Client证书，对于不识别或不允许的证书直接拒绝建立连接。

### 2. 操作步骤

#### 2.1 生成根证书
根证书的作用是作为一个有公信力的CA（证书颁发机构），用来为其他主体颁发证书。
根证书可以验证从自己这颁发出去的证书。

>最好单独保存根证书，一旦泄露，其管理的整个域都不再安全。

在这里，根证书将在后续步骤用于颁发server以及client证书。

>颁发server和client证书时还可以使用自签名的方式，但这种方式无法验证证书身份，不安全。

创建根证书的步骤如下：
```shell
# -keyout 指定输出的私钥文件；-out 指定输出的证书文件
# -days 3650 指定根证书有效期10年
# -subj 指定主体名
# -addext 表示添加扩展项，其中 subjectAltName 表示指定备用主体名
$ openssl req -x509 -newkey rsa:2048 -nodes -days 3650 -keyout ca-key.pem -out ca-cert.pem \
  -subj "/CN=x.microsvc"
  
# 查看根证书内容
$ openssl x509 -in ca-cert.pem -noout -text
...
# 查看证书有效期
$ openssl x509 -in ca-cert.pem -noout -dates                                                                                            1 ↵
notBefore=Sep  2 10:31:52 2023 GMT
notAfter=Aug 30 10:31:52 2033 GMT
```

#### 2.2 为微服务生成server和client证书

- server证书：为指定server颁发的证书，提供给client验证其身份
- client证书：为指定client颁发的证书，提供给server验证其身份

在标准的通信架构中，一般要为每个server生成对应的server证书和client证书，而每个证书都会伴随一个秘钥，
但这样在大型微服务项目中会需要管理大量证书，这会造成较高的维护复杂度。

为了降低维护复杂度，这里建议生成一个通用的server&client证书，当证书发生泄露时，重新签发server&client证书分发至所有服务，
然后重启所有服务即可。

>对于某些Server，可能只允许特定Client访问，可以在为这些特定Client生成证书时，设置特定而非通用的CN属性（Common Name），
> 然后在Server端建立TLS连接时进行CN验证即可。


**这里就是双向认证，提供了较高的安全性。**

>这样的做法也一定程度上降低了微服务之间的通信安全性，需要架构者进行权衡。
> 相比于为每个服务生成证书，更建议使用Server Mesh架构，将服务间通信的加密、授权工作转移给Sidecar容器，以减少单个微服务复杂性。

1. 生成server证书

```shell
# 先生成server私钥和csr（证书签名请求）
# 注意subject指定了一个泛域名，这样才能适配不同服务。同时在代码中用以访问每个服务的虚拟域名也必须是相同的格式，否则不能匹配证书
# 这里的 .default.svc.cluster.local 表示一个k8s集群内通配的域名
$ openssl req -newkey rsa:2048 -nodes -keyout server-key.pem -out server-req.pem \
  -subj "/CN=*.default.svc.cluster.local"

# 再使用ca私钥颁发server证书; 对于server证书，需要指定SAN（subjectAltName），用来替代CN，这个强制要求貌似来自Go的加密库
# -- 这里SAN指定了两个HOST（满足域名格式），在client使用域名连接server时会自动验证，如果使用IP连接则会验证失败！
# -- 可以同时指定IP：subjectAltName=DNS:xxx.default.svc.cluster.local, IP:127.0.0.1
# -- 可以使用通配符匹配多个域名：subjectAltName=DNS:*.default.svc.cluster.local    (但不能匹配多个IP)
# -- 这里添加 IP:127.0.0.1 是为了方便 Dev 环境使用
$ openssl x509 -req -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial \
  -in server-req.pem -out server-cert.pem \
  -days 365 \
  -extfile <(printf "subjectAltName=DNS:*.default.svc.cluster.local, IP:127.0.0.1")

$ ls
ca-cert.pem     ca-cert.srl     ca-key.pem      server-cert.pem server-key.pem  server-req.pem
```

2. 生成client证书

```shell
# 先生成client私钥和csr（证书签名请求）
# 这里的CN有意设置为一个server会拒绝的值，server将会验证下面的 subjectAltName
$ openssl req -newkey rsa:2048 -nodes -keyout client-key.pem -out client-req.pem \
  -subj "/CN=unknown.microsvc"

# 再使用ca证书和私钥颁发client证书
# 这里多指定的 DNS 表示这个为这个client证书添加特定服务的授权（当 CN 验证失败时，server就会验证这个DNS），这就实现了开头所说的某些server仅对特定client授权
$ openssl x509 -req -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial \
  -in client-req.pem -out client-cert.pem -days 365 \
  -extfile <(printf "subjectAltName=DNS:user.client.microsvc,DNS:admin.client.microsvc")

$ ls
ca-cert.pem ca-cert.srl ca-key.pem client-cert.pem client-key.pem client-req.pem server-cert.pem server-key.pem server-req.pem
```

>注意：特定client证书只允许给特定client使用。

#### 2.3 配置微服务以使用证书
在每个微服务的配置中，指定生成的私钥和证书文件，以便它们可以在 TLS 握手期间使用。

具体方式请在代码中搜索：`NewGRPCClient`, `newGRPCServer`

### 3. 安全保管证书
对于服务使用到的server、client以及root证书以及对应的私钥文件，属于非常重要的安全凭据。
一旦泄漏，就大大降低了微服务通信的安全性，需要及时重新生成这些凭据，并更新所有服务，这个过程费时费力。

所以要保管好这些凭据，**不能将它们托管在代码仓库中**，而是单独存放。

另外，在部署时也不能简单拷贝证书文件到服务器上，而是使用其他安全的方式加载，
比如k8s的secret。
