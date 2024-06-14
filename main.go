package main

import (
	"github.com/NumberMan1/MMO-server/bootstrap"
	"github.com/NumberMan1/MMO-server/config/define"
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/inventory/item"
	"github.com/NumberMan1/MMO-server/mgr"
	"github.com/NumberMan1/MMO-server/model"
	"github.com/NumberMan1/MMO-server/service"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer"
	"github.com/NumberMan1/common/summer/timeunit"
)

func initServices() {
	bootstrap.Init("config/config.yaml")
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
	service.GetBattleServiceInstance().Start()
	logger.SLCDebug("战斗服务启动完成")
	service.GetChatServiceInstance().Start()
	logger.SLCDebug("聊天服务启动完成")
	summer.GetScheduleInstance().Start()
	logger.SLCDebug("定时器服务启动完成")

	//space := model.GetSpaceManagerInstance().GetSpace(2)
	//mon := space.MonsterManager.Create(1002, 3, vector3.NewVector3(270038, 0, 322005), vector3.Zero3())
	//mon.AI = model.NewMonsterAI(mon)

	summer.GetScheduleInstance().AddTask(func() {
		mgr.GetEntityManagerInstance().Update()
		model.GetSpaceManagerInstance().Update()
	}, timeunit.Milliseconds, 20, 0)

	model.CreateItemEntity(model.GetSpaceManagerInstance().GetSpace(1), item.NewItemByItemId(1001, 10, 0),
		vector3.NewVector3(0, 0, 0), vector3.Zero3())
	model.CreateItemEntity(model.GetSpaceManagerInstance().GetSpace(1), item.NewItemByItemId(1002, 5, 0),
		vector3.NewVector3(3000, 0, 3000), vector3.Zero3())
	model.CreateItemEntity(model.GetSpaceManagerInstance().GetSpace(1), item.NewItemByItemId(1003, 1, 0),
		vector3.NewVector3(6000, 0, 6000), vector3.Zero3())

	//传送门1：新手村=>森林
	gate1 := model.NewGate(1, 4001001, vector3.NewVector3(10000, 0, 10000), vector3.Zero3())
	gate1.SetName("传送门-森林入口")
	gate1.SetTarget(model.GetSpaceManagerInstance().GetSpace(2), vector3.NewVector3(354947, 1660, 308498))
	//传送门2：森林=>新手村
	gate3 := model.NewGate(2, 4001001, vector3.NewVector3(346318, 1870, 319313), vector3.Zero3())
	gate3.SetName("传送门-新手村")
	gate3.SetTarget(model.GetSpaceManagerInstance().GetSpace(1), vector3.NewVector3(0, 0, 0))
	//山贼附近
	gate2 := model.NewGate(1, 4001001, vector3.NewVector3(15000, 0, 10000), vector3.Zero3())
	gate2.SetName("传送门-山贼")
	gate2.SetTarget(model.GetSpaceManagerInstance().GetSpace(2), vector3.NewVector3(263442, 5457, 306462))
}

func main() {
	initServices()
	select {}
}
