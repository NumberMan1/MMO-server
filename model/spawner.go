package model

import (
	"github.com/NumberMan1/MMO-server/config/define"
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/common/logger"
	"github.com/NumberMan1/common/summer/timeunit"
	"regexp"
	"strconv"
)

func parsePoint(text string) *vector3.Vector3 {
	pattern := `\d+`
	compile := regexp.MustCompile(pattern)
	match := compile.FindAllString(text, -1)
	if len(match) != 0 {
		x, _ := strconv.ParseFloat(match[0], 64)
		y, _ := strconv.ParseFloat(match[1], 64)
		z, _ := strconv.ParseFloat(match[2], 64)
		return vector3.NewVector3(x, y, z)
	}
	return vector3.Zero3()
}

type Spawner struct {
	Define     *define.SpawnDefine
	Space      *Space
	mon        *Monster
	pos        *vector3.Vector3 //刷怪位置
	dir        *vector3.Vector3 //刷怪方向
	reviving   bool             //是否正在复活倒计时
	reviveTime float64          //复活时间
}

func NewSpawner(define *define.SpawnDefine, space *Space) *Spawner {
	s := &Spawner{
		Define: define,
		Space:  space,
		//pos:    parsePoint(define.Pos),
		//dir:    parsePoint(define.Dir),
		pos: vector3.NewVector3(float64(define.Pos[0]), float64(define.Pos[1]), float64(define.Pos[2])),
		dir: vector3.NewVector3(float64(define.Dir[0]), float64(define.Dir[1]), float64(define.Dir[2])),
	}
	logger.SLCDebug("New Spawner:场景[%v],坐标[%v],单位类型[%v]，周期[%v]秒", space.Name, s.Pos, define.TID, define.Period)
	s.spawn()
	return s
}

func (s *Spawner) Pos() *vector3.Vector3 {
	return s.pos
}

func (s *Spawner) Dir() *vector3.Vector3 {
	return s.dir
}

func (s *Spawner) spawn() {
	s.mon = s.Space.MonsterManager.Create(s.Define.TID, s.Define.Level, s.Pos(), s.Dir())
}

func (s *Spawner) Update() {
	if s.mon != nil && s.mon.IsDeath() && !s.reviving {
		s.reviveTime = timeunit.Time + float64(s.Define.Period)
		s.reviving = true
	}
	if s.reviving && s.reviveTime <= timeunit.Time {
		s.reviving = false
		if s.mon != nil {
			s.mon.Revive()
		}
	}
}
