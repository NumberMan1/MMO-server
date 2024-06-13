package bootstrap

import (
	"github.com/NumberMan1/MMO-server/config"
	"github.com/NumberMan1/MMO-server/database"
	"github.com/NumberMan1/common/global"
)

func Init(configPath string) {
	global.Init(configPath)
	database.Init(configPath)
	config.ServerInit(configPath)
}
