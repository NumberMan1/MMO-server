package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type ServerSysConfig struct {
	Server struct {
		Host        string `yaml:"host" mapstructure:"host"`
		Port        int    `yaml:"port" mapstructure:"port"`
		WorkerCount int    `yaml:"worker_count" mapstructure:"worker_count"`
	} `yaml:"server" mapstructure:"server"`
}

var (
	ServerConfig ServerSysConfig
)

func ServerInit() {
	//读取yaml
	config := ServerSysConfig{}
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("Error:", err)
		panic("加载配置出错")
	}

	ServerConfig = config
}
