package model

import (
	"github.com/NumberMan1/common/summer/protocol/gen/proto"
	"github.com/NumberMan1/common/summer/vector3"
)

type Actor struct {
	*Entity
	Id    int
	Name  string
	Space *Space
	Info  *proto.NCharacter
	/*
		int32 id = 1;
		int32 type_id = 2; //角色类型
		int32 entity_id = 3;
		string name = 4;
		int32 level = 5;
		int64 exp = 6;
		int32 spaceId = 7;
		int64 gold = 8;
		NEntity entity = 9;
	*/
}

func NewActor(position, direction vector3.Vector3) *Actor {
	return &Actor{Entity: NewEntity(position, direction)}
}
