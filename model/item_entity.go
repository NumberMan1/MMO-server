package model

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/inventory/item"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
)

type ItemEntity struct {
	*Actor
	item *item.Item
}

// CreateItemEntity 在场景里创建物品
func CreateItemEntity(space *Space, item1 *item.Item, pos, dir *vector3.Vector3) *ItemEntity {
	entity := NewItemEntity(proto.EntityType_Item, 0, 0, pos, dir)
	entity.SetItem(item1)
	entity.Info().ItemInfo = entity.Item().ItemInfo()
	GetEntityManagerInstance().AddEntity(space.Id, entity)
	space.EntityEnter(entity)
	return entity
}

func NewItemEntity(t proto.EntityType, tid, level int, position, direction *vector3.Vector3) *ItemEntity {
	return &ItemEntity{Actor: NewActor(t, tid, level, position, direction)}
}

func (i *ItemEntity) Item() *item.Item {
	return i.item
}

func (i *ItemEntity) SetItem(item *item.Item) {
	i.item = item
}
