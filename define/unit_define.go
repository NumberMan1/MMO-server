package define

// 1. 每个 Sheet 形成一个 Struct 定义, Sheet 的名称作为 Struct 的名称
// 2. 表格约定：第一行是变量名称，第二行是变量类型

type UnitDefine struct {
	TID       int     `json:"TID" bson:"tid"`              // 单位类型
	Name      string  `json:"Name" bson:"name"`            // 名称
	Resource  string  `json:"Resource" bson:"resource"`    // 模型资源
	Kind      string  `json:"Kind" bson:"kind"`            // 类别
	Decs      string  `json:"Decs" bson:"decs"`            // 介绍
	Speed     int     `json:"Speed" bson:"speed"`          // 移动速度
	HPMax     float32 `json:"HPMax" bson:"hp_max"`         // 生命上限
	MPMax     float32 `json:"MPMax" bson:"mp_max"`         // 法力上限
	InitLevel int     `json:"InitLevel" bson:"init_level"` // 初始等级
	Energy    int     `json:"Energy" bson:"energy"`        //活力
	AD        float32 `json:"AD" bson:"ad"`                // 物攻
	AP        float32 `json:"AP" bson:"ap"`                // 魔攻
	DEF       float32 `json:"DEF" bson:"def"`              // 物防
	MDEF      float32 `json:"MDEF" bson:"mdef"`            // 魔防
	CRI       float32 `json:"CRI" bson:"cri"`              // 暴击率
	CRD       float32 `json:"CRD" bson:"crd"`              // 暴击伤害
	HitRate   float32 `json:"HitRate" bson:"hit_rate"`     // 命中率
	DodgeRate float32 `json:"DodgeRate" bson:"dodge_rate"` // 闪避率
	STR       float32 `json:"STR" bson:"str"`              // 力量
	INT       float32 `json:"INT" bson:"int"`              // 智力
	AGI       float32 `json:"AGI" bson:"agi"`              // 敏捷
	GSTR      float32 `json:"GSTR" bson:"gstr"`            // 力量成长
	GINT      float32 `json:"GINT" bson:"gint"`            // 智力成长
	GAGI      float32 `json:"GAGI" bson:"gagi"`            // 敏捷成长
	AI        string  `json:"AI" bson:"ai"`                // AI名称
}

func (u *UnitDefine) GetId() int {
	return u.TID
}
