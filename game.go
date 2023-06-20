package main

import (
	"ebijam23/states"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	States  []states.State
	Manager ResourceManager
	Cursor  Cursor
}

func (g *Game) State() states.State {
	return g.States[len(g.States)-1]
}

func (g *Game) PushState(state states.State) {
	g.States = append(g.States, state)
	state.Init(states.Context{
		Manager:      &g.Manager,
		StateMachine: g,
		Cursor:       &g.Cursor,
	})
}

func (g *Game) PopState() {
	if len(g.States) == 0 {
		return
	}
	g.States = g.States[:len(g.States)-1]
}

func (g *Game) Init() error {
	g.Cursor.image = g.Manager.GetAs("images", "hand-normal", (*ebiten.Image)(nil)).(*ebiten.Image)
	g.Cursor.Enable()

	return nil
}

func (g *Game) Update() error {
	if (ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyEnter)) || inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	if state := g.State(); state != nil {
		return state.Update(states.Context{
			Manager:      &g.Manager,
			StateMachine: g,
			Cursor:       &g.Cursor,
		})
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if state := g.State(); state != nil {
		state.Draw(screen)
	}
	if g.Cursor.Enabled() {
		x, y := ebiten.CursorPosition()
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(g.Cursor.image, opts)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f TPS: %0.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 360
}
