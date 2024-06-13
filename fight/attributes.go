package fight

import "encoding/json"

// Attributes 定义实体的各种属性
type Attributes struct {
	Speed     float32 // 速度
	HPMax     float32 // 最大生命值
	MPMax     float32 // 最大魔法值
	AD        float32 // 物理攻击力
	AP        float32 // 魔法攻击力
	DEF       float32 // 物理防御力
	MDEF      float32 // 魔法防御力
	CRI       float32 // 暴击率
	CRD       float32 // 暴击伤害
	STR       float32 // 力量
	INT       float32 // 智力
	AGI       float32 // 敏捷
	HitRate   float32 // 命中率
	DodgeRate float32 // 闪避率
	HpRegen   float32 // 生命恢复
	HpSteal   float32 // 伤害吸血%
}

// Add 方法增加来自另一个 Attributes 对象的属性
func (a *Attributes) Add(data *Attributes) {
	a.Speed += data.Speed
	a.HPMax += data.HPMax
	a.MPMax += data.MPMax
	a.AD += data.AD
	a.AP += data.AP
	a.DEF += data.DEF
	a.MDEF += data.MDEF
	a.CRI += data.CRI
	a.CRD += data.CRD
	a.STR += data.STR
	a.INT += data.INT
	a.AGI += data.AGI
	a.HitRate += data.HitRate
	a.DodgeRate += data.DodgeRate
	a.HpRegen += data.HpRegen
	a.HpSteal += data.HpSteal
}

// Sub 方法减少来自另一个 Attributes 对象的属性
func (a *Attributes) Sub(data *Attributes) {
	a.Speed -= data.Speed
	a.HPMax -= data.HPMax
	a.MPMax -= data.MPMax
	a.AD -= data.AD
	a.AP -= data.AP
	a.DEF -= data.DEF
	a.MDEF -= data.MDEF
	a.CRI -= data.CRI
	a.CRD -= data.CRD
	a.STR -= data.STR
	a.INT -= data.INT
	a.AGI -= data.AGI
	a.HitRate -= data.HitRate
	a.DodgeRate -= data.DodgeRate
	a.HpRegen -= data.HpRegen
	a.HpSteal -= data.HpSteal
}

// Reset 方法将所有属性重置为零
func (a *Attributes) Reset() {
	a.Speed = 0
	a.HPMax = 0
	a.MPMax = 0
	a.AD = 0
	a.AP = 0
	a.DEF = 0
	a.MDEF = 0
	a.CRI = 0
	a.CRD = 0
	a.STR = 0
	a.INT = 0
	a.AGI = 0
	a.HitRate = 0
	a.DodgeRate = 0
	a.HpRegen = 0
	a.HpSteal = 0
}

// ToString 方法提供 Attributes 对象的 JSON 表示
func (a *Attributes) ToString() string {
	data, _ := json.Marshal(a)
	return string(data)
}
