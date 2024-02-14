package battle

type AttributeData struct {
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
}

// Add 增加属性
func (ad *AttributeData) Add(data AttributeData) {
	ad.Speed += data.Speed
	ad.HPMax += data.HPMax
	ad.MPMax += data.MPMax
	ad.AD += data.AD
	ad.AP += data.AP
	ad.DEF += data.DEF
	ad.MDEF += data.MDEF
	ad.CRI += data.CRI
	ad.CRD += data.CRD
	ad.STR += data.STR
	ad.INT += data.INT
	ad.AGI += data.AGI
	ad.HitRate += data.HitRate
	ad.DodgeRate += data.DodgeRate
}

// Sub 减少属性
func (ad *AttributeData) Sub(data AttributeData) {
	ad.Speed -= data.Speed
	ad.HPMax -= data.HPMax
	ad.MPMax -= data.MPMax
	ad.AD -= data.AD
	ad.AP -= data.AP
	ad.DEF -= data.DEF
	ad.MDEF -= data.MDEF
	ad.CRI -= data.CRI
	ad.CRD -= data.CRD
	ad.STR -= data.STR
	ad.INT -= data.INT
	ad.AGI -= data.AGI
	ad.HitRate -= data.HitRate
	ad.DodgeRate -= data.DodgeRate
}

// Reset 重置属性
func (ad *AttributeData) Reset() {
	ad.Speed = 0
	ad.HPMax = 0
	ad.MPMax = 0
	ad.AD = 0
	ad.AP = 0
	ad.DEF = 0
	ad.MDEF = 0
	ad.CRI = 0
	ad.CRD = 0
	ad.STR = 0
	ad.INT = 0
	ad.AGI = 0
	ad.HitRate = 0
	ad.DodgeRate = 0
}
