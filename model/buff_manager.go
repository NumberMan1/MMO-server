package model

import (
	"github.com/NumberMan1/MMO-server/protocol/gen/proto"
	"github.com/NumberMan1/common/summer/timeunit"
)

// BuffManager 每个角色都有自己的Buff管理器
type BuffManager struct {
	owner    IActor
	buffs    []Buff
	Observer func(Buff)
}

// NewBuffManager creates a new BuffManager
func NewBuffManager(owner IActor) *BuffManager {
	return &BuffManager{
		owner: owner,
		buffs: []Buff{},
	}
}

// AddBuff 在目标身上挂buff
func AddBuff[T Buff](bm *BuffManager, provider IActor, level int, constructor func() T) {
	var temp01 []T
	for _, buff := range bm.buffs {
		if t, ok := buff.(T); ok {
			temp01 = append(temp01, t)
		}
	}

	if len(temp01) == 0 {
		addNewBuff[T](bm, provider, level, constructor)
	} else {
		switch temp01[0].BuffConflict() {
		case proto.BuffConflict_Separate:
			temp := true
			for _, item := range temp01 {
				if item.Provider() == provider {
					item.SetCurrentLevel(item.CurrentLevel() + level)
					temp = false
					continue
				}
			}
			if temp {
				addNewBuff[T](bm, provider, level, constructor)
			}
		case proto.BuffConflict_Combine:
			temp01[0].SetCurrentLevel(temp01[0].CurrentLevel() + level)
		case proto.BuffConflict_Cover:
			bm.RemoveBuff(temp01[0])
			addNewBuff[T](bm, provider, level, constructor)
		}
	}
}

// FindBuff 获得单位身上指定类型的buff的列表
func FindBuff[T Buff](bm *BuffManager) []T {
	var result []T
	for _, buff := range bm.buffs {
		if t, ok := buff.(T); ok {
			result = append(result, t)
		}
	}
	return result
}

// RemoveBuff 移除单位身上指定的一个buff
func (bm *BuffManager) RemoveBuff(buff Buff) bool {
	for i, item := range bm.buffs {
		if item == buff {
			item.SetCurrentLevel(0)
			item.OnLost()
			bm.buffs = append(bm.buffs[:i], bm.buffs[i+1:]...)
			return true
		}
	}
	return false
}

// addNewBuff adds a new buff to the Target
func addNewBuff[T Buff](bm *BuffManager, provider IActor, level int, constructor func() T) {
	buff := constructor()
	err := buff.Init(bm.owner, provider)
	if err != nil {
		return
	}
	bm.buffs = append(bm.buffs, buff)
	buff.SetResidualDuration(buff.MaxDuration())
	buff.SetCurrentLevel(level)
	buff.OnGet()
	if bm.Observer != nil {
		bm.Observer(buff)
	}
}

// OnUpdate updates all buffs each frame
func (bm *BuffManager) OnUpdate(delta float64) {
	//所有Buff执行Update
	for _, item := range bm.buffs {
		if item.CurrentLevel() > 0 && item.Owner() != nil {
			item.OnUpdate(delta)
		}
	}

	//降低持续时间,清理无用buff
	for i := len(bm.buffs) - 1; i >= 0; i-- {
		buff := bm.buffs[i]
		//如果等级为0，则移除
		if buff.CurrentLevel() == 0 {
			bm.RemoveBuff(buff)
			continue
		}
		//如果持续时间为0，则降级,
		//降级后如果等级为0则移除，否则刷新持续时间
		if buff.ResidualDuration() == 0 {
			buff.SetCurrentLevel(buff.CurrentLevel() - buff.Demotion())
			if buff.CurrentLevel() == 0 {
				bm.RemoveBuff(buff)
				continue
			} else {
				buff.SetResidualDuration(buff.MaxDuration())
			}
		}
		//降低持续时间
		buff.SetResidualDuration(buff.ResidualDuration() - float32(timeunit.DeltaTime))
	}
}
