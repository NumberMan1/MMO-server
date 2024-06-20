package item

import (
	"github.com/NumberMan1/MMO-server/config/define"
	"github.com/NumberMan1/MMO-server/fight"
	"github.com/NumberMan1/MMO-server/protocol/gen/proto"
)

// Equipment 装备
type Equipment struct {
	*Item
	//装备类型
	equipType proto.EquipsType
	//装备的战斗属性
	attr *fight.Attributes
}

// Attr 获取装备的战斗属性
func (e *Equipment) Attr() *fight.Attributes {
	return e.attr
}

// EquipType 获取装备类型
func (e *Equipment) EquipType() proto.EquipsType {
	return e.equipType
}

// SetEquipType 设置装备类型
func (e *Equipment) SetEquipType(equipType proto.EquipsType) {
	e.equipType = equipType
}

func NewEquipmentByInfo(itemInfo *proto.ItemInfo) *Equipment {
	e := &Equipment{Item: NewItemByInfo(itemInfo)}
	e.SetEquipType(e.parseEquipType(e.Def().ItemType))
	e.loadAttributes()
	return e
}

func NewEquipmentByDefine(itemDef *define.ItemDefine, position int) *Equipment {
	e := &Equipment{Item: NewItem(itemDef, 1, position)}
	e.SetEquipType(e.parseEquipType(itemDef.EquipsType))
	e.loadAttributes()
	return e
}

// loadAttributes 加载装备属性
func (e *Equipment) loadAttributes() {
	e.attr = &fight.Attributes{
		Speed:     e.Def().Speed,
		HPMax:     e.Def().HPMax,
		MPMax:     e.Def().MPMax,
		AD:        e.Def().AD,
		AP:        e.Def().AP,
		DEF:       e.Def().DEF,
		MDEF:      e.Def().MDEF,
		CRI:       e.Def().CRI,
		CRD:       e.Def().CRD,
		STR:       e.Def().STR,
		INT:       e.Def().INT,
		AGI:       e.Def().AGI,
		HitRate:   e.Def().HitRate,
		DodgeRate: e.Def().DodgeRate,
		HpRegen:   e.Def().HpRegen,
		HpSteal:   e.Def().HpSteal,
	}
}

func (e *Equipment) parseEquipType(value string) proto.EquipsType {
	var t proto.EquipsType
	switch value {
	case "武器":
		t = proto.EquipsType_Weapon
	case "胸甲":
		t = proto.EquipsType_Chest
	case "腰带":
		t = proto.EquipsType_Belt
	case "裤子":
		t = proto.EquipsType_Legs
	case "鞋子":
		t = proto.EquipsType_Boots
	case "戒指":
		t = proto.EquipsType_Ring
	case "项链":
		t = proto.EquipsType_Neck
	case "翅膀":
		t = proto.EquipsType_Wings
	}
	return t
}
