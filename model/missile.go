package model

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/fight"
	"github.com/NumberMan1/common/logger"
)

// Missile 投射物
type Missile struct {
	//投射物所在场景
	space *Space
	//所属技能
	skill *Skill
	//追击目标
	target fight.ISCObject
	//初始位置
	initPos *vector3.Vector3
	//飞弹当前位置
	Position *vector3.Vector3
}

func NewMissile(skill *Skill, initPos *vector3.Vector3, target fight.ISCObject) *Missile {
	m := &Missile{
		space:    skill.Owner.Space(),
		skill:    skill,
		target:   target,
		initPos:  initPos,
		Position: initPos,
	}
	logger.SLCInfo("Position:%v", m.Position)
	return m
}

func (m *Missile) FightMgr() *FightMgr {
	return m.space.FightMgr
}

func (m *Missile) OnUpdate(dt float64) {
	a := m.Position
	b := m.target.GetPosition()
	direction := vector3.Normalize3(vector3.Sub3(b, a))
	dist := float64(m.skill.Def.MissileSpeed) * dt
	if dist > vector3.GetDistance(a, b) {
		m.Position = b
		m.skill.OnHit(m.target)
		for i, missile := range m.FightMgr().Missiles {
			if missile == m {
				m.FightMgr().Missiles = append(m.FightMgr().Missiles[:i], m.FightMgr().Missiles[i+1:]...)
				break
			}
		}
	} else {
		direction.Multiply(dist)
		m.Position = vector3.Add3(m.Position, direction)
	}
}
