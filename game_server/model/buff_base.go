package model

import (
	"errors"
	"github.com/NumberMan1/MMO-server/game_common/protocol/gen/proto"
	"github.com/NumberMan1/MMO-server/game_server/config/define"
	"math"
)

// BuffBase 是 Buff 的基础类
type BuffBase struct {
	//Buff的初始持续时间
	maxDuration float32
	//Buff的时间流失速度
	timeScale float32
	//buff的最大堆叠层数，最小为1，最大为2147483647
	maxLevel int
	// Buff的类型，分为正面、负面、中立三种
	buffType proto.BuffType
	// 当两个不同单位向同一个单位施加同一个buff时的冲突处理
	// 分为三种：
	// combine,合并为一个buff，叠层（提高等级）
	// separate,独立存在
	// cover, 覆盖，后者覆盖前者
	buffConflict proto.BuffConflict
	// 可否被驱散
	dispellable bool
	// Buff对外显示的名称
	name string
	// Buff的介绍文本
	description string
	// 图标资源的路径
	demotion int
	// 每次Buff持续时间结束时降低的等级，一般降低1级或者降低为0级。
	iconPath string
	// Buff的当前等级
	currentLevel int
	// Buff的当前剩余时间
	residualDuration float32
	initialized      bool

	//此buff的持有者
	owner IActor
	//此Buff提供者
	provider IActor
	def      *define.BuffDefine
}

// Owner 此buff的持有者
func (b *BuffBase) Owner() IActor {
	return b.owner
}

// setOwner 设置此buff的持有者
func (b *BuffBase) setOwner(owner IActor) {
	b.owner = owner
}

// Provider 此Buff提供者
func (b *BuffBase) Provider() IActor {
	return b.provider
}

// setProvider 设置此Buff提供者
func (b *BuffBase) setProvider(provider IActor) {
	b.provider = provider
}

func (b *BuffBase) Def() *define.BuffDefine {
	return b.def
}

func (b *BuffBase) setDef(def *define.BuffDefine) {
	b.def = def
}

// NewBuffBase 创建一个新的 BuffBase 实例
func NewBuffBase() *BuffBase {
	return &BuffBase{
		maxDuration:      3,
		timeScale:        1,
		maxLevel:         1,
		buffType:         proto.BuffType_None,
		buffConflict:     proto.BuffConflict_Cover,
		dispellable:      true,
		name:             "默认名称",
		description:      "这个Buff没有介绍",
		demotion:         1,
		iconPath:         "",
		currentLevel:     0,
		residualDuration: 3,
		initialized:      false,
	}
}

func (b *BuffBase) BID() int {
	return b.def.BID
}

// MaxDuration 获取初始持续时间
func (b *BuffBase) MaxDuration() float32 {
	return b.maxDuration
}

// setMaxDuration 设置Buff的初始持续时间
func (b *BuffBase) setMaxDuration(value float32) {
	b.maxDuration = float32(math.Max(0, float64(value)))
}

// TimeScale 获取或设置Buff的时间流失速度，最小为0，最大为10
func (b *BuffBase) TimeScale() float32 {
	return b.timeScale
}

// SetTimeScale 设置Buff的时间流失速度，最小为0，最大为10
func (b *BuffBase) SetTimeScale(value float32) {
	b.timeScale = float32(math.Min(math.Max(0, float64(value)), 10))
}

// MaxLevel 获取Buff的最大堆叠层数，最小为1，最大为2147483647
func (b *BuffBase) MaxLevel() int {
	return b.maxLevel
}

// setMaxLevel 设置Buff的最大堆叠层数，最小为1，最大为2147483647
func (b *BuffBase) setMaxLevel(value int) {
	b.maxLevel = int(math.Max(1, float64(value)))
}

// BuffType 获取Buff的类型
func (b *BuffBase) BuffType() proto.BuffType {
	return b.buffType
}

// setBuffType 设置Buff的类型
func (b *BuffBase) setBuffType(value proto.BuffType) {
	b.buffType = value
}

// BuffConflict 获取Buff冲突处理方式
func (b *BuffBase) BuffConflict() proto.BuffConflict {
	return b.buffConflict
}

// setBuffConflict 设置Buff冲突处理方式
func (b *BuffBase) setBuffConflict(value proto.BuffConflict) {
	b.buffConflict = value
}

// Dispellable 获取Buff是否可驱散
func (b *BuffBase) Dispellable() bool {
	return b.dispellable
}

// setDispellable 设置Buff是否可驱散
func (b *BuffBase) setDispellable(value bool) {
	b.dispellable = value
}

// Name 获取Buff的名称
func (b *BuffBase) Name() string {
	return b.name
}

// setName 设置Buff的名称
func (b *BuffBase) setName(value string) {
	b.name = value
}

// Description 获取Buff的介绍文本
func (b *BuffBase) Description() string {
	return b.description
}

// setDescription 设置Buff的介绍文本
func (b *BuffBase) setDescription(value string) {
	b.description = value
}

// IconPath 获取Buff的图标资源路径
func (b *BuffBase) IconPath() string {
	return b.iconPath
}

// setIconPath 设置Buff的图标资源路径
func (b *BuffBase) setIconPath(value string) {
	b.iconPath = value
}

// Demotion 获取每次Buff持续时间结束时降低的等级
func (b *BuffBase) Demotion() int {
	return b.demotion
}

// setDemotion 设置每次Buff持续时间结束时降低的等级
func (b *BuffBase) setDemotion(value int) {
	b.demotion = int(math.Min(math.Max(0, float64(value)), float64(b.MaxLevel())))
}

// CurrentLevel 获取Buff的当前等级
func (b *BuffBase) CurrentLevel() int {
	return b.currentLevel
}

// SetCurrentLevel 设置Buff的当前等级
func (b *BuffBase) SetCurrentLevel(value int) {
	change := int(math.Min(math.Max(0, float64(value)), float64(b.MaxLevel()))) - b.currentLevel
	b.OnLevelChange(change)
	b.currentLevel += change
}

// ResidualDuration 获取Buff的当前剩余时间
func (b *BuffBase) ResidualDuration() float32 {
	return b.residualDuration
}

// SetResidualDuration 设置Buff的当前剩余时间
func (b *BuffBase) SetResidualDuration(value float32) {
	b.residualDuration = float32(math.Max(0, float64(value)))
}

// OnGet 当 owner 获得此 BuffBase 时触发
func (b *BuffBase) OnGet() {}

// OnLost 当 owner 失去此 BuffBase 时触发
func (b *BuffBase) OnLost() {}

// OnUpdate 每帧更新，由 BuffManager 调用
func (b *BuffBase) OnUpdate(delta float64) {}

// OnLevelChange 当等级改变时调用
func (b *BuffBase) OnLevelChange(change int) {}

// Init 初始化 BuffBase
func (b *BuffBase) Init(owner, provider IActor) error {
	if b.initialized {
		return errors.New("不能对已经初始化的 BuffBase 再次初始化")
	}
	if ActorIsNil(owner) || ActorIsNil(provider) {
		return errors.New("初始化值不能为空")
	}
	b.owner = owner
	b.provider = provider
	b.initialized = true

	def := b.GetBuffDefine()
	if def != nil {
		switch def.BuffType {
		case "正增益":
			b.buffType = proto.BuffType_Buff
		case "负增益":
			b.buffType = proto.BuffType_Debuff
		default:
			b.buffType = proto.BuffType_None
		}
		switch def.BuffConflict {
		case "合并":
			b.buffConflict = proto.BuffConflict_Combine
		case "独立":
			b.buffConflict = proto.BuffConflict_Separate
		case "覆盖":
			b.buffConflict = proto.BuffConflict_Cover
		default:
			return errors.New("buff BuffConflict Not Found" + def.BuffConflict)
		}
		b.setDef(def)
		b.setName(def.Name)
		b.setDescription(def.Description)
		b.setIconPath(def.IconPath)
		b.setMaxDuration(def.MaxDuration)
		b.setMaxLevel(def.MaxLevel)
		b.setDispellable(def.Dispellable)
		b.setDemotion(def.Demotion)
		b.SetTimeScale(def.TimeScale)
	}

	return nil
}

// GetBuffDefine 获取 BuffBase 的定义，需要在具体实现中提供
func (b *BuffBase) GetBuffDefine() *define.BuffDefine {
	return nil
}
