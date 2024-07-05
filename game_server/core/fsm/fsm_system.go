package fsm

type System[T any] struct {
	dict           map[string]IState[T]
	CurrentStateId string
	CurrentState   IState[T]
	//共享参数Params
	P T
}

func NewSystem[T any](param T) *System[T] {
	return &System[T]{P: param, dict: map[string]IState[T]{}}
}

// AddState 添加状态
func (s *System[T]) AddState(id string, state IState[T]) {
	if len(s.CurrentStateId) == 0 {
		s.CurrentStateId = id
		s.CurrentState = state
	}
	s.dict[id] = state
	state.SetFsm(s)
}

// RemoveState 移除状态
func (s *System[T]) RemoveState(id string) {
	_, ok := s.dict[id]
	if ok {
		delete(s.dict, id)
	}
}

// ChangeState 切换状态
func (s *System[T]) ChangeState(id string) {
	if s.CurrentStateId == id {
		return
	}
	v, ok := s.dict[id]
	if !ok {
		return
	}
	if s.CurrentState != nil {
		s.CurrentState.OnLeave()
	}
	s.CurrentStateId = id
	s.CurrentState = v
	s.CurrentState.OnEnter()
}

func (s *System[T]) Update() {
	if s.CurrentState != nil {
		s.CurrentState.OnUpdate()
	}
}
