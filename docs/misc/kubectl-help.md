# kubectl 介绍

本项目基于kubernetes，因此使用kubectl对节点进行调试是很有帮助的，这里对使用到的kubectl命令进行简要介绍


## 查看所有pods命令
通常在部署bgpemu使用拓扑文件创建完成后，查看节点是否成功创建
```bash
kubectl get pods -A
```

## 进入指定命名空间的指定pod的指定容器的shell中
```bash
kubectl exec -it -n < 命名空间 > < pod名称 > -c < 容器名称 > -- sh
# -n 命名空间，此处默认bgp
# -c 容器名称
# 示例 kubectl exec -it -n bgp r1-0 -c r1-0-apnic -- sh
```