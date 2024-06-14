package bootstrap

import (
	"github.com/NumberMan1/MMO-server/config"
	"github.com/NumberMan1/MMO-server/database"
	"github.com/NumberMan1/common/global"
	"github.com/NumberMan1/common/summer/proto_helper"
)

func Init(configPath string) {
	//解析yaml
	global.Init(configPath)
	database.Init(configPath)
	config.ServerInit(configPath)
	//初始化pb字典
	proto_helper.Init()
}
