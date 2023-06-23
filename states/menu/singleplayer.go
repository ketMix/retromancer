package menu

import (
	"ebijam23/resources"
	"ebijam23/states"
	"ebijam23/states/game"
	"math/rand"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SinglePlayer struct {
	clickSound     *resources.Sound
	items          []resources.MenuItem
	hats           []string
	hatIndex       int
	hatItem        *resources.SpriteItem
	hatText        *resources.TextItem
	controllerItem *resources.SpriteItem
	useController  bool
	//
	localPlayers []*game.LocalPlayer
	overlay      game.Overlay
}

func (s *SinglePlayer) Init(ctx states.Context) error {
	s.overlay.Init(ctx)
	//
	s.clickSound = ctx.Manager.GetAs("sounds", "click", (*resources.Sound)(nil)).(*resources.Sound)

	// Set up local player.
	s.localPlayers = append(s.localPlayers, game.NewLocalPlayer())

	// Load in our hats.
	s.hats = ctx.Manager.GetNamesWithPrefix("images", "hat-")
	s.hatIndex = int(rand.Int31n(int32(len(s.hats))))

	centerX := 320.0
	leftX := centerX - 50.0
	rightX := centerX + 50.0
	x := 320.0
	y := 30.0

	s.items = append(s.items, &resources.TextItem{
		X:    centerX,
		Y:    y,
		Text: ctx.L("Hat"),
		Callback: func() bool {
			return false
		},
	})

	y += 30.0

	s.items = append(s.items, &resources.SpriteItem{
		X:      leftX,
		Y:      y,
		Sprite: resources.NewSprite(ctx.Manager.Get("images", "arrow-left").(*ebiten.Image)),
		Callback: func() bool {
			s.clickSound.Play(1.0)
			s.hatIndex--
			if s.hatIndex < 0 {
				s.hatIndex = len(s.hats) - 1
			}
			s.SyncHat(ctx)
			return false
		},
	})
	s.items[len(s.items)-1].(*resources.SpriteItem).Sprite.Centered = true

	s.items = append(s.items, &resources.SpriteItem{
		X:      centerX,
		Y:      y,
		Sprite: resources.NewSprite(ctx.Manager.Get("images", s.hats[s.hatIndex]).(*ebiten.Image)),
		Callback: func() bool {
			return false
		},
	})
	s.hatItem = s.items[len(s.items)-1].(*resources.SpriteItem)
	s.hatItem.Sprite.Centered = true
	s.items[len(s.items)-1].(*resources.SpriteItem).Sprite.Scale = 2.0

	s.items = append(s.items, &resources.TextItem{
		X:    centerX,
		Y:    y + 20,
		Text: "",
		Callback: func() bool {
			return false
		},
	})
	s.hatText = s.items[len(s.items)-1].(*resources.TextItem)

	s.SyncHat(ctx)

	s.items = append(s.items, &resources.SpriteItem{
		X:      rightX,
		Y:      y,
		Sprite: resources.NewSprite(ctx.Manager.Get("images", "arrow-right").(*ebiten.Image)),
		Callback: func() bool {
			s.clickSound.Play(1.0)
			s.hatIndex++
			if s.hatIndex >= len(s.hats) {
				s.hatIndex = 0
			}
			s.SyncHat(ctx)
			return false
		},
	})
	s.items[len(s.items)-1].(*resources.SpriteItem).Sprite.Centered = true

	y += s.items[len(s.items)-1].(*resources.SpriteItem).Sprite.Height() + 60

	// Controller
	s.items = append(s.items, &resources.TextItem{
		X:    centerX,
		Y:    y,
		Text: ctx.L("Input"),
		Callback: func() bool {
			return false
		},
	})

	y += 30.0

	s.items = append(s.items, &resources.SpriteItem{
		X:      leftX,
		Y:      y,
		Sprite: resources.NewSprite(ctx.Manager.Get("images", "arrow-left").(*ebiten.Image)),
		Callback: func() bool {
			s.clickSound.Play(1.0)
			s.useController = !s.useController
			s.SyncController(ctx)
			return false
		},
	})
	s.items[len(s.items)-1].(*resources.SpriteItem).Sprite.Centered = true

	s.items = append(s.items, &resources.SpriteItem{
		X:      centerX,
		Y:      y,
		Sprite: resources.NewSprite(ctx.Manager.Get("images", "keyboard").(*ebiten.Image)),
		Callback: func() bool {
			return false
		},
	})
	s.controllerItem = s.items[len(s.items)-1].(*resources.SpriteItem)
	s.items[len(s.items)-1].(*resources.SpriteItem).Sprite.Centered = true

	s.items = append(s.items, &resources.SpriteItem{
		X:      rightX,
		Y:      y,
		Sprite: resources.NewSprite(ctx.Manager.Get("images", "arrow-right").(*ebiten.Image)),
		Callback: func() bool {
			s.clickSound.Play(1.0)
			s.useController = !s.useController
			s.SyncController(ctx)
			return false
		},
	})
	s.items[len(s.items)-1].(*resources.SpriteItem).Sprite.Centered = true

	// Start and stuff
	x = 320.0
	y = 320.0
	s.items = append(s.items, &resources.TextItem{
		X:    x,
		Y:    y,
		Text: ctx.L("Start"),
		Callback: func() bool {
			s.clickSound.Play(1.0)
			ctx.StateMachine.PopState()
			ctx.StateMachine.PushState(&game.World{
				StartingMap: "start",
				Players: []game.Player{
					s.localPlayers[0],
				},
			})
			return true
		},
	})
	y -= 50 + 16

	return nil
}

func (s *SinglePlayer) SyncHat(ctx states.Context) {
	s.hatItem.Sprite.SetImage(ctx.Manager.Get("images", s.hats[s.hatIndex]).(*ebiten.Image))

	s.hatText.Text = strings.TrimPrefix(s.hats[s.hatIndex], "hat-")

	s.localPlayers[0].SetHat(s.hats[s.hatIndex])
}

func (s *SinglePlayer) SyncController(ctx states.Context) {
	if s.useController && len(ebiten.AppendGamepadIDs(nil)) == 0 {
		s.useController = false
	}

	if s.useController {
		gamepadIDs := ebiten.AppendGamepadIDs(nil)
		s.localPlayers[0].GamepadID = int(gamepadIDs[0])

		s.controllerItem.Sprite.SetImage(ctx.Manager.Get("images", "controller").(*ebiten.Image))
	} else {
		s.localPlayers[0].GamepadID = 0
		s.controllerItem.Sprite.SetImage(ctx.Manager.Get("images", "keyboard").(*ebiten.Image))
	}
}

func (s *SinglePlayer) Finalize(ctx states.Context) error {
	return nil
}

func (s *SinglePlayer) Enter(ctx states.Context) error {
	return nil
}

func (s *SinglePlayer) Update(ctx states.Context) error {
	s.overlay.Update(ctx)
	x, y := ebiten.CursorPosition()
	for _, m := range s.items {
		m.CheckState(float64(x), float64(y))
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		for _, m := range s.items {
			if m.Hovered() {
				if m.Activate() {
					return nil
				}
			}
		}
	}
	return nil
}

func (s *SinglePlayer) Draw(ctx states.DrawContext) {
	for _, m := range s.items {
		m.Draw(ctx)
	}
	s.overlay.Draw(ctx)
}
