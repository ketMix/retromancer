package game

import (
	"github.com/ketMix/retromancer/states"
)

type WorldState interface {
	Enter(w *World, ctx states.Context)
	Leave(w *World, ctx states.Context)
	Tick(w *World, ctx states.Context)
	Draw(w *World, ctx states.DrawContext)
}
