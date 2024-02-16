package model

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/common/summer"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
	"github.com/NumberMan1/common/summer/timeunit"
	"math"
	"math/rand"
	"time"
)

var (
	y1000 = vector3.NewVector3(0, 1000, 0)
)

type Monster struct {
	*Actor
	AI           IAIBase
	Target       IActor          //目标
	MoveTarget   vector3.Vector3 //移动目标位置
	MovePosition vector3.Vector3 //当前移动位置
	InitPosition vector3.Vector3 //出生点
	rand         *rand.Rand
}

func NewMonster(tid, level int, position, direction vector3.Vector3) *Monster {
	m := &Monster{
		Actor:        NewActor(proto.EntityType_Monster, tid, level, position, direction),
		InitPosition: position,
	}
	m.SetState(proto.EntityState_IDLE)
	m.rand = rand.New(rand.NewSource(time.Now().Unix()))
	// 位置同步
	summer.GetScheduleInstance().AddTask(func() {
		if m.State() != proto.EntityState_MOVE {
			return
		}
		es := &proto.NEntitySync{
			Entity: m.EntityData(),
			State:  m.State(),
		}
		m.Space().UpdateEntity(es)
	}, timeunit.Milliseconds, 150, 0)
	//设置AI对象
	switch m.Define().AI {
	case "Monster":
		m.AI = NewMonsterAI(m)
	}
	return m
}

func (m *Monster) MoveTo(target vector3.Vector3) {
	if m.State() == proto.EntityState_IDLE {
		m.SetState(proto.EntityState_MOVE)
	}
	if m.MoveTarget != target {
		m.MoveTarget = target
		m.MovePosition = m.Position()
		dir := vector3.Normalize3(vector3.Sub3(m.MoveTarget, m.MovePosition))
		m.SetDirection(vector3.Dot(m.LookRotation(dir), y1000))
		//广播消息
		es := &proto.NEntitySync{
			Entity: m.EntityData(),
			State:  m.State(),
		}
		m.Space().UpdateEntity(es)
	}
}

func (m *Monster) StopMove() {
	m.SetState(proto.EntityState_IDLE)
	m.MovePosition = m.MoveTarget
	//广播消息
	es := &proto.NEntitySync{
		Entity: m.EntityData(),
		State:  m.State(),
	}
	m.Space().UpdateEntity(es)
}

func (m *Monster) Update() {
	if m.AI != nil {
		m.AI.Update()
	}
	if m.State() == proto.EntityState_MOVE {
		//移动方向
		dir := vector3.Normalize3(vector3.Sub3(m.MoveTarget, m.MovePosition))
		m.SetDirection(vector3.Dot(m.LookRotation(dir), y1000))
		dist := timeunit.Time
		if vector3.GetDistance(m.MoveTarget, m.MovePosition) < dist {
			m.StopMove()
		} else {
			dir.Multiply(dist)
			m.MovePosition.Add(dir)
		}
		m.SetPosition(m.MovePosition)
	}
}

// LookRotation 方向向量转欧拉角
func (m *Monster) LookRotation(fromDir vector3.Vector3) vector3.Vector3 {
	rad2Deg := 57.29578
	eulerAngles := vector3.NewVector3(0, 0, 0)
	eulerAngles.X = math.Acos(math.Sqrt((fromDir.X*fromDir.X+fromDir.Z*fromDir.Z)/(fromDir.X*fromDir.X+fromDir.Y*fromDir.Y+fromDir.Z*fromDir.Z))) * rad2Deg
	if fromDir.Y > 0 {
		eulerAngles.X = 360 - eulerAngles.X
	}
	//AngleY = arc tan(x/z)
	eulerAngles.Y = math.Atan2(fromDir.Z, fromDir.X) * rad2Deg
	if eulerAngles.Y < 0 {
		eulerAngles.Y += 180
	}
	if fromDir.X < 0 {
		eulerAngles.Y += 180
	}
	//AngleZ = 0
	eulerAngles.Z = 0
	return eulerAngles
}

// RandomPointWithBirth 计算出生点附近的随机坐标
func (m *Monster) RandomPointWithBirth(r float64) vector3.Vector3 {
	x := m.rand.Float64()*2 - 1
	z := m.rand.Float64()*2 - 1
	dir := vector3.Normalize3(vector3.NewVector3(x, 0, z))
	dir.Multiply(r)
	dir.Multiply(m.rand.Float64())
	return vector3.Add3(m.InitPosition, dir)
}
