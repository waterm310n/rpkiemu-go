# openssl介绍
openssl获取自签名证书
```bash
openssl req -new \ 
        # -new 生成新的证书请求以及私钥，默认为1024比特。
        -newkey rsa:4096 -keyout issuer.key \ 
        # -newkey rsa:bits
        # -keyout filename 指定生成的私钥文件,默认文件名privkey.pem
        -x509 -out issuer.crt \
        # -x509 生成自签名证书
        # -out 指定输出文件名
        -days 3650 -nodes -subj "/C=NL/L=Amsterdam/O=NLnet Labs"
        # -days 指定自签名证书的有效期限
        # -nodes 全称no des ,当作为参数给出时，这意味着OpenSSL不会加密PKCS＃12文件中的私钥。
        # -subj arg 用于指定生成的证书请求的用户信息，或者处理证书请求时用指定参数替换。
```
数字证书中主题(Subject)中字段的含义
|字段名|字段值|
|--|--|
|公用名称 (Common Name)|简称：CN 字段，对于 SSL 证书，一般为网站域名；而对于代码签名证书则为申请单位名称；而对于客户端证书则为证书申请者的姓名；|
|单位名称 (Organization Name)|简称：O 字段，对于 SSL 证书，一般为网站域名；而对于代码签名证书则为申请单位名称；而对于客户端单位证书则为证书申请者所在单位名称；|
|所在城市 (Locality)|简称：L 字段 |
|所在省份 (State/Provice)|简称：S 字段 |
|所在国家 (Country)|简称：C 字段，只能是国家字母缩写，如中国：CN |

其他还有一些字段没用到就不展示了

openssl创建证书请求
```bash
openssl req -new \
        -newkey rsa:4096 -keyout subject.csr \
        -out subject.crt \
        -days 3650 -nodes -subj "/C=NL/L=Amsterdam/O=NLnet Labs/CN=nginx.krill.testbed"
# 与一条命令的差别在于没有使用-x509选项
```
openssl对证书请求签
```bash
openssl x509 \
        -in subject.csr -req -out subject.crt -extfile subject.ext \
        # -in filename 指定输入文件名
        # -req 输入为证书请求，需要进行处理
        # -out filename 指定输出文件名
        # -extfile filename 指定包含证书扩展项的文件名，如果没有，那么生成的证书将没有任何扩展项
        -CA issuer.crt -CAkey issuer.key -CAcreateserial \
        # -CA arg  设置CA文件，必须为PEM格式
        # -CAkey arg  设置CA私钥文件，必须为PEM格式
        # -CAcreateserial 如果序证书列号文件，则生成
        -days 3650
        # -days arg 设置证书有效期
```
openssl查看证书信息
```bash
openssl x509 -in issuer.crt -noout -text 
# X509命令是一个多用途的证书工具。它可以显示证书信息、转换证书格式、签名证书请求以及改变证书的信任设置等
# -in filename 指定输入文件名
# -text 打印证书信息
# -noout 不显示输入文件的内容
```

## 查看ROA的命令
```bash
openssl asn1parse -in roa_file_name --inform DER
```

# X509拓展字段

##  X509v3 extensions
|字段名|字段值|
|--|--|
|Subject Key Identifier|唯一标记了当前证书中的公钥|
|Authority Key Identifier|唯一标记了该证书签发私钥对应的公钥|
举个例子，一个设备当中可能有好多个证书，但是每一个app可能只需要对应证书中的公钥，那么就可以将这一串ID内置在app中，证书也不需要完整解析，先去看看对应Key Identifier字段是不是匹配即可。通常计算方式为计算公钥的SHA1值，本例子中的Subject Key Identifier即为这种方式：
subjectAltName 是 X.509 version 3 的一个扩展项，该扩展项用于标记和界定证书持有者的身份

# 浏览器如何验证证书
当浏览器使用HTTPS连接到您的服务器时，他们会检查以确保您的SSL证书与地址栏中的主机名称匹配。

浏览器有三种找到匹配的方法：

1.主机名（在地址栏中）与证书主题(Subject)中的通用名称(Common Name)完全匹配。

2.主机名称与通配符通用名称相匹配。例如，www.example.com匹配通用名称* .example.com。

3.主机名 在主题备用名称(SAN: Subject Alternative Name)字段中列出

# 参考资料
[证书各个字段的含义](https://www.cnblogs.com/iiiiher/p/8085698.html)

[openssl命令](https://www.openssl.net.cn/docs/32.html)

[X.509系列（一）：X.509 v3格式下的证书](https://www.cnblogs.com/xiaoxi-jinchen/p/15434662.html)