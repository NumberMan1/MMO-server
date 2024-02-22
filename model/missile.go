package model

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/fight"
	"github.com/NumberMan1/common/logger"
)

type Missile struct {
	skill    *Skill
	target   fight.ISCObject
	initPos  *vector3.Vector3
	Position *vector3.Vector3
}

func NewMissile(skill *Skill, initPos *vector3.Vector3, target fight.ISCObject) *Missile {
	m := &Missile{
		skill:    skill,
		target:   target,
		initPos:  initPos,
		Position: initPos,
	}
	logger.SLCInfo("Position:%v", m.Position)
	return m
}

func (m *Missile) FightMgr() *FightMgr {
	return m.skill.FightMgr()
}

func (m *Missile) OnUpdate(dt float64) {
	a := m.Position
	b := m.target.GetPosition()
	direction := vector3.Normalize3(vector3.Sub3(b, a))
	dist := float64(m.skill.Def.MissileSpeed) * dt
	if dist > vector3.GetDistance(a, b) {
		m.Position = b
		m.skill.OnHit(m.target)
		for e := m.FightMgr().Missiles.Front(); e != nil; e = e.Next() {
			if e.Value.(*Missile) == m {
				m.FightMgr().Missiles.Remove(e)
				break
			}
		}
	} else {
		direction.Multiply(dist)
		m.Position = vector3.Add3(m.Position, direction)
	}
}
