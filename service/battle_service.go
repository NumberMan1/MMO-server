package service

import (
	"github.com/NumberMan1/MMO-server/model"
	"github.com/NumberMan1/MMO-server/protocol/gen/proto"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ns/singleton"
	"github.com/NumberMan1/common/summer/network"
	"github.com/NumberMan1/common/summer/network/message_router"
)

var (
	singleBattleService = singleton.Singleton{}
)

type BattleService struct {
}

func GetBattleServiceInstance() *BattleService {
	result, _ := singleton.GetOrDo[*BattleService](&singleBattleService, func() (*BattleService, error) {
		return &BattleService{}, nil
	})
	return result
}

func (bs *BattleService) Start() {
	network.GetMessageRouterInstance().Subscribe("proto.SpellRequest", message_router.MessageHandler{Op: bs.spellRequest})
}

func (bs *BattleService) spellRequest(msg message_router.Msg) {
	req := msg.Message.(*proto.SpellRequest)
	logger.SLCInfo("技能施法请求：%v", req)
	session := msg.Sender.(network.Connection).Get("Session").(*model.Session)
	chr := session.Character
	if chr.EntityId() != int(req.GetInfo().GetCasterId()) {
		logger.SLCError("施法者ID错误")
		return
	}
	//加入战斗管理器
	chr.Space().FightMgr.CastQueue.Push(req.GetInfo())
}
