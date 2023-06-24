package menu

import (
	"ebijam23/resources"
	"ebijam23/states"
	"image/color"
	"time"
)

type Intro struct {
	vfx resources.VFXList
}

func (i *Intro) Init(ctx states.Context) error {
	ctx.MusicPlayer.Play(ctx.Manager.GetAs("songs", "title-intro", (*resources.Song)(nil)).(states.Song))
	return nil
}
func (i *Intro) Enter(ctx states.Context) error {
	i.vfx.SetMode(resources.Sequential)
	i.vfx.Add(&resources.Fade{
		Duration: 1 * time.Second,
	})
	x := 320.0
	y := 200.0
	// TODO: Make this text actually good.
	i.vfx.Add(&resources.Text{
		Text:         ctx.L("MenuIntro1"),
		InDuration:   1 * time.Second,
		HoldDuration: 2 * time.Second,
		OutDuration:  1 * time.Second,
		X:            x,
		Y:            y,
	})
	i.vfx.Add(&resources.Text{
		Text:         ctx.L("MenuIntro2"),
		InDuration:   1 * time.Second,
		HoldDuration: 1 * time.Second,
		OutDuration:  1 * time.Second,
		X:            x,
		Y:            y,
	})
	i.vfx.Add(&resources.Text{
		Text:         ctx.L("MenuIntro3"),
		InDuration:   1 * time.Second,
		HoldDuration: 2 * time.Second,
		OutDuration:  1 * time.Second,
		X:            x,
		Y:            y,
	})
	i.vfx.Add(&resources.Text{
		Text:         "",
		InDuration:   0,
		HoldDuration: 2 * time.Second,
		OutDuration:  0,
		X:            x,
		Y:            y,
	})
	return nil
}

func (i *Intro) Finalize(ctx states.Context) error {
	return nil
}

func (i *Intro) Update(ctx states.Context) error {
	if !i.vfx.Empty() {
		return nil
	}
	ctx.StateMachine.PopState()
	return nil
}

func (i *Intro) Draw(ctx states.DrawContext) {
	i.vfx.Process(ctx, nil)

	ctx.Text.SetColor(color.NRGBA{0xff, 0xff, 0xff, 0x66})
	ctx.Text.Draw(ctx.Screen, ctx.L("SkipIntro"), ctx.Screen.Bounds().Max.X/2, ctx.Screen.Bounds().Max.Y-int(ctx.Text.Utils().GetLineHeight()))
}
