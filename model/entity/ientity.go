package entity

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/protocol/gen/proto"
)

type IEntity interface {
	PositionTime() float64
	Speed() int
	SetSpeed(speed int)
	Position() *vector3.Vector3
	SetPosition(position *vector3.Vector3)
	Direction() *vector3.Vector3
	SetDirection(direction *vector3.Vector3)
	EntityId() int
	EntityData() *proto.NetEntity
	SetEntityData(entity *proto.NetEntity)
	Update()
}
