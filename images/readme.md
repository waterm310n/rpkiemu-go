# 介绍
本文件夹包含了一些本项目所使用到的docker镜像打包
打包的脚本来自[rpki-deploy](https://github.com/NLnetLabs/rpki-deploy)

|镜像名|作用|
|---|---|
|krill:local|负责搭建CA架构并签署roa|
|routinator:local|负责验证krill的数据|
|nginx:local|负责代理krill以支持rrdp|
|rsyncd:local|负责代理krill以支持rsync服务|

# 前置要求

## 镜像构建
```bash
docker buildx build -t krill:local ./krill
docker buildx build -t nginx:local ./nginx
docker buildx build -t rsyncd:local ./rsyncd
docker buildx build -t routinator:local ./routinator
docker buildx build -t gobgp:local ./gobgp
```
# 运行流程

## 使用host网络自建
运行krill
```bash
docker run -d  --network host --name krill\
 -v krill_data:/var/krill/data/ \
 -v krill_rsync:/var/krill/data/repo/rsync \
 -e KRILL_CLI_TOKEN=krillTestBed \
 -v ./krill.conf:/var/krill/data/krill.conf \
 krill:local
```
然后运行nginx
```bash
docker run -d --network host --name krill.testbed nginx:local
```
接着运行rsyncd
```bash
docker run -d --network host --name rsyncd \
 -v krill_rsync:/share:ro \
 rsyncd:local
```
运行routinator
```bash
docker run -d --network host --name routinator  -e SRC_TAL=ta.tal  routinator:local
```

### 注意事项
以上四个docker都以主机的网络连接。
需要为主机修改`/etc/hosts`，将`krill`解析为本机地址

## 使用bridge网络自建
```bash
docker network create rpki_test
```
运行带有仓库的krill端
```bash
docker run -d --network rpki_test --name krill\
 -v krill_data:/var/krill/data/ \
 -v krill_rsync:/var/krill/data/repo/rsync \
 -v ./krill/krill.conf:/var/krill/data/krill.conf \
 -v krill_data_share:/tmp \
 -e KRILL_CLI_TOKEN=krillTestBed \
 nlnetlabs/krill:v0.12.1
```
假设有其他的krill运行
```bash
docker run -d --network rpki_test --name krill2\
 -v krill2_data:/var/krill/data/ \
 -v krill2_rsync:/var/krill/data/repo/rsync \
 -v krill_data_share:/tmp \
 -v ./krill/krill2.conf:/var/krill/data/krill.conf \
 -e KRILL_CLI_TOKEN=krillTestBed2 \
 nlnetlabs/krill:v0.12.1
```
运行nginx
```bash
docker run -d --network rpki_test --name nginx.krill.testbed nginx:local
```
运行rsync
```bash
docker run -d --network rpki_test --name rsyncd.krill.testbed \
 -e SRC_CER=https://nginx.krill.testbed/ta/ta.cer\
 -v krill_rsync:/share:ro \
 rsyncd:local
```
运行routinator
```bash
docker run -d --network rpki_test --name routinator -e SRC_TALS=https://krill:3000/ta/ta.tal  -p 9556:9556 -p 3323:3323 routinator:local

```
```bash
docker run -it --network rpki_test --rm alpine:3.17.2 #这里使用busybox的话会连接不上，存在问题。
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