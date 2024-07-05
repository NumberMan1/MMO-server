package model

import (
	"errors"
	"github.com/NumberMan1/MMO-server/game_common"
	"github.com/NumberMan1/MMO-server/game_common/protocol/gen/proto"
	"github.com/NumberMan1/MMO-server/game_server/core/vector3"
	"github.com/NumberMan1/MMO-server/game_server/mgr"
	"time"
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
	game_common.GetScheduleInstance().AddTask(g.teleport, time.Millisecond*500, 0)
	sp := GetSpaceManagerInstance().GetSpace(spaceId)
	if sp != nil {
		sp.EntityEnter(g)
	}
	return g
}

// teleport 执行传送
func (g *Gate) teleport() error {
	if g.Space() == nil || g.TargetSpace() == nil {
		return errors.New("传送门的所在地图或传送的目标地图为空")
	}
	units := RangeUnit(g.Position(), g.Space().Id, teleportGateRange)
	for e := units.Front(); e != nil; e = e.Next() {
		if chr, ok := e.Value.(*Character); ok {
			TeleportSpace(g.TargetSpace(), g.TargetPosition(), vector3.Zero3(), chr)
		}
	}
	return nil
}

// SetTarget 设置传送目标
func (g *Gate) SetTarget(targetSpace *Space, targetPosition *vector3.Vector3) {
	g.SetTargetSpace(targetSpace)
	g.SetTargetPosition(targetPosition)
}
