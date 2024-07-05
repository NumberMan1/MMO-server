package model

import (
	"github.com/NumberMan1/MMO-server/game_common/protocol/gen/proto"
	define2 "github.com/NumberMan1/MMO-server/game_server/config/define"
	"github.com/NumberMan1/MMO-server/game_server/model/entity"
)

// IActor 接口定义
// IActor 代表一个游戏中的角色，其行为和属性由接口中定义的方法来描述
type IActor interface {
	// IEntity interface 的继承
	entity.IEntity

	// BuffManager 获取该角色的buff管理器
	BuffManager() *BuffManager
	// SetBuffManager 设置该角色的buff管理器
	SetBuffManager(buffManager *BuffManager)

	// SetHp 设置角色当前血量
	SetHp(hp float32)

	// SetMp 设置角色当前法力值
	SetMp(mp float32)

	// SetLevel 设置角色等级
	SetLevel(level int32)

	// SetExp 设置角色经验值
	SetExp(exp int64)

	// SetGold 设置角色金币数量
	SetGold(gold int64)

	// Level 获取角色等级
	Level() int

	// Exp 获取角色经验值
	Exp() int

	// Gold 获取角色金币数量
	Gold() int

	// HPMax 获取角色最大 HP
	HPMax() int

	// MPMax 获取角色最大 MP
	MPMax() int

	// UnitState 获取角色的单位状态
	UnitState() proto.UnitState

	// SetUnitState 设置角色的单位状态
	SetUnitState(unitState proto.UnitState)

	// SkillMgr 获取角色的技能管理器
	SkillMgr() *SkillManager

	// SetSkillMgr 设置角色的技能管理器
	SetSkillMgr(skillMgr *SkillManager)

	// Spell 获取角色的魔法
	Spell() *Spell

	// SetSpell 设置角色的魔法
	SetSpell(spell *Spell)

	// Hp 获取角色当前血量
	Hp() float32

	// Mp 获取角色当前法力值
	Mp() float32

	// Id 获取角色 ID
	Id() int

	// Name 获取角色名称
	Name() string

	// Type 获取角色类型
	Type() proto.EntityType

	// SetId 设置角色 Id
	SetId(v int)

	// SetName 设置角色名称
	SetName(v string)

	// SetType 设置角色类型
	SetType(v proto.EntityType)

	// State 获取角色的当前状态
	State() proto.EntityState

	// SetState 设置角色的当前状态
	SetState(state proto.EntityState)

	// Space 获取角色所在的地图
	Space() *Space

	// SetSpace 设置角色所在的地图
	SetSpace(space *Space)

	// Info 获取角色的基础信息
	Info() *proto.NetActor

	// SetInfo 设置角色的基础信息
	SetInfo(info *proto.NetActor)

	// Define 获取角色的定义数据
	Define() *define2.UnitDefine

	// SetDefine 设置角色的定义数据
	SetDefine(define *define2.UnitDefine)

	// Attr 获取角色的属性集合
	Attr() *AttributesAssembly

	// SetAttr 设置角色的属性集合
	SetAttr(attr *AttributesAssembly)

	// IsDeath 判断角色是否已死亡
	IsDeath() bool

	// Revive 复活角色
	Revive()

	// Update 更新角色的状态
	Update()

	// Die 处理角色死亡
	Die(killerID int)

	// OnBeforeDie 角色死前的处理
	OnBeforeDie(killerID int)

	// OnAfterDie 角色死后的处理
	OnAfterDie(killerID int)

	// RecvDamage 处理角色收到的伤害
	RecvDamage(dmg *proto.Damage)

	// SetAndUpdateHp 设置并更新角色当前血量
	SetAndUpdateHp(hp float32)

	// SetAndUpdateMp 设置并更新角色当前法力值
	SetAndUpdateMp(mp float32)

	// SetAndUpdateState 设置并更新角色的单位状态
	SetAndUpdateState(unitState proto.UnitState)

	// SetAndUpdateGolds 设置并更新角色的金币数量
	SetAndUpdateGolds(value int64)

	// SetAndUpdateExp 设置并更新角色的经验值
	SetAndUpdateExp(value int64)

	// SetAndUpdateLevel 设置并更新角色的等级
	SetAndUpdateLevel(value int)

	// SyncSpeed 通知客户端：Speed 变化
	SyncSpeed(value int)

	// SyncHpMax 通知客户端：HPMax 变化
	SyncHpMax(value float32)

	// SyncMpMax 通知客户端：MPMax 变化
	SyncMpMax(value float32)
}
