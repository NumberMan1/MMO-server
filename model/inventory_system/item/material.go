package item

import "github.com/NumberMan1/MMO-server/protocol/gen/proto"

// Material 材料
type Material struct {
	*Item
}

func NewMaterial(itemInfo *proto.ItemInfo) *Material {
	return &Material{Item: NewItemByInfo(itemInfo)}
}
