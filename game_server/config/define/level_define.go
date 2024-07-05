package define

// LevelDefine 定义了一个等级的属性结构体
type LevelDefine struct {
	Level    int   `json:"Level" bson:"level"`        // 等级
	ExpLimit int64 `json:"ExpLimit" bson:"exp_limit"` // 升级需要经验
	BaseExp  int   `json:"BaseExp" bson:"base_exp"`   // 标准单位经验
	PlanTime int   `json:"PlanTime" bson:"plan_time"` // 经验时间
}

func (l *LevelDefine) GetId() int {
	return l.Level
}
