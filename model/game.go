package model

import (
	"container/list"
	"github.com/NumberMan1/MMO-server/core/vector3"
)

func GetUnit(entityId int) IActor {
	return GetEntityManagerInstance().GetEntity(entityId).(IActor)
}

// RangeUnit 查找范围内的人物
// 从spaceId的地图的position位置查找r范围的人物 []model.IActor
func RangeUnit(position *vector3.Vector3, spaceId, r int) *list.List {
	return GetEntityList(GetEntityManagerInstance(), spaceId, func(t IActor) bool {
		return vector3.GetDistance(position, t.Position()) <= float64(r)
	})
}
