package model

import (
	"container/list"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ns"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
)

// FightMgr 战斗管理器
type FightMgr struct {
	space *Space
	//投射物列表
	Missiles *list.List

	//技能施法队列
	CastQueue *ns.TSQueue[*proto.CastInfo]
	//等待广播的队列
	SpellQueue *ns.TSQueue[*proto.CastInfo]
	//等待广播的伤害队列
	DamageQueue *ns.TSQueue[*proto.Damage]
	//角色属性变动
	PropertyUpdateQueue *ns.TSQueue[*proto.PropertyUpdate]

	//施法响应对象，每帧发送一次
	spellResponse *proto.SpellResponse
	//伤害消息，每帧发送一次
	damageResponse *proto.DamageResponse
	//角色属性变化
	propertyUpdateResponse *proto.PropertyUpdateResponse
}

func NewFightMgr(space *Space) *FightMgr {
	return &FightMgr{
		space:                  space,
		Missiles:               list.New(),
		CastQueue:              ns.NewTSQueue[*proto.CastInfo](),
		SpellQueue:             ns.NewTSQueue[*proto.CastInfo](),
		DamageQueue:            ns.NewTSQueue[*proto.Damage](),
		PropertyUpdateQueue:    ns.NewTSQueue[*proto.PropertyUpdate](),
		spellResponse:          &proto.SpellResponse{CastList: make([]*proto.CastInfo, 0)},
		damageResponse:         &proto.DamageResponse{List: make([]*proto.Damage, 0)},
		propertyUpdateResponse: &proto.PropertyUpdateResponse{List: make([]*proto.PropertyUpdate, 0)},
	}
}

func (fm *FightMgr) OnUpdate(delta float64) {
	for !fm.CastQueue.Empty() {
		cast := fm.CastQueue.Pop()
		logger.SLCInfo("执行施法:%v", cast)
		fm.RunCast(cast)
	}
	for e := fm.Missiles.Front(); e != nil; e = e.Next() {
		e.Value.(*Missile).OnUpdate(delta)
	}
	fm.broadcastSpell()
	fm.broadcastDamage()
	fm.broadcastProperties()
}

func (fm *FightMgr) broadcastProperties() {
	for !fm.PropertyUpdateQueue.Empty() {
		item := fm.PropertyUpdateQueue.Pop()
		fm.propertyUpdateResponse.List = append(fm.propertyUpdateResponse.List, item)
	}
	if len(fm.propertyUpdateResponse.List) > 0 {
		fm.space.Broadcast(fm.propertyUpdateResponse)
		fm.propertyUpdateResponse.List = make([]*proto.PropertyUpdate, 0)
	}
}

func (fm *FightMgr) broadcastDamage() {
	for !fm.DamageQueue.Empty() {
		item := fm.DamageQueue.Pop()
		fm.damageResponse.List = append(fm.damageResponse.List, item)
	}
	if len(fm.damageResponse.List) > 0 {
		fm.space.Broadcast(fm.damageResponse)
		fm.damageResponse.List = make([]*proto.Damage, 0)
	}
}

// 广播施法信息
func (fm *FightMgr) broadcastSpell() {
	for !fm.SpellQueue.Empty() {
		item := fm.SpellQueue.Pop()
		fm.spellResponse.CastList = append(fm.spellResponse.CastList, item)
	}
	if len(fm.spellResponse.CastList) > 0 {
		fm.space.Broadcast(fm.spellResponse)
		fm.spellResponse.CastList = make([]*proto.CastInfo, 0)
	}
}

func (fm *FightMgr) RunCast(info *proto.CastInfo) {
	caster := GetEntityManagerInstance().GetEntity(int(info.CasterId)).(IActor)
	if caster == nil {
		logger.SLCInfo("RunCast: Caster is null %v", info.CasterId)
		return
	}
	caster.Spell().RunCast(info)
}
