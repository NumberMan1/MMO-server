package define

// 1. 每个 Sheet 形成一个 Struct 定义, Sheet 的名称作为 Struct 的名称
// 2. 表格约定：第一行是变量名称，第二行是变量类型

type SkillDefine struct {
	ID           int       `json:"ID"`           // 编号
	TID          int       `json:"TID"`          // 单位类型
	Code         int       `json:"Code"`         // 技能码
	Name         string    `json:"Name"`         // 技能名称
	Description  string    `json:"Description"`  // 技能描述
	Level        int       `json:"Level"`        // 技能等级
	MaxLevel     int       `json:"MaxLevel"`     // 技能上限
	Type         string    `json:"Type"`         // 类别
	Icon         string    `json:"Icon"`         // 技能图标
	TargetType   string    `json:"TargetType"`   // 目标类型
	CD           float32   `json:"CD"`           // 冷却时间
	SpellRange   int       `json:"SpellRange"`   // 施法距离
	CastTime     float32   `json:"CastTime"`     // 施法前摇
	Cost         int       `json:"Cost"`         // 魔法消耗
	AnimName     string    `json:"AnimName"`     // 施法动作
	ReqLevel     int       `json:"ReqLevel"`     // 等级要求
	Missile      string    `json:"Missile"`      // 投射物
	MissileSpeed int       `json:"MissileSpeed"` // 投射速度
	HitArt       string    `json:"HitArt"`       // 击中效果
	Area         int       `json:"Area"`         // 影响区域
	HitTime      []float32 `json:"HitTime"`      // 命中时间
	BUFF         []int     `json:"BUFF"`         // 附加效果
	AD           float32   `json:"AD"`           // 物理攻击
	AP           float32   `json:"AP"`           // 法术攻击
	ADC          float32   `json:"ADC"`          // 物攻加成
	APC          float32   `json:"APC"`          // 法攻加成
}
