package fight

import (
	"github.com/NumberMan1/MMO-server/game_server/core/vector3"
	"github.com/NumberMan1/MMO-server/game_server/model/entity"
)

type ISCObject interface {
	GetRealObj() any
	GetId() int
	GetPosition() *vector3.Vector3
	GetDirection() *vector3.Vector3
}

// SCObject Server-Client-Object，可以代表一个人或者位置
type SCObject struct {
	realObj any
}

func NewSCObject(realObj any) *SCObject {
	return &SCObject{realObj: realObj}
}

func (sco *SCObject) GetRealObj() any {
	return sco.realObj
}

func (sco *SCObject) GetId() int {
	return 0
}

func (sco *SCObject) GetPosition() *vector3.Vector3 {
	return vector3.Zero3()
}

func (sco *SCObject) GetDirection() *vector3.Vector3 {
	return vector3.Zero3()
}

// SCEntity 定义SCEntity类，继承自SCObject
type SCEntity struct {
	*SCObject
}

func NewSCEntity(realObj entity.IEntity) *SCEntity {
	return &SCEntity{
		SCObject: NewSCObject(realObj),
	}
}

func (sce *SCEntity) getObj() entity.IEntity {
	return sce.GetRealObj().(entity.IEntity)
}

func (sce *SCEntity) GetId() int {
	return sce.getObj().EntityId()
}

func (sce *SCEntity) GetPosition() *vector3.Vector3 {
	return sce.getObj().Position()
}

func (sce *SCEntity) GetDirection() *vector3.Vector3 {
	return sce.getObj().Direction()
}

type SCPosition struct {
	*SCObject
}

func NewSCPosition(realObj *vector3.Vector3) *SCPosition {
	return &SCPosition{
		SCObject: NewSCObject(realObj),
	}
}

func (SCP *SCPosition) GetPosition() *vector3.Vector3 {
	return SCP.GetRealObj().(*vector3.Vector3)
}
