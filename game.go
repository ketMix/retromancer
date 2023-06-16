package main

import (
	"ebijam23/states"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	States  []states.State
	Manager ResourceManager
}

func (g *Game) State() states.State {
	return g.States[len(g.States)-1]
}

func (g *Game) PushState(state states.State) {
	g.States = append(g.States, state)
	state.Init(states.Context{
		Manager:      &g.Manager,
		StateMachine: g,
	})
}

func (g *Game) PopState() {
	if len(g.States) == 0 {
		return
	}
	g.States = g.States[:len(g.States)-1]
}

func (g *Game) Update() error {
	if state := g.State(); state != nil {
		return state.Update(states.Context{
			Manager:      &g.Manager,
			StateMachine: g,
		})
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if state := g.State(); state != nil {
		state.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 360
}
