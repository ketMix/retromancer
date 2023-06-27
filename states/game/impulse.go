package game

type ImpulseSet struct {
	Move        *ImpulseMove
	Interaction Impulse
}

type Impulse interface {
	Cost() int
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
	return 2
}

type ImpulseShield struct {
}

func (i ImpulseShield) Cost() int {
	return 4
}

type ImpulseShoot struct {
	X, Y float64 // X and Y are the current cursor coordinates.
}

func (i ImpulseShoot) Cost() int {
	return 6
}
