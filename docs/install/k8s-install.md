# Kubernetes配置

## 前置准备

使用一台 `ubuntu2004 server`的虚拟机

下为vmware workstation pro下的一台机器ip与系统数据

| hostname | ip              | system             |

| -------- | --------------- | ------------------ |

| master   | 192.168.106.143 | ubuntu20.04 server |

### 关闭防火墙

```bash

#查看防火墙状态

sudoufwstatus

#如果是inactive无需执行下面命令 

sudosystemctlstopfirewalld.service

sudosystemctldisablefirewalld.service

```

### 关闭selinux

`ubuntu2204`无需此操作

`centos7`需执行下面命令

```bash

getenforce

cat/etc/selinux/config

sudosetenforce0

sudosed-i's/^SELINUX=enforcing$/SELINUX=permissive/'/etc/selinux/config

cat/etc/selinux/config

```

### 关闭swap

```bash

master@master:~$free-m

               total        used        free      shared  buff/cache   available

Mem:            3876         328        2942           1         604        3323

Swap:              0           0           0

#可以看到并没有使用swap

#如果有的话编辑/etc/fstab，将下面这一行注释

#/swap.img      none    swap    sw      0       0

#可以使用该命令`sudo sed -i 's/.*swap.*/#&/' /etc/fstab`

```

### 修改内核参数以安装k8s依赖

```bash

cat<<EOF| sudo tee /etc/modules-load.d/k8s.conf

overlay

br_netfilter

EOF

sudomodprobeoverlay

sudomodprobebr_netfilter

```

```bash

# 设置所需的 sysctl 参数，参数在重新启动后保持不变

cat<<EOF| sudo tee /etc/sysctl.d/k8s.conf

net.bridge.bridge-nf-call-iptables  = 1

net.bridge.bridge-nf-call-ip6tables = 1

net.ipv4.ip_forward                 = 1

EOF

# 应用 sysctl 参数而不重新启动

sudosysctl--system

```

```bash

# 通过运行以下指令确认 br_netfilter 和 overlay 模块被加载：

lsmod|grepbr_netfilter

lsmod|grepoverlay

```

### 修改apt源

将 `apt`源全部修改为清华开源镜像源（其他源也可以）

1.备份旧的镜像源

```bash

sudocp/etc/apt/sources.list/etc/apt/sources.list.bak

```

2.修改 `/etc/apt/sources.list`如下

```

# 默认注释了源码镜像以提高 apt update 速度，如有需要可自行取消注释

deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ jammy main restricted universe multiverse

# deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ jammy main restricted universe multiverse

deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ jammy-updates main restricted universe multiverse

# deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ jammy-updates main restricted universe multiverse

deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ jammy-backports main restricted universe multiverse

# deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ jammy-backports main restricted universe multiverse

# deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ jammy-security main restricted universe multiverse

# # deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ jammy-security main restricted universe multiverse

deb http://security.ubuntu.com/ubuntu/ jammy-security main restricted universe multiverse

# deb-src http://security.ubuntu.com/ubuntu/ jammy-security main restricted universe multiverse

# 预发布软件源，不建议启用

# deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ jammy-proposed main restricted universe multiverse

# # deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ jammy-proposed main restricted universe multiverse

```

3.更新系统软件

```bash

sudoaptupdate && sudoaptupgrade

```

## 安装docker

为每台虚拟机都安装 `docker`

#### 设置apt仓库

1. 更新 `apt`包索引并安装包，以允许 `apt`通过 `HTTPS` 使用仓库:

```bash

sudoapt-getupdate

sudoapt-getinstallca-certificatescurlgnupg

```

2. 添加 `docker`的GPG密钥

```bash

sudoinstall-m0755-d/etc/apt/keyrings

curl-fsSLhttps://download.docker.com/linux/ubuntu/gpg|sudogpg--dearmor-o/etc/apt/keyrings/docker.gpg

sudochmoda+r/etc/apt/keyrings/docker.gpg

```

3. 使用下列命令配置仓库

```bash

echo\

  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \

  "$(. /etc/os-release && echo "$VERSION_CODENAME")"stable" | \

  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

```

### 安装docker引擎

1. 更新 `apt`包索引并允许 `apt`可以使用https连接仓库

```bash

sudoapt-getupdate

```

2. 安装 `Docker Engine`, `containerd`,和 `Docker Compose`

```bash

sudoapt-getinstalldocker-cedocker-ce-clicontainerd.iodocker-buildx-plugindocker-compose-plugin

```

```bash

sudosystemctlstopcontainerd.service

sudocp/etc/containerd/config.toml/etc/containerd/config.toml.bak

sudocontainerdconfigdefault>$HOME/config.toml

sudocp$HOME/config.toml/etc/containerd/config.toml

# 修改 /etc/containerd/config.toml 文件后，要将 docker、containerd 停止后，再启动

sudosed-i"s#registry.k8s.io/pause#registry.aliyuncs.com/google_containers/pause#g"/etc/containerd/config.toml

# https://kubernetes.io/zh-cn/docs/setup/production-environment/container-runtimes/#containerd-systemd

# 确保 /etc/containerd/config.toml 中的 disabled_plugins 内不存在 cri

sudosed-i"s#SystemdCgroup = false#SystemdCgroup = true#g"/etc/containerd/config.toml


# containerd 忽略证书验证的配置

#      [plugins."io.containerd.grpc.v1.cri".registry.configs]

#        [plugins."io.containerd.grpc.v1.cri".registry.configs."192.168.0.12:8001".tls]

#          insecure_skip_verify = true



sudosystemctlenable--nowcontainerd.service

# sudo systemctl status containerd.service


# sudo systemctl status docker.service

sudosystemctlstartdocker.service

# sudo systemctl status docker.service

sudosystemctlenabledocker.service

sudosystemctlenabledocker.socket

sudosystemctllist-unit-files|grepdocker


sudomkdir-p/etc/docker


sudotee/etc/docker/daemon.json<<-'EOF'

{

  "registry-mirrors": [

  "https://registry.docker-cn.com",

  "http://hub-mirror.c.163.com",

  "https://docker.mirrors.ustc.edu.cn"

  ],

  "exec-opts": ["native.cgroupdriver=systemd"],

  "proxies": {

    "http-proxy": "http://192.168.106.140:7890",

    "https-proxy": "http://192.168.106.140:7890",

    "no-proxy": "*127.0.0.0/8"

  }

}

EOF


sudosystemctldaemon-reload

sudosystemctlrestartdocker

```

更多详细内容可见[Install Docker Engine on Ubuntu | Docker Documentation](https://docs.docker.com/engine/install/ubuntu/)

与[Linux post-installation steps for Docker Engine | Docker Documentation](https://docs.docker.com/engine/install/linux-postinstall/)

## 安装kubeadm，kubectl，kubelet

- `kubeadm`：用来初始化集群的指令。
- `kubelet`：在集群中的每个节点上用来启动 Pod 和容器等。
- `kubectl`：用来与集群通信的命令行工具。

这里使用阿里云开源镜像站安装，下列的命令请在 `root`权限下运行

```bash

apt-getupdate && apt-getinstall-yapt-transport-https

curlhttps://mirrors.aliyun.com/kubernetes/apt/doc/apt-key.gpg|apt-keyadd-

cat<<EOF>/etc/apt/sources.list.d/kubernetes.list

deb https://mirrors.aliyun.com/kubernetes/apt/ kubernetes-xenial main

EOF

apt-getupdate

apt-getinstall-ykubelet=1.26.3-00kubeadm=1.26.3-00kubectl=1.26.3-00

```

## 启用集群

```bash

kubeadmconfigimagespull--image-repositoryregistry.aliyuncs.com/google_containers--cri-socketunix:///var/run/cri-dockerd.sock

```

在master节点上启用

```bash

kubeadminit\

--image-repository=registry.aliyuncs.com/google_containers 

```

## 配置[flannel](https://github.com/flannel-io/flannel) cni

```bash
kubectl apply -f https://github.com/flannel-io/flannel/releases/latest/download/kube-flannel.yml

```
##### node节点加入（如果有多台机器）

在node节点上执行

```bash

kubeadmjoin192.168.106.140:6443--tokenssy6d5.44hrz8njkbajcgju\

        --discovery-token-ca-cert-hash sha256:419b066b5857aedad83ad0bedf725520ee95dcf7e8b733f66750813737103706

```

## 重置集群

```bash

kubeadmreset--cri-socketunix:///var/run/cri-dockerd.sock

sudorm-f/etc/cni/net.d/10-flannel.conflist

sudoiptables-F && sudoiptables-tnat-F && sudoiptables-tmangle-F && sudoiptables-X

sudorm-f$HOME/.kube/config

```

## 问题集

征求中

## 参考

参考[kubernetes（k8s）安装](https://www.yuque.com/xuxiaowei-com-cn/gitlab-k8s/k8s-install)
