package define

import "encoding/json"

// BuffDefine 定义了一个Buff的属性结构体
type BuffDefine struct {
	BID          int     `json:"BID" bson:"bid"`                    // 编号
	Name         string  `json:"Name" bson:"name"`                  // 名称
	Description  string  `json:"Description" bson:"description"`    // 介绍
	IconPath     string  `json:"IconPath" bson:"icon_path"`         // 图标路径
	MaxDuration  float32 `json:"MaxDuration" bson:"max_duration"`   // 持续时间
	MaxLevel     int     `json:"MaxLevel" bson:"max_level"`         // 堆叠上限
	BuffType     string  `json:"BuffType" bson:"buff_type"`         // 种类
	BuffConflict string  `json:"BuffConflict" bson:"buff_conflict"` // 叠加方式
	Dispellable  bool    `json:"Dispellable" bson:"dispellable"`    // 是否可驱散
	Demotion     int     `json:"Demotion" bson:"demotion"`          // 降级
	TimeScale    float32 `json:"TimeScale" bson:"time_scale"`       // 时间速率
}

// UnmarshalJSON Custom unmarshalling to handle the specific requirements
func (b *BuffDefine) UnmarshalJSON(data []byte) error {
	type Alias BuffDefine
	aux := &struct {
		Dispellable int `json:"Dispellable"`
		*Alias
	}{
		Alias: (*Alias)(b),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Convert Dispellable from int to bool
	b.Dispellable = aux.Dispellable == 1
	return nil
}

func (b *BuffDefine) GetId() int {
	return b.BID
}
