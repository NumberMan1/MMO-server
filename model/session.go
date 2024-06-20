package model

import (
	"github.com/NumberMan1/MMO-server/database"
	"github.com/NumberMan1/common/global/variable"
	"github.com/NumberMan1/common/ns"
	"github.com/NumberMan1/common/summer/network"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"net"
	"time"
)

// Session 用户会话类
type Session struct {
	// 当前登录的角色
	Character *Character
	// 数据库玩家信息
	DbPlayer *database.DbPlayer
	//如果网络连接断开，把消息临时缓存在这里
	buffer ns.TSQueue[proto.Message]
	Id     string
	//心跳时间
	HeartTime time.Time
	//网络连接
	Conn network.Connection
}

func NewSession(id string) *Session {
	return &Session{Id: id}
}

// Space 当前所在地图
func (s *Session) Space() *Space {
	if s.Character != nil {
		return s.Character.Space()
	}
	return nil
}

// isConnected 检测连接状态
func isConnected(conn net.Conn) bool {
	// 发送一个空的写操作来检测连接状态
	_, err := conn.Write([]byte{})
	return err == nil
}

// Send 发送消息
func (s *Session) Send(msg proto.Message) {
	if s.Conn != nil && isConnected(s.Conn.Socket()) {
		var m proto.Message
		for !s.buffer.Empty() {
			m = s.buffer.Pop()
			s.Conn.Send(m)
			if s.Character != nil {
				variable.Log.Info("补发消息",
					zap.Int32("eid", int32(s.Character.EntityId())),
					zap.Any("data", msg))
			}
		}
		s.Conn.Send(msg)
	} else {
		s.buffer.Push(msg)
	}
}

// Leave Session离开游戏
func (s *Session) Leave() {
	variable.Log.Info("Session离开", zap.String("sid", s.Id))
	if s.Character != nil {
		if s.Character.Space() != nil {
			s.Character.Space().EntityLeave(s.Character)
		}
		GetCharacterManagerInstance().RemoveCharacter(s.Character.Id())
	}

}
