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
	States      []states.State
	Flags       Flags
	Text        *etxt.Renderer
	Resources   ResourceManager
	Localizer   Localizer
	Cursor      Cursor
	MusicPlayer MusicPlayer
	Difficulty  states.Difficulty
}

func (g *Game) CreateContext() states.Context {
	return states.Context{
		StateMachine: g,
		Cursor:       &g.Cursor,
		MusicPlayer:  &g.MusicPlayer,
		Difficulty:   g.Difficulty,
		L:            &g.Localizer,
		R:            &g.Resources,
	}
}

func (g *Game) State() states.State {
	return g.States[len(g.States)-1]
}

func (g *Game) PushState(state states.State) {
	g.States = append(g.States, state)
	ctx := g.CreateContext()
	state.Init(ctx)
}

func (g *Game) PopState(v interface{}) {
	if len(g.States) == 0 {
		return
	}
	ctx := g.CreateContext()
	g.States[len(g.States)-1].Finalize(ctx)
	g.States = g.States[:len(g.States)-1]
	g.States[len(g.States)-1].Enter(ctx, v)
}

func (g *Game) Init() error {
	g.Cursor.image = g.Resources.GetAs("images", "hand-normal", (*ebiten.Image)(nil)).(*ebiten.Image)
	g.Cursor.Enable()

	g.Text = etxt.NewRenderer()

	return nil
}

func (g *Game) Update() error {
	if (ebiten.IsKeyPressed(ebiten.KeyAlt) && inpututil.IsKeyJustPressed(ebiten.KeyEnter)) || inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	// Check if music needs to loop, etc.
	g.MusicPlayer.Update()

	if state := g.State(); state != nil {
		return state.Update(g.CreateContext())
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if state := g.State(); state != nil {
		state.Draw(states.DrawContext{
			Screen: screen,
			Text:   g.Text,
			L:      &g.Localizer,
			R:      &g.Resources,
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
