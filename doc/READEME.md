# 初始化环境

## docker创建基本的mysql和mongo容器

创建mongo

```sh
docker run --hostname=2bda484650a9 --env=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin --env=HOME=/data/db --volume=/data/configdb --volume=/data/db -p 27017:27017 --restart=no --label='org.opencontainers.image.ref.name=ubuntu' --label='org.opencontainers.image.version=22.04' --runtime=runc -d mongo

```

创建mysql

```sh
docker run --hostname=ac818ccf9f17 --env=MYSQL_ROOT_PASSWORD=root --env=PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin --env=GOSU_VERSION=1.16 --env=MYSQL_MAJOR=innovation --volume=/var/lib/mysql -p 3306:3306 -p 30060:33060 --restart=no --runtime=runc -d mysql:latest

```

## 使用脚本初始环境

运行db_init.sh