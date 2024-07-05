package model

import (
	"github.com/NumberMan1/MMO-server/game_common/protocol/gen/proto"
	"github.com/NumberMan1/MMO-server/game_server/core/vector3"
	"github.com/NumberMan1/MMO-server/game_server/model/fight"
	"github.com/NumberMan1/common/logger"
)

type Spell struct {
	owner IActor
}

func NewSpell(owner IActor) *Spell {
	return &Spell{owner: owner}
}

func (s *Spell) Owner() IActor {
	return s.owner
}

// Intonate 吟唱技能
func (s *Spell) Intonate(skill *Skill) {

}

// RunCast 统一施法
func (s *Spell) RunCast(info *proto.CastInfo) {
	sk := s.Owner().SkillMgr().GetSkill(int(info.SkillId))
	if sk.IsUnitTarget() {
		s.SpellTarget(int(info.SkillId), int(info.TargetId))
	} else if sk.IsPointTarget() {
		loc := vector3.NewVector3(float64(info.TargetLoc.X), float64(info.TargetLoc.Y), float64(info.TargetLoc.Z))
		s.SpellPosition(int(info.SkillId), loc)
	} else if sk.IsNoneTarget() {
		s.SpellNoTarget(int(info.SkillId))
	}
}

// SpellNoTarget 施放无目标技能
func (s *Spell) SpellNoTarget(skillId int) {
	logger.SLCInfo("Spell::SpellNoTarget():Caster[%v]:Skill[%v]", s.Owner().EntityId(), skillId)
	//检查技能
	sk := s.Owner().SkillMgr().GetSkill(skillId)
	if sk == nil {
		logger.SLCInfo("Spell::SpellNoTarget():owner[%v]:Skill=%v not found", s.Owner().EntityId(), skillId)
		return
	}
	//执行技能
	sco := fight.NewSCEntity(s.Owner())
	res := sk.CanUse(sco)
	if res != proto.CastResult_Success {
		logger.SLCInfo("Cast Fail Skill %v %v", sk.Def.GetId(), res)
		s.OnSpellFailure(skillId, res)
		return
	}
	sk.Use(sco)
	info := &proto.CastInfo{
		CasterId: int32(s.Owner().EntityId()),
		SkillId:  int32(skillId),
	}
	s.Owner().Space().FightMgr.SpellQueue.Push(info)
}

// SpellTarget 施放单位目标技能
func (s *Spell) SpellTarget(skillId int, targetId int) {
	logger.SLCInfo("Spell::SpellTarget():Caster[%v]:Skill[%v]:Target[%v]", s.Owner().EntityId(), skillId, targetId)
	//检查技能
	sk := s.Owner().SkillMgr().GetSkill(skillId)
	if sk == nil {
		logger.SLCInfo("Spell::SpellTarget():owner[%v]:Skill=%v not found", s.Owner().EntityId(), skillId)
		return
	}
	//检查目标
	target := GetUnit(targetId)
	if ActorIsNil(target) {
		logger.SLCInfo("Spell::SpellTarget():owner[%v]:Target=%v not found", s.Owner().EntityId(), targetId)
		return
	}
	//执行技能
	sco := fight.NewSCEntity(target)
	res := sk.CanUse(sco)
	if res != proto.CastResult_Success {
		logger.SLCInfo("Cast Fail Skill %v %v", sk.Def.GetId(), res)
		s.OnSpellFailure(skillId, res)
		return
	}
	sk.Use(sco)
	info := &proto.CastInfo{
		CasterId: int32(s.Owner().EntityId()),
		TargetId: int32(targetId),
		SkillId:  int32(skillId),
	}
	s.Owner().Space().FightMgr.SpellQueue.Push(info)
}

// SpellPosition 施放点目标技能
func (s *Spell) SpellPosition(skillId int, position *vector3.Vector3) {
	logger.SLCInfo("Spell::SpellPosition():Caster[%v]:Skill[%v]:Pos[%v]", s.Owner().EntityId(), skillId, *position)
	//sco := NewSCObject(NewSCPosition(position))
}

// OnSpellFailure 通知玩家技能失败
func (s *Spell) OnSpellFailure(skillId int, reason proto.CastResult) {
	if chr, ok := s.Owner().(*Character); ok {
		resp := &proto.SpellFailResponse{
			CasterId: int32(s.Owner().EntityId()),
			SkillId:  int32(skillId),
			Reason:   reason,
		}
		chr.Session.Send(resp)
	}
}
