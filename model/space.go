package model

import (
	"github.com/NumberMan1/MMO-server/define"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer/network"
	pt "github.com/NumberMan1/common/summer/protocol/gen/proto"
)

type Space struct {
	Id   int
	Name string
	Def  define.SpaceDefine
	//当前场景中全部的角色
	characterDict map[int]*Character
	connCharacter map[network.Connection]*Character
	//当前场景中的野怪 <MonsterId,Monster>
	monsterDict    map[int]*Monster
	MonsterManager *MonsterManager
	SpawnManager   *SpawnManager
}

func NewSpace(def define.SpaceDefine) *Space {
	s := &Space{
		Id:             def.SID,
		Name:           def.Name,
		Def:            def,
		characterDict:  map[int]*Character{},
		connCharacter:  map[network.Connection]*Character{},
		monsterDict:    map[int]*Monster{},
		MonsterManager: NewMonsterManager(),
		SpawnManager:   NewSpawnManager(),
	}
	s.MonsterManager.Init(s)
	s.SpawnManager.Init(s)
	return s
}

// CharacterJoin 角色进入场景
func (s *Space) CharacterJoin(conn network.Connection, character *Character) {
	logger.SLCInfo("角色进入场景:%d", character.Id)
	conn.Get("Session").(*Session).Character = character //把角色存入Session
	character.OnEnterSpace(s)
	character.Conn = conn
	s.characterDict[character.Id()] = character
	_, ok := s.connCharacter[conn]
	if ok == false {
		s.connCharacter[conn] = character
	}
	//把新进入的角色广播给其他玩家
	character.Info().Entity = character.EntityData()
	response := &pt.SpaceCharactersEnterResponse{
		SpaceId:       int32(s.Id),
		CharacterList: make([]*pt.NCharacter, 0),
	}
	response.CharacterList = append(response.CharacterList, character.Info())
	for _, v := range s.characterDict {
		if v.Conn != conn {
			v.Conn.Send(response)
		}
	}
	//新上线的角色需要获取全部角色
	for _, v := range s.characterDict {
		if v.Conn == conn {
			continue
		}
		response.CharacterList = make([]*pt.NCharacter, 0)
		response.CharacterList = append(response.CharacterList, v.Info())
		conn.Send(response)
	}
	for _, v := range s.monsterDict {
		response.CharacterList = append(response.CharacterList, v.Info())
	}
	conn.Send(response)
}

// CharacterLeave 角色离开地图
// 客户端离线、切换地图
func (s *Space) CharacterLeave(conn network.Connection, character *Character) {
	logger.SLCInfo("角色离开场景:%d", character.EntityId())
	delete(s.characterDict, character.Id())
	response := &pt.SpaceCharacterLeaveResponse{
		EntityId: int32(character.EntityId()),
	}
	for _, v := range s.characterDict {
		v.Conn.Send(response)
	}
}

// UpdateEntity 广播更新Entity信息
func (s *Space) UpdateEntity(sync *pt.NEntitySync) {
	for _, v := range s.characterDict {
		if v.EntityId() == int(sync.Entity.Id) {
			v.SetEntityData(sync.GetEntity())
		} else {
			response := &pt.SpaceEntitySyncResponse{EntitySync: sync}
			v.Conn.Send(response)
		}
	}
}

// MonsterEnter 怪物进入场景
func (s *Space) MonsterEnter(mon *Monster) {
	s.monsterDict[mon.Id()] = mon
	mon.OnEnterSpace(s)
	resp := &pt.SpaceCharactersEnterResponse{
		SpaceId:       int32(s.Id),
		CharacterList: make([]*pt.NCharacter, 0),
	}
	resp.CharacterList = append(resp.CharacterList, mon.Info())
	for _, v := range s.characterDict {
		v.Conn.Send(resp)
	}
}
func (s *Space) Update() {
	s.SpawnManager.Update()
}
