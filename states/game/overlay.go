package game

import (
	"ebijam23/resources"
	"ebijam23/states"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Overlay struct {
	soundButton *resources.Sprite
}

func (o *Overlay) Init(ctx states.Context) error {
	o.soundButton = resources.NewSprite(ctx.Manager.GetAs("images", "sound-high", (*ebiten.Image)(nil)).(*ebiten.Image))
	o.soundButton.X = 640 - o.soundButton.Width() - 8
	o.soundButton.Y = 8

	o.Sync(ctx)

	return nil
}

func (o *Overlay) Draw(ctx states.DrawContext) {
	o.soundButton.Draw(ctx)
}

func (o *Overlay) Sync(ctx states.Context) {
	if resources.Volume == 0.0 {
		o.soundButton.SetImage(ctx.Manager.GetAs("images", "sound-none", (*ebiten.Image)(nil)).(*ebiten.Image))
	} else if resources.Volume == 0.5 {
		o.soundButton.SetImage(ctx.Manager.GetAs("images", "sound-low", (*ebiten.Image)(nil)).(*ebiten.Image))
	} else {
		o.soundButton.SetImage(ctx.Manager.GetAs("images", "sound-high", (*ebiten.Image)(nil)).(*ebiten.Image))
	}
}

func (o *Overlay) Update(ctx states.Context) error {
	x, y := ebiten.CursorPosition()
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if o.soundButton.Hit(float64(x), float64(y)) {
			// toggle 'em, bois
			if resources.Volume == 0.0 {
				resources.Volume = 0.5
			} else if resources.Volume == 0.5 {
				resources.Volume = 1.0
			} else {
				resources.Volume = 0.0
			}
		}
		o.Sync(ctx)
	}

	return nil
}
