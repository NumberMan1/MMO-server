package model

// AIBase AI基础类
type AIBase struct {
	owner *Monster
}

func (b *AIBase) Owner() *Monster {
	return b.owner
}

func (b *AIBase) SetOwner(owner *Monster) {
	b.owner = owner
}

func NewBase(owner *Monster) *AIBase {
	return &AIBase{owner: owner}
}
