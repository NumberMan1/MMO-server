package model

import (
	"github.com/NumberMan1/MMO-server/define"
)

type Skill struct {
	Owner  IActor
	Define define.SkillDefine
}
