# 安装说明

## golang安装

```bash
#安装go
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.1.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin">> ~/.bashrc
#配置七牛云go代理
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```

## k8s安装
本项目基于1.26.3的kubernetes版本
见[k8s安装](k8s-install.md)

### 配置非root权限用户使用docker与kubectl

非root权限用户使用docker
```bash
sudo usermod -aG docker $USER
newgrp docker #激活改变
```

非root权限用户使用kubectl
```bash
#非root权限用户登录下
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.kubeconfig $HOME/.kube/config
sudo chown <username>:<username> $HOME/.kube/config
```

## kne安装
见[kne安装](kne-install.md)

## bgpemu安装
见[bgpemu安装](https://github.com/V3rgilius/bgpemu/blob/master/docs/setup.md)

## 镜像安装
TODO