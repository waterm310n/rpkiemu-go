# 数据生成

生成的数据呈现树形结构,其文件结构大致如下所示

```
.
├── mripe-ncc-ta
├── mripe-ncc-ta_children
│   ├── m102ed0852ea4b4700eae91a42e9f7e0fe53497e6
│   ├── m102ed0852ea4b4700eae91a42e9f7e0fe53497e6_children
│   ├── ma4250f3c598917bc119f7ddd423595bb7251181c
│   ├── ma4250f3c598917bc119f7ddd423595bb7251181c_children
│   │   ├── m1jEEXgH74HRy26tqz6rGi2p4R28
│   │   ├── mCnL-b5z5_hDy2tw1BsZhkUG21hY
│   │   ├── mNtmoRKu7irFWj9OV4aAdDNTs5lU
│   │   ├── mjJSdv3GcxMvclDgsa6pNVv1_v1w
│   │   ├── mjt86YvmixlvlDi5KNkl_9WRWSkk
│   │   └── mxEVPS6F2oPWaqya7YmwaYMAFK4s
│   └── maa2304f03c4e807dae51fd33008c93692b973c40
│       ├── mzpyS25FDQF186nbbBJlhIVgSs3E
│       └── mzxwWMquEOVl0Q9l4xgbK1-rubso
└── ...
```


## 启动mysql容器
```bash
docker run --net=host -p 3306:3306 --name mysql-docker -v ~/mysql/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 --privileged=true --restart=always -d mysql
```
请手动**向mysql中注入RPKI数据**

## 获取RPKI数据（层次结构与ROAs）
首先需要配置`rpkiemu-go.json`文件,样例如下
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
    "kubeConfig":"/home/master/.kube/config", //k8s配置文件
    ...
}
```
然后运行下面的命令生成数据
```bash
$ rpkiemu-go ca generate -d < 数据目录 > 
# 默认情况会在执行命令所在的目录创建examples目录作为数据目录
```
## 将RPKI数据与BGP数据结合
```bash
$ rpkiemu-go ca adapt -i < rpki发布点文件 > -t < bgp拓扑文件 > -o < 带有rpki发布点的bgp拓扑文件 >
```