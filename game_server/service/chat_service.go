package service

import (
	"fmt"
	"github.com/NumberMan1/MMO-server/game_common/protocol/gen/proto"
	"github.com/NumberMan1/MMO-server/game_server/config/define"
	"github.com/NumberMan1/MMO-server/game_server/core/vector3"
	"github.com/NumberMan1/MMO-server/game_server/model"
	item2 "github.com/NumberMan1/MMO-server/game_server/model/inventory_system/item"
	"github.com/NumberMan1/common/ns/singleton"
	"github.com/NumberMan1/common/summer/network"
	"github.com/NumberMan1/common/summer/network/message_router"
)

var (
	singleChatService = singleton.Singleton{}
)

type ChatService struct {
}

func GetChatServiceInstance() *ChatService {
	result, _ := singleton.GetOrDo[*ChatService](&singleChatService, func() (*ChatService, error) {
		return &ChatService{}, nil
	})
	return result
}

func (cs *ChatService) Start() {
	network.GetMessageRouterInstance().Subscribe("proto.ChatRequest", message_router.MessageHandler{Op: cs.chatRequest})
}

func (cs *ChatService) chatRequest(msg message_router.Msg) {
	//获取当前主角对象
	session := msg.Sender.(network.Connection).Get("Session").(*model.Session)
	message := msg.Message.(*proto.ChatRequest)
	chr := session.Character
	//广播聊天消息
	resp := &proto.ChatResponse{
		SenderId:  int32(chr.EntityId()),
		TextValue: message.TextValue,
	}
	chr.Space().Broadcast(resp)

	//var sd *define2.SpaceDefine
	//for _, v := range define2.GetDataManagerInstance().Spaces {
	//	if v.Name == message.TextValue {
	//		sd = v
	//		break
	//	}
	//}
	//var sp *model.Space
	//if sd != nil {
	//	sp = model.GetSpaceManagerInstance().GetSpace(sd.GetId())
	//	switch message.TextValue {
	//	case "新手村":
	//		model.TeleportSpace(sp, vector3.Zero3(), vector3.Zero3(), chr)
	//	}
	//} else {
	//	switch message.TextValue {
	//	case "森林":
	//		model.TeleportSpace(model.GetSpaceManagerInstance().GetSpace(2), vector3.NewVector3(354947, 1660, 308498), vector3.Zero3(), chr)
	//	case "山贼":
	//		model.TeleportSpace(model.GetSpaceManagerInstance().GetSpace(2), vector3.NewVector3(263442, 5457, 306462), vector3.Zero3(), chr)
	//	}
	//}
	instance := model.GetSpaceManagerInstance()
	switch message.TextValue {
	case "新手村":
		sp := instance.GetSpace(1)
		model.TeleportSpace(sp, vector3.NewVector3(0, 0, 0), vector3.Zero3(), chr)
	case "森林":
		sp := instance.GetSpace(2)
		model.TeleportSpace(sp, vector3.NewVector3(354947, 1660, 308498), vector3.Zero3(), chr)
	case "山贼":
		sp := instance.GetSpace(2)
		model.TeleportSpace(sp, vector3.NewVector3(263442, 5457, 306462), vector3.Zero3(), chr)
	case "-bull":
		sp := instance.GetSpace(2)
		model.TeleportSpace(sp, vector3.NewVector3(225597, 29220, 393920), vector3.Zero3(), chr)
	case "-attach":
		//创建鞋子
		def := define.GetDataManagerInstance().Items[1005]
		item := item2.NewEquipmentByDefine(def, 0)
		//穿在身上
		chr.EquipsManager.Attach(item)
		//刷新背包
		chr.SendInventory(true, false, false)
	case "-detach":
		chr.EquipsManager.Detach(proto.EquipsType_Boots)
		chr.SendInventory(true, false, false)
	case "-sudu3":
		chr.SyncSpeed(3000)
	case "-sudu8":
		chr.SyncSpeed(8000)
	default:
		fmt.Println("Unrecognized command")
	}
}
