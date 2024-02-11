package model

import (
	"github.com/NumberMan1/MMO-server/database"
	"github.com/NumberMan1/common/summer/network"
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
	"github.com/NumberMan1/common/summer/vector3"
)

// Character 角色
type Character struct {
	*Actor
	//当前角色的客户端连接
	Conn network.Connection
	//当前角色对应的数据库对象
	Data database.DbCharacter
}

func NewCharacter(position, direction vector3.Vector3) *Character {
	return &Character{Actor: NewActor(position, direction)}
}

func CharacterFromDbCharacter(character database.DbCharacter) *Character {
	c := &Character{
		Actor: NewActor(vector3.NewVector3(float64(character.X), float64(character.Y), float64(character.Z)), vector3.Zero3()),
		Data:  character,
	}
	c.Id = int(character.ID)
	c.Name = character.Name
	c.Info = &proto.NCharacter{
		Id:       int32(character.ID),
		TypeId:   int32(character.JobId),
		EntityId: 0,
		Name:     character.Name,
		Level:    int32(character.Level),
		Exp:      int64(character.Exp),
		SpaceId:  int32(character.SpaceId),
		Gold:     character.Gold,
		Entity:   nil,
		Hp:       int32(character.Hp),
		Mp:       int32(character.Mp),
	}
	c.Data = character
	c.SetSpeed(3000)
	return c
}
