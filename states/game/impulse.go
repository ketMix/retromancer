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

type ImpulseDeflect struct {
	X, Y float64 // X and Y are the current cursor coordinates.
}
