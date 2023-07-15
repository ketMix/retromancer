package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"image/color"
	"time"
)

type WorldStateBegin struct {
	vfx resources.VFXList
}

func (w *WorldStateBegin) Enter(s *World, ctx states.Context) {
	w.vfx.SetMode(resources.Sequential)
	w.vfx.Add(&resources.Fade{
		Duration: 1 * time.Second,
	})
	x := 320.0
	y := 200.0
	// TODO: Make this text actually good.
	w.vfx.Add(&resources.Text{
		Text:         ctx.L.Get("Intro1"),
		InDuration:   1 * time.Second,
		HoldDuration: 2 * time.Second,
		OutDuration:  1 * time.Second,
		X:            x,
		Y:            y,
	})
	w.vfx.Add(&resources.Text{
		Text:         ctx.L.Get("Intro2"),
		InDuration:   1 * time.Second,
		HoldDuration: 1 * time.Second,
		OutDuration:  1 * time.Second,
		X:            x,
		Y:            y,
	})
	w.vfx.Add(&resources.Text{
		Text:         ctx.L.Get("Intro3"),
		InDuration:   1 * time.Second,
		HoldDuration: 2 * time.Second,
		OutDuration:  1 * time.Second,
		X:            x,
		Y:            y,
	})
	w.vfx.Add(&resources.Text{
		Text:         ctx.L.Get("Intro4"),
		InDuration:   1 * time.Second,
		HoldDuration: 2 * time.Second,
		OutDuration:  1 * time.Second,
		X:            x,
		Y:            y,
	})

}

func (w *WorldStateBegin) Leave(s *World, ctx states.Context) {
}

func (w *WorldStateBegin) Tick(s *World, ctx states.Context) {
	skip := false

	if s.DoPlayersShareThought(ResetThought{}) || s.DoPlayersShareThought(QuitThought{}) {
		skip = true
	}

	if !skip && !w.vfx.Empty() {
		return
	}

	s.PopState(ctx)
	s.PushState(&WorldStateLive{}, ctx)
}

func (w *WorldStateBegin) Draw(s *World, ctx states.DrawContext) {
	w.vfx.Process(ctx, nil)

	ctx.Text.SetColor(color.NRGBA{0xff, 0xff, 0xff, 0x66})
	ctx.Text.Draw(ctx.Screen, ctx.L.Get("SkipIntro"), ctx.Screen.Bounds().Max.X/2, ctx.Screen.Bounds().Max.Y-int(ctx.Text.Utils().GetLineHeight()))
}
