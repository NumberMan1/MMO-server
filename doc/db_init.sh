#!/bin/bash

# 定义变量
MONGO_CONTAINER="mongo"
MYSQL_CONTAINER="mysql"
MONGO_PORT=27017
MONGO_USER="root"
MONGO_PASSWORD="root"
MYSQL_PORT=3306
MYSQL_USER="root"
MYSQL_PASSWORD="root"
DB_NAME="mmo_game"

# 创建 MongoDB 架构
echo "Creating $DB_NAME database in MongoDB..."
docker run --name=mongo --env=MONGO_INITDB_ROOT_USERNAME=$MONGO_USER --env=MONGO_INITDB_ROOT_PASSWORD=$MONGO_PASSWORD --env=MONGO_INITDB_DATABASE=$DB_NAME -p 27017:27017 -d mongo:latest

# 创建 MySQL 架构
echo "Creating $DB_NAME database in MySQL..."
docker run --name=mysql --env=MYSQL_ROOT_PASSWORD=$MYSQL_PASSWORD --env=MYSQL_DATABASE=$DB_NAME -p 3306:3306 -p 30060:33060 -d mysql:latest

echo "Databases created successfully."
