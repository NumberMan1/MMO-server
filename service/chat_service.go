package service

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/define"
	"github.com/NumberMan1/MMO-server/model"
	"github.com/NumberMan1/common/ns/singleton"
	"github.com/NumberMan1/common/summer/network"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
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
	network.GetMessageRouterInstance().Subscribe("proto.ChatRequest", network.MessageHandler{Op: cs.chatRequest})
}

func (cs *ChatService) chatRequest(msg network.Msg) {
	//获取当前主角对象
	session := msg.Sender.Get("Session").(*model.Session)
	message := msg.Message.(*proto.ChatRequest)
	chr := session.Character
	//广播聊天消息
	resp := &proto.ChatResponse{
		SenderId:  int32(chr.EntityId()),
		TextValue: message.TextValue,
	}
	chr.Space().Broadcast(resp)

	var sd *define.SpaceDefine
	for _, v := range define.GetDataManagerInstance().Spaces {
		if v.Name == message.TextValue {
			sd = v
			break
		}
	}
	var sp *model.Space
	if sd != nil {
		sp = model.GetSpaceManagerInstance().GetSpace(sd.GetId())
		switch message.TextValue {
		case "新手村":
			chr.TelportSpace(sp, vector3.Zero3(), vector3.Zero3(), chr)
		}
	} else {
		switch message.TextValue {
		case "森林":
			chr.TelportSpace(model.GetSpaceManagerInstance().GetSpace(2), vector3.NewVector3(354947, 1660, 308498), vector3.Zero3(), chr)
		case "山贼":
			chr.TelportSpace(model.GetSpaceManagerInstance().GetSpace(2), vector3.NewVector3(263442, 5457, 306462), vector3.Zero3(), chr)
		}
	}
}
