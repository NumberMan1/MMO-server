package model

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/define"
	"github.com/NumberMan1/common/logger"
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
	//connCharacter map[network.Connection]*Character
	//当前场景中的全部演员 <entityId,IActor>
	actorDict      map[int]IActor
	MonsterManager *MonsterManager
	SpawnManager   *SpawnManager
}

func NewSpace(def *define.SpaceDefine) *Space {
	s := &Space{
		Id:            def.SID,
		Name:          def.Name,
		Def:           def,
		characterDict: map[int]*Character{},
		//connCharacter:  map[network.Connection]*Character{},
		actorDict:      map[int]IActor{},
		MonsterManager: NewMonsterManager(),
		SpawnManager:   NewSpawnManager(),
	}
	s.FightMgr = NewFightMgr(s)
	s.MonsterManager.Init(s)
	s.SpawnManager.Init(s)
	return s
}

// CharacterJoin 主角进入场景
func (s *Space) CharacterJoin(character *Character) {
	logger.SLCInfo("玩家进入场景 Chr[%v],Entity[%v]", character.Id(), character.EntityId())
	//character.OnEnterSpace(s, character)
	//记录到主角字典
	s.characterDict[character.EntityId()] = character
	character.SetSpace(s)
	//拉取附近的演员信息
	rsp := &pt.SpaceEnterResponse{
		Character: character.Info(),
		List:      make([]*pt.NetActor, 0),
	}
	for _, c := range s.actorDict {
		rsp.List = append(rsp.List, c.Info())
	}
	character.Conn.Send(rsp)
}

// EntityLeave 演员离开场景
// 客户端离线、切换地图
func (s *Space) EntityLeave(actor IActor) {
	logger.SLCInfo("角色离开场景:%d", actor.EntityId())
	delete(s.actorDict, actor.EntityId())
	response := &pt.SpaceCharacterLeaveResponse{
		EntityId: int32(actor.EntityId()),
	}
	s.Broadcast(response)
	//如果是主角
	if chr, ok := actor.(*Character); ok {
		delete(s.characterDict, chr.EntityId())
	}
}

// Teleport 同场景传送
func (s *Space) Teleport(actor IActor, pos, dir *vector3.Vector3) {
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

// EntityEnter 演员进入场景
func (s *Space) EntityEnter(actor IActor) {
	logger.SLCInfo("角色进入场景:eid=%v", actor.EntityId())
	s.actorDict[actor.EntityId()] = actor
	actor.OnEnterSpace(s, actor)
	resp := &pt.SpaceCharactersEnterResponse{
		SpaceId:       int32(s.Id), //场景ID
		CharacterList: make([]*pt.NetActor, 0),
	}
	resp.CharacterList = append(resp.CharacterList, actor.Info())
	s.Broadcast(resp)
	//如果是主角
	if chr, ok := actor.(*Character); ok {
		s.CharacterJoin(chr)
	}
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
