package states

type State interface {
	Init(ctx Context) error
	Finalize(ctx Context) error
	Enter(ctx Context, v interface{}) error
	Update(ctx Context) error
	Draw(ctx DrawContext)
}

type StateMachine interface {
	PushState(state State)
	PopState(v interface{})
}
