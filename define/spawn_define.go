package define

// 1. 每个 Sheet 形成一个 Struct 定义, Sheet 的名称作为 Struct 的名称
// 2. 表格约定：第一行是变量名称，第二行是变量类型

type SpawnDefine struct {
	ID      int    `json:"ID" bson:"id"`            // ID
	SpaceId int    `json:"SpaceId" bson:"space_id"` // 场景ID
	Pos     string `json:"Pos" bson:"pos"`          // 刷怪位置
	Dir     string `json:"Dir" bson:"dir"`          // 刷怪方向
	TID     int    `json:"TID" bson:"tid"`          // 单位类型
	Level   int    `json:"Level" bson:"level"`      // 单位等级
	Period  int    `json:"Period" bson:"period"`    // 刷怪周期（秒）

}

func (s *SpawnDefine) GetId() int {
	return s.ID
}
