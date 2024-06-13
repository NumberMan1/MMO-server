package database

import (
	"context"
	"fmt"
	"github.com/NumberMan1/common"
	"github.com/NumberMan1/common/logger"
	mongobrocker "github.com/NumberMan1/common/mongo"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
)

// mongodb可配置的数据库
var (
	MongoDbClient *mongobrocker.Client
)

type sysConfig struct {
	Mongodb struct {
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		MinPoolSize int    `yaml:"min_pool_size"`
		MaxPoolSize int    `yaml:"max_pool_size"`
	} `yaml:"mongodb"`
}

func Init(configPath string) {
	ctx := context.Background()
	//读取yaml
	file, _ := os.Open(configPath)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.SLoggerConsole.Error("Error closing")
		}
	}(file)
	decoder := yaml.NewDecoder(file)
	config := sysConfig{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error:", err)
		panic("加载配置出错")
	}
	url := "mongodb://" + config.Mongodb.Host + ":" + strconv.FormatInt(int64(config.Mongodb.Port), 10)
	fmt.Println(url)
	//创建mongo客户端
	MongoDbClient = &mongobrocker.Client{
		BaseComponent: common.NewBaseComponent(),
		RealCli: mongobrocker.NewClient(ctx, &mongobrocker.Config{
			URI:         url,
			MinPoolSize: 3,
			MaxPoolSize: 3000,
		}),
	}
}
