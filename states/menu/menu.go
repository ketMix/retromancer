package menu

import (
	"ebijam23/resources"
	"ebijam23/states"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Menu struct {
	logo, sp, mp, quit, ballpit *resources.Sprite
	sprites                     resources.Sprites
}

func (m *Menu) Init(ctx states.Context) error {
	x := 320.0
	y := 25.0
	m.logo = resources.NewSprite(ctx.Manager.GetAs("images", "logo", (*ebiten.Image)(nil)).(*ebiten.Image))
	m.logo.X = x - m.logo.Width()/2
	m.logo.Y = y
	y += m.logo.Height() + 100
	m.sp = resources.NewSprite(ctx.Manager.GetAs("images", "sp", (*ebiten.Image)(nil)).(*ebiten.Image))
	m.sp.X = x - m.sp.Width()/2
	m.sp.Y = y
	y += m.sp.Height() + 16
	m.mp = resources.NewSprite(ctx.Manager.GetAs("images", "mp", (*ebiten.Image)(nil)).(*ebiten.Image))
	m.mp.X = x - m.mp.Width()/2
	m.mp.Y = y
	y += m.mp.Height() + 16
	m.quit = resources.NewSprite(ctx.Manager.GetAs("images", "quit", (*ebiten.Image)(nil)).(*ebiten.Image))
	m.quit.X = x - m.quit.Width()/2
	m.quit.Y = y
	y += m.quit.Height() + 16
	m.ballpit = resources.NewSprite(ctx.Manager.GetAs("images", "ballpit", (*ebiten.Image)(nil)).(*ebiten.Image))
	m.ballpit.X = x - m.quit.Width()/2
	m.ballpit.Y = y
	y += m.ballpit.Height() + 16
	m.sprites = append(m.sprites, m.sp, m.mp, m.quit, m.ballpit)
	return nil
}

func (m *Menu) Finalize(ctx states.Context) error {
	return nil
}

func (m *Menu) Update(ctx states.Context) error {
	x, y := ebiten.CursorPosition()

	for _, sprite := range m.sprites {
		if sprite.Hit(float64(x), float64(y)) {
			sprite.Options.ColorScale.Reset()
			if sprite == m.quit {
				sprite.Options.ColorScale.Scale(1.0, 0.25, 0.25, 1.0)
			} else {
				sprite.Options.ColorScale.Scale(0.25, 0.75, 1.0, 1.0)
			}
		} else {
			sprite.Options.ColorScale.Reset()
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if m.sp.Hit(float64(x), float64(y)) {
			ctx.StateMachine.PushState(&SinglePlayer{})
		} else if m.mp.Hit(float64(x), float64(y)) {
			ctx.StateMachine.PushState(&MultiPlayer{})
		} else if m.quit.Hit(float64(x), float64(y)) {
			return states.ErrQuitGame
		} else if m.ballpit.Hit(float64(x), float64(y)) {
			ctx.StateMachine.PushState(&Ballpit{})
		}

	}
	return nil
}

func (m *Menu) Draw(screen *ebiten.Image) {
	m.logo.Draw(screen)
	for _, sprite := range m.sprites {
		sprite.Draw(screen)
	}
}
