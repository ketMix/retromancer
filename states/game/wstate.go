package game

import (
	"ebijam23/states"
)

type WorldState interface {
	Enter(w *World, ctx states.Context)
	Leave(w *World, ctx states.Context)
	Tick(w *World, ctx states.Context)
	Draw(w *World, ctx states.DrawContext)
}
