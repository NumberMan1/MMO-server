package model

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
)

// MonsterManager 每个地图都有怪物管理器
type MonsterManager struct {
	// <EntityId,Monster>
	dict  map[int]*Monster
	space *Space
}

func NewMonsterManager() *MonsterManager {
	return &MonsterManager{dict: map[int]*Monster{}}
}

func (mm *MonsterManager) Init(space *Space) {
	mm.space = space
}

func (mm *MonsterManager) Create(tid, level int, pos, dir *vector3.Vector3) *Monster {
	monster := NewMonster(tid, level, pos, dir)
	GetEntityManagerInstance().AddEntity(mm.space.Id, monster)
	monster.Info().SpaceId = int32(mm.space.Id)
	monster.Info().GetEntity().Id = int32(monster.EntityId())
	mm.dict[monster.EntityId()] = monster
	monster.SetId(monster.EntityId())
	mm.space.EntityEnter(monster)
	return monster
}
