package model

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/inventory_system/item"
	"github.com/NumberMan1/MMO-server/mgr"
	"github.com/NumberMan1/MMO-server/protocol/gen/proto"
	"github.com/NumberMan1/common/global/variable"
	"go.uber.org/zap"
)

// ItemEntity 场景里的物品
type ItemEntity struct {
	*Actor
	//真正的物品对象
	item item.IItem
}

// CreateItemEntity 在场景里创建物品
func CreateItemEntity(space *Space, item1 item.IItem, pos, dir *vector3.Vector3) *ItemEntity {
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
	it, err := item.CreateItemByItemId(itemId, itemAmount, 0)
	if err != nil {
		variable.Log.Error("CreateItemEntityById", zap.Error(err))
		return nil
	}
	return CreateItemEntity(sp, it, pos, dir)
}

func NewItemEntity(t proto.EntityType, tid, level int, position, direction *vector3.Vector3) *ItemEntity {
	return &ItemEntity{Actor: NewActor(t, tid, level, position, direction)}
}

func (i *ItemEntity) Item() item.IItem {
	return i.item
}

func (i *ItemEntity) SetItem(item item.IItem) {
	i.item = item
}
