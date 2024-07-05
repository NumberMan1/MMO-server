package model

import (
	"github.com/NumberMan1/MMO-server/game_common/protocol/gen/proto"
	"github.com/NumberMan1/MMO-server/game_server/core/vector3"
	"github.com/NumberMan1/MMO-server/game_server/mgr"
	item2 "github.com/NumberMan1/MMO-server/game_server/model/inventory_system/item"
	"github.com/NumberMan1/common/global/variable"
	"go.uber.org/zap"
)

// ItemEntity 场景里的物品
type ItemEntity struct {
	*Actor
	//真正的物品对象
	item item2.IItem
}

// CreateItemEntity 在场景里创建物品
func CreateItemEntity(space *Space, item1 item2.IItem, pos, dir *vector3.Vector3) *ItemEntity {
	entity := NewItemEntity(proto.EntityType_Item, 0, 0, pos, dir)
	entity.SetItem(item1)
	entity.Info().ItemInfo = entity.Item().ItemInfo()
	mgr.GetEntityManagerInstance().AddEntity(space.Id, entity)
	space.EntityEnter(entity)
	return entity
}

// CreateItemEntityById 通过id在场景里创建物品
func CreateItemEntityById(spaceId, itemId, itemAmount int, pos, dir *vector3.Vector3) *ItemEntity {
	sp := GetSpaceManagerInstance().GetSpace(spaceId)
	it, err := item2.CreateItemByItemId(itemId, itemAmount, 0)
	if err != nil {
		variable.Log.Error("CreateItemEntityById", zap.Error(err))
		return nil
	}
	return CreateItemEntity(sp, it, pos, dir)
}

func NewItemEntity(t proto.EntityType, tid, level int, position, direction *vector3.Vector3) *ItemEntity {
	return &ItemEntity{Actor: NewActor(t, tid, level, position, direction)}
}

func (i *ItemEntity) Item() item2.IItem {
	return i.item
}

func (i *ItemEntity) SetItem(item item2.IItem) {
	i.item = item
}
