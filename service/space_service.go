package service

import (
	"fmt"
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/mgr"
	"github.com/NumberMan1/MMO-server/model"
	"github.com/NumberMan1/MMO-server/protocol/gen/proto"
	"github.com/NumberMan1/common/ns/singleton"
	"github.com/NumberMan1/common/summer/network"
	"github.com/NumberMan1/common/summer/network/message_router"
	"math"
)

var (
	singleSpaceService = singleton.Singleton{}
)

// SpaceService 地图服务
type SpaceService struct {
}

func GetSpaceServiceInstance() *SpaceService {
	instance, _ := singleton.GetOrDo[*SpaceService](&singleSpaceService, func() (*SpaceService, error) {
		return &SpaceService{}, nil
	})
	return instance
}

func (ss *SpaceService) Start() {
	//初始化地图
	model.GetSpaceManagerInstance().Init()
	//位置同步请求
	network.GetMessageRouterInstance().Subscribe("proto.SpaceEntitySyncRequest", message_router.MessageHandler{Op: ss.spaceEntitySyncRequest})
}

func (ss *SpaceService) GetSpace(id int) *model.Space {
	return model.GetSpaceManagerInstance().GetSpace(id)
}

func (ss *SpaceService) spaceEntitySyncRequest(msg message_router.Msg) {
	//获取当前角色所在的地图
	sp := msg.Sender.(network.Connection).Get("Session").(*model.Session).Space()
	if sp == nil {
		return
	}
	//同步请求信息
	netEntity := msg.Message.(*proto.SpaceEntitySyncRequest).EntitySync.Entity
	netV3 := vector3.NewVector3(float64(netEntity.Position.X), float64(netEntity.Position.Y), float64(netEntity.Position.Z))
	//服务端实际的角色信息
	serEntity := mgr.GetEntityManagerInstance().GetEntity(int(netEntity.Id))
	serV3 := vector3.NewVector3(serEntity.Position().X, serEntity.Position().Y, serEntity.Position().Z)
	//计算距离
	distance := vector3.GetDistance(netV3, serV3)
	//使用服务器移动速度
	netEntity.Speed = int32(serEntity.Speed())
	//计算时间差
	dt := min(serEntity.PositionTime(), 1.0)
	//计算限额
	limit := float64(serEntity.Speed()) * dt * 3
	fmt.Printf("距离%v，阈值%v，间隔%v\n", distance, limit, dt)
	if math.IsNaN(distance) || distance > limit {
		//拉回原位置
		resp := &proto.SpaceEntitySyncResponse{
			EntitySync: &proto.NetEntitySync{
				Entity: serEntity.EntityData(),
				Force:  true,
			},
		}
		msg.Sender.(network.Connection).Send(resp)
		return
	}

	//广播同步信息
	sp.UpdateEntity(msg.Message.(*proto.SpaceEntitySyncRequest).EntitySync)
}
