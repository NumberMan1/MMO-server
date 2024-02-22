package model

import (
	"github.com/NumberMan1/MMO-server/core/vector3"
	"github.com/NumberMan1/MMO-server/database"
	"github.com/NumberMan1/common/summer/network"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
)

// Character 角色
type Character struct {
	*Actor
	//当前角色的客户端连接
	Conn network.Connection
	//当前角色对应的数据库对象
	Data *database.DbCharacter
}

func NewCharacter(dbCharacter *database.DbCharacter) *Character {
	c := &Character{
		Actor: NewActor(proto.EntityType_Character, dbCharacter.JobId, dbCharacter.Level,
			vector3.NewVector3(float64(dbCharacter.X), float64(dbCharacter.Y), float64(dbCharacter.Z)),
			vector3.Zero3()),
	}
	c.SetId(int(dbCharacter.ID))
	c.SetName(dbCharacter.Name)
	c.Info().Id = int32(dbCharacter.ID)
	c.Info().Name = dbCharacter.Name
	c.Info().Tid = int32(dbCharacter.JobId) //单位类型
	c.Info().Level = int32(dbCharacter.Level)
	c.Info().Exp = int64(dbCharacter.Exp)
	c.Info().SpaceId = int32(dbCharacter.SpaceId)
	c.Info().Gold = dbCharacter.Gold
	c.Info().Hp = float32(dbCharacter.Hp)
	c.Info().Mp = float32(dbCharacter.Mp)
	c.Data = dbCharacter
	return c
}

// CharacterId 玩家角色唯一ID
func (c *Character) CharacterId() int {
	return int(c.Data.ID)
}
