package model

import (
	"container/list"
	"github.com/NumberMan1/MMO-server/define"
)

// SpawnManager 刷怪管理器
type SpawnManager struct {
	Rules *list.List
	Space *Space
}

func NewSpawnManager() *SpawnManager {
	return &SpawnManager{Rules: list.New()}
}

func (sm *SpawnManager) Init(space *Space) {
	sm.Space = space
	//根据当前场景加载对应的规则
	rules := make([]define.SpawnDefine, 0)
	for _, v := range define.GetDataManagerInstance().Spawns {
		if v.SpaceId == space.Id {
			rules = append(rules, v)
		}
	}
	for _, v := range rules {
		sm.Rules.PushBack(NewSpawner(v, space))
	}
}

func (sm *SpawnManager) Update() {
	for e := sm.Rules.Front(); e != nil; e = e.Next() {
		s := e.Value.(*Spawner)
		s.Update()
	}
}
