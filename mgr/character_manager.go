package mgr

import (
	"github.com/NumberMan1/MMO-server/database"
	"github.com/NumberMan1/MMO-server/model"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/ns/singleton"
	"github.com/NumberMan1/common/summer"
	"github.com/NumberMan1/common/summer/timeunit"
	"sync"
)

var (
	singleCharacterManager = singleton.Singleton{}
)

// CharacterManager 统一管理全部的角色（创建，移除，获取）
type CharacterManager struct {
	//游戏里全部的角色 <ChrId,ChrObj>
	characters sync.Map
}

func GetCharacterManagerInstance() *CharacterManager {
	result, _ := singleton.GetOrDo[*CharacterManager](&singleCharacterManager, func() (*CharacterManager, error) {
		return &CharacterManager{
			characters: sync.Map{},
		}, nil
	})
	summer.GetScheduleInstance().AddTask(result.save, timeunit.Seconds, 2, 0)
	return result
}

func (cm *CharacterManager) CreateCharacter(dbChr database.DbCharacter) *model.Character {
	character := model.CharacterFromDbCharacter(dbChr)
	cm.characters.Store(character.Id, character)
	GetEntityManagerInstance().AddEntity(dbChr.SpaceId, character.Entity)
	return character
}

func (cm *CharacterManager) RemoveCharacter(chrId int) {
	character, ok := cm.characters.Load(chrId)
	if ok {
		cm.characters.Delete(chrId)
		chr := character.(*model.Character)
		GetEntityManagerInstance().RemoveEntity(chr.Data.SpaceId, chr.Entity)
	}
}

func (cm *CharacterManager) GetCharacter(chrId int) *model.Character {
	c, ok := cm.characters.Load(chrId)
	if ok {
		return c.(*model.Character)
	} else {
		return nil
	}
}

func (cm *CharacterManager) Clear() {
	cm.characters = sync.Map{}
}

func (cm *CharacterManager) save() {
	cm.characters.Range(func(key, value any) bool {
		logger.SLCDebug("save character:%v", value.(*model.Character).Data)
		database.OrmDb.Save(value.(*model.Character).Data)
		return true
	})
}
