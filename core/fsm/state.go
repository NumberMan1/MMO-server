package fsm

type IState[T any] interface {
	Fsm() *System[T]
	SetFsm(fsm *System[T])
	P() T
	SetP(p T)
	OnEnter()
	OnUpdate()
	OnLeave()
}

// State 状态基础类
type State[T any] struct {
	fsm *System[T]
}

func NewState[T any]() *State[T] {
	return &State[T]{fsm: &System[T]{}}
}

func (t *State[T]) Fsm() *System[T] {
	return t.fsm
}

func (t *State[T]) SetFsm(fsm *System[T]) {
	t.fsm = fsm
}

func (t *State[T]) P() T {
	return t.fsm.P
}

func (t *State[T]) SetP(p T) {
	t.fsm.P = p
}

func (t *State[T]) OnEnter() {

}

func (t *State[T]) OnUpdate() {

}

func (t *State[T]) OnLeave() {

}
