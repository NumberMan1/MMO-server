package model

import (
	define2 "github.com/NumberMan1/MMO-server/config/define"
)

// SpawnManager 刷怪管理器
type SpawnManager struct {
	Rules []*Spawner
	Space *Space
}

func NewSpawnManager() *SpawnManager {
	return &SpawnManager{Rules: make([]*Spawner, 0)}
}

func (sm *SpawnManager) Init(space *Space) {
	sm.Space = space
	//根据当前场景加载对应的规则
	rules := make([]*define2.SpawnDefine, 0)
	for _, v := range define2.GetDataManagerInstance().Spawns {
		if v.SpaceId == space.Id {
			rules = append(rules, v)
		}
	}
	for _, v := range rules {
		sm.Rules = append(sm.Rules, NewSpawner(v, space))
	}
}

func (sm *SpawnManager) Update() {
	for _, rule := range sm.Rules {
		rule.Update()
	}
}
