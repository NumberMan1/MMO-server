package model

import (
	define2 "github.com/NumberMan1/MMO-server/config/define"
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/fight"
	"github.com/NumberMan1/MMO-server/protocol/gen/proto"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer"
	"github.com/NumberMan1/common/summer/core"
	"github.com/NumberMan1/common/summer/timeunit"
	"math"
	"math/rand"
	"slices"
)

// Stage 技能阶段
type Stage int

const (
	Stage_None Stage = iota
	Stage_Intonate
	Stage_Active
	Stage_Coolding
)

type Skill struct {
	IsPassive bool                 // 是否被动技能
	Def       *define2.SkillDefine // 技能设定
	Owner     IActor               // 技能归属者
	Cd        float64              // 冷却计时，0代表技能可用
	time      float64              // 技能运行时间
	State     Stage                // 当前技能状态
	notCrit   int                  // 未暴击次数
	target    fight.ISCObject
	//rand      *rand.Rand
}

func NewSkill(owner IActor, skid int) *Skill {
	s := &Skill{
		Owner: owner,
		Def:   define2.GetDataManagerInstance().Skills[skid],
		//rand:  rand.New(rand.NewSource(time.Now().Unix())),
	}
	//伤害延迟默认float[]={0}
	if len(s.Def.HitDelay) == 0 {
		s.Def.HitDelay = []float32{0}
	}

	return s
}

func (s *Skill) Target() fight.ISCObject {
	return s.target
}

// 强制触发暴击的次数
func (s *Skill) forceCritAfter() int {
	return int((float32(100))/s.Owner.Attr().Final.CRI) + 2
}

func (s *Skill) IsUnitTarget() bool {
	return s.Def.TargetType == "单位"
}

func (s *Skill) IsPointTarget() bool {
	return s.Def.TargetType == "点"
}

func (s *Skill) IsNoneTarget() bool {
	return s.Def.TargetType == "None"
}

func (s *Skill) IsNormal() bool {
	return s.Def.Type == "普通攻击"
}

func (s *Skill) FightMgr() *FightMgr {
	return s.Owner.Space().FightMgr
}

func (s *Skill) Update() {
	if s.State == Stage_None && core.Equal(s.Cd, 0) {
		return
	}
	if s.Cd > 0 {
		s.Cd -= timeunit.DeltaTime
		return
	}
	if s.Cd < 0 {
		s.Cd = 0
	}
	if s.time < 0.0001 {
		s.onIntonate()
	}
	s.time += timeunit.DeltaTime
	if s.State == Stage_Intonate && s.time >= float64(s.Def.IntonateTime) {
		s.State = Stage_Active
		s.Cd = float64(s.Def.CD)
		s.OnActive()
	}
	if s.State == Stage_Active {
		if s.time >= float64(s.Def.IntonateTime+slices.Max(s.Def.HitDelay)) {
			s.State = Stage_Coolding
		}
	}
	if s.State == Stage_Coolding {
		if core.Equal(s.Cd, 0) {
			s.time = 0
			s.State = Stage_None
			s.OnFinish()
		}
	}
}

func (s *Skill) onIntonate() {
	logger.SLCInfo("技能蓄力：Owner[%v],Skill[%v]", s.Owner.EntityId(), s.Def.Name)
}

func (s *Skill) OnActive() {
	logger.SLCInfo("技能激活：Owner[%v],Skill[%v]", s.Owner.EntityId(), s.Def.Name)
	//如果是投射物
	if s.Def.IsMissile {
		missile := NewMissile(s, s.Owner.Position(), s.target)
		s.FightMgr().Missiles = append(s.FightMgr().Missiles, missile)
	} else {
		//如果不是投射物
		logger.SLCInfo("Def.HitDelay.Length=%v", len(s.Def.HitDelay))
		for _, v := range s.Def.HitDelay {
			summer.GetScheduleInstance().AddTask(s.hitTrigger, timeunit.Milliseconds, int(v*1000), 1)
		}
	}
}

// 触发延迟伤害
func (s *Skill) hitTrigger() {
	logger.SLCInfo("hitTrigger：Owner[%v],Skill[%v]", s.Owner.EntityId(), s.Def.Name)
	s.OnHit(s.target)
}

// OnHit 技能打到目标
func (s *Skill) OnHit(sco fight.ISCObject) {
	logger.SLCInfo("OnHit：Owner[%v],Skill[%v]，SCO[%v]", s.Owner.EntityId(), s.Def.Name, sco)
	//单体伤害
	if s.Def.Area == 0 {
		if ob, ok := sco.(*fight.SCEntity); ok {
			actor := ob.GetRealObj().(IActor)
			s.takeDamage(actor)
		}
	} else {
		//范围伤害
		logger.SLCInfo("范围伤害：Space[%v],Center[%v],Area[%v]", s.Owner.EntityId(), sco.GetPosition(), s.Def.Area)
		l := RangeUnit(sco.GetPosition(), s.Owner.Space().Id, s.Def.Area)
		for e := l.Front(); e != nil; e = e.Next() {
			actor := e.Value.(IActor)
			s.takeDamage(actor)
		}
	}
}

// 对目标造成伤害
func (s *Skill) takeDamage(target IActor) {
	if target.IsDeath() || target == s.Owner {
		return
	}
	logger.SLCInfo("Skill:TakeDamage:Atker[%v],Target[%v]", s.Owner.EntityId(), target.EntityId())
	//计算伤害数值、暴击、闪避、
	//扣除目标HP、广播通知

	//伤害=攻击[攻]×(1-护甲[守]/(护甲[守]+400+85×等级[攻]))
	a := s.Owner.Attr().Final //攻击者属性
	b := target.Attr().Final  //被攻击者属性
	//伤害信息
	dmg := &proto.Damage{
		AttackerId: int32(s.Owner.EntityId()),
		TargetId:   int32(target.EntityId()),
		SkillId:    int32(s.Def.GetId()),
	}
	//技能的物攻和法攻
	ad := s.Def.AD + a.AD*s.Def.ADC
	ap := s.Def.AP + a.AP*s.Def.APC
	//计算伤害
	ads := ad * (1 - b.DEF/(b.DEF+400+85*float32(s.Owner.Info().Level)))
	aps := ap * (1 - b.MDEF/(b.MDEF+400+85*float32(s.Owner.Info().Level)))
	logger.SLCInfo("ads=%v , aps=%v", ads, aps)
	dmg.Amount = ads + aps
	//计算暴击
	s.notCrit += 1
	//randCri := s.rand.Float32()
	randCri := rand.Float32()
	cri := a.CRI * 0.01
	logger.SLCInfo("暴击计算：%v / %v | [%v/%v]", randCri, cri, s.notCrit, s.forceCritAfter())
	if randCri < cri || s.notCrit > s.forceCritAfter() {
		s.notCrit = 0
		dmg.IsCrit = true
		dmg.Amount *= float32(math.Max(float64(a.CRD), 100) * 0.01)
	}
	//计算闪避
	hitRate := (a.HitRate - b.DodgeRate) * 0.01
	logger.SLCInfo("闪避计算：%v %v %v", hitRate, a.HitRate, b.DodgeRate)
	//if s.rand.Float32() > hitRate {
	if rand.Float32() > hitRate {
		dmg.IsMiss = true
		dmg.Amount = 0
	}
	target.RecvDamage(dmg)
}

func (s *Skill) OnFinish() {
	logger.SLCInfo("技能结束：Owner[%v],Skill[%v]", s.Owner.EntityId(), s.Def.Name)
}

// CanUse 检查技能是否可用
func (s *Skill) CanUse(sco fight.ISCObject) proto.CastResult {
	if s.IsPassive { //被动技能
		return proto.CastResult_IsPassive
	} else if s.Owner.Mp() < float32(s.Def.Cost) { //MP不足
		return proto.CastResult_MpLack
	} else if s.State != Stage_None { //正在进行
		return proto.CastResult_Running
	} else if s.Cd != 0 { //冷却中
		return proto.CastResult_Cooldown
	} else if s.Owner.IsDeath() { //Entity已经死亡
		return proto.CastResult_EntityDead
	} else if s, ok := sco.(*fight.SCEntity); ok && s.GetRealObj().(IActor).IsDeath() { //目标已经死亡
		return proto.CastResult_EntityDead
	}
	//施法者和目标的距离
	logger.SLCInfo("2者位置：%v , %v", s.Owner.Position(), sco.GetPosition())
	dist := vector3.GetDistance(s.Owner.Position(), sco.GetPosition())
	logger.SLCInfo("施法者和目标的距离 %v, %v", dist, s.Def.SpellRange)
	if dist > float64(s.Def.SpellRange) {
		return proto.CastResult_OutOfRange
	}
	return proto.CastResult_Success
}

// Use 使用技能
func (s *Skill) Use(sco fight.ISCObject) proto.CastResult {
	s.target = sco
	s.time = 0
	s.State = Stage_Intonate
	return proto.CastResult_Success
}
