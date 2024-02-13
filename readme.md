# rpkiemu-go

[![许可证](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE) ![Go](https://img.shields.io/badge/Language-go-blue)

`rpkiemu-go`是一个go语言实现的，基于容器技术的，RPKI网络仿真系统

## 目录
- [rpkiemu-go](#rpkiemu-go)
  - [目录](#目录)
  - [安装](#安装)
  - [使用](#使用)
  - [致谢](#致谢)
  - [许可证](#许可证)

## 安装
见[安装](docs/install/install.md)

## 使用
本项目的配置文件为`.rpkiemu-go.json`，关于其中的所有选项介绍见[配置说明](docs/configDoc.md)

### 数据生成

```bash
rpkiemu-go ca generate 
```

### RPKI证书体系创建
见[rpkiCreate](docs/manual/rpkiCreate.md)

### RPKI恶意攻击模拟
见[rpkiAttack](docs/manual/rpkiAttack.md)
## 致谢
感谢这个B[V3rgilius/bgpemu (github.com)](https://github.com/V3rgilius/bgpemu)

## 许可证
本项目采用 MIT 许可证。详细信息请查阅 [LICENSE](LICENSE) 文件。