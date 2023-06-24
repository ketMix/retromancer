package menu

import (
	"ebijam23/resources"
	"ebijam23/states"
	"ebijam23/states/game"
	"fmt"
	"math/rand"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// PlayerEntry represents the lobby entry for a player.
type PlayerEntry struct {
	clickSound *resources.Sound
	//
	items []resources.MenuItem
	//
	hats     []string
	hatIndex int
	hatTitle *resources.TextItem
	hatItem  *resources.SpriteItem
	hatLeft  *resources.SpriteItem
	hatRight *resources.SpriteItem
	hatText  *resources.TextItem
	//
	controllerTitle  *resources.TextItem
	controllerItem   *resources.SpriteItem
	controllerLeft   *resources.SpriteItem
	controllerRight  *resources.SpriteItem
	controllerIdText *resources.TextItem
	controllerId     int
	//
	startText *resources.TextItem
	//
	waitingText *resources.TextItem
	//
	player game.Player
}

func (e *PlayerEntry) SetPlayer(player game.Player) {
	e.player = player
}

func (e *PlayerEntry) SyncHat(ctx states.Context) {
	e.hatItem.Sprite.SetImage(ctx.Manager.Get("images", e.hats[e.hatIndex]).(*ebiten.Image))

	e.hatText.Text = strings.TrimPrefix(e.hats[e.hatIndex], "hat-")

	if e.player != nil {
		e.player.SetHat(e.hats[e.hatIndex])
	}
}

func (e *PlayerEntry) SyncController(ctx states.Context) {
	if player, ok := e.player.(*game.LocalPlayer); ok {
		if e.controllerId == 0 {
			player.GamepadID = 0
			e.controllerItem.Sprite.SetImage(ctx.Manager.Get("images", "keyboard").(*ebiten.Image))
			e.controllerIdText.Text = ""
		} else {
			player.GamepadID = e.controllerId
			e.controllerItem.Sprite.SetImage(ctx.Manager.Get("images", "controller").(*ebiten.Image))
			e.controllerIdText.Text = fmt.Sprintf("%d", e.controllerId)
		}
	}
}

func (e *PlayerEntry) SetController(dir int) {
	controllers := ebiten.AppendGamepadIDs(nil)

	next := e.controllerId + dir

	if next > len(controllers) {
		next = 0
		e.controllerId = 0
	} else if next < 0 {
		next = len(controllers)
		e.controllerId = next
	} else {
		e.controllerId = next
	}
}

func (e *PlayerEntry) Init(ctx states.Context) error {
	e.clickSound = ctx.Manager.GetAs("sounds", "click", (*resources.Sound)(nil)).(*resources.Sound)

	e.hats = ctx.Manager.GetNamesWithPrefix("images", "hat-")
	e.hatIndex = int(rand.Int31n(int32(len(e.hats))))

	e.hatTitle = &resources.TextItem{
		Text: ctx.L("Hat"),
		Callback: func() bool {
			return false
		},
	}

	e.hatLeft = &resources.SpriteItem{
		Sprite: resources.NewSprite(ctx.Manager.Get("images", "arrow-left").(*ebiten.Image)),
		Callback: func() bool {
			e.hatIndex--
			if e.hatIndex < 0 {
				e.hatIndex = len(e.hats) - 1
			}
			e.hatItem.Sprite = resources.NewSprite(ctx.Manager.Get("images", e.hats[e.hatIndex]).(*ebiten.Image))
			e.clickSound.Play(1.0)
			e.SyncHat(ctx)
			return false
		},
	}
	e.hatLeft.Sprite.Centered = true

	e.hatItem = &resources.SpriteItem{
		Sprite: resources.NewSprite(ctx.Manager.Get("images", e.hats[e.hatIndex]).(*ebiten.Image)),
		Callback: func() bool {
			return false
		},
	}
	e.hatItem.Sprite.Centered = true
	e.hatItem.Sprite.Scale = 2.0

	e.hatRight = &resources.SpriteItem{
		Sprite: resources.NewSprite(ctx.Manager.Get("images", "arrow-right").(*ebiten.Image)),
		Callback: func() bool {
			e.hatIndex++
			if e.hatIndex >= len(e.hats) {
				e.hatIndex = 0
			}
			e.clickSound.Play(1.0)
			e.SyncHat(ctx)
			return false
		},
	}
	e.hatRight.Sprite.Centered = true

	e.hatText = &resources.TextItem{
		Text: "",
		Callback: func() bool {
			return false
		},
	}

	e.SyncHat(ctx)

	// Controller
	e.controllerTitle = &resources.TextItem{
		Text: ctx.L("Input"),
		Callback: func() bool {
			return false
		},
	}
	e.controllerLeft = &resources.SpriteItem{
		Sprite: resources.NewSprite(ctx.Manager.Get("images", "arrow-left").(*ebiten.Image)),
		Callback: func() bool {
			e.clickSound.Play(1.0)
			e.SetController(-1)
			e.SyncController(ctx)
			return false
		},
	}
	e.controllerLeft.Sprite.Centered = true
	e.controllerItem = &resources.SpriteItem{
		Sprite: resources.NewSprite(ctx.Manager.Get("images", "keyboard").(*ebiten.Image)),
		Callback: func() bool {
			return false
		},
	}
	e.controllerItem.Sprite.Centered = true
	e.controllerRight = &resources.SpriteItem{
		Sprite: resources.NewSprite(ctx.Manager.Get("images", "arrow-right").(*ebiten.Image)),
		Callback: func() bool {
			e.clickSound.Play(1.0)
			e.SetController(1)
			e.SyncController(ctx)
			return false
		},
	}
	e.controllerRight.Sprite.Centered = true

	e.controllerIdText = &resources.TextItem{
		Text: "",
		Callback: func() bool {
			return false
		},
	}

	// Other controls
	e.startText = &resources.TextItem{
		Text: ctx.L("Start"),
		Callback: func() bool {
			e.clickSound.Play(1.0)
			// FIXME: Need to agree w/ players to start (or assume host has full control).
			ctx.StateMachine.PushState(&game.World{
				StartingMap: "start",
				Players: []game.Player{
					e.player,
				},
			})
			return false
		},
	}

	e.waitingText = &resources.TextItem{
		Text: ctx.L("WaitingPlayer"),
		Callback: func() bool {
			return false
		},
	}

	e.items = append(e.items, e.hatTitle, e.hatLeft, e.hatItem, e.hatRight, e.hatText, e.controllerTitle, e.controllerLeft, e.controllerItem, e.controllerRight, e.startText, e.controllerIdText)

	return nil
}

func (e *PlayerEntry) Update(ctx states.Context, offsetX float64) error {
	// This is inefficent, but we reposition the entry stuff every frame.
	centerX := 320.0 - offsetX
	leftX := centerX - 50.0
	rightX := centerX + 50.0

	x := centerX
	y := 60.0

	e.hatText.X = x
	e.hatText.Y = y

	y += 30.0

	e.hatLeft.X = leftX
	e.hatLeft.Y = y

	e.hatItem.X = x
	e.hatItem.Y = y

	e.hatRight.X = rightX
	e.hatRight.Y = y

	y += 30.0

	e.hatTitle.X = x
	e.hatTitle.Y = y

	y += 60.0

	e.controllerTitle.X = x
	e.controllerTitle.Y = y

	y += 30.0

	e.controllerLeft.X = leftX
	e.controllerLeft.Y = y

	e.controllerItem.X = x
	e.controllerItem.Y = y

	e.controllerIdText.X = x + 20
	e.controllerIdText.Y = y

	e.controllerRight.X = rightX
	e.controllerRight.Y = y

	y += 60.0

	e.startText.X = x
	e.startText.Y = y

	//
	e.waitingText.X = centerX
	e.waitingText.Y = 200.0

	// Check for collisions.
	mx, my := ebiten.CursorPosition()
	for _, item := range e.items {
		item.CheckState(float64(mx), float64(my))
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		for _, item := range e.items {
			if item.Hovered() {
				if item.Activate() {
					return nil
				}
			}
		}
	}

	return nil
}

func (e *PlayerEntry) Draw(ctx states.DrawContext) {
	if e.player == nil {
		e.waitingText.Draw(ctx)
		return
	}
	for _, item := range e.items {
		item.Draw(ctx)
	}
}