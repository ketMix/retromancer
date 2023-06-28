package menu

import (
	"ebijam23/resources"
	"ebijam23/states"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Loading struct {
	vfx   resources.VFXList
	ticks int
}

func (l *Loading) Init(ctx states.Context) error {
	x, y := ebiten.WindowSize()
	x /= 4
	y /= 4

	l.vfx.Add(&resources.Text{
		Text:         ctx.L("Loading"),
		InDuration:   0,
		HoldDuration: 5000 * time.Second,
		OutDuration:  0,
		X:            float64(x),
		Y:            float64(y),
	})
	return nil
}

func (l *Loading) Enter(ctx states.Context) error {
	return nil
}

func (l *Loading) Finalize(ctx states.Context) error {
	return nil
}

func (l *Loading) Update(ctx states.Context) error {
	l.ticks++
	if l.ticks > 20 {
		ctx.SetLocale(ctx.Locale(), true)

		// Pop the loading and the GPT question screen
		ctx.StateMachine.PopState()
		ctx.StateMachine.PopState()
	}
	return nil
}

func (l *Loading) Draw(ctx states.DrawContext) {
	l.vfx.Process(ctx, nil)
}
