package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"image/color"

	"github.com/tinne26/etxt"
)

type WorldStateDead struct {
}

func (w *WorldStateDead) Enter(s *World, ctx states.Context) {
}

func (w *WorldStateDead) Leave(s *World, ctx states.Context) {
}

func (w *WorldStateDead) Tick(s *World, ctx states.Context) {
	for _, actor := range s.activeMap.actors {
		actor.Update()
	}

	// Might as well still process particles.
	for _, p := range s.activeMap.particles {
		p.Update()
	}

	if s.DoPlayersShareThought(ResetThought{}) {
		s.PopState(ctx)
		s.PushState(&WorldStateLive{}, ctx)
		s.ResetActiveMap(ctx)
	} else if s.DoPlayersShareThought(QuitThought{}) {
		ctx.StateMachine.PopState()
	}
}

func (w *WorldStateDead) Draw(s *World, ctx states.DrawContext) {
	s.activeMap.Draw(ctx)

	ctx.Text.SetAlign(etxt.YCenter | etxt.XCenter)
	x := ctx.Screen.Bounds().Max.X / 2
	y := float64(ctx.Screen.Bounds().Max.Y / 2)
	y -= ctx.Text.Utils().GetLineHeight() / 2
	// Death title
	{
		ctx.Text.SetScale(2.0)
		ctx.Text.SetColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
		resources.DrawTextOutline(ctx.Text, ctx.Screen, ctx.L("Morte"), x, int(y), 2)
		ctx.Text.SetColor(color.RGBA{0xff, 0x00, 0x00, 0xff})
		ctx.Text.Draw(ctx.Screen, ctx.L("Morte"), x, int(y))
	}
	y += ctx.Text.Utils().GetLineHeight()
	// Subtitles
	{
		ctx.Text.SetScale(1.0)
		ctx.Text.SetColor(color.Black)
		resources.DrawTextOutline(ctx.Text, ctx.Screen, ctx.L("ResetRoom"), x, int(y), 1)
		ctx.Text.SetColor(color.White)
		ctx.Text.Draw(ctx.Screen, ctx.L("ResetRoom"), x, int(y))
		y += ctx.Text.Utils().GetLineHeight()
		ctx.Text.SetColor(color.Black)
		resources.DrawTextOutline(ctx.Text, ctx.Screen, ctx.L("Quit"), x, int(y), 1)
		ctx.Text.SetColor(color.White)
		ctx.Text.Draw(ctx.Screen, ctx.L("Quit"), ctx.Screen.Bounds().Max.X/2, int(y))
	}
}
