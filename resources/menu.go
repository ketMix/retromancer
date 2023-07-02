package resources

import (
	"ebijam23/states"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tinne26/etxt"
	"github.com/tinne26/etxt/fract"
)

type MenuItem interface {
	Hovered() bool
	Draw(ctx states.DrawContext)
	CheckState(x, y float64) bool
	Activate() bool
	Hidden() bool
	SetHidden(bool)
}

type SpriteItem struct {
	X, Y     float64
	Sprite   *Sprite
	hovered  bool
	Callback func() bool
}

func (s *SpriteItem) Hidden() bool {
	return s.Sprite.Hidden
}

func (s *SpriteItem) SetHidden(h bool) {
	s.Sprite.Hidden = h
}

func (s *SpriteItem) CheckState(x, y float64) bool {
	if s.Sprite.Hidden {
		return false
	}

	s.hovered = s.Sprite.Hit(x, y)
	return s.hovered
}

func (s *SpriteItem) Hovered() bool {
	return !s.Sprite.Hidden && s.hovered
}

func (s *SpriteItem) Activate() bool {
	return s.Callback()
}

func (s *SpriteItem) Draw(ctx states.DrawContext) {
	s.Sprite.X = s.X
	s.Sprite.Y = s.Y
	s.Sprite.Draw(ctx)
}

type TextItem struct {
	X, Y            float64
	renderRect      fract.Rect
	hovered         bool
	Text            string
	Underline       bool
	Callback        func() bool
	SelfRefCallback *func(*TextItem) bool
	hidden          bool
}

func (t *TextItem) Hidden() bool {
	return t.hidden
}

func (t *TextItem) SetHidden(h bool) {
	t.hidden = h
}

func (t *TextItem) CheckState(x, y float64) bool {
	if t.hidden {
		return false
	}
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
	return !t.hidden && t.hovered
}

func (t *TextItem) Activate() bool {
	if t.SelfRefCallback != nil {
		return (*t.SelfRefCallback)(t)
	}
	return t.Callback()
}

func (t *TextItem) Draw(ctx states.DrawContext) {
	if t.hidden {
		return
	}
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

	ctx.Text.SetColor(color.NRGBA{0xff, 0xff, 0xff, 0xff})
	ctx.Text.Draw(ctx.Screen, t.Text, int(t.X), int(t.Y))

	if t.Underline {
		vector.StrokeLine(ctx.Screen, t.renderRect.Min.X.ToFloat32(), t.renderRect.Max.Y.ToFloat32(), t.renderRect.Max.X.ToFloat32(), t.renderRect.Max.Y.ToFloat32(), 1, color.NRGBA{0xff, 0xff, 0xff, 0xff}, false)
	}
}

type ButtonItem struct {
	X, Y       float64
	renderRect fract.Rect
	hovered    bool
	hidden     bool
	Text       string
	Callback   func() bool
}

func (t *ButtonItem) Hidden() bool {
	return t.hidden
}

func (t *ButtonItem) SetHidden(h bool) {
	t.hidden = h
}

func (t *ButtonItem) CheckState(x, y float64) bool {
	if t.hidden {
		return false
	}
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

func (t *ButtonItem) Hovered() bool {
	return !t.hidden && t.hovered
}

func (t *ButtonItem) Activate() bool {
	return t.Callback()
}

func (t *ButtonItem) Draw(ctx states.DrawContext) {
	if t.hidden {
		return
	}

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

	x1 := float32(t.X) - t.renderRect.Width().ToFloat32()/2 - 4
	x2 := float32(t.X) + t.renderRect.Width().ToFloat32()/2 + 4
	y1 := float32(t.Y) - t.renderRect.Height().ToFloat32()/2 - 4
	y2 := float32(t.Y) + t.renderRect.Height().ToFloat32()/2 + 4
	c := color.NRGBA{0xff, 0xff, 0xff, 0x80}

	if t.hovered {
		c = color.NRGBA{0xff, 0xff, 0xff, 0xff}
	}

	vector.StrokeLine(ctx.Screen, x1, y1, x2, y1, 1, c, false)
	vector.StrokeLine(ctx.Screen, x1, y1, x1, y2, 1, c, false)
	vector.StrokeLine(ctx.Screen, x2, y1, x2, y2, 1, c, false)
	vector.StrokeLine(ctx.Screen, x1, y2, x2, y2, 1, c, false)

	ctx.Text.SetColor(color.NRGBA{0xff, 0xff, 0xff, 0xff})
	ctx.Text.Draw(ctx.Screen, t.Text, int(t.X), int(t.Y))
}

type InputItem struct {
	X, Y        float64
	Width       float64
	Placeholder string
	renderRect  fract.Rect
	hovered     bool
	hidden      bool
	active      bool
	Text        string
	Callback    func() bool
}

func (t *InputItem) CheckState(x, y float64) bool {
	if t.hidden {
		return false
	}
	x1 := t.X - t.Width/2
	x2 := t.X + t.Width/2
	y1 := t.Y - t.renderRect.Height().ToFloat64()/2 - 4
	y2 := t.Y + t.renderRect.Height().ToFloat64()/2 + 4
	if x >= x1 && x <= x2 && y >= y1 && y <= y2 {
		t.hovered = true
	} else {
		t.hovered = false
	}
	return t.hovered
}

func (t *InputItem) Hidden() bool {
	return t.hidden
}

func (t *InputItem) SetHidden(h bool) {
	t.hidden = h
}

func (t *InputItem) Hovered() bool {
	return !t.hidden && t.hovered
}

func (t *InputItem) Activate() bool {
	t.active = true
	return t.Callback()
}

func (t *InputItem) Deactivate() {
	t.active = false
}

func (t *InputItem) IsActive() bool {
	return t.active
}

func (t *InputItem) Update() {
	if t.hidden {
		return
	}
	if t.active {
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			if len(t.Text) > 0 {
				t.Text = t.Text[:len(t.Text)-1]
			}
		} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			t.active = false
		} else if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyV) {
			t.Text += ReadClipboard()
		} else if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyC) {
			WriteClipboard(t.Text)
		} else {
			runes := ebiten.AppendInputChars(nil)
			t.Text += string(runes)
		}
	}
}

func (t *InputItem) Draw(ctx states.DrawContext) {
	if t.hidden {
		return
	}

	ctx.Text.SetAlign(etxt.YCenter | etxt.XCenter)
	txt := t.Text
	if txt == "" && t.Placeholder != "" {
		ctx.Text.SetColor(color.NRGBA{0x80, 0x80, 0x80, 0xff})
		txt = t.Placeholder
	} else {
		ctx.Text.SetColor(color.NRGBA{0xff, 0xff, 0xff, 0xff})
	}
	if txt == "" {
		txt = " "
	}
	t.renderRect = ctx.Text.Measure(txt)
	t.renderRect = t.renderRect.AddInts(int(t.X), int(t.Y))

	align := ctx.Text.GetAlign()
	if align&etxt.YCenter != 0 {
		t.renderRect = t.renderRect.AddInts(0, -t.renderRect.Height().ToInt()/2)
	}
	if align&etxt.XCenter != 0 {
		t.renderRect = t.renderRect.AddInts(-t.renderRect.Width().ToInt()/2, 0)
	}

	x1 := float32(t.X) - float32(t.Width)/2
	x2 := float32(t.X) + float32(t.Width)/2
	y2 := float32(t.Y) + t.renderRect.Height().ToFloat32()/2 + 4
	c := color.NRGBA{0xff, 0xff, 0xff, 0x80}

	if t.hovered {
		c = color.NRGBA{0xff, 0xff, 0xff, 0xff}
	}

	vector.StrokeLine(ctx.Screen, x1, y2, x2, y2, 1, c, false)

	ctx.Text.Draw(ctx.Screen, txt, int(t.X), int(t.Y))
}
