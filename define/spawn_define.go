package define

// 1. 每个 Sheet 形成一个 Struct 定义, Sheet 的名称作为 Struct 的名称
// 2. 表格约定：第一行是变量名称，第二行是变量类型

type SpawnDefine struct {
	ID      int    `json:"ID"`      // ID
	SpaceId int    `json:"SpaceId"` // 场景ID
	Pos     string `json:"Pos"`     // 刷怪位置
	Dir     string `json:"Dir"`     // 刷怪方向
	TID     int    `json:"TID"`     // 单位类型
	Level   int    `json:"Level"`   // 单位等级
	Period  int    `json:"Period"`  // 刷怪周期（秒）

}
