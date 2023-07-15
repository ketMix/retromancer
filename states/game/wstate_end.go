package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
)

type WorldStateEnd struct {
	vfx       resources.VFXList
	npcSprite *resources.Sprite
}

func (w *WorldStateEnd) Enter(s *World, ctx states.Context) {
	w.npcSprite = resources.NewAnimatedSprite([]*ebiten.Image{
		ctx.R.GetAs("images", "npc-saved1", (*ebiten.Image)(nil)).(*ebiten.Image),
		ctx.R.GetAs("images", "npc-saved2", (*ebiten.Image)(nil)).(*ebiten.Image),
		ctx.R.GetAs("images", "npc-saved3", (*ebiten.Image)(nil)).(*ebiten.Image),
		ctx.R.GetAs("images", "npc-saved4", (*ebiten.Image)(nil)).(*ebiten.Image),
	})
	w.npcSprite.Loop = true
	w.npcSprite.Centered = true
	w.npcSprite.Framerate = 4

	w.vfx.SetMode(resources.Sequential)
	w.vfx.Add(&resources.Fade{
		Duration: 1 * time.Second,
	})
	x := 320.0
	y := 200.0

	if len(s.savedNPCs) < 14 {
		// Bad ending
		w.vfx.Add(&resources.Text{
			Text:         ctx.L.Get("Outro1"),
			Scale:        1.0,
			InDuration:   1 * time.Second,
			HoldDuration: 2 * time.Second,
			OutDuration:  1 * time.Second,
			X:            x,
			Y:            y,
		})
		w.vfx.Add(&resources.Text{
			Text:         ctx.L.Get("Outro2"),
			Scale:        1.0,
			InDuration:   1 * time.Second,
			HoldDuration: 2 * time.Second,
			OutDuration:  1 * time.Second,
			X:            x,
			Y:            y,
		})
		w.vfx.Add(&resources.Text{
			Text:         ctx.L.Get("Outro3"),
			Scale:        1.0,
			InDuration:   1 * time.Second,
			HoldDuration: 2 * time.Second,
			OutDuration:  1 * time.Second,
			X:            x,
			Y:            y,
		})
		w.vfx.Add(&resources.Text{
			Text:         ctx.L.Get("Outro4"),
			Scale:        1.0,
			InDuration:   1 * time.Second,
			HoldDuration: 2 * time.Second,
			OutDuration:  1 * time.Second,
			X:            x,
			Y:            y,
		})
		w.vfx.Add(&resources.Text{
			Text:         ctx.L.Get("Outro5"),
			InDuration:   1 * time.Second,
			HoldDuration: 3 * time.Second,
			OutDuration:  1 * time.Second,
			X:            x,
			Y:            y,
		})
		w.vfx.Add(&resources.Text{
			Text:         ctx.L.Get("Outro6"),
			InDuration:   1 * time.Second,
			HoldDuration: 3 * time.Second,
			OutDuration:  1 * time.Second,
			X:            x,
			Y:            y,
		})
		w.vfx.Add(&resources.Text{
			Text:         ctx.L.Get("Outro7"),
			InDuration:   3 * time.Second,
			HoldDuration: 2 * time.Second,
			OutDuration:  4 * time.Second,
			X:            x,
			Y:            y,
		})

	} else {
		// Good ending
		w.vfx.Add(&resources.Text{
			Text:         ctx.L.Get("Good1"),
			Scale:        1.0,
			InDuration:   1 * time.Second,
			HoldDuration: 2 * time.Second,
			OutDuration:  1 * time.Second,
			X:            x,
			Y:            y,
		})
		w.vfx.Add(&resources.Text{
			Text:         ctx.L.Get("Good2"),
			Scale:        1.0,
			InDuration:   1 * time.Second,
			HoldDuration: 3 * time.Second,
			OutDuration:  1 * time.Second,
			X:            x,
			Y:            y,
		})
		w.vfx.Add(&resources.Text{
			Text:         ctx.L.Get("Good3"),
			Scale:        1.0,
			InDuration:   1 * time.Second,
			HoldDuration: 2 * time.Second,
			OutDuration:  1 * time.Second,
			X:            x,
			Y:            y,
		})
		w.vfx.Add(&resources.Text{
			Text:         ctx.L.Get("Good4"),
			Scale:        1.0,
			InDuration:   1 * time.Second,
			HoldDuration: 2 * time.Second,
			OutDuration:  1 * time.Second,
			X:            x,
			Y:            y,
		})
		w.vfx.Add(&resources.Text{
			Text:         ctx.L.Get("Good5"),
			Scale:        1.0,
			InDuration:   1 * time.Second,
			HoldDuration: 2 * time.Second,
			OutDuration:  1 * time.Second,
			X:            x,
			Y:            y,
		})
	}
}

func (w *WorldStateEnd) Leave(s *World, ctx states.Context) {
}

func (w *WorldStateEnd) Tick(s *World, ctx states.Context) {
	skip := false

	if !w.vfx.Empty() {
		return
	} else {
		w.npcSprite.Update()
		if s.DoPlayersShareThought(ResetThought{}) || s.DoPlayersShareThought(QuitThought{}) {
			skip = true
		}
	}

	if !skip {
		return
	}

	ctx.StateMachine.PopState(len(s.savedNPCs) >= 14)
}

func (w *WorldStateEnd) Draw(s *World, ctx states.DrawContext) {
	w.vfx.Process(ctx, nil)

	if w.vfx.Empty() {
		ctx.Text.SetScale(1.0)
		ctx.Text.SetColor(color.White)
		ctx.Text.SetAlign(etxt.XCenter | etxt.YCenter)
		ctx.Text.Draw(ctx.Screen, ctx.L.Get("SavedNPCs"), ctx.Screen.Bounds().Dx()/2, 100)
		for i := 0; i < len(s.savedNPCs); i++ {
			x := (i % 13) * 20
			y := (i / 13) * 20
			w.npcSprite.X = float64(200 + x)
			w.npcSprite.Y = float64(130 + y)
			w.npcSprite.Draw(ctx)
		}

		ctx.Text.SetColor(color.NRGBA{0xff, 0xff, 0xff, 0x66})
		ctx.Text.Draw(ctx.Screen, ctx.L.Get("SkipOutro"), ctx.Screen.Bounds().Max.X/2, ctx.Screen.Bounds().Max.Y-int(ctx.Text.Utils().GetLineHeight()))

		return
	}
}
