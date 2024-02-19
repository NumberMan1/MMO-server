package model

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
	"time"
)

// Entity 在MMO世界进行同步的实体
type Entity struct {
	speed      int              // //移动速度
	position   *vector3.Vector3 //位置
	direction  *vector3.Vector3 //方向
	netObj     *proto.NetEntity //网络对象
	lastUpdate int64            //最后一次更新位置的时间戳
}

// PositionTime 距离上次位置更新的间隔（秒）
func (e *Entity) PositionTime() float64 {
	return float64(time.Now().UnixMilli()-e.lastUpdate) * 0.001
}

func (e *Entity) Speed() int {
	return e.speed
}

func (e *Entity) SetSpeed(speed int) {
	e.speed = speed
	e.netObj.Speed = int32(speed)
}

func (e *Entity) Position() *vector3.Vector3 {
	return e.position
}

func (e *Entity) SetPosition(position *vector3.Vector3) {
	e.position = position
	e.netObj.Position = &proto.Vec3{
		X: int32(position.X),
		Y: int32(position.Y),
		Z: int32(position.Z),
	}
	e.lastUpdate = time.Now().UnixMilli()
}

func (e *Entity) Direction() *vector3.Vector3 {
	return e.direction
}

func (e *Entity) SetDirection(direction *vector3.Vector3) {
	e.direction = direction
	e.netObj.Direction = &proto.Vec3{
		X: int32(direction.X),
		Y: int32(direction.Y),
		Z: int32(direction.Z),
	}
}

func NewEntity(position, direction *vector3.Vector3) *Entity {
	e := &Entity{netObj: &proto.NetEntity{}}
	e.SetPosition(position)
	e.SetDirection(direction)
	return e
}

func (e *Entity) EntityId() int {
	return int(e.netObj.Id)
}

func (e *Entity) EntityData() *proto.NetEntity {
	return e.netObj
}

func (e *Entity) SetEntityData(entity *proto.NetEntity) {
	e.SetPosition(vector3.NewVector3(float64(entity.Position.X), float64(entity.Position.Y), float64(entity.Position.Z)))
	e.SetDirection(vector3.NewVector3(float64(entity.Direction.X), float64(entity.Direction.Y), float64(entity.Direction.Z)))
	e.SetSpeed(int(entity.Speed))
}

func (e *Entity) Update() {

}
