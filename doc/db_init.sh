#!/bin/bash

DB_NAME="mmo_game"

# MongoDB 部分
echo "Creating MongoDB database '$DB_NAME'..."
sudo docker exec -i mongo mongo <<EOF
use $DB_NAME
db.createCollection("$DB_NAME")
EOF
echo "MongoDB database '$DB_NAME' created successfully."

# MySQL 部分
MYSQL_USER="root"
MYSQL_PASSWORD="root"
MYSQL_HOST="127.0.0.1"
MYSQL_PORT="3306"

echo "Creating MySQL database '$DB_NAME'..."
sudo docker exec -i mysql mysql -u$MYSQL_USER -p$MYSQL_PASSWORD -h$MYSQL_HOST -P$MYSQL_PORT <<EOF
CREATE DATABASE $DB_NAME;
EOF
echo "MySQL database '$DB_NAME' created successfully."
