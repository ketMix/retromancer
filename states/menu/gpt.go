package menu

import (
	"ebijam23/resources"
	"ebijam23/states"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GPTOptions struct {
	back, checkKey, generate *resources.ButtonItem
	key, locale, style       *resources.InputItem
	buttons                  []*resources.ButtonItem
	inputs                   []*resources.InputItem
	click                    *resources.Sound
	gptKeyIsValid            bool
	gptIsActive              bool
}

func (s *GPTOptions) Init(ctx states.Context) error {
	s.gptKeyIsValid = ctx.L.CheckGPTKey()
	s.click = ctx.R.GetAs("sounds", "click", (*resources.Sound)(nil)).(*resources.Sound)

	x, y := 1280, 720
	centerX := float64(x / 4)
	width := float64(350)
	s.key = &resources.InputItem{
		X:           centerX,
		Y:           float64(y/20) * 6,
		Width:       width,
		Text:        ctx.L.GetGPTKey(),
		Placeholder: ctx.L.Get("OpenAI API Key"),
		Callback: func() bool {
			return false
		},
	}
	s.style = &resources.InputItem{
		X:           centerX,
		Y:           float64(y/20) * 7,
		Width:       width,
		Placeholder: ctx.L.Get("Writing Style"),
		Text:        ctx.L.GetGPTStyle(),
		Callback: func() bool {
			return false
		},
	}
	s.locale = &resources.InputItem{
		X:           centerX,
		Y:           float64(y/20) * 8,
		Width:       100,
		Text:        ctx.L.Locale(),
		Placeholder: ctx.L.Get("Language"),
		Callback: func() bool {
			return false
		},
	}

	s.checkKey = &resources.ButtonItem{
		X:    centerX - width/4,
		Y:    335,
		Text: ctx.L.Get("Check Key"),
		Callback: func() bool {
			s.click.Play(1.0)
			if s.key.Text != "" {
				ctx.L.SetGPTKey(s.key.Text)
			}
			if ctx.L.CheckGPTKey() {
				s.gptKeyIsValid = true
			} else {
				s.gptKeyIsValid = false
			}
			return true
		},
	}

	s.generate = &resources.ButtonItem{
		X:    centerX + width/4,
		Y:    335,
		Text: ctx.L.Get("Generate"),
		Callback: func() bool {
			s.click.Play(1.0)
			ctx.L.SetGPTKey(s.key.Text)
			if ctx.L.CheckGPTKey() {
				s.gptKeyIsValid = true
				ctx.L.SetGPTStyle(s.style.Text)
				ctx.L.SetLocale(s.locale.Text, true)
			} else {
				s.gptKeyIsValid = false
			}
			return true
		},
	}
	s.back = &resources.ButtonItem{
		Text: ctx.L.Get("Back"),
		X:    30,
		Y:    335,
		Callback: func() bool {
			s.click.Play(1.0)
			ctx.StateMachine.PopState(nil)
			return false
		},
	}
	s.inputs = append(s.inputs, s.key, s.style, s.locale)
	s.buttons = append(s.buttons, s.checkKey, s.generate, s.back)
	return nil
}

func (s *GPTOptions) Finalize(ctx states.Context) error {
	return nil
}

func (s *GPTOptions) Enter(ctx states.Context, v interface{}) error {
	return nil
}

func (s *GPTOptions) Update(ctx states.Context) error {
	s.gptIsActive = ctx.L.GPTIsActive()
	s.generate.Text = ctx.L.Get("Generate")
	s.back.Text = ctx.L.Get("Back")
	s.locale.Placeholder = ctx.L.Get("Language")
	s.style.Placeholder = ctx.L.Get("Writing Style")

	x, y := ebiten.CursorPosition()
	for _, m := range s.buttons {
		m.CheckState(float64(x), float64(y))
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		for _, m := range s.buttons {
			if m.Hovered() {
				if m.Activate() {
					return nil
				}
			}
		}
	}
	for _, m := range s.inputs {
		m.CheckState(float64(x), float64(y))
		m.Update()
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		for _, m := range s.inputs {
			if m.Hovered() {
				if m.Activate() {
					return nil
				}
			} else {
				m.Deactivate()
			}
		}
	}
	return nil
}

func (s *GPTOptions) Draw(ctx states.DrawContext) {
	for _, m := range s.buttons {
		m.Draw(ctx)
	}
	for _, m := range s.inputs {
		m.Draw(ctx)
	}
	x, y := 1280, 720

	if !s.gptKeyIsValid {
		ctx.Text.SetScale(1.0)
		ctx.Text.SetColor(color.RGBA{255, 0, 0, 255})
		ctx.Text.Draw(ctx.Screen, ctx.L.Get("GPTKeyNotValid"), x/4, (y / 24))
	} else {
		ctx.Text.SetScale(1.0)
		ctx.Text.SetColor(color.RGBA{0, 255, 0, 255})
		ctx.Text.Draw(ctx.Screen, ctx.L.Get("GPTKeyValid"), x/4, (y / 24))
	}

	if !s.gptIsActive {
		ctx.Text.SetScale(1.0)
		ctx.Text.SetColor(color.RGBA{255, 0, 0, 255})
		ctx.Text.Draw(ctx.Screen, ctx.L.Get("GPTNotActive"), x/4, (y/24)*2)
	}

	// Draw instruction
	text := strings.Split(ctx.L.Get("GPTInstructions"), "\n")
	splitText := make([]string, 0)
	maxLen := 60
	for _, line := range text {
		if len(line) <= maxLen {
			splitText = append(splitText, line)
			continue
		}
		// split text into lines that are less than maxLen
		// split on spaces
		words := strings.Split(line, " ")
		currentLine := ""
		for _, word := range words {
			if len(currentLine)+len(word) > maxLen {
				splitText = append(splitText, currentLine)
				currentLine = ""
			}
			currentLine += word + " "
		}
		splitText = append(splitText, currentLine)
	}
	y = (y/20)*4 - (len(splitText)/2)*int(ctx.Text.Utils().GetLineHeight())
	for _, line := range splitText {
		{
			ctx.Text.SetScale(1.0)
			ctx.Text.SetColor(color.White)
			ctx.Text.Draw(ctx.Screen, line, x/4, y)
		}
		y += int(ctx.Text.Utils().GetLineHeight())
	}
}
