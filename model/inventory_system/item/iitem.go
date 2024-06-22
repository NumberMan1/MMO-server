package item

import (
	"github.com/NumberMan1/MMO-server/config/define"
	"github.com/NumberMan1/MMO-server/protocol/gen/proto"
)

// IItem 物品接口
type IItem interface {
	// Amount 获取数量
	Amount() int
	// SetAmount 设置数量
	SetAmount(amount int)
	// Position 获取所处位置
	Position() int
	// SetPosition 设置所处位置
	SetPosition(position int)
	// ItemInfo 获取物品信息
	ItemInfo() *proto.ItemInfo
	// init 初始化物品
	init(id int, name string, itemType ItemType, quality Quality, description string, capacity, buyPrice, sellPrice int, sprite string)
	// Def 获取物品定义
	Def() *define.ItemDefine
	// Id 获取物品ID
	Id() int
	// SetId 设置物品ID
	SetId(id int)
	// Name 获取物品名称
	Name() string
	// SetName 设置物品名称
	SetName(name string)
	// ItemType 获取物品种类
	ItemType() ItemType
	// SetItemType 设置物品种类
	SetItemType(itemType ItemType)
	// Quality 获取物品品质
	Quality() Quality
	// SetQuality 设置物品品质
	SetQuality(quality Quality)
	// Description 获取物品描述
	Description() string
	// SetDescription 设置物品描述
	SetDescription(description string)
	// Capacity 获取物品叠加数量上限
	Capacity() int
	// SetCapacity 设置物品叠加数量上限
	SetCapacity(capacity int)
	// BuyPrice 获取物品买入价格
	BuyPrice() int
	// SetBuyPrice 设置物品买入价格
	SetBuyPrice(buyPrice int)
	// SellingPrice 获取物品卖出价格
	SellingPrice() int
	// SetSellingPrice 设置物品卖出价格
	SetSellingPrice(sellingPrice int)
	// Sprite 获取物品图片路径
	Sprite() string
	// SetSprite 设置物品图片路径
	SetSprite(sprite string)
}
