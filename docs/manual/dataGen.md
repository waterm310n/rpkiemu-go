# 数据生成

启动mysql容器
```bash
docker run --net=host -p 3306:3306 --name mysql-docker -v ~/mysql/data:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 --privileged=true --restart=always -d mysql
```