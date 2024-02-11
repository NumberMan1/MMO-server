package main

import (
	"github.com/NumberMan1/MMO-server/mgr"
	"github.com/NumberMan1/MMO-server/service"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer"
)

func initServices() {
	//加载JSON配置文件
	mgr.GetDataManagerInstance().Init()
	netService := service.NewNetService()
	netService.Start()
	logger.SLCDebug("网络服务启动完成")
	service.GetSpaceServiceInstance().Start()
	logger.SLCDebug("地图服务启动完成")
	service.GetUserServiceInstance().Start()
	logger.SLCDebug("玩家服务启动完成")
	summer.GetScheduleInstance().Start()
	logger.SLCDebug("定时器服务启动完成")
}

func main() {
	initServices()
	select {}
}
