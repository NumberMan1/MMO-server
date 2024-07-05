package model

import (
	"github.com/NumberMan1/MMO-server/game_common/protocol/gen/proto"
	"github.com/NumberMan1/MMO-server/game_server/core/vector3"
	"github.com/NumberMan1/MMO-server/game_server/database"
	"github.com/NumberMan1/MMO-server/game_server/model/inventory_system/item"
)

// Character 角色
type Character struct {
	*Actor
	//当前角色的客户端连接
	Session *Session
	//当前角色对应的数据库对象
	Data *database.DbCharacter
	//背包
	Knapsack *Inventory
	//装备管理器
	EquipsManager *EquipsManager
}

func NewCharacter(dbCharacter *database.DbCharacter) *Character {
	c := &Character{
		Actor: NewActor(proto.EntityType_Character, dbCharacter.JobId, dbCharacter.Level,
			vector3.NewVector3(float64(dbCharacter.X), float64(dbCharacter.Y), float64(dbCharacter.Z)),
			vector3.Zero3()),
	}
	c.SetId(int(dbCharacter.ID))
	c.SetName(dbCharacter.Name)
	c.Info().Id = int32(dbCharacter.ID)
	c.Info().Name = dbCharacter.Name
	c.Info().Tid = int32(dbCharacter.JobId) //单位类型
	c.Info().Level = int32(dbCharacter.Level)
	c.Info().Exp = int64(dbCharacter.Exp)
	c.Info().SpaceId = int32(dbCharacter.SpaceId)
	c.Info().Gold = dbCharacter.Gold
	c.Info().Hp = float32(dbCharacter.Hp)
	c.Info().Mp = float32(dbCharacter.Mp)
	c.Data = dbCharacter
	//创建背包
	c.Knapsack = NewInventory(c)
	c.Knapsack.Init(c.Data.Knapsack)
	//初始化装备
	c.EquipsManager = NewEquipsManager(c)
	c.EquipsManager.Init(c.Data.EquipsData)
	return c
}

// CharacterId 玩家角色唯一ID
func (c *Character) CharacterId() int {
	return int(c.Data.ID)
}

// UseItem 使用物品
func (c *Character) UseItem(slotIndex int) {
	item1, ok := c.Knapsack.TrySlotItem(slotIndex)
	if !ok {
		return
	}
	if item1.ItemType() != item.ItemType_Consumable {
		return
	}
	item1.SetAmount(item1.Amount() - 1)
	if item1.Amount() <= 0 {
		c.Knapsack.SetItem(slotIndex, nil)
	}
	//发送消息
	c.SendInventory(true, false, false)
	//物品效果
	if item1.Id() == 1001 {
		c.SetAndUpdateHp(c.Hp() + 50)
	}
	if item1.Id() == 1002 {
		c.SetAndUpdateMp(c.Mp() + 50)
	}
}

// SendInventory 发送背包到客户端
func (c *Character) SendInventory(isKnapsack, isStorage, isEquips bool) {
	rsp := &proto.InventoryResponse{
		EntityId: int32(c.EntityId()),
	}
	if isKnapsack {
		rsp.KnapsackInfo = c.Knapsack.InventoryInfo()
	}
	if isStorage {

	}
	if isEquips {

	}
	c.Session.Send(rsp)
}
