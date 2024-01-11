package menu

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/ketMix/retromancer/states"
	"github.com/ketMix/retromancer/states/game"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/ketMix/retromancer/resources"
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
	controllerIndex  int
	useController    bool
	//
	difficulties []states.Difficulty
	diffIndex    int
	diffTitle    *resources.TextItem
	diffLeft     *resources.SpriteItem
	diffRight    *resources.SpriteItem
	diffItem     *resources.TextItem
	//
	startText *resources.TextItem
	//
	waitingText *resources.TextItem
	//
	id     uint32 // ID. Only set if it is a networked player.
	player game.Player
}

func (e *PlayerEntry) SetPlayer(player game.Player) {
	e.player = player
}

func (e *PlayerEntry) SyncHat(ctx states.Context) {
	e.hatItem.Sprite.SetImage(ctx.R.Get("images", e.hats[e.hatIndex]).(*ebiten.Image))

	e.hatText.Text = strings.TrimPrefix(e.hats[e.hatIndex], "hat-")

	if e.player != nil {
		e.player.SetHat(e.hats[e.hatIndex])
	}
}

func (e *PlayerEntry) SyncController(ctx states.Context) {
	controllers := resources.GetFunctionalGamepads()

	if player, ok := e.player.(*game.LocalPlayer); ok {
		if !e.useController {
			player.GamepadID = -1
			e.controllerItem.Sprite.SetImage(ctx.R.Get("images", "keyboard").(*ebiten.Image))
			e.controllerIdText.Text = ""
		} else {
			player.GamepadID = controllers[e.controllerIndex]
			player.GamepadMap = resources.GetBestGamemap(player.GamepadID)
			e.controllerItem.Sprite.SetImage(ctx.R.Get("images", "controller").(*ebiten.Image))
			e.controllerIdText.Text = fmt.Sprintf("%d", player.GamepadID)
		}
	}
}

func (e *PlayerEntry) SyncDifficulty(ctx states.Context, i int) {
	e.diffItem.Text = ctx.L.Get(string(e.difficulties[i]))
}

func (e *PlayerEntry) SetController(dir int) {
	controllers := resources.GetFunctionalGamepads()

	if len(controllers) == 0 {
		e.useController = false
		return
	}

	next := e.controllerIndex + dir

	if e.useController {
		if next >= len(controllers) {
			e.useController = false
		} else if next < 0 {
			e.useController = false
		} else {
			e.controllerIndex = next
		}
	} else {
		if dir == -1 {
			e.controllerIndex = len(controllers) - 1
			e.useController = true
		} else {
			e.controllerIndex = 0
			e.useController = true
		}
	}
}

func (e *PlayerEntry) Init(s *Lobby, ctx states.Context) error {
	e.clickSound = ctx.R.GetAs("sounds", "click", (*resources.Sound)(nil)).(*resources.Sound)

	// Dem hats
	e.hats = ctx.R.GetNamesWithPrefix("images", "hat-")
	e.hatIndex = int(rand.Int31n(int32(len(e.hats))))

	e.hatTitle = &resources.TextItem{
		Text: ctx.L.Get("Hat"),
		Callback: func() bool {
			return false
		},
	}

	e.hatLeft = &resources.SpriteItem{
		Sprite: resources.NewSprite(ctx.R.Get("images", "arrow-left").(*ebiten.Image)),
		Callback: func() bool {
			e.hatIndex--
			if e.hatIndex < 0 {
				e.hatIndex = len(e.hats) - 1
			}
			e.clickSound.Play(1.0)
			e.SyncHat(ctx)
			s.SendToNetPlayers(HatMessage{Hat: e.hatIndex})
			return false
		},
	}
	e.hatLeft.Sprite.Centered = true

	e.hatItem = &resources.SpriteItem{
		Sprite: resources.NewSprite(ctx.R.Get("images", e.hats[e.hatIndex]).(*ebiten.Image)),
		Callback: func() bool {
			return false
		},
	}
	e.hatItem.Sprite.Centered = true
	e.hatItem.Sprite.Scale = 2.0

	e.hatRight = &resources.SpriteItem{
		Sprite: resources.NewSprite(ctx.R.Get("images", "arrow-right").(*ebiten.Image)),
		Callback: func() bool {
			e.hatIndex++
			if e.hatIndex >= len(e.hats) {
				e.hatIndex = 0
			}
			e.clickSound.Play(1.0)
			e.SyncHat(ctx)
			s.SendToNetPlayers(HatMessage{Hat: e.hatIndex})
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

	hatItems := []resources.MenuItem{
		e.hatTitle,
		e.hatLeft,
		e.hatItem,
		e.hatRight,
		e.hatText,
	}
	e.SyncHat(ctx)

	// Controller
	e.controllerTitle = &resources.TextItem{
		Text: ctx.L.Get("Input"),
		Callback: func() bool {
			return false
		},
	}
	e.controllerLeft = &resources.SpriteItem{
		Sprite: resources.NewSprite(ctx.R.Get("images", "arrow-left").(*ebiten.Image)),
		Callback: func() bool {
			e.clickSound.Play(1.0)
			e.SetController(-1)
			e.SyncController(ctx)
			return false
		},
	}
	e.controllerLeft.Sprite.Centered = true
	e.controllerItem = &resources.SpriteItem{
		Sprite: resources.NewSprite(ctx.R.Get("images", "keyboard").(*ebiten.Image)),
		Callback: func() bool {
			return false
		},
	}
	e.controllerItem.Sprite.Centered = true
	e.controllerRight = &resources.SpriteItem{
		Sprite: resources.NewSprite(ctx.R.Get("images", "arrow-right").(*ebiten.Image)),
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
	controllerItems := []resources.MenuItem{
		e.controllerTitle,
		e.controllerLeft,
		e.controllerItem,
		e.controllerRight,
		e.controllerIdText,
	}

	// Difficulty
	e.difficulties = []states.Difficulty{states.DifficultyEasy, states.DifficultyNormal, states.DifficultyHard}
	e.diffIndex = 1

	e.diffTitle = &resources.TextItem{
		Text: ctx.L.Get("Difficulty"),
	}

	e.diffLeft = &resources.SpriteItem{
		Sprite: resources.NewSprite(ctx.R.Get("images", "arrow-left").(*ebiten.Image)),
		Callback: func() bool {
			e.diffIndex--
			e.clickSound.Play(1.0)
			if e.diffIndex < 0 {
				e.diffIndex = 0
			}
			e.SyncDifficulty(ctx, e.diffIndex)
			return false
		},
	}
	e.diffLeft.Sprite.Centered = true

	e.diffItem = &resources.TextItem{
		Text: ctx.L.Get(string(e.difficulties[e.diffIndex])),
	}

	e.diffRight = &resources.SpriteItem{
		Sprite: resources.NewSprite(ctx.R.Get("images", "arrow-right").(*ebiten.Image)),
		Callback: func() bool {
			e.diffIndex++
			e.clickSound.Play(1.0)
			if e.diffIndex > len(e.difficulties)-1 {
				e.diffIndex = len(e.difficulties) - 1
			}
			e.SyncDifficulty(ctx, e.diffIndex)
			return false
		},
	}
	e.diffRight.Sprite.Centered = true

	difficultyItems := []resources.MenuItem{
		e.diffTitle,
		e.diffLeft,
		e.diffRight,
		e.diffItem,
	}

	// Other controls
	e.startText = &resources.TextItem{
		Text: ctx.L.Get("Start"),
		Callback: func() bool {
			e.clickSound.Play(1.0)
			if !s.net.Running || s.net.Hosting {
				s.shouldStart = true
				s.difficulty = e.difficulties[e.diffIndex]
			}
			return false
		},
	}

	e.waitingText = &resources.TextItem{
		Text: ctx.L.Get("WaitingPlayer"),
		Callback: func() bool {
			return false
		},
	}

	e.items = append(e.items, e.startText)
	e.items = append(e.items, hatItems...)
	e.items = append(e.items, controllerItems...)
	e.items = append(e.items, difficultyItems...)
	return nil
}

func (e *PlayerEntry) Update(ctx states.Context, offsetX float64) error {
	// This is inefficent, but we reposition the entry stuff every frame.
	centerX := 320.0 - offsetX
	leftX := centerX - 50.0
	rightX := centerX + 50.0

	x := centerX
	y := 60.0

	// Hats
	e.hatTitle.X = x
	e.hatTitle.Y = y

	y += 30.0

	e.hatLeft.X = leftX
	e.hatLeft.Y = y

	e.hatItem.X = x
	e.hatItem.Y = y

	e.hatRight.X = rightX
	e.hatRight.Y = y

	y += 30.0
	e.hatText.X = x
	e.hatText.Y = y

	y += 50.0

	// Input
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

	y += 50.0

	// Difficulty
	e.diffTitle.X = x
	e.diffTitle.Y = y

	y += 30.0
	e.diffLeft.X = leftX
	e.diffLeft.Y = y

	e.diffItem.X = x
	e.diffItem.Y = y

	e.diffRight.X = rightX
	e.diffRight.Y = y

	y += 50.0
	//
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
