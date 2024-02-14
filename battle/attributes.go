package battle

import "github.com/NumberMan1/MMO-server/define"

type Attributes struct {
	Basic AttributeData //基础属性（初始+成长）
	Equip AttributeData //装备属性
	Buffs AttributeData //Buff属性
	Final AttributeData //最终属性
}

func (a *Attributes) Init(define define.UnitDefine, level int) {
	a.Basic = AttributeData{}
	a.Equip = AttributeData{}
	a.Buffs = AttributeData{}
	a.Final = AttributeData{}
	//初始化属性
	initial := AttributeData{
		Speed: float32(define.Speed),
		HPMax: define.HPMax,
		MPMax: define.MPMax,
		AD:    define.AD,
		AP:    define.AP,
		DEF:   define.DEF,
		MDEF:  define.MDEF,
		CRI:   define.CRI,
		CRD:   define.CRD,
		STR:   define.STR,
		INT:   define.INT,
		AGI:   define.AGI,
	}
	//成长属性
	growth := AttributeData{
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
	extra := AttributeData{
		HPMax: a.Final.STR * 5,
		AP:    a.Final.INT * 1.5,
	}
	a.Final.Add(extra)
}
