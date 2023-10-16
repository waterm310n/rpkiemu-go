# CA方搭建

## 搭建CA层次结构
```bash
$ rpkiemu-go ca create -d < 数据目录 >
# 默认情况会在执行命令所在的目录使用examples目录作为数据目录
```

## 发布ROAS
```bash
$ rpkiemu-go ca publish -d < 数据目录 >
# 默认情况会在执行命令所在的目录使用examples目录作为数据目录
```