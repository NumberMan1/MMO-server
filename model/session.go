package model

import (
	"github.com/NumberMan1/MMO-server/database"
)

// Session 用户会话类
type Session struct {
	// 当前登录的角色
	Character *Character
	// 数据库玩家信息
	DbPlayer *database.DbPlayer
}

func NewSession() *Session {
	return &Session{}
}

// Space 当前所在地图
func (s *Session) Space() *Space {
	if s.Character != nil {
		return s.Character.Space()
	}
	return nil
}
