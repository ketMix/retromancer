package menu

import (
	"ebijam23/resources"
	"ebijam23/states"
	"ebijam23/states/game"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type SinglePlayer struct {
	clickSound      *resources.Sound
	items           []resources.MenuItem
	multiplayerItem *resources.ButtonItem
	joinItem        *resources.ButtonItem
	hostItem        *resources.ButtonItem
	backItem        *resources.TextItem
	lobbyItem       *resources.InputItem
	playerEntries   [2]*PlayerEntry
	overlay         game.Overlay
}

func (s *SinglePlayer) Init(ctx states.Context) error {
	s.overlay.Init(ctx)
	//
	s.clickSound = ctx.Manager.GetAs("sounds", "click", (*resources.Sound)(nil)).(*resources.Sound)

	s.playerEntries[0] = &PlayerEntry{
		player: game.NewLocalPlayer(),
	}
	// Dummy player entry for 2nd player.
	s.playerEntries[1] = &PlayerEntry{}

	for _, e := range s.playerEntries {
		e.Init(ctx)
	}

	s.multiplayerItem = &resources.ButtonItem{
		Text: ctx.L("Multiplayer"),
		X:    500,
		Y:    20,
		Callback: func() bool {
			s.clickSound.Play(1.0)
			for i, item := range s.items {
				if item == s.multiplayerItem {
					s.items = append(s.items[:i], s.items[i+1:]...)
					s.items = append(s.items, s.lobbyItem, s.joinItem, s.hostItem)
					break
				}
			}
			return false
		},
	}

	s.lobbyItem = &resources.InputItem{
		X:     350,
		Y:     20,
		Width: 150,
		Callback: func() bool {
			return false
		},
	}

	s.joinItem = &resources.ButtonItem{
		Text: ctx.L("Host"),
		X:    450,
		Y:    20,
		Callback: func() bool {
			s.clickSound.Play(1.0)
			return false
		},
	}

	s.hostItem = &resources.ButtonItem{
		Text: ctx.L("Join"),
		X:    450 + 50,
		Y:    20,
		Callback: func() bool {
			s.clickSound.Play(1.0)
			return false
		},
	}

	s.backItem = &resources.TextItem{
		Text: ctx.L("Back"),
		X:    30,
		Y:    335,
		Callback: func() bool {
			s.clickSound.Play(1.0)
			ctx.StateMachine.PopState()
			return false
		},
	}
	s.items = append(s.items, s.backItem, s.multiplayerItem)

	return nil
}

func (s *SinglePlayer) Finalize(ctx states.Context) error {
	return nil
}

func (s *SinglePlayer) Enter(ctx states.Context) error {
	return nil
}

func (s *SinglePlayer) Update(ctx states.Context) error {
	s.overlay.Update(ctx)

	s.lobbyItem.Update()

	// Check for controller button hit to activate player 2.
	for _, gamepadID := range ebiten.AppendGamepadIDs(nil) {
		if inpututil.IsGamepadButtonJustPressed(gamepadID, ebiten.GamepadButton9) {
			pl := game.NewLocalPlayer()
			s.playerEntries[1].player = pl
			s.playerEntries[1].controllerId = int(gamepadID)
			s.playerEntries[1].SyncController(ctx)
			pl.GamepadID = int(gamepadID)
		}
	}

	x := -(len(s.playerEntries) - 1) * 150 / 2
	for i, e := range s.playerEntries {
		e.Update(ctx, float64(x+i*150))
	}

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
	ctx.Text.SetColor(color.White)
	for _, e := range s.playerEntries {
		e.Draw(ctx)
	}
	for _, m := range s.items {
		m.Draw(ctx)
	}
	s.overlay.Draw(ctx)
}
