package database

import (
	"context"
	"fmt"
	"github.com/NumberMan1/common"
	"github.com/NumberMan1/common/logger"
	mongobrocker "github.com/NumberMan1/common/mongo"
	"github.com/NumberMan1/common/ormdb"
	"gorm.io/gorm"
)

// 可配置的数据库参数
var (
	Host             string
	Port             string
	User             string
	Password         string
	OrmConnectionStr string
	OrmDb            *gorm.DB
	MongoDbClient    *mongobrocker.Client
)

func init() {
	Host = "127.0.0.1"
	Port = "3406"
	User = "root"
	Password = "root"
	OrmConnectionStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/game?charset=utf8&parseTime=True&loc=Local&timeout=1000ms", User, Password, Host, Port)
	var err error
	OrmDb, err = ormdb.ConnectToDB("mysql", OrmConnectionStr)
	if err != nil {
		logger.SLCError(err.Error())
	}
	err = OrmDb.AutoMigrate(&DbPlayer{})
	if err != nil {
		logger.SLCError(err.Error())
	}
	err = OrmDb.AutoMigrate(&DbCharacter{})
	if err != nil {
		logger.SLCError(err.Error())
	}
	ctx := context.Background()
	MongoDbClient = &mongobrocker.Client{
		BaseComponent: common.NewBaseComponent(),
		RealCli: mongobrocker.NewClient(ctx, &mongobrocker.Config{
			URI:         "mongodb://localhost:26017",
			MinPoolSize: 3,
			MaxPoolSize: 3000,
		}),
	}
}
