package resources

import (
	"ebijam23/states"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprites []*Sprite

type Sprite struct {
	image            *ebiten.Image
	images           []*ebiten.Image
	frame            int
	Framerate        int
	elapsed          int
	Loop             bool
	X, Y             float64
	lastX, lastY     float64
	interpX, interpY float64
	Scale            float64
	Flipped          bool
	Interpolate      bool
	Centered         bool
	Hidden           bool
	Options          ebiten.DrawImageOptions
	vfx              VFXList
}

func (s *Sprite) Width() float64 {
	if s.Scale != 0 {
		return float64(s.image.Bounds().Dx()) * s.Scale
	}
	return float64(s.image.Bounds().Dx())
}

func (s *Sprite) Height() float64 {
	if s.Scale != 0 {
		return float64(s.image.Bounds().Dy()) * s.Scale
	}
	return float64(s.image.Bounds().Dy())
}

func (s *Sprite) SetXY(x, y float64) {
	s.lastX = s.X
	s.lastY = s.Y
	s.X = x
	s.Y = y
}

func (s *Sprite) Reset() {
	if len(s.images) > 0 {
		s.frame = 0
		s.image = s.images[0]
	}
}

func (s *Sprite) Frame() int {
	return s.frame
}

func (s *Sprite) SetFrame(frame int) {
	if frame >= len(s.images) {
		frame = len(s.images) - 1
	} else if frame < 0 {
		frame = 0
	}
	s.image = s.images[frame]
	s.frame = frame
}

func (s *Sprite) AddImage(image *ebiten.Image) {
	s.images = append(s.images, image)
}

func (s *Sprite) SetImage(image *ebiten.Image) {
	if len(s.images) == 0 {
		s.images = append(s.images, s.image)
	} else {
		s.images[s.frame] = s.image
	}
	s.image = image
}

func (s *Sprite) Image() *ebiten.Image {
	return s.image
}
func (s *Sprite) Images() []*ebiten.Image {
	return s.images
}

func (s *Sprite) Update() {
	if s.Framerate > 0 && len(s.images) > 0 {
		s.elapsed++
		if s.elapsed >= s.Framerate {
			s.elapsed = 0
			if s.frame+1 >= len(s.images) {
				if s.Loop {
					s.frame = 0
				} else {
					return
				}
			} else {
				s.frame++
			}
			s.image = s.images[s.frame]
		}
	}
}

func (s *Sprite) Draw(ctx states.DrawContext) {
	if s.Hidden {
		return
	}
	s.DrawWithOptions(ctx, &ebiten.DrawImageOptions{})
}

func (s *Sprite) DrawWithOptions(ctx states.DrawContext, opts *ebiten.DrawImageOptions) {
	s.vfx.Process(ctx, opts)
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

	if s.Flipped {
		s.Options.GeoM.Scale(-1, 1)
		s.Options.GeoM.Translate(s.Width(), 0)
	}

	if s.Scale != 0 {
		s.Options.GeoM.Scale(s.Scale, s.Scale)
	}

	if s.Centered {
		s.Options.GeoM.Translate(s.interpX-s.Width()/2, s.interpY-s.Height()/2)
	} else {
		s.Options.GeoM.Translate(s.interpX, s.interpY)
	}

	s.Options.GeoM.Concat(opts.GeoM)

	if opts.ColorScale.A() != 1 || opts.ColorScale.R() != 1 || opts.ColorScale.G() != 1 || opts.ColorScale.B() != 1 {
		s.Options.ColorScale = opts.ColorScale
	} else {
		s.Options.ColorScale.Reset()
	}

	ctx.Screen.DrawImage(s.image, &s.Options)
}

func (s *Sprite) Hit(x, y float64) bool {
	if s.Centered {
		x += s.Width() / 2
		y += s.Height() / 2
	}
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

func NewAnimatedSprite(images []*ebiten.Image) *Sprite {
	return &Sprite{
		image:  images[0],
		images: images,
	}
}
