package item

import (
	"errors"
	"github.com/NumberMan1/MMO-server/game_common/protocol/gen/proto"
	define2 "github.com/NumberMan1/MMO-server/game_server/config/define"
)

// ItemIsNil 判断物品接口变量是否为nil
func ItemIsNil(item IItem) bool {
	if item == nil {
		return true
	}
	if it, ok := item.(*Consumable); ok && it == nil {
		return true
	}
	if it, ok := item.(*Material); ok && it == nil {
		return true
	}
	if it, ok := item.(*Equipment); ok && it == nil {
		return true
	}
	return false
}

func NewItem(def *define2.ItemDefine, amount, position int) *Item {
	i := &Item{def: def, amount: amount, position: position}
	i.init(def.ID, def.Name, ItemType_Material, Quality_Common, def.Description, def.Capicity, def.BuyPrice, def.SellPrice, def.Icon)
	switch def.ItemType {
	case "消耗品":
		i.SetItemType(ItemType_Consumable)
	case "道具":
		i.SetItemType(ItemType_Material)
	case "装备":
		i.SetItemType(ItemType_Equipment)
	}
	switch def.Quality {
	case "普通":
		i.SetQuality(Quality_Common)
	case "非凡":
		i.SetQuality(Quality_Uncommon)
	case "稀有":
		i.SetQuality(Quality_Rare)
	case "史诗":
		i.SetQuality(Quality_Epic)
	case "传说":
		i.SetQuality(Quality_Legendary)
	case "神器":
		i.SetQuality(Quality_Artifact)
	}
	return i
}

func NewItemByInfo(info *proto.ItemInfo) *Item {
	def := define2.GetDataManagerInstance().Items[int(info.ItemId)]
	i := NewItem(def, 1, 0)
	i.itemInfo = info
	i.SetAmount(int(info.Amount))
	i.SetPosition(int(info.Position))
	return i
}

// CreateItemByItemId 通过物品id统一创建物品的方法
func CreateItemByItemId(itemId, amount, position int) (IItem, error) {
	info := proto.ItemInfo{
		ItemId:   int32(itemId),
		Amount:   int32(amount),
		Position: int32(position),
	}
	return CreateItemByItemInfo(&info)
}

// CreateItemByItemInfo 通过物品信息统一创建物品的方法
func CreateItemByItemInfo(info *proto.ItemInfo) (IItem, error) {
	def := define2.GetDataManagerInstance().Items[int(info.ItemId)]
	var i IItem
	switch def.ItemType {
	case "消耗品":
		i = NewConsumable(info)
	case "道具":
		i = NewMaterial(info)
	case "装备":
		i = NewEquipmentByInfo(info)
	default:
		return nil, errors.New("物品初始化失败:unknown item type")
	}
	return i, nil
}

// Item 物品基类
type Item struct {
	id           int      //物品ID
	name         string   //物品名称
	itemType     ItemType //物品种类
	quality      Quality  //物品品质
	description  string   //物品描述
	capacity     int      //物品叠加数量上限
	buyPrice     int      //物品买入价格
	sellingPrice int      //物品卖出价格
	sprite       string   //存放物品的图片路径，通过Resources加载
	def          *define2.ItemDefine
	amount       int //数量
	position     int //所处位置
	itemInfo     *proto.ItemInfo
}

// Amount 获取数量
func (i *Item) Amount() int {
	return i.amount
}

// SetAmount 设置数量
func (i *Item) SetAmount(amount int) {
	i.amount = amount
}

// Position 获取所处位置
func (i *Item) Position() int {
	return i.position
}

// SetPosition 设置所处位置
func (i *Item) SetPosition(position int) {
	i.position = position
}

func (i *Item) ItemInfo() *proto.ItemInfo {
	if i.itemInfo == nil {
		i.itemInfo = &proto.ItemInfo{ItemId: int32(i.Id())}
	}
	i.itemInfo.Amount = int32(i.Amount())
	i.itemInfo.Position = int32(i.Position())
	return i.itemInfo
}

func (i *Item) init(id int, name string, itemType ItemType, quality Quality, description string, capacity, buyPrice, sellPrice int, sprite string) {
	i.SetId(id)
	i.SetName(name)
	i.SetItemType(itemType)
	i.SetQuality(quality)
	i.SetDescription(description)
	i.SetCapacity(capacity)
	i.SetBuyPrice(buyPrice)
	i.SetSellingPrice(sellPrice)
	i.SetSprite(sprite)
}

func (i *Item) Def() *define2.ItemDefine {
	return i.def // 获取物品定义
}

// Id 获取物品ID
func (i *Item) Id() int {
	return i.id
}

// SetId 设置物品ID
func (i *Item) SetId(id int) {
	i.id = id
}

// Name 获取物品名称
func (i *Item) Name() string {
	return i.name
}

// SetName 设置物品名称
func (i *Item) SetName(name string) {
	i.name = name
}

// ItemType 获取物品种类
func (i *Item) ItemType() ItemType {
	return i.itemType
}

// SetItemType 设置物品种类
func (i *Item) SetItemType(itemType ItemType) {
	i.itemType = itemType
}

// Quality 获取物品品质
func (i *Item) Quality() Quality {
	return i.quality
}

// SetQuality 设置物品品质
func (i *Item) SetQuality(quality Quality) {
	i.quality = quality
}

// Description 获取物品描述
func (i *Item) Description() string {
	return i.description
}

// SetDescription 设置物品描述
func (i *Item) SetDescription(description string) {
	i.description = description
}

// Capacity 获取物品叠加数量上限
func (i *Item) Capacity() int {
	return i.capacity
}

// SetCapacity 设置物品叠加数量上限
func (i *Item) SetCapacity(capacity int) {
	i.capacity = capacity
}

// BuyPrice 获取物品买入价格
func (i *Item) BuyPrice() int {
	return i.buyPrice
}

// SetBuyPrice 设置物品买入价格
func (i *Item) SetBuyPrice(buyPrice int) {
	i.buyPrice = buyPrice
}

// SellingPrice 获取物品卖出价格
func (i *Item) SellingPrice() int {
	return i.sellingPrice
}

// SetSellingPrice 设置物品卖出价格
func (i *Item) SetSellingPrice(sellingPrice int) {
	i.sellingPrice = sellingPrice
}

// Sprite 获取物品图片路径
func (i *Item) Sprite() string {
	return i.sprite
}

// SetSprite 设置物品图片路径
func (i *Item) SetSprite(sprite string) {
	i.sprite = sprite
}
