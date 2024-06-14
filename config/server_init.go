package config

import (
	"fmt"
	"github.com/NumberMan1/common/logger"
	"gopkg.in/yaml.v3"
	"os"
)

type ServerSysConfig struct {
	Server struct {
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		WorkerCount int    `yaml:"worker_count"`
	} `yaml:"server"`
}

var (
	ServerConfig ServerSysConfig
)

func ServerInit(configPath string) {
	//读取yaml
	file, _ := os.Open(configPath)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.SLoggerConsole.Error("Error closing")
		}
	}(file)
	decoder := yaml.NewDecoder(file)
	config := ServerSysConfig{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error:", err)
		panic("加载配置出错")
	}

	ServerConfig = config
}
