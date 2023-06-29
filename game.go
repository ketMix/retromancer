package main

import (
	"ebijam23/states"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tinne26/etxt"
)

type Game struct {
	States      []states.State
	Flags       Flags
	Text        *etxt.Renderer
	Manager     ResourceManager
	Localizer   Localizer
	Cursor      Cursor
	MusicPlayer MusicPlayer
}

func (g *Game) State() states.State {
	return g.States[len(g.States)-1]
}

func (g *Game) PushState(state states.State) {
	g.States = append(g.States, state)
	ctx := states.Context{
		Manager:      &g.Manager,
		L:            g.Localizer.Get,
		Locale:       g.Localizer.Locale,
		SetLocale:    g.Localizer.SetLocale,
		SetGPTStyle:  g.Localizer.SetGPTStyle,
		CheckGPTKey:  g.Localizer.CheckGPTKey,
		GPTIsActive:  g.Localizer.GPTIsActive,
		StateMachine: g,
		Cursor:       &g.Cursor,
		MusicPlayer:  &g.MusicPlayer,
	}
	state.Init(ctx)
}

func (g *Game) PopState(v interface{}) {
	if len(g.States) == 0 {
		return
	}
	ctx := states.Context{
		Manager:      &g.Manager,
		L:            g.Localizer.Get,
		Locale:       g.Localizer.Locale,
		SetLocale:    g.Localizer.SetLocale,
		SetGPTStyle:  g.Localizer.SetGPTStyle,
		CheckGPTKey:  g.Localizer.CheckGPTKey,
		GPTIsActive:  g.Localizer.GPTIsActive,
		StateMachine: g,
		Cursor:       &g.Cursor,
		MusicPlayer:  &g.MusicPlayer,
	}
	g.States[len(g.States)-1].Finalize(ctx)
	g.States = g.States[:len(g.States)-1]
	g.States[len(g.States)-1].Enter(ctx, v)
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

	// Check if music needs to loop, etc.
	g.MusicPlayer.Update()

	if state := g.State(); state != nil {
		return state.Update(states.Context{
			Manager:      &g.Manager,
			L:            g.Localizer.Get,
			Locale:       g.Localizer.Locale,
			SetLocale:    g.Localizer.SetLocale,
			SetGPTStyle:  g.Localizer.SetGPTStyle,
			CheckGPTKey:  g.Localizer.CheckGPTKey,
			GPTIsActive:  g.Localizer.GPTIsActive,
			StateMachine: g,
			Cursor:       &g.Cursor,
			MusicPlayer:  &g.MusicPlayer,
		})
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if state := g.State(); state != nil {
		state.Draw(states.DrawContext{
			L:         g.Localizer.Get,
			Locale:    g.Localizer.Locale,
			SetLocale: g.Localizer.SetLocale,
			Screen:    screen,
			Text:      g.Text,
		})
	}
	if g.Cursor.Enabled() {
		x, y := ebiten.CursorPosition()
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x), float64(y))
		opts.GeoM.Translate(-float64(g.Cursor.image.Bounds().Dx())/2, -float64(g.Cursor.image.Bounds().Dy())/2)
		screen.DrawImage(g.Cursor.image, opts)
	}
	//ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f TPS: %0.2f", ebiten.ActualFPS(), ebiten.ActualTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 360
}
