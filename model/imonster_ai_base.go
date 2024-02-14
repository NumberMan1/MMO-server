package model

type IAIBase interface {
	Owner() *Monster
	SetOwner(owner *Monster)
	Update()
}
