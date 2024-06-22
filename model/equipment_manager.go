package model

import (
	"github.com/NumberMan1/MMO-server/model/inventory_system/item"
	proto2 "github.com/NumberMan1/MMO-server/protocol/gen/proto"
	"github.com/NumberMan1/common/global/variable"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"sync"
)

// EquipsManager 装备管理器，每个主角都有
type EquipsManager struct {
	Chr *Character
	//装备数据 proto.EquipsType:*item.Equipment
	dict          sync.Map
	inventoryInfo *proto2.InventoryInfo
}

func NewEquipsManager(owner *Character) *EquipsManager {
	return &EquipsManager{Chr: owner}
}

// Attach 穿戴装备
func (em *EquipsManager) Attach(equip *item.Equipment) bool {
	if !em.Detach(equip.EquipType()) {
		return false
	}
	em.dict.Store(equip.EquipType(), equip)
	variable.Log.Info("EquipsManager:Attach", zap.Any("InventoryInfo", em.InventoryInfo()))
	//添加属性
	em.Chr.Attr().Equip.Add(equip.Attr())
	em.Chr.Attr().Reload()
	//广播
	em.SendEquipsToClient()
	return true
}

// Detach 卸下装备
func (em *EquipsManager) Detach(equipsType proto2.EquipsType) bool {
	if value, ok := em.dict.Load(equipsType); ok {
		equip := value.(*item.Equipment)
		index := em.Chr.Knapsack.FindEmptyIndex()
		if index == -1 {
			return false
		}
		em.dict.Delete(equipsType)
		em.Chr.Knapsack.SetItem(index, equip)
		em.Chr.Attr().Equip.Sub(equip.Attr())
		em.Chr.Attr().Reload()
		em.SendEquipsToClient()
	}
	variable.Log.Info("EquipsManager:Detach", zap.Any("InventoryInfo", em.InventoryInfo()))
	return true
}

func (em *EquipsManager) InventoryInfo() *proto2.InventoryInfo {
	if em.inventoryInfo == nil {
		em.inventoryInfo = &proto2.InventoryInfo{List: make([]*proto2.ItemInfo, 0)}
	}
	em.inventoryInfo.Capacity = 0
	em.inventoryInfo.List = make([]*proto2.ItemInfo, 0)
	em.dict.Range(func(key, value any) bool {
		equip := value.(*item.Equipment)
		em.inventoryInfo.List = append(em.inventoryInfo.List, equip.ItemInfo())
		return true
	})
	return em.inventoryInfo
}

// Init 初始化装备
func (em *EquipsManager) Init(bytes []byte) {
	if bytes == nil || len(bytes) == 0 {
		return
	}
	inv := &proto2.InventoryInfo{}
	err := proto.Unmarshal(bytes, inv)
	if err != nil {
		variable.Log.Error("EquipsManager:Init", zap.Error(err))
		return
	}
	variable.Log.Info("EquipsManager:Init", zap.Any("InventoryInfo", inv))
	for _, itemInfo := range inv.List {
		//还原装备对象
		equip := item.NewEquipmentByInfo(itemInfo)
		//穿戴装备
		em.Attach(equip)
	}
	//更新角色装备信息
	em.Chr.Info().EquipsList = make([]*proto2.ItemInfo, 0)
	for _, itemInfo := range em.InventoryInfo().List {
		em.Chr.Info().EquipsList = append(em.Chr.Info().EquipsList, itemInfo)
	}
}

// SendEquipsToClient 发送装备给客户端
func (em *EquipsManager) SendEquipsToClient() {
	//更新角色装备信息
	em.Chr.Info().EquipsList = make([]*proto2.ItemInfo, 0)
	for _, itemInfo := range em.InventoryInfo().List {
		em.Chr.Info().EquipsList = append(em.Chr.Info().EquipsList, itemInfo)
	}
	//广播
	resp := &proto2.EquipsResponse{EntityId: int32(em.Chr.EntityId()), EquipsList: em.InventoryInfo().List}
	if em.Chr != nil {
		if em.Chr.Space() != nil {
			em.Chr.Space().Broadcast(resp)
		}
	}
}
