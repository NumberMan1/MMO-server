package model

import (
	"github.com/NumberMan1/MMO-server/core/fsm"
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
	"github.com/NumberMan1/common/summer/timeunit"
	"math/rand"
	"reflect"
	"time"
)

type Param struct {
	Owner *Monster
	//视野范围
	ViewRange int
	//相对于出生点的活动范围
	WalkRange int
	//追击范围
	ChaseRange int
	Rand       *rand.Rand
}

func NewParam() *Param {
	return &Param{
		Owner:      nil,
		ViewRange:  8000,
		WalkRange:  8000,
		ChaseRange: 12000,
		Rand:       rand.New(rand.NewSource(time.Now().Unix())),
	}
}

// WalkState 巡逻状态
type WalkState struct {
	*fsm.State[*Param]
	lastTime float64
	//等待时间
	waitTime float64
}

func NewWalkState() *WalkState {
	return &WalkState{
		State:    fsm.NewState[*Param](),
		lastTime: timeunit.Time,
		waitTime: 10,
	}
}

func (ws *WalkState) OnEnter() {
	ws.State.P().Owner.StopMove()
}

func (ws *WalkState) OnUpdate() {
	mon := ws.P().Owner
	//查询 8000 范围内的玩家
	chr := GetNearest[*Character](GetEntityManagerInstance(), mon.Space().Id, ws.P().ViewRange, mon.Position())
	if !reflect.ValueOf(chr).IsZero() {
		mon.Target = chr
		ws.Fsm().ChangeState("chase")
		return
	}
	if mon.State() == proto.EntityState_IDLE {
		if ws.lastTime+ws.waitTime < timeunit.Time {
			ws.lastTime = timeunit.Time
			ws.waitTime = (ws.P().Rand.Float64() * 20) + 10
			//移动到随机位置
			target := mon.RandomPointWithBirth(float64(ws.P().WalkRange))
			mon.MoveTo(target)
		}
	}
}

// ChaseState 追击状态
type ChaseState struct {
	*fsm.State[*Param]
}

func NewChaseState() *ChaseState {
	return &ChaseState{State: fsm.NewState[*Param]()}
}

func (cs *ChaseState) OnUpdate() {
	mon := cs.P().Owner
	if mon.Target == nil || mon.Target.IsDeath() || !GetEntityManagerInstance().Exist(mon.Target.Id()) {
		mon.Target = nil
		cs.Fsm().ChangeState("walk")
		return
	}
	//自身与出生点的距离
	m := vector3.GetDistance(mon.InitPosition, mon.Position())
	//自身和目标的距离
	n := vector3.GetDistance(mon.Position(), mon.Target.Position())
	if m > float64(cs.P().ChaseRange) || n > float64(cs.P().ViewRange) {
		//返回出生点
		cs.Fsm().ChangeState("goback")
		return
	}
	if n < 1200 {
		if mon.State() == proto.EntityState_MOVE {
			mon.StopMove()
		}
		logger.SLCInfo("发起攻击")
	} else {
		mon.MoveTo(mon.Target.Position())
	}
	return
}

// GoBackState 返回状态
type GoBackState struct {
	*fsm.State[*Param]
}

func NewGoBackState() *GoBackState {
	return &GoBackState{State: fsm.NewState[*Param]()}
}

func (gbs *GoBackState) OnEnter() {
	gbs.P().Owner.MoveTo(gbs.P().Owner.InitPosition)
}

func (gbs *GoBackState) OnUpdate() {
	mon := gbs.P().Owner
	if vector3.GetDistance(mon.InitPosition, mon.Position()) > 100 {
		gbs.Fsm().ChangeState("walk")
	}
}

type MonsterAI struct {
	*AIBase
	fsmSystem *fsm.System[*Param]
}

func NewMonsterAI(owner *Monster) *MonsterAI {
	param := NewParam()
	param.Owner = owner
	a := &MonsterAI{AIBase: NewBase(owner), fsmSystem: fsm.NewSystem[*Param](param)}
	a.fsmSystem.AddState("walk", NewWalkState())
	a.fsmSystem.AddState("chase", NewChaseState())
	a.fsmSystem.AddState("goback", NewGoBackState())
	return a
}

func (a *MonsterAI) Update() {
	if a.fsmSystem != nil {
		a.fsmSystem.Update()
	}
}
