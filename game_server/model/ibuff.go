package model

import (
	"github.com/NumberMan1/MMO-server/game_common/protocol/gen/proto"
	"github.com/NumberMan1/MMO-server/game_server/config/define"
)

// Buff 接口定义了 Buff 的基本行为
type Buff interface {
	// Owner 获取此 Buff 的持有者
	Owner() IActor
	// Provider 获取此 Buff 的提供者
	Provider() IActor
	// Def 获取 Buff 的定义
	Def() *define.BuffDefine
	// BID 获取 Buff 的 ID
	BID() int
	// MaxDuration 获取 Buff 的初始持续时间
	MaxDuration() float32
	// TimeScale 获取 Buff 的时间流失速度
	TimeScale() float32
	// SetTimeScale 设置 Buff 的时间流失速度
	SetTimeScale(value float32)
	// MaxLevel 获取 Buff 的最大堆叠层数
	MaxLevel() int
	// BuffType 获取 Buff 的类型
	BuffType() proto.BuffType
	// BuffConflict 获取 Buff 的冲突处理方式
	BuffConflict() proto.BuffConflict
	// Dispellable 获取 Buff 是否可驱散
	Dispellable() bool
	// Name 获取 Buff 的名称
	Name() string
	// Description 获取 Buff 的介绍文本
	Description() string
	// IconPath 获取 Buff 的图标资源路径
	IconPath() string
	// Demotion 获取每次 Buff 持续时间结束时降低的等级
	Demotion() int
	// CurrentLevel 获取 Buff 的当前等级
	CurrentLevel() int
	// SetCurrentLevel 设置 Buff 的当前等级
	SetCurrentLevel(value int)
	// ResidualDuration 获取 Buff 的当前剩余时间
	ResidualDuration() float32
	// SetResidualDuration 设置 Buff 的当前剩余时间
	SetResidualDuration(value float32)
	// OnGet 当持有者获得此 Buff 时触发
	OnGet()
	// OnLost 当持有者失去此 Buff 时触发
	OnLost()
	// OnUpdate 每帧更新，由 BuffManager 调用
	OnUpdate(delta float64)
	// OnLevelChange 当等级改变时调用
	OnLevelChange(change int)
	// Init 初始化 Buff
	Init(owner, provider IActor) error
	// GetBuffDefine 获取 Buff 的定义
	GetBuffDefine() *define.BuffDefine
}
