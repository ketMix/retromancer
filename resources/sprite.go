package resources

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprites []*Sprite

type Sprite struct {
	image            *ebiten.Image
	X, Y             float64
	lastX, lastY     float64
	interpX, interpY float64
	Interpolate      bool
	Centered         bool
	Options          ebiten.DrawImageOptions
}

func (s *Sprite) Width() float64 {
	return float64(s.image.Bounds().Dx())
}

func (s *Sprite) Height() float64 {
	return float64(s.image.Bounds().Dy())
}

func (s *Sprite) SetXY(x, y float64) {
	s.lastX = s.X
	s.lastY = s.Y
	s.X = x
	s.Y = y
}

func (s *Sprite) SetImage(image *ebiten.Image) {
	s.image = image
}

func (s *Sprite) Image() *ebiten.Image {
	return s.image
}

func (s *Sprite) Draw(screen *ebiten.Image) {
	// hmmmmm
	/*lx := (s.X * (1.0 - 0.1)) + (s.interpX * 0.1)
	ly := (s.Y * (1.0 - 0.1)) + (s.interpY * 0.1)
	s.interpX = lx
	s.interpY = ly
	fmt.Println(s.X, s.interpX, s.lastX)*/
	// FIXME
	// Alternatively, we could be dumb:
	if s.Interpolate {
		if s.interpX < s.X {
			s.interpX++
		} else if s.interpX > s.X {
			s.interpX--
		}
		if math.Abs(s.interpX-s.X) < 1 {
			s.interpX = s.X
		}
		if s.interpY < s.Y {
			s.interpY++
		} else if s.interpY > s.Y {
			s.interpY--
		}
		if math.Abs(s.interpY-s.Y) < 1 {
			s.interpY = s.Y
		}
	} else {
		s.interpX = s.X
		s.interpY = s.Y
	}
	s.Options.GeoM.Reset()
	if s.Centered {
		s.Options.GeoM.Translate(s.interpX-s.Width()/2, s.interpY-s.Height()/2)
	} else {
		s.Options.GeoM.Translate(s.interpX, s.interpY)
	}
	screen.DrawImage(s.image, &s.Options)
}

func (s *Sprite) Hit(x, y float64) bool {
	if x < s.X || x > s.X+s.Width() {
		return false
	}
	if y < s.Y || y > s.Y+s.Height() {
		return false
	}
	return true
}

func NewSprite(image *ebiten.Image) *Sprite {
	return &Sprite{
		image: image,
	}
}
