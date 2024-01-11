package menu

import (
	"github.com/ketMix/retromancer/states"
	"github.com/ketMix/retromancer/states/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/ketMix/retromancer/resources"
)

type Jukebox struct {
	click    *resources.Sound
	items    []resources.MenuItem
	logo     *resources.Sprite
	backItem *resources.TextItem
	overlay  game.Overlay
	sprites  []*resources.Sprite
}

func (j *Jukebox) Init(ctx states.Context) error {
	j.overlay.Init(ctx)
	j.click = ctx.R.GetAs("sounds", "click", (*resources.Sound)(nil)).(*resources.Sound)

	j.logo = resources.NewSprite(ctx.R.GetAs("images", "logo", (*ebiten.Image)(nil)).(*ebiten.Image))
	j.logo.Centered = true
	j.logo.X = 320
	j.logo.Y = j.logo.Height()/2 - 30

	// Back button
	j.backItem = &resources.TextItem{
		Text: ctx.L.Get("Back"),
		X:    30,
		Y:    335,
		Callback: func() bool {
			j.click.Play(1.0)
			ctx.StateMachine.PopState(nil)
			return false
		},
	}
	j.items = append(j.items, j.backItem)

	// Create the song items
	x := 320.0
	y := j.logo.Y + j.logo.Height()/2

	// Get all songs
	songs := ctx.R.GetNamesWithPrefix("songs", "")
	for i, song := range songs {
		playSongFn := func(t *resources.TextItem) bool {
			j.click.Play(1.0)
			ctx.MusicPlayer.Play(ctx.R.Get("songs", t.Text).(*resources.Song))
			return true
		}
		item := &resources.TextItem{
			Text:            song,
			X:               x,
			Y:               y + float64(i)*20,
			Underline:       true,
			SelfRefCallback: &playSongFn,
		}
		j.items = append(j.items, item)
	}

	// Create the sprites
	y = j.logo.Y + j.logo.Height()/2
	batBossSprite := resources.NewAnimatedSpriteFromName(ctx.R, "bat-boss-alive")
	batBossSprite.Framerate = 5
	batBossSprite.Loop = true
	batBossSprite.SetXY(x-250, y)

	skellBossHeadSprite := resources.NewAnimatedSpriteFromName(ctx.R, "skell-boss-head-alive")
	skellBossHeadSprite.Framerate = 5
	skellBossHeadSprite.Loop = true
	skellBossHeadSprite.SetXY(x-150, y+100)

	batBossRedSprite := resources.NewAnimatedSpriteFromName(ctx.R, "bat-boss-red-alive")
	batBossRedSprite.Framerate = 5
	batBossRedSprite.Loop = true
	batBossRedSprite.SetXY(x+200, y)

	skellBossBodySprite := resources.NewAnimatedSpriteFromName(ctx.R, "skell-boss-body-alive")
	skellBossBodySprite.Framerate = 5
	skellBossBodySprite.Loop = true
	skellBossBodySprite.SetXY(x+100, y+65)

	j.sprites = append(j.sprites,
		batBossSprite,
		batBossRedSprite,
		skellBossBodySprite,
		skellBossHeadSprite,
	)
	return nil
}

func (j *Jukebox) Enter(ctx states.Context, v interface{}) error {
	return nil
}

func (j *Jukebox) Finalize(ctx states.Context) error {
	return nil
}

func (j *Jukebox) Update(ctx states.Context) error {
	x, y := ebiten.CursorPosition()

	for _, s := range j.sprites {
		s.Update()
	}
	for _, b := range j.items {
		b.CheckState(float64(x), float64(y))
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		for _, b := range j.items {
			if b.Hovered() {
				if b.Activate() {
					return nil
				}
			}
		}
	}

	j.overlay.Update(ctx)

	return nil
}

func (j *Jukebox) Draw(ctx states.DrawContext) {
	j.logo.Draw(ctx)

	for _, s := range j.sprites {
		s.Draw(ctx)
	}

	for _, b := range j.items {
		b.Draw(ctx)
	}
	j.overlay.Draw(ctx)
}
