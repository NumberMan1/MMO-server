编译后用于初始化mongo数据的命令如下

```shell
./程序名 init --config 配置文件路径
```

```shell
#示例命令
./build_init init --config game_server/config/config.yaml
```

查看是否成功,打印mongo的数据

```shell
./程序名
```

```shell
#示例命令
./build_init --config game_server/config/config.yaml
```



