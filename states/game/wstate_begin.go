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

func (w *WorldStateBegin) Enter(s *World) {
	w.vfx.SetMode(resources.Sequential)
	w.vfx.Add(&resources.Fade{
		Duration: 1 * time.Second,
	})
	x := 320.0
	y := 200.0
	// TODO: Make this text actually good.
	w.vfx.Add(&resources.Text{
		Text:         "Beneaf the wisper of time...",
		InDuration:   1 * time.Second,
		HoldDuration: 2 * time.Second,
		OutDuration:  1 * time.Second,
		X:            x,
		Y:            y,
	})
	w.vfx.Add(&resources.Text{
		Text:         "da forgetne echo has tarted to",
		InDuration:   1 * time.Second,
		HoldDuration: 1 * time.Second,
		OutDuration:  1 * time.Second,
		X:            x,
		Y:            y,
	})
	w.vfx.Add(&resources.Text{
		Text:         "ripl agin",
		InDuration:   1 * time.Second,
		HoldDuration: 2 * time.Second,
		OutDuration:  1 * time.Second,
		X:            x,
		Y:            y,
	})
	w.vfx.Add(&resources.Text{
		Text:         "... u wakem ups",
		InDuration:   1 * time.Second,
		HoldDuration: 2 * time.Second,
		OutDuration:  1 * time.Second,
		X:            x,
		Y:            y,
	})

}

func (w *WorldStateBegin) Leave(s *World) {
}

func (w *WorldStateBegin) Tick(s *World, ctx states.Context) {
	skip := false

	if s.DoPlayersShareThought(ResetThought{}) || s.DoPlayersShareThought(QuitThought{}) {
		skip = true
	}

	if !skip && !w.vfx.Empty() {
		return
	}

	s.PopState()
	s.PushState(&WorldStateLive{})
}

func (w *WorldStateBegin) Draw(s *World, ctx states.DrawContext) {
	w.vfx.Process(ctx, nil)

	ctx.Text.SetColor(color.NRGBA{0xff, 0xff, 0xff, 0x66})
	ctx.Text.Draw(ctx.Screen, "Press <Enter> or <Escape> to skip", ctx.Screen.Bounds().Max.X/2, ctx.Screen.Bounds().Max.Y-int(ctx.Text.Utils().GetLineHeight()))
}
