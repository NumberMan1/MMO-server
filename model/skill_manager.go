package model

import (
	"container/list"
	"github.com/NumberMan1/MMO-server/define"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
)

type SkillManager struct {
	//归属的角色
	owner IActor
	//技能列表
	Skills *list.List
}

func NewSkillManager(owner IActor) *SkillManager {
	sm := &SkillManager{
		owner:  owner,
		Skills: list.New(),
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
		if v.TID == job {
			s = append(s, v.ID)
		}
	}
	sm.loadSkill(s)
}

func (sm *SkillManager) loadSkill(ids []int) {
	for _, id := range ids {
		sm.owner.Info().Skills = append(sm.owner.Info().Skills, &proto.SkillInfo{Id: int32(id)})
		sk := NewSkill(sm.owner, id)
		sm.Skills.PushBack(sk)
		logger.SLCInfo("角色[%v]加载技能[%v-%v]", sm.owner.Name(), sk.Def.GetId(), sk.Def.Name)
	}
}

func (sm *SkillManager) GetSkill(id int) *Skill {
	for e := sm.Skills.Front(); e != nil; e = e.Next() {
		if e.Value.(*Skill).Def.GetId() == id {
			return e.Value.(*Skill)
		}
	}
	return nil
}

func (sm *SkillManager) Update() {
	for e := sm.Skills.Front(); e != nil; e = e.Next() {
		e.Value.(*Skill).Update()
	}
}
