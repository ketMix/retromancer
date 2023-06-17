package states

type Context struct {
	Manager      ResourceManager
	StateMachine StateMachine
	Cursor       Cursor
}
