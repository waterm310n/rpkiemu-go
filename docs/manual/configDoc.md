# 配置选项说明

## databaseConfig

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
|键| 说明| 例子|
|--|--|--|
|kubeConfig|k8s配置文件位置| /home/master/.kube/config |
