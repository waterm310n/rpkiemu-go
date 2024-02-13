# 介绍
本文件夹包含了一些本项目所使用到的docker镜像打包
打包的脚本来自[rpki-deploy](https://github.com/NLnetLabs/rpki-deploy)

|镜像名|作用|
|---|---|
|krill:local|负责搭建CA架构并签署roa|
|routinator:local|负责验证krill的数据|
|nginx:local|负责代理krill以支持rrdp|
|rsyncd:local|负责代理krill以支持rsync服务|

## 镜像构建
```bash
docker buildx build -t krill:local ./krill
docker buildx build -t rsyncd:local ./rsyncd
docker buildx build -t routinator:local ./routinator
docker buildx build -t gobgp:local ./gobgp
```

# 帮助指令
批量删除docker中的<none>镜像
```bash
docker rmi $(docker images|grep "^<none>"|awk '{print $3}')
```
批量删除docker中未运行的镜像
```bash
docker rm $(docker ps -qf status=exited)
```
批量删除运行中的docker镜像
```bash
docker stop $(docker ps -q) && docker rm $(docker ps -qf status=exited)
```