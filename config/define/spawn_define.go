package define

import (
	"encoding/json"
	"strconv"
	"strings"
)

// 1. 每个 Sheet 形成一个 Struct 定义, Sheet 的名称作为 Struct 的名称
// 2. 表格约定：第一行是变量名称，第二行是变量类型

type SpawnDefine struct {
	ID      int `json:"ID" bson:"id"`            // ID
	SpaceId int `json:"SpaceId" bson:"space_id"` // 场景ID
	//Pos     string `json:"Pos" bson:"pos"`          // 刷怪位置
	//Dir     string `json:"Dir" bson:"dir"`          // 刷怪方向
	Pos    []int `json:"Pos" bson:"pos"`       // 刷怪位置
	Dir    []int `json:"Dir" bson:"dir"`       // 刷怪方向
	TID    int   `json:"Job" bson:"tid"`       // 单位类型
	Level  int   `json:"Level" bson:"level"`   // 单位等级
	Period int   `json:"Period" bson:"period"` // 刷怪周期（秒）
	Count  int   `json:"Count" bson:"count"`   // 刷怪数量
	Range  int   `json:"Range" bson:"range"`   // 随机距离
}

// UnmarshalJSON Custom unmarshalling to handle the specific requirements
func (s *SpawnDefine) UnmarshalJSON(data []byte) error {
	type Alias SpawnDefine
	aux := &struct {
		Pos string `json:"Pos"`
		Dir string `json:"Dir"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Convert Pos from string to []int
	if aux.Pos != "" {
		var pos []int
		posStr := strings.Trim(aux.Pos, "[]")
		for _, val := range strings.Split(posStr, ",") {
			num, err := strconv.Atoi(strings.TrimSpace(val))
			if err != nil {
				return err
			}
			pos = append(pos, num)
		}
		s.Pos = pos
	}

	// Convert Dir from string to []int
	if aux.Dir != "" {
		var dir []int
		dirStr := strings.Trim(aux.Dir, "[]")
		for _, val := range strings.Split(dirStr, ",") {
			num, err := strconv.Atoi(strings.TrimSpace(val))
			if err != nil {
				return err
			}
			dir = append(dir, num)
		}
		s.Dir = dir
	}
	return nil
}

func (s *SpawnDefine) GetId() int {
	return s.ID
}
