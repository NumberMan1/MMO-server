package model

import (
	"github.com/NumberMan1/MMO-server/fight"
	"github.com/NumberMan1/common/logger"
)

type AttributesAssembly struct {
	Basic *fight.Attributes //基础属性（初始+成长）
	Equip *fight.Attributes //装备属性
	Buffs *fight.Attributes //Buff属性
	Final *fight.Attributes //最终属性
	owner IActor
}

func NewAttributesAssembly() *AttributesAssembly {
	return &AttributesAssembly{
		Basic: &fight.Attributes{},
		Equip: &fight.Attributes{},
		Buffs: &fight.Attributes{},
		Final: &fight.Attributes{},
	}
}

func (aa *AttributesAssembly) Init(actor IActor) {
	aa.owner = actor
	unitDefine := aa.owner.Define()

	aa.Equip.Reset()
	aa.Buffs.Reset()
	aa.Final.Reset()
	//基础属性
	aa.Basic = &fight.Attributes{
		Speed:     float32(unitDefine.Speed),
		HPMax:     unitDefine.HPMax,
		MPMax:     unitDefine.MPMax,
		AD:        unitDefine.AD,
		AP:        unitDefine.AP,
		DEF:       unitDefine.DEF,
		MDEF:      unitDefine.MDEF,
		CRI:       unitDefine.CRI,
		CRD:       unitDefine.CRD,
		STR:       unitDefine.STR,
		INT:       unitDefine.INT,
		AGI:       unitDefine.AGI,
		HitRate:   unitDefine.HitRate,
		DodgeRate: unitDefine.DodgeRate,
		HpRegen:   unitDefine.HpRegen,
		HpSteal:   unitDefine.HpSteal,
	}
	aa.Reload()
}

// Reload 重新加载
func (aa *AttributesAssembly) Reload() {
	define := aa.owner.Define()
	level := aa.owner.Info().Level

	//成长属性
	growth := &fight.Attributes{
		STR: define.GSTR * float32(level), // 力量成长
		INT: define.GINT * float32(level), // 智力成长
		AGI: define.GAGI * float32(level), // 敏捷成长
	}

	//todo 处理装备和buff
	//合并到最终属性
	aa.Final.Reset()
	aa.Final.Add(aa.Basic)
	aa.Final.Add(growth)
	aa.Final.Add(aa.Equip)
	aa.Final.Add(aa.Buffs)
	//附加属性
	extra := &fight.Attributes{
		HPMax: aa.Final.STR * 5,   //力量加生命上限
		AP:    aa.Final.INT * 1.5, //智力加法攻
	}
	aa.Final.Add(extra)
	logger.SLCInfo("最终属性：%v", aa.Final)
	//赋值与同步
	aa.owner.SetSpeed(int(aa.Final.Speed))
	aa.owner.SyncHpMax(aa.Final.HPMax)
	aa.owner.SyncMpMax(aa.Final.MPMax)
}
