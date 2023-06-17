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
	return 4
}

type ActorActions struct {
	Actor   Actor
	Actions []Action
}

type Action interface {
}

type ActionMove struct {
	X, Y float64
}

type ActionReflect struct {
	X, Y float64
}

type ActionDeflect struct {
	Direction float64
}
