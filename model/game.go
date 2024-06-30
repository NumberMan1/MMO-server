package model

import (
	"container/list"
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/mgr"
	"math"
)

// ActorIsNil 判断演员接口变量是否为nil
func ActorIsNil(actor IActor) bool {
	if actor == nil {
		return true
	}
	if ac, ok := actor.(*Character); ok && ac == nil {
		return true
	}
	if am, ok := actor.(*Monster); ok && am == nil {
		return true
	}
	return false
}

func GetUnit(entityId int) IActor {
	return mgr.GetEntityManagerInstance().GetEntity(entityId).(IActor)
}

// RangeUnit 查找范围内的人物
// 从spaceId的地图的position位置查找r范围的人物 []model.IActor
func RangeUnit(position *vector3.Vector3, spaceId, r int) *list.List {
	return mgr.GetEntityList(mgr.GetEntityManagerInstance(), spaceId, func(t IActor) bool {
		dist := vector3.GetDistance(position, t.Position())
		return !math.IsNaN(dist) && dist <= float64(r)
	})
}
