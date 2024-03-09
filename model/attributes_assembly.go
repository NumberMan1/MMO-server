package model

import (
	"github.com/NumberMan1/MMO-server/fight"
	"github.com/NumberMan1/common/logger"
)

type AttributesAssembly struct {
	initial *fight.Attributes //初始属性(来自于mongo的属性)
	growth  *fight.Attributes //成长属性
	//Basic *fight.Attributes //基础属性（初始+成长）
	equip *fight.Attributes //装备属性
	buffs *fight.Attributes //Buff属性
	Final *fight.Attributes //最终属性
	actor IActor
}

func NewAttributesAssembly() *AttributesAssembly {
	return &AttributesAssembly{
		initial: &fight.Attributes{},
		growth:  &fight.Attributes{},
		equip:   &fight.Attributes{},
		buffs:   &fight.Attributes{},
		Final:   &fight.Attributes{},
	}
}

func (aa *AttributesAssembly) Init(actor IActor) {
	aa.actor = actor

}

// Reload 重新加载
func (aa *AttributesAssembly) Reload() {
	aa.growth.Reset()
	aa.equip.Reset()
	aa.buffs.Reset()
	aa.Final.Reset()

	define := aa.actor.Define()
	level := aa.actor.Info().Level
	//初始化属性
	aa.initial.Speed = float32(define.Speed)
	aa.initial.HPMax = define.HPMax
	aa.initial.MPMax = define.MPMax
	aa.initial.AD = define.AD
	aa.initial.AP = define.AP
	aa.initial.DEF = define.DEF
	aa.initial.MDEF = define.MDEF
	aa.initial.CRI = define.CRI
	aa.initial.CRD = define.CRD
	aa.initial.STR = define.STR
	aa.initial.INT = define.INT
	aa.initial.AGI = define.AGI
	aa.initial.HitRate = define.HitRate
	aa.initial.DodgeRate = define.DodgeRate
	aa.initial.HpRegen = define.HpRegen
	aa.initial.HpSteal = define.HpSteal

	//成长属性
	aa.growth.STR = define.GSTR * float32(level) // 力量成长
	aa.growth.INT = define.GINT * float32(level) // 智力成长
	aa.growth.AGI = define.GAGI * float32(level) // 敏捷成长

	////基础属性（初始+成长）
	//a.Basic.Add(initial)
	//a.Basic.Add(growth)

	//todo 处理装备和buff
	//合并到最终属性"
	aa.Final.Add(aa.initial)
	aa.Final.Add(aa.growth)
	aa.Final.Add(aa.equip)
	aa.Final.Add(aa.buffs)
	//附加属性
	extra := &fight.Attributes{
		HPMax: aa.Final.STR * 5,   //力量加生命上限
		AP:    aa.Final.INT * 1.5, //智力加法攻
	}
	aa.Final.Add(extra)
	logger.SLCInfo("最终属性：%v", aa.Final)
	//赋值与同步
	aa.actor.SetSpeed(int(aa.Final.Speed))
	aa.actor.Info().Hpmax = aa.Final.HPMax
	aa.actor.Info().Mpmax = aa.Final.MPMax
	aa.actor.OnHpMaxChanged(aa.Final.HPMax)
	aa.actor.OnMpMaxChanged(aa.Final.MPMax)
}
