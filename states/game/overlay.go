package game

import (
	"fmt"

	"github.com/ketMix/retromancer/states"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/ketMix/retromancer/resources"
)

type Overlay struct {
	soundButton *resources.Sprite
	musicButton *resources.Sprite
	localeIcon  *resources.Sprite
}

func (o *Overlay) Init(ctx states.Context) error {
	o.soundButton = resources.NewSprite(ctx.R.GetAs("images", "sound-high", (*ebiten.Image)(nil)).(*ebiten.Image))
	o.soundButton.X = 640 - o.soundButton.Width() - 8
	o.soundButton.Y = 8

	o.musicButton = resources.NewSprite(ctx.R.GetAs("images", "music-high", (*ebiten.Image)(nil)).(*ebiten.Image))
	o.musicButton.X = 640 - o.musicButton.Width() - 8
	o.musicButton.Y = 8 + o.soundButton.Height() + 8

	o.localeIcon = resources.NewSprite(ctx.R.GetAs("images", fmt.Sprintf("flag-%s", ctx.L.Locale()), (*ebiten.Image)(nil)).(*ebiten.Image))
	o.localeIcon.X = 640 - o.localeIcon.Width() - 8
	o.localeIcon.Y = 8 + o.soundButton.Height() + 8 + o.musicButton.Height() + 8

	o.Sync(ctx)

	return nil
}

func (o *Overlay) Draw(ctx states.DrawContext) {
	o.soundButton.Draw(ctx)
	o.musicButton.Draw(ctx)
	o.localeIcon.Draw(ctx)
}

func (o *Overlay) Sync(ctx states.Context) {
	if resources.Volume == 0.0 {
		o.soundButton.SetImage(ctx.R.GetAs("images", "sound-none", (*ebiten.Image)(nil)).(*ebiten.Image))
	} else if resources.Volume == 0.5 {
		o.soundButton.SetImage(ctx.R.GetAs("images", "sound-low", (*ebiten.Image)(nil)).(*ebiten.Image))
	} else {
		o.soundButton.SetImage(ctx.R.GetAs("images", "sound-high", (*ebiten.Image)(nil)).(*ebiten.Image))
	}

	if ctx.MusicPlayer.Volume() == 0.0 {
		o.musicButton.SetImage(ctx.R.GetAs("images", "music-none", (*ebiten.Image)(nil)).(*ebiten.Image))
	} else if ctx.MusicPlayer.Volume() == 0.25 {
		o.musicButton.SetImage(ctx.R.GetAs("images", "music-low", (*ebiten.Image)(nil)).(*ebiten.Image))
	} else {
		o.musicButton.SetImage(ctx.R.GetAs("images", "music-high", (*ebiten.Image)(nil)).(*ebiten.Image))
	}

	if ctx.L.GPTIsActive() {
		o.localeIcon.SetImage(ctx.R.GetAs("images", "flag-gpt", (*ebiten.Image)(nil)).(*ebiten.Image))
	} else {
		o.localeIcon.SetImage(ctx.R.GetAs("images", fmt.Sprintf("flag-%s", ctx.L.Locale()), (*ebiten.Image)(nil)).(*ebiten.Image))
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
		if o.musicButton.Hit(float64(x), float64(y)) {
			// toggle 'em, gois
			if ctx.MusicPlayer.Volume() == 0.0 {
				ctx.MusicPlayer.SetVolume(0.25)
			} else if ctx.MusicPlayer.Volume() == 0.25 {
				ctx.MusicPlayer.SetVolume(0.5)
			} else {
				ctx.MusicPlayer.SetVolume(0.0)
			}
		}
		if o.localeIcon.Hit(float64(x), float64(y)) {
			locales := ctx.R.GetNamesWithPrefix("locales", "")
			var localeIndex int
			for i, l := range locales {
				if l == ctx.L.Locale() {
					localeIndex = i
					break
				}
			}
			if localeIndex+1 >= len(locales) {
				localeIndex = 0
			} else {
				localeIndex++
			}
			ctx.L.SetLocale(locales[localeIndex], false)
		}
		o.Sync(ctx)
	}

	return nil
}
