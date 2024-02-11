package model

import "github.com/NumberMan1/common/summer/vector3"

type Monster struct {
	*Actor
}

func NewMonster(position, direction vector3.Vector3) *Monster {
	return &Monster{Actor: NewActor(position, direction)}
}
