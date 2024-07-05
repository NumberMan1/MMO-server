package item

import (
	"github.com/NumberMan1/MMO-server/game_common/protocol/gen/proto"
)

// Consumable 消耗品
type Consumable struct {
	*Item
}

func NewConsumable(itemInfo *proto.ItemInfo) *Consumable {
	return &Consumable{Item: NewItemByInfo(itemInfo)}
}
