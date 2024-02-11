package service

import "github.com/NumberMan1/MMO-server/model"

// Session 用户会话类
type Session struct {
	// 当前登录的角色
	Character *model.Character
	// 当前所在地图
	Space *model.Space
}
