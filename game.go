package main

import (
	"ebijam23/states"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tinne26/etxt"
)

type Game struct {
	States    []states.State
	Text      *etxt.Renderer
	Manager   ResourceManager
	Localizer Localizer
	Cursor    Cursor
}

func (g *Game) State() states.State {
	return g.States[len(g.States)-1]
}

func (g *Game) PushState(state states.State) {
	g.States = append(g.States, state)
	state.Init(states.Context{
		Manager:      &g.Manager,
		L:            g.Localizer.Get,
		StateMachine: g,
		Cursor:       &g.Cursor,
	})
}

func (g *Game) PopState() {
	if len(g.States) == 0 {
		return
	}
	g.States[len(g.States)-1].Finalize(states.Context{
		Manager:      &g.Manager,
		L:            g.Localizer.Get,
		StateMachine: g,
		Cursor:       &g.Cursor,
	})
	g.States = g.States[:len(g.States)-1]
}

func (g *Game) Init() error {
	g.Cursor.image = g.Manager.GetAs("images", "hand-normal", (*ebiten.Image)(nil)).(*ebiten.Image)
	g.Cursor.Enable()

	g.Text = etxt.NewRenderer()

	return nil
}

func (g *Game) Update() error {
	if (ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyEnter)) || inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	if state := g.State(); state != nil {
		return state.Update(states.Context{
			Manager:      &g.Manager,
			L:            g.Localizer.Get,
			StateMachine: g,
			Cursor:       &g.Cursor,
		})
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if state := g.State(); state != nil {
		state.Draw(states.DrawContext{
			L:      g.Localizer.Get,
			Screen: screen,
			Text:   g.Text,
		})
	}
	if g.Cursor.Enabled() {
		x, y := ebiten.CursorPosition()
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x), float64(y))
		opts.GeoM.Translate(-float64(g.Cursor.image.Bounds().Dx())/2, -float64(g.Cursor.image.Bounds().Dy())/2)
		screen.DrawImage(g.Cursor.image, opts)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f TPS: %0.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 360
}
