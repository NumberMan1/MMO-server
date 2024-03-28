package model

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/define"
	"github.com/NumberMan1/MMO-server/inventory/item"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer/proto_helper"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
	"sync"
)

type Inventory struct {
	chr           *Character
	capacity      int
	itemMap       *sync.Map
	hasChanged    bool
	inventoryInfo *proto.InventoryInfo
}

func (i *Inventory) InventoryInfo() *proto.InventoryInfo {
	if i.inventoryInfo == nil {
		i.inventoryInfo = &proto.InventoryInfo{}
	}
	i.inventoryInfo.Capacity = int32(i.Capacity())
	//如果数据有变化
	if i.hasChanged {
		//重新拼装proto对象
		i.inventoryInfo.List = make([]*proto.ItemInfo, 0)
		i.itemMap.Range(func(key, value any) bool {
			i.inventoryInfo.List = append(i.inventoryInfo.List, value.(*item.Item).ItemInfo())
			return true
		})
	}
	return i.inventoryInfo
}

func (i *Inventory) setChr(chr *Character) {
	i.chr = chr
}

func (i *Inventory) setCapacity(capacity int) {
	i.capacity = capacity
}

func NewInventory(chr *Character) *Inventory {
	return &Inventory{chr: chr, itemMap: &sync.Map{}}
}

func (i *Inventory) Chr() *Character {
	return i.chr
}

func (i *Inventory) Capacity() int {
	return i.capacity
}

func (i *Inventory) Init(data []byte) {
	//默认背包
	if data == nil {
		i.setCapacity(10)
	} else { //数据还原
		msg, err := proto_helper.Parse([]byte("proto.InventoryInfo"), data)
		if err != nil {
			panic(err)
		}
		inv := msg.(*proto.InventoryInfo)
		logger.SLCInfo("数据还原: %v", inv)
		i.setCapacity(int(inv.Capacity))
		//创建物品
		for _, itemInfo := range inv.List {
			i.SetItem(int(itemInfo.Position), item.NewItemByItemId(int(itemInfo.ItemId), int(itemInfo.Amount), int(itemInfo.Position)))
		}
	}
}

// SetItem 设置插槽的物品，插槽索引从0开始
func (i *Inventory) SetItem(slotIndex int, item2 *item.Item) bool {
	//如果索引大于容量则停止设置
	if slotIndex >= i.Capacity() {
		return false
	}
	//标记数据发生变化
	i.hasChanged = true
	//清空当前插槽
	if item2 == nil {
		value, ok := i.itemMap.Load(slotIndex)
		if ok {
			i.itemMap.Delete(slotIndex)
			value.(*item.Item).Position = -1
		}
		return true
	}
	//设置插槽物品
	i.itemMap.Store(slotIndex, item2)
	item2.Position = slotIndex
	return true
}

// TrySlotItem 获取插槽物品
func (i *Inventory) TrySlotItem(slotIndex int) (*item.Item, bool) {
	value, ok := i.itemMap.Load(slotIndex)
	return value.(*item.Item), ok
}

// AddItem 物品加入库存
func (i *Inventory) AddItem(itemId, amount int) bool {
	//验证物品类型是否存在
	if _, ok := define.GetDataManagerInstance().Items[itemId]; !ok {
		logger.SLCInfo("物品id不存在: %v", itemId)
		return false
	}
	def := define.GetDataManagerInstance().Items[itemId]
	//检查剩余空间
	if i.calculateMaxRemainingQuantity(itemId) < amount {
		return false
	}
	//amount代表期望添加的数量，大于0则一直循环
	for amount > 0 {
		//查找id相同且未满的物品
		sameItem := i.findSameItemAndNotFull(itemId)
		if sameItem != nil {
			//本次可以处理的数量
			current := min(amount, sameItem.Capacity()-sameItem.Amount)
			sameItem.Amount += current
			amount -= current
		} else {
			//查找背包空闲的插槽索引，-1代表背包已满
			index := i.findEmptyIndex()
			if index > -1 {
				//本次可处理的数量
				current := min(amount, def.Capicity)
				i.SetItem(index, item.NewItem(def, current, index))
				amount -= current
			} else {
				logger.SLCDebug("没有空的物品槽")
				return false
			}
		}
	}
	return true
}

// Exchange 交换物品位置，该索引是插槽索引
func (i *Inventory) Exchange(originSlotIndex, targetSlotIndex int) bool {
	//交换id不能相同
	if originSlotIndex == targetSlotIndex {
		return false
	}
	//交换id不能小于零
	if originSlotIndex < 0 || targetSlotIndex < 0 {
		return false
	}
	//查找原始插槽物品
	if _, ok := i.itemMap.Load(originSlotIndex); !ok {
		return false
	}
	t1, _ := i.itemMap.Load(originSlotIndex)
	item1 := t1.(*item.Item)
	logger.SLCInfo("item1 = %v", item1.Name())
	//查找目标插槽物品，如果为空，直接放置
	var item2 *item.Item
	if t2, ok := i.itemMap.Load(targetSlotIndex); !ok {
		i.SetItem(originSlotIndex, nil)
		i.SetItem(targetSlotIndex, item1)
	} else {
		//如果物品类型相同
		item2 = t2.(*item.Item)
		if item1.Id() == item2.Id() {
			//可移动的数量
			num := min(item2.Capacity()-item2.Amount, item1.Amount)
			// 如果原始物品数量小于等于可移动数量，将原始物品全部移动到目标插槽)
			if item1.Amount <= num {
				item2.Amount += item1.Amount
				i.SetItem(originSlotIndex, nil)
			} else {
				// 否则，不移动物品只修改数量
				item1.Amount -= num
				item2.Amount += num
			}
		} else {
			//如果类型不同则交换位置
			i.SetItem(originSlotIndex, item2)
			i.SetItem(targetSlotIndex, item1)
		}
	}
	return true
}

// RemoveItem 移除指定数量的物品，无视插槽位置
func (i *Inventory) RemoveItem(itemId, amount int) int {
	removedAmount := 0
	for amount > 0 {
		findSameItem := i.findSameItem(itemId)
		if findSameItem == nil {
			break
		}
		//判断要移除的数量是否大于物品的当前数量
		currentAmount := min(amount, findSameItem.Amount)
		findSameItem.Amount -= currentAmount
		removedAmount += currentAmount
		amount -= currentAmount
		//清空物品槽
		if findSameItem.Amount == 0 {
			i.SetItem(findSameItem.Position, nil)
		}
	}
	return removedAmount
}

// Discard 丢弃指定插槽的物品，返回实际丢弃的数量
func (i *Inventory) Discard(slotIndex, amount int) int {
	if amount < 1 {
		return 0
	}
	if _, ok := i.itemMap.Load(slotIndex); !ok {
		return 0
	}
	t1, _ := i.itemMap.Load(slotIndex)
	item1 := t1.(*item.Item)
	//只丢弃一部分
	if amount < item1.Amount {
		item1.Amount -= amount
		newItem := item.NewItemByItemId(item1.Id(), amount, 0)
		CreateItemEntity(i.Chr().Space(), newItem, i.Chr().Position(), vector3.Zero3())
		return amount
	}
	//全额丢弃
	i.SetItem(slotIndex, nil)
	CreateItemEntity(i.Chr().Space(), item1, i.Chr().Position(), vector3.Zero3())
	return item1.Amount
}

// calculateMaxRemainingQuantity 计算背包里还能放多少个这样的物品
func (i *Inventory) calculateMaxRemainingQuantity(itemId int) int {
	//检查物品类型是否存在
	def, ok := define.GetDataManagerInstance().Items[itemId]
	if !ok {
		logger.SLCInfo("物品id不存在: %v", itemId)
		return 0
	}
	//记录可用数量
	quantity := 0
	//遍历全部插槽
	for index := 0; index < i.Capacity(); index++ {
		//如果插槽有物品
		if t1, ok := i.itemMap.Load(index); ok {
			item1 := t1.(*item.Item)
			//如果物品类型相同
			if item1.Id() == itemId {
				quantity += item1.Capacity() - item1.Amount
			}
		} else {
			quantity += def.Capicity
		}
	}
	logger.SLCDebug("Inventory：Entity[%v] 物品[%v]还能放入[%v]个", i.Chr().EntityId(), def.Name, quantity)
	return quantity
}

// findSameItem 查找ID相同的物品
func (i *Inventory) findSameItem(itemId int) *item.Item {
	var item1 *item.Item
	i.itemMap.Range(func(key, value any) bool {
		if value != nil {
			if item2 := value.(*item.Item); item2.Id() == itemId {
				item1 = item2
				return false
			}
		}
		return true
	})
	return item1
}

// findSameItemAndNotFull 查找ID相同且未满的物品
func (i *Inventory) findSameItemAndNotFull(itemId int) *item.Item {
	var item1 *item.Item
	i.itemMap.Range(func(key, value any) bool {
		if value != nil {
			if item2 := value.(*item.Item); item2.Id() == itemId && item2.Amount < item2.Capacity() {
				item1 = item2
				return false
			}
		}
		return true
	})
	return item1
}

// findEmptyIndex 查找空的插槽位置
func (i *Inventory) findEmptyIndex() int {
	for index := 0; index < i.Capacity(); index++ {
		if _, ok := i.itemMap.Load(index); !ok {
			return index
		}
	}
	return -1
}
