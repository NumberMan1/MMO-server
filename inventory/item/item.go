package item

import (
	define2 "github.com/NumberMan1/MMO-server/config/define"
	"github.com/NumberMan1/MMO-server/protocol/gen/proto"
)

type ItemType int

const (
	ItemType_Material   ItemType = iota //材料&道具
	ItemType_Consumable                 //消耗品
	ItemType_Equipment                  //武器&装备
)

type Quality int

const (
	Quality_Common    = iota //普通
	Quality_Uncommon         //非凡
	Quality_Rare             //稀有
	Quality_Epic             //史诗
	Quality_Legendary        //传说
	Quality_Artifact         //神器
)

// Item 物品基类
type Item struct {
	id           int
	name         string
	itemType     ItemType
	quality      Quality
	description  string
	capacity     int
	buyPrice     int
	sellingPrice int
	sprite       string //存放物品的图片路径，通过Resources加载
	def          *define2.ItemDefine
	Amount       int //数量
	Position     int //所处位置
	itemInfo     *proto.ItemInfo
}

func (i *Item) ItemInfo() *proto.ItemInfo {
	if i.itemInfo == nil {
		i.itemInfo = &proto.ItemInfo{ItemId: int32(i.Id())}
	}
	i.itemInfo.Amount = int32(i.Amount)
	i.itemInfo.Position = int32(i.Position)
	return i.itemInfo
}

func NewItem(def *define2.ItemDefine, amount, position int) *Item {
	i := &Item{def: def, Amount: amount, Position: position}
	i.SetId(def.GetId())
	i.SetName(def.Name)
	i.SetDescription(def.Description)
	i.SetCapacity(def.Capicity)
	i.SetBuyPrice(def.BuyPrice)
	i.SetSellingPrice(def.SellPrice)
	i.SetSprite(def.Icon)
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

func NewItemByItemId(itemId, amount, position int) *Item {
	def := define2.GetDataManagerInstance().Items[itemId]
	return NewItem(def, amount, position)
}

func (i *Item) Def() *define2.ItemDefine {
	return i.def
}

func (i *Item) Id() int {
	return i.id
}

func (i *Item) SetId(id int) {
	i.id = id
}

func (i *Item) Name() string {
	return i.name
}

func (i *Item) SetName(name string) {
	i.name = name
}

func (i *Item) ItemType() ItemType {
	return i.itemType
}

func (i *Item) SetItemType(itemType ItemType) {
	i.itemType = itemType
}

func (i *Item) Quality() Quality {
	return i.quality
}

func (i *Item) SetQuality(quality Quality) {
	i.quality = quality
}

func (i *Item) Description() string {
	return i.description
}

func (i *Item) SetDescription(description string) {
	i.description = description
}

func (i *Item) Capacity() int {
	return i.capacity
}

func (i *Item) SetCapacity(capacity int) {
	i.capacity = capacity
}

func (i *Item) BuyPrice() int {
	return i.buyPrice
}

func (i *Item) SetBuyPrice(buyPrice int) {
	i.buyPrice = buyPrice
}

func (i *Item) SellingPrice() int {
	return i.sellingPrice
}

func (i *Item) SetSellingPrice(sellingPrice int) {
	i.sellingPrice = sellingPrice
}

func (i *Item) Sprite() string {
	return i.sprite
}

func (i *Item) SetSprite(sprite string) {
	i.sprite = sprite
}
