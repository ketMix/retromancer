package resources

import "github.com/hajimehoshi/ebiten/v2"

type Sprites []*Sprite

type Sprite struct {
	image   *ebiten.Image
	X, Y    float64
	Options ebiten.DrawImageOptions
}

func (s *Sprite) Width() float64 {
	return float64(s.image.Bounds().Dx())
}

func (s *Sprite) Height() float64 {
	return float64(s.image.Bounds().Dy())
}

func (s *Sprite) Image() *ebiten.Image {
	return s.image
}

func (s *Sprite) Draw(screen *ebiten.Image) {
	s.Options.GeoM.Reset()
	s.Options.GeoM.Translate(s.X, s.Y)
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
