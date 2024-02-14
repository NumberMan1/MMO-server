package model

import (
	"github.com/NumberMan1/MMO-server/battle"
	"github.com/NumberMan1/MMO-server/define"
	"github.com/NumberMan1/MMO-server/model/core"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
)

type IActor interface {
	core.IEntity
	State() proto.EntityState
	SetState(state proto.EntityState)
	Space() *Space
	SetSpace(space *Space)
	Info() *proto.NCharacter
	SetInfo(info *proto.NCharacter)
	Define() define.UnitDefine
	SetDefine(define define.UnitDefine)
	Attr() *battle.Attributes
	SetAttr(attr *battle.Attributes)
	IsDeath() bool
	SetIsDeath(isDeath bool)
	Id() int
	Name() string
	Type() proto.EntityType
	SetId(v int)
	SetName(v string)
	SetType(v proto.EntityType)
	OnEnterSpace(space *Space)
	Revive()
}
