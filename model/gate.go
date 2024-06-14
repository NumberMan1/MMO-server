package model

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/mgr"
	"github.com/NumberMan1/MMO-server/protocol/gen/proto"
	"github.com/NumberMan1/common/summer"
	"github.com/NumberMan1/common/summer/timeunit"
)

const (
	teleportGateRange = 2000
)

// Gate 传送门
type Gate struct {
	*Actor
	//传送到哪个场景
	targetSpace *Space
	//传送到什么位置
	targetPosition *vector3.Vector3
}

func (g *Gate) TargetSpace() *Space {
	return g.targetSpace
}

func (g *Gate) SetTargetSpace(targetSpace *Space) {
	g.targetSpace = targetSpace
}

func (g *Gate) TargetPosition() *vector3.Vector3 {
	return g.targetPosition
}

func (g *Gate) SetTargetPosition(targetPosition *vector3.Vector3) {
	g.targetPosition = targetPosition
}

func NewGate(spaceId, tid int, position, direction *vector3.Vector3) *Gate {
	g := &Gate{Actor: NewActor(proto.EntityType_Gate, tid, 0, position, direction)}
	mgr.GetEntityManagerInstance().AddEntity(spaceId, g)
	summer.GetScheduleInstance().AddTask(g.teleport, timeunit.Milliseconds, 500, 0)
	sp := GetSpaceManagerInstance().GetSpace(spaceId)
	if sp != nil {
		sp.EntityEnter(g)
	}
	return g
}

// teleport 执行传送
func (g *Gate) teleport() {
	if g.Space() == nil || g.TargetSpace() == nil {
		return
	}
	units := RangeUnit(g.Position(), g.Space().Id, teleportGateRange)
	for e := units.Front(); e != nil; e = e.Next() {
		if chr, ok := e.Value.(*Character); ok {
			TeleportSpace(g.TargetSpace(), g.TargetPosition(), vector3.Zero3(), chr)
		}
	}
}

// SetTarget 设置传送目标
func (g *Gate) SetTarget(targetSpace *Space, targetPosition *vector3.Vector3) {
	g.SetTargetSpace(targetSpace)
	g.SetTargetPosition(targetPosition)
}
