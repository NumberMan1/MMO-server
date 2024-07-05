# MMO游戏服务器(基于github.com/NumberMan1/common实现) 持续更新

## 需要的第三方组件

需要mongodb,mysql,docker

## 游戏通用组件(game_common)

其它游戏服务所需的功能如协议,rpc等

## 游戏服务(game_server)

具体配置在config/config.yaml里
先执行doc/db_init.sh初始化数据库,然后执行/config/init/data_init.go初始化mongodb数据

运行后main.go服务后会创建storage目录,里面是服务的具体日志

视频展示https://www.bilibili.com/video/BV1QW421N7J2/?spm_id_from=333.999.0.0

目前忙于工作,有时间就会更新项目