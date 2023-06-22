package menu

import (
	"ebijam23/resources"
	"ebijam23/states"
	"ebijam23/states/game"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tinne26/etxt"
	"github.com/tinne26/etxt/fract"
)

type SinglePlayer struct {
	items    []MenuItem
	hats     []string
	hatIndex int
	hatItem  *SpriteItem
}

type MenuItem interface {
	Hovered() bool
	Draw(ctx states.DrawContext)
	CheckState(x, y float64) bool
	Callback() bool
}

type SpriteItem struct {
	X, Y     float64
	Sprite   *resources.Sprite
	hovered  bool
	callback func() bool
}

func (s *SpriteItem) CheckState(x, y float64) bool {
	s.hovered = s.Sprite.Hit(x, y)
	return s.hovered
}

func (s *SpriteItem) Hovered() bool {
	return s.hovered
}

func (s *SpriteItem) Callback() bool {
	return s.callback()
}

func (s *SpriteItem) Draw(ctx states.DrawContext) {
	s.Sprite.X = s.X
	s.Sprite.Y = s.Y
	s.Sprite.Draw(ctx.Screen)
}

type TextItem struct {
	X, Y       float64
	renderRect fract.Rect
	hovered    bool
	Text       string
	callback   func() bool
}

func (t *TextItem) CheckState(x, y float64) bool {
	x1 := t.renderRect.Min.X.ToFloat64()
	y1 := t.renderRect.Min.Y.ToFloat64()
	x2 := t.renderRect.Max.X.ToFloat64()
	y2 := t.renderRect.Max.Y.ToFloat64()
	if x >= x1 && x <= x2 && y >= y1 && y <= y2 {
		t.hovered = true
	} else {
		t.hovered = false
	}
	return t.hovered
}

func (t *TextItem) Hovered() bool {
	return t.hovered
}

func (t *TextItem) Callback() bool {
	return t.callback()
}

func (t *TextItem) Init(ctx states.Context) error {

	return nil
}

func (t *TextItem) Draw(ctx states.DrawContext) {
	ctx.Text.SetAlign(etxt.YCenter | etxt.XCenter)
	t.renderRect = ctx.Text.Measure(t.Text)
	t.renderRect = t.renderRect.AddInts(int(t.X), int(t.Y))

	align := ctx.Text.GetAlign()
	if align&etxt.YCenter != 0 {
		t.renderRect = t.renderRect.AddInts(0, -t.renderRect.Height().ToInt()/2)
	}
	if align&etxt.XCenter != 0 {
		t.renderRect = t.renderRect.AddInts(-t.renderRect.Width().ToInt()/2, 0)
	}

	ctx.Text.Draw(ctx.Screen, t.Text, int(t.X), int(t.Y))
}

func (s *SinglePlayer) Init(ctx states.Context) error {
	// Load in our hats.
	s.hats = ctx.Manager.GetNamesWithPrefix("images", "hat-")
	s.hatIndex = int(rand.Int31n(int32(len(s.hats))))

	x := 320.0
	s.items = append(s.items, &SpriteItem{
		X:      x,
		Y:      30,
		Sprite: resources.NewSprite(ctx.Manager.Get("images", "arrow-left").(*ebiten.Image)),
		callback: func() bool {
			s.hatIndex--
			if s.hatIndex < 0 {
				s.hatIndex = len(s.hats) - 1
			}
			s.hatItem.Sprite.SetImage(ctx.Manager.Get("images", s.hats[s.hatIndex]).(*ebiten.Image))
			return false
		},
	})

	x += s.items[len(s.items)-1].(*SpriteItem).Sprite.Width() + 5

	s.items = append(s.items, &SpriteItem{
		X:      x,
		Y:      35,
		Sprite: resources.NewSprite(ctx.Manager.Get("images", s.hats[s.hatIndex]).(*ebiten.Image)),
		callback: func() bool {
			return false
		},
	})
	s.hatItem = s.items[len(s.items)-1].(*SpriteItem)
	s.items[len(s.items)-1].(*SpriteItem).Sprite.Scale = 2.0

	x += s.items[len(s.items)-1].(*SpriteItem).Sprite.Width() + 5

	s.items = append(s.items, &SpriteItem{
		X:      x,
		Y:      30,
		Sprite: resources.NewSprite(ctx.Manager.Get("images", "arrow-right").(*ebiten.Image)),
		callback: func() bool {
			s.hatIndex++
			if s.hatIndex >= len(s.hats) {
				s.hatIndex = 0
			}
			s.hatItem.Sprite.SetImage(ctx.Manager.Get("images", s.hats[s.hatIndex]).(*ebiten.Image))
			return false
		},
	})

	//
	x = 320.0
	y := 320.0
	s.items = append(s.items, &TextItem{
		X:    x,
		Y:    y,
		Text: "Start",
		callback: func() bool {
			ctx.StateMachine.PopState()
			ctx.StateMachine.PushState(&game.World{
				StartingMap: "start",
				Players: []game.Player{
					&game.LocalPlayer{},
				},
			})
			return true
		},
	})
	y -= 50 + 16

	return nil
}

func (s *SinglePlayer) Finalize(ctx states.Context) error {
	return nil
}

func (s *SinglePlayer) Update(ctx states.Context) error {
	x, y := ebiten.CursorPosition()
	for _, m := range s.items {
		m.CheckState(float64(x), float64(y))
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		for _, m := range s.items {
			if m.Hovered() {
				if m.Callback() {
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
}
