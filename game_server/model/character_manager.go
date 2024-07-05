package model

import (
	"github.com/NumberMan1/MMO-server/game_common"
	"github.com/NumberMan1/MMO-server/game_server/database"
	"github.com/NumberMan1/MMO-server/game_server/mgr"
	"github.com/NumberMan1/common/global/variable"
	"github.com/NumberMan1/common/ns/singleton"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

var (
	singleCharacterManager = singleton.Singleton{}
)

// CharacterManager 统一管理全部的角色（创建，移除，获取）
type CharacterManager struct {
	//游戏里全部的角色 <ChrId,ChrObj>
	characters *sync.Map
}

func init() {
	result := GetCharacterManagerInstance()
	//每隔5秒保存Data到数据库
	game_common.GetScheduleInstance().AddTask(result.save, 5*time.Second, 0)
}

func GetCharacterManagerInstance() *CharacterManager {
	result, _ := singleton.GetOrDo[*CharacterManager](&singleCharacterManager, func() (*CharacterManager, error) {
		return &CharacterManager{
			characters: &sync.Map{},
		}, nil
	})
	return result
}

func (cm *CharacterManager) CreateCharacter(dbChr *database.DbCharacter) *Character {
	character := NewCharacter(dbChr)
	cm.characters.Store(character.Id(), character)
	mgr.GetEntityManagerInstance().AddEntity(dbChr.SpaceId, character)
	return character
}

func (cm *CharacterManager) RemoveCharacter(chrId int) {
	character, ok := cm.characters.Load(chrId)
	if ok {
		cm.characters.Delete(chrId)
		chr := character.(*Character)
		mgr.GetEntityManagerInstance().RemoveEntity(chr.Data.SpaceId, chr)
	}
}

func (cm *CharacterManager) GetCharacter(chrId int) *Character {
	c, ok := cm.characters.Load(chrId)
	if ok {
		return c.(*Character)
	} else {
		return nil
	}
}

func (cm *CharacterManager) Clear() {
	cm.characters = &sync.Map{}
}

func (cm *CharacterManager) save() error {
	cm.characters.Range(func(key, value any) bool {
		chr := value.(*Character)
		chr.Data.X = int(chr.Position().X)
		chr.Data.Y = int(chr.Position().Y)
		chr.Data.Z = int(chr.Position().Z)
		chr.Data.JobId = int(chr.Info().Tid)
		chr.Data.Hp = int(chr.Info().Hp)
		chr.Data.Mp = int(chr.Info().Mp)
		chr.Data.Exp = int(chr.Info().Exp)
		chr.Data.Level = int(chr.Info().Level)
		chr.Data.Gold = chr.Info().Gold
		chr.Data.SpaceId = int(chr.Info().SpaceId)
		bs, _ := proto.Marshal(chr.Knapsack.InventoryInfo())
		chr.Data.Knapsack = bs
		bs, _ = proto.Marshal(chr.EquipsManager.InventoryInfo())
		chr.Data.EquipsData = bs
		variable.GDb.Save(chr.Data)
		return true
	})
	return nil
}
