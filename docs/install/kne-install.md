# kne安装

## 创建docker network

```bash
docker network create multinode
```

## 安装make

```bash
sudo apt install make
```

## 安装kne
```bash
make install
```
## 运行kne

```bash
kne deploy deploy/kne/external-multinode.yaml
kubectl taint node master node-role.kubernetes.io/control-plane-
```