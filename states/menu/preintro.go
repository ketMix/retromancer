package menu

import (
	"ebijam23/resources"
	"ebijam23/states"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type PreIntro struct {
	buttons []*resources.ButtonItem
	text    resources.TextItem
}

func (p *PreIntro) Init(ctx states.Context) error {
	if !ctx.L.CheckGPTKey() {
		ctx.StateMachine.PopState(nil)
	}
	x, y := 1280, 720
	centerX := x / 4
	centerY := y / 4
	p.buttons = append(p.buttons,
		&resources.ButtonItem{
			Text: ctx.L.Get("Yes"),
			X:    float64(centerX) - 25,
			Y:    float64(centerY) + 100,
			Callback: func() bool {
				ctx.StateMachine.PushState(&Loading{})
				return false
			},
		},
		&resources.ButtonItem{
			Text: ctx.L.Get("No"),
			X:    float64(centerX) + 25,
			Y:    float64(centerY) + 100,
			Callback: func() bool {
				ctx.StateMachine.PopState(nil)
				return false
			},
		},
	)
	text := make([]string, 0)
	text = append(text,
		"This game has the potential to utilize ChatGPT to provide a custom player experience.",
		"The style of all the game's text can be adjusted from an option in the main menu.",
		"This is an experimental feature and may not work as expected.",
		"Generating the text for this game takes roughly ~40 seconds.",
		"",
		"",
		"",
		"Would you like to enable ChatGPT?",
	)
	p.text = resources.TextItem{
		X:    float64(centerX),
		Y:    float64(centerY),
		Text: strings.Join(text, "\n"),
	}
	return nil
}

func (p *PreIntro) Enter(ctx states.Context, v interface{}) error {
	return nil
}

func (p *PreIntro) Finalize(ctx states.Context) error {
	return nil
}

func (p *PreIntro) Update(ctx states.Context) error {
	x, y := ebiten.CursorPosition()
	for _, m := range p.buttons {
		m.CheckState(float64(x), float64(y))
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		for _, m := range p.buttons {
			if m.Hovered() {
				if m.Activate() {
					return nil
				}
			}
		}
	}
	return nil
}

func (p *PreIntro) Draw(ctx states.DrawContext) {
	for _, m := range p.buttons {
		m.Draw(ctx)
	}
	p.text.Draw(ctx)
}
