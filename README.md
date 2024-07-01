# MMO游戏服务器(基于github.com/NumberMan1/common实现) 持续更新

## 第三方工具

mongodb,mysql,docker(可选)

## 配置文件

具体配置在config/config.yaml里

## 运行服务端

先执行doc/db_init.sh初始化数据库,然后执行/config/init/data_init.go初始化mongodb数据

运行后main.go服务后会创建storage目录,里面是服务的具体日志

视频展示https://www.bilibili.com/video/BV1QW421N7J2/?spm_id_from=333.999.0.0

## 客户端获取

[NumberMan1/MMO-client (github.com)](https://github.com/NumberMan1/MMO-client)对应分支的客户端

目前忙于工作,有时间就会更新项目