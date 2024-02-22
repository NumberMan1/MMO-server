package model

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/define"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer/network"
	pt "github.com/NumberMan1/common/summer/protocol/gen/proto"
	"github.com/NumberMan1/common/summer/timeunit"
	"google.golang.org/protobuf/proto"
)

type Space struct {
	Id       int
	Name     string
	Def      *define.SpaceDefine
	FightMgr *FightMgr
	//当前场景中全部的角色
	characterDict map[int]*Character
	connCharacter map[network.Connection]*Character
	//当前场景中的野怪 <MonsterId,Monster>
	monsterDict    map[int]*Monster
	MonsterManager *MonsterManager
	SpawnManager   *SpawnManager
}

func NewSpace(def *define.SpaceDefine) *Space {
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
	s.FightMgr = NewFightMgr(s)
	s.MonsterManager.Init(s)
	s.SpawnManager.Init(s)
	return s
}

// CharacterJoin 角色进入场景
func (s *Space) CharacterJoin(character *Character) {
	logger.SLCInfo("角色进入场景:%d", character.Id())
	character.OnEnterSpace(s, character)
	s.characterDict[character.Id()] = character
	_, ok := s.connCharacter[character.Conn]
	if ok == false {
		s.connCharacter[character.Conn] = character
	}
	//把新进入的角色广播给其他玩家
	character.Info().Entity = character.EntityData()
	response := &pt.SpaceCharactersEnterResponse{
		SpaceId:       int32(s.Id),
		CharacterList: make([]*pt.NetActor, 0),
	}
	response.CharacterList = append(response.CharacterList, character.Info())
	for _, v := range s.characterDict {
		if v.Conn != character.Conn {
			v.Conn.Send(response)
		}
	}
	//新上线的角色需要获取全部角色
	ser := &pt.SpaceEnterResponse{
		Character: character.Info(),
		List:      make([]*pt.NetActor, 0),
	}
	for _, v := range s.characterDict {
		if v.Conn == character.Conn {
			continue
		}
		ser.List = append(ser.List, v.Info())
	}
	for _, v := range s.monsterDict {
		ser.List = append(ser.List, v.Info())
	}
	character.Conn.Send(ser)
}

// CharacterLeave 角色离开地图
// 客户端离线、切换地图
func (s *Space) CharacterLeave(character *Character) {
	logger.SLCInfo("角色离开场景:%d", character.EntityId())
	delete(s.characterDict, character.Id())
	response := &pt.SpaceCharacterLeaveResponse{
		EntityId: int32(character.EntityId()),
	}
	for _, v := range s.characterDict {
		v.Conn.Send(response)
	}
}

// Telport 同场景传送
func (s *Space) Telport(actor IActor, pos, dir *vector3.Vector3) {
	actor.SetPosition(pos)
	actor.SetDirection(dir)
	resp := &pt.SpaceEntitySyncResponse{
		EntitySync: &pt.NetEntitySync{
			Entity: actor.EntityData(),
			Force:  true,
		},
	}
	s.Broadcast(resp)
}

// UpdateEntity 广播更新Entity信息
func (s *Space) UpdateEntity(sync *pt.NetEntitySync) {
	//logger.SLCDebug("%v", sync)
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
	mon.OnEnterSpace(s, mon)
	resp := &pt.SpaceCharactersEnterResponse{
		SpaceId:       int32(s.Id), //场景ID
		CharacterList: make([]*pt.NetActor, 0),
	}
	resp.CharacterList = append(resp.CharacterList, mon.Info())
	s.Broadcast(resp)
}

// Broadcast 广播Proto消息给场景的全体玩家
func (s *Space) Broadcast(message proto.Message) {
	for _, v := range s.characterDict {
		v.Conn.Send(message)
	}
}

func (s *Space) Update() {
	s.SpawnManager.Update()
	s.FightMgr.OnUpdate(timeunit.DeltaTime)
}
