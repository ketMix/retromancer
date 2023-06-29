package menu

import (
	"ebijam23/resources"
	"ebijam23/states"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Intro struct {
	vfx resources.VFXList
}

func (i *Intro) Init(ctx states.Context) error {
	return nil
}

func (i *Intro) Enter(ctx states.Context, v interface{}) error {
	ctx.MusicPlayer.Play(ctx.Manager.GetAs("songs", "title-intro", (*resources.Song)(nil)).(states.Song))

	i.vfx.SetMode(resources.Sequential)
	i.vfx.Add(&resources.Fade{
		Duration: 1 * time.Second,
	})
	x := 275.0
	y := 150.0
	xOffset := 50.0
	yOffset := 30.0

	// There's a little difference in timing if loading intro first
	// which we do when the key is invalid
	introHold := 900 * time.Millisecond
	if !ctx.CheckGPTKey() {
		introHold = 50 * time.Millisecond
	}

	i.vfx.Add(&resources.Text{
		Text:         "",
		HoldDuration: introHold,
	})
	i.vfx.Add(&resources.Text{
		Text:         ctx.L("MenuIntro1"),
		InDuration:   1450 * time.Millisecond,
		HoldDuration: 800 * time.Millisecond,
		OutDuration:  1350 * time.Millisecond,
		X:            x,
		Y:            y,
	})
	x += xOffset
	y += yOffset
	i.vfx.Add(&resources.Text{
		Text:         ctx.L("MenuIntro2"),
		InDuration:   1450 * time.Millisecond,
		HoldDuration: 800 * time.Millisecond,
		OutDuration:  1350 * time.Millisecond,
		X:            x,
		Y:            y,
	})
	x += xOffset
	y += yOffset
	i.vfx.Add(&resources.Text{
		Text:         ctx.L("MenuIntro3"),
		InDuration:   1450 * time.Millisecond,
		HoldDuration: 1250 * time.Millisecond,
		OutDuration:  3000 * time.Millisecond,
		X:            x,
		Y:            y,
	})
	return nil
}

func (i *Intro) Finalize(ctx states.Context) error {
	return nil
}

func (i *Intro) Update(ctx states.Context) error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeyEscape) || i.vfx.Empty() {
		ctx.StateMachine.PopState(nil)
	}

	return nil
}

func (i *Intro) Draw(ctx states.DrawContext) {
	i.vfx.Process(ctx, nil)

	ctx.Text.SetColor(color.NRGBA{0xff, 0xff, 0xff, 0x66})
	ctx.Text.Draw(ctx.Screen, ctx.L("SkipIntro"), ctx.Screen.Bounds().Max.X/2, ctx.Screen.Bounds().Max.Y-int(ctx.Text.Utils().GetLineHeight()))
}
