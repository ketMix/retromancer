package menu

import (
	"ebijam23/resources"
	"ebijam23/states"
	"ebijam23/states/game"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Menu struct {
	bg1, bg2         *resources.Sprite
	bg1logo, bg2logo *resources.Sprite

	play, credits, gpt *resources.TextItem
	sprites            resources.Sprites
	buttons            []*resources.TextItem
	click              *resources.Sound
	overlay            game.Overlay

	firstVfx  resources.VFXList
	secondVfx resources.VFXList
}

func (m *Menu) Init(ctx states.Context) error {
	m.overlay.Init(ctx)

	m.bg1 = resources.NewSprite(ctx.Manager.GetAs("images", "bg-1", (*ebiten.Image)(nil)).(*ebiten.Image))
	m.bg1.Hidden = true
	m.bg1logo = resources.NewSprite(ctx.Manager.GetAs("images", "bg-logo-1", (*ebiten.Image)(nil)).(*ebiten.Image))
	m.bg1logo.Hidden = true
	m.bg2 = resources.NewSprite(ctx.Manager.GetAs("images", "bg-2", (*ebiten.Image)(nil)).(*ebiten.Image))
	m.bg2logo = resources.NewSprite(ctx.Manager.GetAs("images", "bg-logo-2", (*ebiten.Image)(nil)).(*ebiten.Image))

	x := 20.0
	y := 20.0

	m.play = &resources.TextItem{
		Text: ctx.L("Play"),
		X:    x,
		Y:    y,
		Callback: func() bool {
			m.click.Play(1.0)
			ctx.StateMachine.PushState(&Lobby{})
			return true
		},
	}

	m.credits = &resources.TextItem{
		Text: ctx.L("Credits"),
		X:    x,
		Y:    y,
		Callback: func() bool {
			m.click.Play(1.0)
			ctx.StateMachine.PushState(&Credits{})
			return true
		},
	}

	m.gpt = &resources.TextItem{
		Text: "GPT Options",
		X:    x,
		Y:    y,
		Callback: func() bool {
			m.click.Play(1.0)
			ctx.StateMachine.PushState(&GPTOptions{})
			return true
		},
	}
	m.sprites = append(m.sprites, m.bg1, m.bg2)
	m.buttons = append(m.buttons, m.play, m.gpt, m.credits)

	m.click = ctx.Manager.GetAs("sounds", "click", (*resources.Sound)(nil)).(*resources.Sound)

	m.firstVfx.Add(&resources.Fade{
		Alpha:        1.0,
		Duration:     1 * time.Second,
		ApplyToImage: true,
	})
	m.secondVfx.Add(&resources.Fade{
		Alpha:        1.0,
		Duration:     1 * time.Second,
		ApplyToImage: true,
	})

	return nil
}

func (m *Menu) Finalize(ctx states.Context) error {
	return nil
}

func (m *Menu) Enter(ctx states.Context) error {
	ctx.MusicPlayer.Play(ctx.Manager.GetAs("songs", "title-menu", (*resources.Song)(nil)).(states.Song))

	return nil
}

func (m *Menu) Update(ctx states.Context) error {
	x, y := ebiten.CursorPosition()

	for _, m := range m.buttons {
		m.CheckState(float64(x), float64(y))
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		for _, m := range m.buttons {
			if m.Hovered() {
				if m.Activate() {
					return nil
				}
			}
		}
	}

	m.overlay.Update(ctx)

	return nil
}

func (m *Menu) Draw(ctx states.DrawContext) {
	if !m.firstVfx.Empty() {
		m.bg2logo.Draw(ctx)
		m.firstVfx.Process(ctx, nil)
		return
	}
	m.bg2.Draw(ctx)

	for _, sprite := range m.sprites {
		sprite.Draw(ctx)
	}

	x := 40
	y := m.bg2logo.Height() + 25
	for _, button := range m.buttons {
		button.X = float64(x)
		button.Y = y
		x += 80
		button.Draw(ctx)
	}

	m.overlay.Draw(ctx)

	m.secondVfx.Process(ctx, nil)
	m.bg2logo.Draw(ctx)
}
