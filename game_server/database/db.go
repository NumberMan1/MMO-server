package database

import (
	"context"
	"fmt"
	"github.com/NumberMan1/common"
	"github.com/NumberMan1/common/global/variable"
	mongobrocker "github.com/NumberMan1/common/mongo"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

type DBModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// mongodb可配置的数据库
var (
	MongoDbClient *mongobrocker.Client
	DatabaseName  string
)

type sysConfig struct {
	Mongodb struct {
		Host        string `yaml:"host" mapstructure:"host"`
		Port        int    `yaml:"port" mapstructure:"port"`
		User        string `yaml:"user" mapstructure:"user"`
		Password    string `yaml:"password" mapstructure:"password"`
		Database    string `yaml:"database" mapstructure:"database"`
		MinPoolSize int    `yaml:"min_pool_size" mapstructure:"min_pool_size"`
		MaxPoolSize int    `yaml:"max_pool_size" mapstructure:"max_pool_size"`
	} `yaml:"mongodb" mapstructure:"mongodb"`
}

func Init() {
	ctx := context.Background()
	//读取yaml
	config := sysConfig{}
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("Error:", err)
		panic("加载配置出错")
	}
	DatabaseName = config.Mongodb.Database
	url := "mongodb://" + config.Mongodb.User + ":" + config.Mongodb.Password + "@" +
		config.Mongodb.Host + ":" + strconv.FormatInt(int64(config.Mongodb.Port), 10)
	//创建mongo客户端
	MongoDbClient = &mongobrocker.Client{
		BaseComponent: common.NewBaseComponent(),
		RealCli: mongobrocker.NewClient(ctx, &mongobrocker.Config{
			URI:         url,
			MinPoolSize: 3,
			MaxPoolSize: 3000,
		}),
	}
	//设置mysql表
	err := variable.GDb.AutoMigrate(DbPlayer{})
	if err != nil {
		panic(err)
	}
	err = variable.GDb.AutoMigrate(DbCharacter{})
	if err != nil {
		panic(err)
	}
}
