package resources

import (
	"ebijam23/states"

	"github.com/tinne26/etxt"
	"github.com/tinne26/etxt/fract"
)

type MenuItem interface {
	Hovered() bool
	Draw(ctx states.DrawContext)
	CheckState(x, y float64) bool
	Activate() bool
}

type SpriteItem struct {
	X, Y     float64
	Sprite   *Sprite
	hovered  bool
	Callback func() bool
}

func (s *SpriteItem) CheckState(x, y float64) bool {
	s.hovered = s.Sprite.Hit(x, y)
	return s.hovered
}

func (s *SpriteItem) Hovered() bool {
	return s.hovered
}

func (s *SpriteItem) Activate() bool {
	return s.Callback()
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
	Callback   func() bool
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

func (t *TextItem) Activate() bool {
	return t.Callback()
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
