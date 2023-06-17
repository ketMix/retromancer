package states

type Cursor interface {
	Enabled() bool
	Enable()
	Disable()
}
