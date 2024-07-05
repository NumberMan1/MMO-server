# 初始化数据库环境

## docker创建基本的mysql和mongo容器

创建mongo

```sh
docker run --name=mongo --env=MONGO_INITDB_ROOT_USERNAME=$MONGO_USER --env=MONGO_INITDB_ROOT_PASSWORD=$MONGO_PASSWORD --env=MONGO_INITDB_DATABASE=$DB_NAME -p 27017:27017 -d mongo:latest
```

创建mysql

```sh
docker run --name=mysql --env=MYSQL_ROOT_PASSWORD=$MYSQL_PASSWORD --env=MYSQL_DATABASE=$DB_NAME -p 3306:3306 -p 30060:33060 -d mysql:latest
```

## 使用脚本初始环境

运行db_init.sh