package model

import (
	"errors"
	"github.com/NumberMan1/MMO-server/game_common"
	"github.com/NumberMan1/MMO-server/game_common/protocol/gen/proto"
	"github.com/NumberMan1/MMO-server/game_server/core"
	"github.com/NumberMan1/MMO-server/game_server/core/vector3"
	"math/rand"
	"time"
)

var (
	y1000 = vector3.NewVector3(0, 1000, 0)
)

type Monster struct {
	*Actor
	AI           IAIBase
	Target       IActor           //目标
	MoveTarget   *vector3.Vector3 //移动目标位置
	MovePosition *vector3.Vector3 //当前移动位置
	InitPosition *vector3.Vector3 //出生点
	//rand         *rand.Rand
}

func NewMonster(tid, level int, position, direction *vector3.Vector3) *Monster {
	m := &Monster{
		Actor:        NewActor(proto.EntityType_Monster, tid, level, position, direction),
		InitPosition: position,
	}
	m.SetState(proto.EntityState_IDLE)
	//m.rand = rand.New(rand.NewSource(time.Now().Unix()))
	// 位置同步
	game_common.GetScheduleInstance().AddTask(func() error {
		if m.State() != proto.EntityState_MOVE || m.IsDeath() {
			return errors.New("怪物无法移动,因为状态不对或已经死亡")
		}
		//广播消息
		es := &proto.NetEntitySync{
			Entity: m.EntityData(),
			State:  m.State(),
		}
		m.Space().UpdateEntity(es)
		return nil
	}, time.Millisecond*150, 0)
	//设置AI对象
	switch m.Define().AI {
	case "Monster":
		m.AI = NewMonsterAI(m)
	}
	return m
}

func (m *Monster) MoveTo(target *vector3.Vector3) {
	if m.State() == proto.EntityState_IDLE {
		m.SetState(proto.EntityState_MOVE)
	}
	if m.MoveTarget != target {
		m.MoveTarget = target
		m.MovePosition = m.Position()
		dir := vector3.Normalize3(vector3.Sub3(m.MoveTarget, m.MovePosition))
		m.SetDirection(vector3.Dot(core.LookRotation(dir), y1000))
		//广播消息
		es := &proto.NetEntitySync{
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
	es := &proto.NetEntitySync{
		Entity: m.EntityData(),
		State:  m.State(),
	}
	m.Space().UpdateEntity(es)
}

func (m *Monster) Update() {
	if m.IsDeath() {
		return
	}
	m.Actor.Update()
	if m.AI != nil {
		m.AI.Update()
	}
	if m.State() == proto.EntityState_MOVE {
		//移动方向
		dir := vector3.Normalize3(vector3.Sub3(m.MoveTarget, m.MovePosition))
		m.SetDirection(vector3.Dot(core.LookRotation(dir), y1000))
		//logger.SLCDebug("-----------------")
		//logger.SLCDebug("%d speed %v", m.EntityId(), m.Speed())
		//logger.SLCDebug("-----------------")
		dist := float64(m.Speed()) * game_common.GetScheduleInstance().Clock().DeltaTime
		if vector3.GetDistance(m.MoveTarget, m.MovePosition) < dist {
			m.StopMove()
		} else {
			dir.Multiply(dist)
			m.MovePosition.Add(dir)
		}
		m.SetPosition(m.MovePosition)
	}
}

// RandomPointWithBirth 计算出生点附近的随机坐标
func (m *Monster) RandomPointWithBirth(r float64) *vector3.Vector3 {
	//x := m.rand.Float64()*2 - 1
	//z := m.rand.Float64()*2 - 1
	x := rand.Float64()*2 - 1
	z := rand.Float64()*2 - 1
	dir := vector3.Normalize3(vector3.NewVector3(x, 0, z))
	dir.Multiply(r)
	//dir.Multiply(m.rand.Float64())
	dir.Multiply(rand.Float64())
	return vector3.Add3(m.InitPosition, dir)
}

func (m *Monster) Attack(target IActor) {
	var sk *Skill = nil
	for _, skill := range m.SkillMgr().Skills {
		if skill.IsNormal() {
			sk = skill
			break
		}
	}
	//eSkill := m.skillMgr.Skills.Front()
	if sk == nil {
		return
	}
	if sk.State != Stage_None {
		return
	}
	m.Spell().SpellTarget(sk.Def.GetId(), target.EntityId())
}
