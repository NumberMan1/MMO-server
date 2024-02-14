package main

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/define"
	"github.com/NumberMan1/MMO-server/model"
	"github.com/NumberMan1/MMO-server/service"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer"
	"github.com/NumberMan1/common/summer/timeunit"
)

func initServices() {
	//加载JSON配置文件
	define.GetDataManagerInstance().Init()
	//网路服务模块
	netService := service.NewNetService()
	netService.Start()
	logger.SLCDebug("网络服务启动完成")
	service.GetSpaceServiceInstance().Start()
	logger.SLCDebug("地图服务启动完成")
	service.GetUserServiceInstance().Start()
	logger.SLCDebug("玩家服务启动完成")
	summer.GetScheduleInstance().Start()
	logger.SLCDebug("定时器服务启动完成")

	space := model.GetSpaceManagerInstance().GetSpace(2)
	mon := space.MonsterManager.Create(1002, 3, vector3.NewVector3(270038, 0, 322005), vector3.Zero3())
	mon.AI = model.NewMonsterAI(mon)

	summer.GetScheduleInstance().AddTask(func() {
		model.GetEntityManagerInstance().Update()
		model.GetSpaceManagerInstance().Update()
	}, timeunit.Milliseconds, 20, 0)
}

func main() {
	initServices()
	select {}
}
