package define

import "encoding/json"

// 1. 每个 Sheet 形成一个 Struct 定义, Sheet 的名称作为 Struct 的名称
// 2. 表格约定：第一行是变量名称，第二行是变量类型

type SkillDefine struct {
	ID           int       `json:"ID" bson:"id"`                      // 编号
	Job          int       `json:"Job" bson:"job"`                    // 单位类型
	Code         int       `json:"Code" bson:"code"`                  // 技能码
	Name         string    `json:"Name" bson:"name"`                  // 技能名称
	Description  string    `json:"Description" bson:"description"`    // 技能描述
	Level        int       `json:"Level" bson:"level"`                // 技能等级
	MaxLevel     int       `json:"MaxLevel" bson:"max_level"`         // 技能上限
	Type         string    `json:"Type" bson:"type"`                  // 类别
	Icon         string    `json:"Icon" bson:"icon"`                  // 技能图标
	TargetType   string    `json:"TargetType" bson:"target_type"`     // 目标类型
	CD           float32   `json:"CD" bson:"cd"`                      // 冷却时间
	SpellRange   int       `json:"SpellRange" bson:"spell_range"`     // 施法距离
	IntonateTime float32   `json:"IntonateTime" bson:"IntonateTime"`  // 施法前摇
	Cost         int       `json:"Cost" bson:"cost"`                  // 魔法消耗
	Anim1        string    `json:"Anim1" bson:"anim_1"`               //前摇动作
	Anim2        string    `json:"AnimName" bson:"anim_2"`            // 激活动作
	ReqLevel     int       `json:"ReqLevel" bson:"req_level"`         // 等级要求
	IsMissile    bool      `json:"IsMissile" bson:"is_missile"`       //是否投射物
	Missile      string    `json:"Missile" bson:"missile"`            // 投射物
	MissileSpeed int       `json:"MissileSpeed" bson:"missile_speed"` // 投射速度
	HitArt       string    `json:"HitArt" bson:"hit_art"`             // 击中效果
	Area         int       `json:"Area" bson:"area"`                  // 影响区域
	Duration     float32   `json:"Duration" bson:"duration"`          //持续时间
	Interval     float32   `json:"Interval" bson:"interval"`          //伤害间隔
	HitDelay     []float32 `json:"HitDelay" bson:"hit_delay"`         // 命中时间
	BUFF         []int     `json:"BUFF" bson:"buff"`                  // 附加效果
	AD           float32   `json:"AD" bson:"ad"`                      // 物理攻击
	AP           float32   `json:"AP" bson:"ap"`                      // 法术攻击
	ADC          float32   `json:"ADC" bson:"adc"`                    // 物攻加成
	APC          float32   `json:"APC" bson:"apc"`                    // 法攻加成
}

// UnmarshalJSON Custom unmarshalling to handle the specific requirements
func (s *SkillDefine) UnmarshalJSON(data []byte) error {
	type Alias SkillDefine
	aux := &struct {
		IsMissile int    `json:"IsMissile"`
		HitDelay  string `json:"HitDelay"`
		BUFF      string `json:"BUFF"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Convert IsMissile from int to bool
	s.IsMissile = aux.IsMissile == 1

	// Convert HitDelay from string to []float32
	if aux.HitDelay != "" {
		var hitDelay []float32
		if err := json.Unmarshal([]byte(aux.HitDelay), &hitDelay); err != nil {
			return err
		}
		s.HitDelay = hitDelay
	}

	// Convert BUFF from string to []int
	if aux.BUFF != "" {
		var buff []int
		if err := json.Unmarshal([]byte(aux.BUFF), &buff); err != nil {
			return err
		}
		s.BUFF = buff
	}
	return nil
}

func (s *SkillDefine) GetId() int {
	return s.ID
}
