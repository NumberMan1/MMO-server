package model

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/define"
	"github.com/NumberMan1/MMO-server/inventory/item"
	"github.com/NumberMan1/MMO-server/model/entity"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ns"
	"github.com/NumberMan1/common/summer/core"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
)

type Actor struct {
	*entity.Entity
	space     *Space
	info      *proto.NetActor
	define    *define.UnitDefine
	state     proto.EntityState
	attr      *AttributesAssembly
	unitState proto.UnitState
	skillMgr  *SkillManager
	spell     *Spell
}

func (a *Actor) Level() int {
	return int(a.Info().Exp)
}

func (a *Actor) Exp() int {
	return int(a.Info().Exp)
}

func (a *Actor) Gold() int {
	return int(a.Info().Gold)
}

func (a *Actor) HPMax() int {
	return int(a.Attr().Final.HPMax)
}

func (a *Actor) MPMax() int {
	return int(a.Attr().Final.MPMax)
}

func (a *Actor) UnitState() proto.UnitState {
	return a.unitState
}

func (a *Actor) SetUnitState(unitState proto.UnitState) {
	a.unitState = unitState
}

func (a *Actor) SkillMgr() *SkillManager {
	return a.skillMgr
}

func (a *Actor) SetSkillMgr(skillMgr *SkillManager) {
	a.skillMgr = skillMgr
}

func (a *Actor) Spell() *Spell {
	return a.spell
}

func (a *Actor) SetSpell(spell *Spell) {
	a.spell = spell
}

func (a *Actor) Hp() float32 {
	return a.Info().Hp
}

func (a *Actor) Mp() float32 {
	return a.Info().Mp
}

func NewActor(t proto.EntityType, tid, level int, position, direction *vector3.Vector3) *Actor {
	a := &Actor{
		Entity: entity.NewEntity(position, direction),
		info: &proto.NetActor{
			Tid:   int32(tid),
			Level: int32(level),
			Type:  t,
		},
		attr: NewAttributesAssembly(),
	}
	a.Info().Entity = a.EntityData()
	if def, ok := define.GetDataManagerInstance().Units[tid]; ok {
		a.define = def
		a.Info().Name = a.define.Name
		a.Info().Hp = a.define.HPMax
		a.Info().Mp = a.define.MPMax
	}
	if a.Type() != proto.EntityType_Item {
		a.SetSkillMgr(NewSkillManager(a))
		a.Attr().Init(a)
		a.SetSpell(NewSpell(a))
	}

	//a.SetSpeed(a.define.Speed)
	return a
}

func (a *Actor) Id() int {
	return int(a.info.Id)
}

func (a *Actor) Name() string {
	return a.info.Name
}

func (a *Actor) Type() proto.EntityType {
	return a.info.Type
}

func (a *Actor) SetId(v int) {
	a.info.Id = int32(v)
}

func (a *Actor) SetName(v string) {
	a.info.Name = v
}

func (a *Actor) SetType(v proto.EntityType) {
	a.info.Type = v
}

func (a *Actor) State() proto.EntityState {
	return a.state
}

func (a *Actor) SetState(state proto.EntityState) {
	a.state = state
}

func (a *Actor) Space() *Space {
	return a.space
}

func (a *Actor) SetSpace(space *Space) {
	a.space = space
}

func (a *Actor) Info() *proto.NetActor {
	return a.info
}

func (a *Actor) SetInfo(info *proto.NetActor) {
	a.info = info
}

func (a *Actor) Define() *define.UnitDefine {
	return a.define
}

func (a *Actor) SetDefine(define *define.UnitDefine) {
	a.define = define
}

func (a *Actor) Attr() *AttributesAssembly {
	return a.attr
}

func (a *Actor) SetAttr(attr *AttributesAssembly) {
	a.attr = attr
}

func (a *Actor) IsDeath() bool {
	return a.unitState == proto.UnitState_DEAD
}

func (a *Actor) OnEnterSpace(space *Space, chr IActor) {
	if a.space != nil && space != nil {
		GetEntityManagerInstance().ChangeSpace(a, a.Space().Id, space.Id)
	}
	a.space = space
	a.info.SpaceId = int32(space.Id)
	if c, ok := chr.(*Character); ok {
		c.Data.SpaceId = space.Id
	}
}

func (a *Actor) Revive() {
	logger.SLCInfo("Actor.Revive:%v", a.EntityId())
	if !a.IsDeath() {
		return
	}
	a.SetAndUpdateHp(a.Attr().Final.HPMax)
	a.SetAndUpdateMp(a.Attr().Final.MPMax)
	a.SetAndUpdateState(proto.UnitState_FREE)
}

func (a *Actor) TeleportSpace(space *Space, pos, dir *vector3.Vector3, chr IActor) {
	if _, ok := chr.(*Character); !ok {
		return
	}
	chrTmp := chr.(*Character)
	if space != a.Space() {
		//1.退出当前场景
		space.EntityLeave(chrTmp)
		//2.设置坐标和方向
		chrTmp.SetPosition(pos)
		chrTmp.SetDirection(dir)
		//3.进入新场景
		space.EntityEnter(chrTmp)
	} else {
		space.Teleport(chrTmp, pos, dir)
	}
}

func (a *Actor) Update() {
	if a.SkillMgr() != nil {
		a.SkillMgr().Update()
	}
}

func (a *Actor) Die(killerID int) {
	if a.IsDeath() {
		return
	}
	a.OnBeforeDie(killerID)
	a.SetAndUpdateHp(0)
	a.SetAndUpdateMp(0)
	a.SetAndUpdateState(proto.UnitState_DEAD)
	a.OnAfterDie(killerID)
}

func (a *Actor) OnBeforeDie(killerID int) {

}

func (a *Actor) OnAfterDie(killerID int) {
	// 物品池
	arr := []int{1001, 1002}
	// 生成一个随机索引
	randIndex := ns.RandInt(0, len(arr))
	// 获取随机索引对应的元素
	itemId := arr[randIndex]
	CreateItemEntity(GetSpaceManagerInstance().GetSpace(a.Space().Id), item.NewItemByItemId(itemId, 5, 0), a.Position(), vector3.Zero3())
	//如果击杀者是玩家，给与奖励
	killer := GetUnit(killerID)
	if killer != nil {
		if chr, ok := killer.(*Character); ok {
			chr.SetAndUpdateLevel(chr.Level() + 1)
			chr.SetAndUpdateGolds(int64(chr.Gold() + 50))
			chr.SetAndUpdateExp(int64(chr.Exp() + 32))
		}
	}
}

func (a *Actor) RecvDamage(dmg *proto.Damage) {
	logger.SLCInfo("Actor:RecvDamage[%v]", dmg)
	//添加广播
	a.Space().FightMgr.DamageQueue.Push(dmg)
	//扣血或者死亡
	if a.Hp() > dmg.Amount {
		a.SetAndUpdateHp(a.Hp() - dmg.Amount)
	} else {
		a.Die(int(dmg.AttackerId))
	}
}

func (a *Actor) SetAndUpdateHp(hp float32) {
	if core.Equal(float64(a.Info().Hp), float64(hp)) {
		return
	}
	if hp <= 0 {
		hp = 0
	}
	if hp > a.Attr().Final.HPMax {
		hp = a.Attr().Final.HPMax
	}
	oldValue := a.Info().Hp
	a.Info().Hp = hp
	po := &proto.PropertyUpdate{
		EntityId: int32(a.EntityId()),
		Property: proto.PropertyUpdate_HP,
		OldValue: &proto.PropertyUpdate_PropertyValue{
			Value: &proto.PropertyUpdate_PropertyValue_FloatValue{
				FloatValue: oldValue,
			},
		},
		NewValue: &proto.PropertyUpdate_PropertyValue{
			Value: &proto.PropertyUpdate_PropertyValue_FloatValue{
				FloatValue: a.Info().Hp,
			},
		},
	}
	a.Space().FightMgr.PropertyUpdateQueue.Push(po)
}

func (a *Actor) SetAndUpdateMp(mp float32) {
	if core.Equal(float64(a.Info().Mp), float64(mp)) {
		return
	}
	if mp <= 0 {
		mp = 0
	}
	if mp > a.Attr().Final.MPMax {
		mp = a.Attr().Final.MPMax
	}
	oldValue := a.Info().Mp
	a.Info().Mp = mp
	po := &proto.PropertyUpdate{
		EntityId: int32(a.EntityId()),
		Property: proto.PropertyUpdate_MP,
		OldValue: &proto.PropertyUpdate_PropertyValue{
			Value: &proto.PropertyUpdate_PropertyValue_FloatValue{
				FloatValue: oldValue,
			},
		},
		NewValue: &proto.PropertyUpdate_PropertyValue{
			Value: &proto.PropertyUpdate_PropertyValue_FloatValue{
				FloatValue: a.Info().Mp,
			},
		},
	}
	a.Space().FightMgr.PropertyUpdateQueue.Push(po)
}

func (a *Actor) SetAndUpdateState(unitState proto.UnitState) {
	if a.UnitState() == unitState {
		return
	}
	oldValue := a.UnitState()
	a.SetUnitState(unitState)
	po := &proto.PropertyUpdate{
		EntityId: int32(a.EntityId()),
		Property: proto.PropertyUpdate_State,
		OldValue: &proto.PropertyUpdate_PropertyValue{
			Value: &proto.PropertyUpdate_PropertyValue_StateValue{
				StateValue: oldValue,
			},
		},
		NewValue: &proto.PropertyUpdate_PropertyValue{
			Value: &proto.PropertyUpdate_PropertyValue_StateValue{
				StateValue: a.UnitState(),
			},
		},
	}
	a.Space().FightMgr.PropertyUpdateQueue.Push(po)
}

// SetAndUpdateGolds 金币
func (a *Actor) SetAndUpdateGolds(value int64) {
	if a.Info().Gold == value {
		return
	}
	oldValue := a.Info().Gold
	a.Info().Gold = value
	rsp := &proto.PropertyUpdate{
		EntityId: int32(a.EntityId()),
		Property: proto.PropertyUpdate_Golds,
		OldValue: &proto.PropertyUpdate_PropertyValue{
			Value: &proto.PropertyUpdate_PropertyValue_LongValue{
				LongValue: oldValue,
			},
		},
		NewValue: &proto.PropertyUpdate_PropertyValue{
			Value: &proto.PropertyUpdate_PropertyValue_LongValue{
				LongValue: a.Info().Gold,
			},
		},
	}
	a.Space().FightMgr.PropertyUpdateQueue.Push(rsp)
}

// SetAndUpdateExp 经验
func (a *Actor) SetAndUpdateExp(value int64) {
	if a.Info().Exp == value {
		return
	}
	oldValue := a.Info().Exp
	a.Info().Exp = value
	rsp := &proto.PropertyUpdate{
		EntityId: int32(a.EntityId()),
		Property: proto.PropertyUpdate_Exp,
		OldValue: &proto.PropertyUpdate_PropertyValue{
			Value: &proto.PropertyUpdate_PropertyValue_LongValue{
				LongValue: oldValue,
			},
		},
		NewValue: &proto.PropertyUpdate_PropertyValue{
			Value: &proto.PropertyUpdate_PropertyValue_LongValue{
				LongValue: a.Info().Exp,
			},
		},
	}
	a.Space().FightMgr.PropertyUpdateQueue.Push(rsp)
}

// SetAndUpdateLevel 等级
func (a *Actor) SetAndUpdateLevel(value int) {
	if int(a.Info().Level) == value {
		return
	}
	oldValue := a.Info().Level
	a.Info().Level = int32(value)
	a.Attr().Reload()
	rsp := &proto.PropertyUpdate{
		EntityId: int32(a.EntityId()),
		Property: proto.PropertyUpdate_Level,
		OldValue: &proto.PropertyUpdate_PropertyValue{
			Value: &proto.PropertyUpdate_PropertyValue_IntValue{
				IntValue: oldValue,
			},
		},
		NewValue: &proto.PropertyUpdate_PropertyValue{
			Value: &proto.PropertyUpdate_PropertyValue_IntValue{
				IntValue: a.Info().Level,
			},
		},
	}
	a.Space().FightMgr.PropertyUpdateQueue.Push(rsp)
}

// OnHpMaxChanged 通知客户端：HPMax变化
func (a *Actor) OnHpMaxChanged(value float32) {
	a.Info().Hpmax = value
	po := &proto.PropertyUpdate{
		EntityId: int32(a.EntityId()),
		Property: proto.PropertyUpdate_HPMax,
		OldValue: &proto.PropertyUpdate_PropertyValue{
			Value: &proto.PropertyUpdate_PropertyValue_FloatValue{
				FloatValue: 0,
			},
		},
		NewValue: &proto.PropertyUpdate_PropertyValue{
			Value: &proto.PropertyUpdate_PropertyValue_FloatValue{
				FloatValue: value,
			},
		},
	}
	if a.Space() != nil {
		if a.Space().FightMgr != nil {
			a.Space().FightMgr.PropertyUpdateQueue.Push(po)
		}
	}
}

// OnMpMaxChanged 通知客户端：MPMax变化
func (a *Actor) OnMpMaxChanged(value float32) {
	a.Info().Mpmax = value
	po := &proto.PropertyUpdate{
		EntityId: int32(a.EntityId()),
		Property: proto.PropertyUpdate_MPMax,
		OldValue: &proto.PropertyUpdate_PropertyValue{
			Value: &proto.PropertyUpdate_PropertyValue_FloatValue{
				FloatValue: 0,
			},
		},
		NewValue: &proto.PropertyUpdate_PropertyValue{
			Value: &proto.PropertyUpdate_PropertyValue_FloatValue{
				FloatValue: value,
			},
		},
	}
	if a.Space() != nil {
		if a.Space().FightMgr != nil {
			a.Space().FightMgr.PropertyUpdateQueue.Push(po)
		}
	}
}
