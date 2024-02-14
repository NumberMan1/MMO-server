package define

// 1. 每个 Sheet 形成一个 Struct 定义, Sheet 的名称作为 Struct 的名称
// 2. 表格约定：第一行是变量名称，第二行是变量类型

type UnitDefine struct {
	TID       int     `json:"TID"`       // 单位类型
	Name      string  `json:"Name"`      // 名称
	Resource  string  `json:"Resource"`  // 模型资源
	Type      string  `json:"Type"`      // 类别
	Decs      string  `json:"Decs"`      // 介绍
	Speed     int     `json:"Speed"`     // 移动速度
	HPMax     float32 `json:"HPMax"`     // 生命上限
	MPMax     float32 `json:"MPMax"`     // 法力上限
	InitLevel int     `json:"InitLevel"` // 初始等级
	AD        float32 `json:"AD"`        // 物攻
	AP        float32 `json:"AP"`        // 魔攻
	DEF       float32 `json:"DEF"`       // 物防
	MDEF      float32 `json:"MDEF"`      // 魔防
	CRI       float32 `json:"CRI"`       // 暴击率
	CRD       float32 `json:"CRD"`       // 暴击伤害
	HitRate   float32 `json:"HitRate"`   // 命中率
	DodgeRate float32 `json:"DodgeRate"` // 闪避率
	STR       float32 `json:"STR"`       // 力量
	INT       float32 `json:"INT"`       // 智力
	AGI       float32 `json:"AGI"`       // 敏捷
	GSTR      float32 `json:"GSTR"`      // 力量成长
	GINT      float32 `json:"GINT"`      // 智力成长
	GAGI      float32 `json:"GAGI"`      // 敏捷成长
	AI        string  `json:"AI"`        // AI名称
}
