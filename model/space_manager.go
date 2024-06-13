package model

import (
	"github.com/NumberMan1/MMO-server/config/define"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ns/singleton"
)

var (
	singleSpaceManager = singleton.Singleton{}
)

type SpaceManager struct {
	//地图字典
	dict map[int]*Space
}

func (sm *SpaceManager) Init() {
	for k, s := range define.GetDataManagerInstance().Spaces {
		sm.dict[k] = NewSpace(s)
		logger.SLCInfo("初始化地图:%s", s.Name)
	}
}

func (sm *SpaceManager) GetSpace(spaceId int) *Space {
	s, ok := sm.dict[spaceId]
	if ok {
		return s
	} else {
		return nil
	}
}

func GetSpaceManagerInstance() *SpaceManager {
	result, _ := singleton.GetOrDo[*SpaceManager](&singleSpaceManager, func() (*SpaceManager, error) {
		return &SpaceManager{
			dict: map[int]*Space{},
		}, nil
	})
	return result
}

func (sm *SpaceManager) Update() {
	for _, s := range sm.dict {
		s.Update()
	}
}
