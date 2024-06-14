package model

import (
	"github.com/NumberMan1/MMO-server/config/define"
	"github.com/NumberMan1/MMO-server/protocol/gen/proto"
	"github.com/NumberMan1/common/logger"
)

type SkillManager struct {
	//归属的角色
	owner IActor
	//技能列表
	Skills []*Skill
}

func NewSkillManager(owner IActor) *SkillManager {
	sm := &SkillManager{
		owner:  owner,
		Skills: make([]*Skill, 0),
	}
	sm.InitSkills()
	return sm
}

// InitSkills 初始化技能信息,通过读取DataManager来加载技能信息
func (sm *SkillManager) InitSkills() {
	job := sm.owner.Define().TID
	sks := define.GetDataManagerInstance().Skills
	s := make([]int, 0)
	for _, v := range sks {
		if v.Job == job {
			s = append(s, v.ID)
		}
	}
	// 移除自定义技能,暂未实现
	for i, v := range s {
		if v == 2003 {
			s = append(s[:i], s[i+1:]...)
		}
	}
	// 加载野怪通用技能
	if job == 1002 || job == 1003 {
		s = append(s, 101)
	}
	sm.loadSkill(s)
}

func (sm *SkillManager) loadSkill(ids []int) {
	for _, id := range ids {
		sm.owner.Info().Skills = append(sm.owner.Info().Skills, &proto.SkillInfo{Id: int32(id)})
		sk := NewSkill(sm.owner, id)
		sm.Skills = append(sm.Skills, sk)
		logger.SLCInfo("角色[%v]加载技能[%v-%v]", sm.owner.Name(), sk.Def.GetId(), sk.Def.Name)
	}
}

func (sm *SkillManager) GetSkill(id int) *Skill {
	for _, skill := range sm.Skills {
		if skill.Def.GetId() == id {
			return skill
		}
	}
	return nil
}

func (sm *SkillManager) Update() {
	for _, sk := range sm.Skills {
		sk.Update()
	}
}
