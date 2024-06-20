package model

import (
	"github.com/NumberMan1/MMO-server/database"
	"github.com/NumberMan1/common/ns/singleton"
	"github.com/google/uuid"
	"sync"
	"time"
)

type SessionManager struct {
	//id-session
	dict sync.Map
}

var instance singleton.Singleton

func GetSessionManagerInstance() *SessionManager {
	r, _ := singleton.GetOrDo[*SessionManager](&instance, func() (*SessionManager, error) {
		return &SessionManager{
			dict: sync.Map{},
		}, nil
	})
	go r.startTimer()
	return r
}

// NewSession 给登录玩家分配Session
func (sm *SessionManager) NewSession(dbPlayer *database.DbPlayer) *Session {
	session := NewSession(uuid.New().String())
	session.DbPlayer = dbPlayer
	sm.dict.Store(session.Id, session)
	return session
}

func (sm *SessionManager) GetSession(sessionId string) *Session {
	if value, ok := sm.dict.Load(sessionId); ok {
		return value.(*Session)
	}
	return nil
}

func (sm *SessionManager) Remove(sessionId string) {
	sm.dict.Delete(sessionId)
}

// GetPlayerSession 查找玩家的Session
func (sm *SessionManager) GetPlayerSession(playerId uint) *Session {
	var result *Session
	sm.dict.Range(func(key, value interface{}) bool {
		session := value.(*Session)
		if session.DbPlayer.ID == playerId {
			result = session
			return false
		}
		return true
	})
	return result
}

func (sm *SessionManager) startTimer() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		sm.CheckSessions()
	}
}

// CheckSessions 检查Session是否还有效
func (sm *SessionManager) CheckSessions() {
	now := time.Now()
	sm.dict.Range(func(key, value interface{}) bool {
		session := value.(*Session)
		//如果没有回应心跳超过10秒就离开
		if now.Sub(session.HeartTime) > time.Second*10 {
			session.Leave()
		}
		return true
	})
}
