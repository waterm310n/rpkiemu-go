# CA方搭建

## BGP部分
使用[bgpemu](https://github.com/V3rgilius/bgpemu)创建BGP网络中的节点
**注：要求拓扑文件说明CA所处的位置并且与本项目的.rpkiemu-go.json文件中的publishPoints对应**
```bash
$ bgpemu topo create < 拓扑文件 >
# 通常拓扑文件以.yaml格式结尾
# 例如bgpemu topo create topo.yaml
```
接着创建BGP网络中的路径
```bash
$ bgpemu lab deploy < 场景文件 >
# 通常拓扑文件以.yaml格式结尾
# 例如bgpemu lab deploy scene.yaml
```

## RPKI部分

### 初始化容器服务
```bash
$ ./rpkiemu-go ca setup
# 根据.rpkiemu-go.json文件中的publishPoints所述，对bgpemu中的相应pod节点中的krill容器进行初始化

```

### 搭建CA层次结构
```bash
$ rpkiemu-go ca create -d < 数据目录 >
# 默认情况会在执行命令所在的目录使用examples目录作为数据目录
```

### 发布ROAS
```bash
$ rpkiemu-go ca publish -d < 数据目录 >
# 默认情况会在执行命令所在的目录使用examples目录作为数据目录
```

### 依赖方初始化
```bash
$ rpkiemu-go rp setup
# 根据.rpkiemu-go.json文件中的relyParties所述，对bgpemu中的相应pod节点中的routinator容器进行初始化
```

### 
## 系统结束运行
```bash
$ bgpemu topo delete < 拓扑文件 >
# 通常拓扑文件以.yaml格式结尾
# 例如bgpemu topo delete topo.yaml
```