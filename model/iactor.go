package model

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/define"
	"github.com/NumberMan1/MMO-server/model/entity"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
)

type IActor interface {
	entity.IEntity
	UnitState() proto.UnitState
	SetUnitState(unitState proto.UnitState)
	SkillMgr() *SkillManager
	SetSkillMgr(skillMgr *SkillManager)
	Spell() *Spell
	SetSpell(spell *Spell)
	Hp() float32
	Mp() float32
	SetHp(hp float32)
	SetMp(mp float32)
	Id() int
	Name() string
	Type() proto.EntityType
	SetId(v int)
	SetName(v string)
	SetType(v proto.EntityType)
	State() proto.EntityState
	SetState(state proto.EntityState)
	Space() *Space
	SetSpace(space *Space)
	Info() *proto.NetActor
	SetInfo(info *proto.NetActor)
	Define() *define.UnitDefine
	SetDefine(define *define.UnitDefine)
	Attr() *AttributesAssembly
	SetAttr(attr *AttributesAssembly)
	IsDeath() bool
	OnEnterSpace(space *Space, chr IActor)
	Revive()
	TelportSpace(space *Space, pos, dir *vector3.Vector3, chr IActor)
	Update()
	Die(killerID int)
	OnBeforeDie(killerID int)
	OnAfterDie(killerID int)
	RecvDamage(dmg *proto.Damage)
}
