package game

type ImpulseSet struct {
	Move        *ImpulseMove
	Interaction Impulse
}

type Impulse interface {
}

type ImpulseMove struct {
	Direction float64
}

type ImpulseReflect struct {
	X, Y float64 // X and Y are the current cursor coordinates.
}

func (i ImpulseReflect) Cost() int {
	return 1
}

type ImpulseDeflect struct {
	X, Y float64 // X and Y are the current cursor coordinates.
}

func (i ImpulseDeflect) Cost() int {
	return 4
}
