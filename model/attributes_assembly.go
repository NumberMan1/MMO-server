package model

import (
	"github.com/NumberMan1/MMO-server/fight"
	"github.com/NumberMan1/common/logger"
)

type AttributesAssembly struct {
	Basic fight.Attributes //基础属性（初始+成长）
	Equip fight.Attributes //装备属性
	Buffs fight.Attributes //Buff属性
	Final fight.Attributes //最终属性
}

func (a *AttributesAssembly) Init(actor IActor) {
	a.Basic = fight.Attributes{}
	a.Equip = fight.Attributes{}
	a.Buffs = fight.Attributes{}
	a.Final = fight.Attributes{}
	define := actor.Define()
	level := actor.Info().Level
	//初始化属性
	initial := fight.Attributes{
		Speed:     float32(define.Speed),
		HPMax:     define.HPMax,
		MPMax:     define.MPMax,
		AD:        define.AD,
		AP:        define.AP,
		DEF:       define.DEF,
		MDEF:      define.MDEF,
		CRI:       define.CRI,
		CRD:       define.CRD,
		STR:       define.STR,
		INT:       define.INT,
		AGI:       define.AGI,
		HitRate:   define.HitRate,
		DodgeRate: define.DodgeRate,
		HpRegen:   define.HpRegen,
		HpSteal:   define.HpSteal,
	}
	//成长属性
	growth := fight.Attributes{
		STR: define.GSTR * float32(level), // 力量成长
		INT: define.GINT * float32(level), // 智力成长
		AGI: define.GAGI * float32(level), // 敏捷成长
	}
	//基础属性（初始+成长）
	a.Basic.Add(initial)
	a.Basic.Add(growth)
	//todo 处理装备和buff
	//合并到最终属性"
	a.Final.Add(a.Basic)
	a.Final.Add(a.Equip)
	a.Final.Add(a.Buffs)
	//附加属性
	extra := fight.Attributes{
		HPMax: a.Final.STR * 5,
		AP:    a.Final.INT * 1.5,
	}
	a.Final.Add(extra)
	logger.SLCInfo("最终属性：%v", a.Final)
}
