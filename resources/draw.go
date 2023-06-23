package resources

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
)

// This isn't a resources, but whatever.

// DrawArc is an inefficient pixel-based arc drawererer
func DrawArc(screen *ebiten.Image, posX, posY float64, radius float64, start, end float64, color color.Color) {
	start += math.Pi / 2
	end += math.Pi / 2

	for i := start; i < end; i += 0.05 {
		x := math.Cos(i)*radius + posX
		y := math.Sin(i)*radius + posY
		screen.Set(int(x), int(y), color)
	}
}

func DrawTextOutline(text *etxt.Renderer, screen *ebiten.Image, str string, x, y int, scale int) {
	// :)
	text.Draw(screen, str, x-scale, y)
	text.Draw(screen, str, x+scale, y)
	text.Draw(screen, str, x, y-scale)
	text.Draw(screen, str, x, y+scale)
}
