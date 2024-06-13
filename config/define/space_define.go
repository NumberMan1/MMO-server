package define

// 1. 每个 Sheet 形成一个 Struct 定义, Sheet 的名称作为 Struct 的名称
// 2. 表格约定：第一行是变量名称，第二行是变量类型

type SpaceDefine struct {
	SID      int    `json:"SID" bson:"sid"`           // 场景编号
	Name     string `json:"Name" bson:"name"`         // 名称
	Resource string `json:"Resource" bson:"resource"` // 资源
	Kind     string `json:"Kind" bson:"kind"`         // 类型
	AllowPK  int    `json:"AllowPK" bson:"allow_pk"`  //允许PK（1允许，0不允许）
}

func (s *SpaceDefine) GetId() int {
	return s.SID
}
