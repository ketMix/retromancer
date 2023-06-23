package game

import (
	"ebijam23/states"
)

type WorldState interface {
	Enter(w *World)
	Tick(w *World, ctx states.Context)
	Draw(w *World, ctx states.DrawContext)
	Leave(w *World)
}
