package model

import (
	"github.com/NumberMan1/MMO-server/battle"
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/define"
	"github.com/NumberMan1/MMO-server/model/core"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
)

type Actor struct {
	*core.Entity
	space   *Space
	info    *proto.NCharacter
	define  define.UnitDefine
	state   proto.EntityState
	attr    *battle.Attributes
	isDeath bool
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

func (a *Actor) Info() *proto.NCharacter {
	return a.info
}

func (a *Actor) SetInfo(info *proto.NCharacter) {
	a.info = info
}

func (a *Actor) Define() define.UnitDefine {
	return a.define
}

func (a *Actor) SetDefine(define define.UnitDefine) {
	a.define = define
}

func (a *Actor) Attr() *battle.Attributes {
	return a.attr
}

func (a *Actor) SetAttr(attr *battle.Attributes) {
	a.attr = attr
}

func (a *Actor) IsDeath() bool {
	return a.isDeath
}

func (a *Actor) SetIsDeath(isDeath bool) {
	a.isDeath = isDeath
}

func NewActor(t proto.EntityType, tid, level int, position, direction vector3.Vector3) *Actor {
	a := &Actor{
		Entity: core.NewEntity(position, direction),
		define: define.GetDataManagerInstance().Units[tid],
		info: &proto.NCharacter{
			Tid:   int32(tid),
			Level: int32(level),
			Type:  t,
		},
		attr: &battle.Attributes{},
	}
	a.info.Name = a.define.Name
	a.info.Entity = a.EntityData()
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

func (a *Actor) OnEnterSpace(space *Space) {
	a.space = space
	a.info.SpaceId = int32(space.Id)
}

func (a *Actor) Revive() {
	a.isDeath = false
}
