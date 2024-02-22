package model

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/define"
	"github.com/NumberMan1/MMO-server/model/entity"
	"github.com/NumberMan1/common/logger"
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

func (a *Actor) SetHp(hp float32) {
	a.Info().Hp = hp
}

func (a *Actor) SetMp(mp float32) {
	a.Info().Mp = mp
}

func NewActor(t proto.EntityType, tid, level int, position, direction *vector3.Vector3) *Actor {
	a := &Actor{
		Entity: entity.NewEntity(position, direction),
		define: define.GetDataManagerInstance().Units[tid],
		info: &proto.NetActor{
			Tid:   int32(tid),
			Level: int32(level),
			Type:  t,
		},
		attr: &AttributesAssembly{},
	}
	a.info.Name = a.define.Name
	a.info.Entity = a.EntityData()
	a.SetSkillMgr(NewSkillManager(a))
	a.SetHp(a.define.HPMax)
	a.SetMp(a.define.MPMax)
	a.Attr().Init(a)
	a.SetSpell(NewSpell(a))
	a.SetSpeed(a.define.Speed)
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
	a.setHp(a.Attr().Final.HPMax)
	a.setMp(a.Attr().Final.MPMax)
	a.setState(proto.UnitState_FREE)
}

func (a *Actor) TelportSpace(space *Space, pos, dir *vector3.Vector3, chr IActor) {
	if _, ok := chr.(*Character); !ok {
		return
	}
	chrTmp := chr.(*Character)
	if space != a.Space() {
		//1.退出当前场景
		space.CharacterLeave(chrTmp)
		//2.设置坐标和方向
		chrTmp.SetPosition(pos)
		chrTmp.SetDirection(dir)
		//3.进入新场景
		space.CharacterJoin(chrTmp)
	} else {
		space.Telport(chrTmp, pos, dir)
	}
}

func (a *Actor) Update() {
	a.SkillMgr().Update()
}

func (a *Actor) Die(killerID int) {
	if a.IsDeath() {
		return
	}
	a.OnBeforeDie(killerID)
	a.setHp(0)
	a.setMp(0)
	a.setState(proto.UnitState_DEAD)
	a.OnAfterDie(killerID)
}

func (a *Actor) OnBeforeDie(killerID int) {

}

func (a *Actor) OnAfterDie(killerID int) {

}

func (a *Actor) RecvDamage(dmg *proto.Damage) {
	logger.SLCInfo("Actor:RecvDamage[%v]", dmg)
	//添加广播
	a.Space().FightMgr.DamageQueue.Push(dmg)
	//扣血或者死亡
	if a.Hp() > dmg.Amount {
		a.setHp(a.Hp() - dmg.Amount)
	} else {
		a.Die(int(dmg.AttackerId))
	}
}

func (a *Actor) setHp(hp float32) {
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

func (a *Actor) setMp(mp float32) {
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

func (a *Actor) setState(unitState proto.UnitState) {
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
