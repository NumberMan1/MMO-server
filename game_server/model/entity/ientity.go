package entity

import (
	"github.com/NumberMan1/MMO-server/game_common/protocol/gen/proto"
	"github.com/NumberMan1/MMO-server/game_server/core/vector3"
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
