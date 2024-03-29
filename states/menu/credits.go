package menu

import (
	"github.com/ketMix/retromancer/resources"

	"github.com/ketMix/retromancer/states"
	"github.com/ketMix/retromancer/states/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Credits struct {
	clickSound *resources.Sound
	items      []resources.MenuItem
	sections   []CreditSection
	jam        *resources.Sprite
	logo       *resources.Sprite
	backItem   *resources.TextItem
	overlay    game.Overlay
}

type CreditSection struct {
	Title  resources.MenuItem
	Person resources.MenuItem
	Hat    *resources.Sprite
}

func (c *Credits) Init(ctx states.Context) error {
	c.overlay.Init(ctx)
	c.clickSound = ctx.R.GetAs("sounds", "click", (*resources.Sound)(nil)).(*resources.Sound)

	c.logo = resources.NewSprite(ctx.R.GetAs("images", "logo", (*ebiten.Image)(nil)).(*ebiten.Image))
	c.logo.Centered = true
	c.logo.X = 320
	c.logo.Y = c.logo.Height()/2 - 30

	c.backItem = &resources.TextItem{
		Text: ctx.L.Get("Back"),
		X:    30,
		Y:    335,
		Callback: func() bool {
			c.clickSound.Play(1.0)
			ctx.StateMachine.PopState(nil)
			return false
		},
	}
	c.items = append(c.items, c.backItem)

	c.jam = resources.NewSprite(ctx.R.GetAs("images", "jam", (*ebiten.Image)(nil)).(*ebiten.Image))
	c.jam.X = 80
	c.jam.Y = c.logo.Y + c.logo.Height()/2 + 20

	x := 320.0
	y := c.logo.Y + c.logo.Height()/2
	section := CreditSection{
		Title: &resources.TextItem{
			Text: ctx.L.Get("Engine, Programming, Art, Lore"),
			X:    x,
			Y:    y,
		},
		Person: &resources.TextItem{
			Text:      "kettek",
			X:         x,
			Y:         y + 20,
			Underline: true,
			Callback: func() bool {
				OpenURL("https://kettek.net")
				return true
			},
		},
		Hat: resources.NewSprite(ctx.R.GetAs("images", "hat-coffee", (*ebiten.Image)(nil)).(*ebiten.Image)),
	}
	section.Hat.Scale = 2.0
	section.Hat.Centered = true
	section.Hat.X = x - 36
	section.Hat.Y = y + 19
	y += 60
	c.sections = append(c.sections, section)

	section = CreditSection{
		Title: &resources.TextItem{
			Text: ctx.L.Get("Programming, Music, SFX, Levels, Lore"),
			X:    x,
			Y:    y,
			Callback: func() bool {
				ctx.StateMachine.PushState(&Jukebox{})
				return true
			},
		},
		Person: &resources.TextItem{
			Text:      "liqMix",
			X:         x,
			Y:         y + 20,
			Underline: true,
			Callback: func() bool {
				OpenURL("https://liq.mx/")
				return true
			},
		},
		Hat: resources.NewSprite(ctx.R.GetAs("images", "hat-pep", (*ebiten.Image)(nil)).(*ebiten.Image)),
	}
	section.Hat.Scale = 2.0
	section.Hat.Centered = true
	section.Hat.X = x - 36
	section.Hat.Y = y + 19
	y += 60
	c.sections = append(c.sections, section)

	section = CreditSection{
		Title: &resources.TextItem{
			Text: ctx.L.Get("Menu Art & Logo"),
			X:    x,
			Y:    y,
		},
		Person: &resources.TextItem{
			Text:      "Amaruuk",
			X:         x,
			Y:         y + 20,
			Underline: true,
			Callback: func() bool {
				OpenURL("https://birdtooth.studio/")
				return true
			},
		},
		Hat: resources.NewSprite(ctx.R.GetAs("images", "hat-egg", (*ebiten.Image)(nil)).(*ebiten.Image)),
	}
	section.Hat.Scale = 2.0
	section.Hat.Centered = true
	section.Hat.X = x - 36
	section.Hat.Y = y + 19
	y += 60
	c.sections = append(c.sections, section)

	return nil
}

func (c *Credits) Enter(ctx states.Context, v interface{}) error {
	return nil
}

func (c *Credits) Finalize(ctx states.Context) error {
	return nil
}

func (c *Credits) Update(ctx states.Context) error {
	x, y := ebiten.CursorPosition()

	for _, b := range c.items {
		b.CheckState(float64(x), float64(y))
	}
	for _, s := range c.sections {
		s.Person.CheckState(float64(x), float64(y))
		s.Title.CheckState(float64(x), float64(y))
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		if c.jam.Hit(float64(x), float64(y)) {
			OpenURL("https://itch.io/jam/ebitengine-game-jam-2023")
			return nil
		}

		for _, b := range c.items {
			if b.Hovered() {
				if b.Activate() {
					return nil
				}
			}
		}
		for _, s := range c.sections {
			if s.Person.Hovered() {
				if s.Person.Activate() {
					return nil
				}
			}
			if s.Title.Hovered() {
				if s.Title.Activate() {
					return nil
				}
			}
		}
	}

	c.overlay.Update(ctx)

	return nil
}

func (c *Credits) Draw(ctx states.DrawContext) {
	c.logo.Draw(ctx)
	c.jam.Draw(ctx)
	for _, b := range c.items {
		b.Draw(ctx)
	}
	for _, s := range c.sections {
		s.Title.Draw(ctx)
		s.Person.Draw(ctx)
		s.Hat.Draw(ctx)
	}
	c.overlay.Draw(ctx)
}
