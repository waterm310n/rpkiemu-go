# 配置选项说明

## databaseConfig
数据库配置文件
|键| 说明| 例子|
|--|--|--|
|host|数据库主机地址|172.19.220.163|
|user|数据库用户名|root|
|password|数据库密码|123456|
|port|数据库端口|3306|
|database|数据库名称|rpki_db|
|tables|使用的数据库表的映射|{"cas": "cas_new","roas":"roas_new"}|
|rirs|查找起点（5大RIR名称）|["RIPE"]|
|ases|查找的目标AS（将其证书链查找出来）|["AS8393"]|
|limitLayer|查找深度|5|

## kubeConfig
kubectl配置文件位置
|键| 说明| 例子|
|--|--|--|
|kubeConfig|k8s配置文件位置| /home/master/.kube/config |

## publishPoints
该配置的目的是：明确CA发布点在bgpemu中对应命名空间下的对应pod中的对应容器名，即CA方在bgp网络中所在的节点位置。
如果不设置该配置选项，那么就无法在bgp环境中创建CA方容器，进而无法创建RPKI证书体系

该配置以Hash表方式组成，其中键值为发布点名称，对应的值类型说明如下
|键| 说明| 值类型| 例子|
|--|--|--|--|
|namespace| 命名空间| 字符串 | "bgp" |
|podName| pod名称 | 字符串| "r5" |
|CAcontainerName| CA容器名称 | 字符串| "r5-ripe" |
|RSYNCDcontainerName| RSYNCD容器名称 | 字符串| "r5-rsyncd" |
|isRIR| 该发布点是否是RIR | 布尔值| true |

## relyParties
该配置的目的是：明确依赖方在bgpemu中对应命名空间下的对应pod中的对应容器名，即依赖方在bgp网络中所在的节点位置。
如果不设置该配置选项，那么就无法在bgp环境中创建依赖方方容器，进而bgp网络中的节点无法通过依赖方获取路由过滤表（VRP）。

该配置同样以Hash表方式组成，其中键值为依赖方名称，对应的值类型说明如下。
|键| 说明| 值类型| 例子|
|--|--|--|--|
|namespace| 命名空间| 字符串 | "bgp" |
|podName| pod名称 | 字符串| "r4" |
|containerName| 容器名称 | 字符串| "r4-routinator" |

## 一份正常的配置文件示例
```json
{
    "databaseConfig":{
        "host": "192.168.106.140", //mysql数据库地址
        "user": "root", //mysql数据库用户
        "password": "123456", //mysql数据库密码
        "port": 3306, //mysql数据库端口
        "database": "rpki_db", // 使用的数据库
        "rirs": ["RIPE"], // 使用哪些rirs
        "tables": {
            "cas": "cas_new", //cas表对应表名
            "roas":"roas_new" //roas表对应表名
        },
        "ases": ["AS8393"], //ASes对应表明
        "limitLayer":5 //层数
    },
    "kubeConfig":"/home/master/.kube/config",//k8s配置文件
    "publishPoints":{
       // 在bgpemu的bgp命名空间下，pod r4 运行CA容器，并且该容器名字为r4-ripe
       "ripe": {
            "namespace": "bgp",
            "pod_name": "r4",
            "ca_container_name": "r4-ripe",
            "rsyncd_container_name": "r4-rsyncd",
            "is_rir": true
        },
       // 在bgpemu的bgp命名空间下，pod r2-0 运行CA容器，并且该容器名字为r2-ca
       "r2-child":{
            "namespace":"bgp",
            "podName":"r2-0",
            "containerName":"r2-ca",
            "isRIR":false
       }
    },
    "relyParties":{
        // 在bgpemu的bgp命名空间下，pod r5 运行RP容器，并且该容器名字为r5-routinator
        "r5-routinator": {
            "namespace": "bgp",
            "pod_name": "r5",
            "ca_container_name": "r5-routinator"
        }
    }
}
```

