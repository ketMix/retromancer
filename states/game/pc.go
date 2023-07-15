package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	playerStartLives = 3
	playerMaxLives   = playerStartLives + 2
)

type PC struct {
	player Player
	saved  *PC
	//
	Arrow                     *resources.Sprite
	Sprite                    *resources.Sprite
	DeathSprite               *resources.Sprite
	Phylactery                *resources.Sprite
	Hat                       *resources.Sprite
	Life                      *resources.Sprite
	shape                     CircleShape
	Hand                      Hand
	Lives                     int
	InvulnerableTicks         int // Ticks the player should be invulnerable for
	TicksSinceLastInteraction int
	Energy                    int
	MaxEnergy                 int
	EnergyRestoreRate         int
	HasDeflect                bool
	HasShield                 bool
	//
	shielding bool
	//
	previousInteraction Action
	//
	impulses    ImpulseSet
	resurrected bool
	//
	momentumX float64
	momentumY float64
	//
	audioPlayer *audio.Player
	currentSfx  *resources.Sound
	reverseSfx  *resources.Sound
	deflectSfx  *resources.Sound
	shieldSfx   *resources.Sound
	hurtSfx     *resources.Sound
}

func (s *World) NewPC(ctx states.Context) *PC {
	pc := &PC{
		Sprite:            resources.NewSprite(ctx.R.GetAs("images", "player", (*ebiten.Image)(nil)).(*ebiten.Image)),
		DeathSprite:       resources.NewSprite(ctx.R.GetAs("images", "player-dead1", (*ebiten.Image)(nil)).(*ebiten.Image)),
		Phylactery:        resources.NewSprite(ctx.R.GetAs("images", "phylactery", (*ebiten.Image)(nil)).(*ebiten.Image)),
		Arrow:             resources.NewSprite(ctx.R.GetAs("images", "direction-arrow", (*ebiten.Image)(nil)).(*ebiten.Image)),
		Life:              resources.NewSprite(ctx.R.GetAs("images", "life", (*ebiten.Image)(nil)).(*ebiten.Image)),
		Energy:            100,
		MaxEnergy:         100,
		EnergyRestoreRate: 2,
		Lives:             playerStartLives,
		reverseSfx:        ctx.R.GetAs("sounds", "reverse-sfx", (*resources.Sound)(nil)).(*resources.Sound),
		deflectSfx:        ctx.R.GetAs("sounds", "deflect-sfx", (*resources.Sound)(nil)).(*resources.Sound),
		shieldSfx:         ctx.R.GetAs("sounds", "shield-sfx", (*resources.Sound)(nil)).(*resources.Sound),
		hurtSfx:           ctx.R.GetAs("sounds", "hurt-sfx", (*resources.Sound)(nil)).(*resources.Sound),
	}

	// FIXME: This shouldn't be hardcoded.
	pc.DeathSprite.Framerate = 2
	pc.DeathSprite.Centered = true
	for i := 1; i <= 11; i++ {
		pc.DeathSprite.AddImage(ctx.R.GetAs("images", fmt.Sprintf("player-dead%d", i), (*ebiten.Image)(nil)).(*ebiten.Image))
	}

	pc.shape.Radius = 2
	//pc.Sprite.Interpolate = true
	pc.Sprite.Centered = true
	pc.Hand.Sprite = resources.NewSprite(ctx.R.GetAs("images", "hand-normal", (*ebiten.Image)(nil)).(*ebiten.Image))
	pc.Hand.Sprite.Centered = true
	pc.Hand.HoverSprite = resources.NewSprite(ctx.R.GetAs("images", "hand-glow", (*ebiten.Image)(nil)).(*ebiten.Image))
	pc.Hand.HoverSprite.Centered = true

	return pc
}

func (p *PC) Save() {
	saved := *p
	saved.DeathSprite.Reset()
	p.saved = &saved
}

func (p *PC) Restore() {
	*p = *p.saved
}

func (p *PC) SetPlayer(player Player) {
	p.player = player
}

func (p *PC) Player() Player {
	return p.player
}

func (p *PC) Update() (actions []Action) {
	if p.Lives < 0 {
		p.DeathSprite.Update()
		return
	}

	p.InvulnerableTicks--

	p.TicksSinceLastInteraction++
	if p.TicksSinceLastInteraction > 20 {
		if p.Energy+p.EnergyRestoreRate <= p.MaxEnergy {
			p.Energy += p.EnergyRestoreRate
		}
	}
	// Do not handle movement until we are resurrected.
	if p.resurrected {
		if p.impulses.Move != nil {
			p.momentumX = 0.3*p.momentumX + 3.7*math.Cos((*p.impulses.Move).Direction)
			p.momentumY = 0.3*p.momentumY + 3.7*math.Sin((*p.impulses.Move).Direction)
			/*x := 5 * math.Cos((*p.impulses.Move).Direction)
			y := 5 * math.Sin((*p.impulses.Move).Direction)
			if x < 0 {
				p.Sprite.Flipped = true
			} else if x > 0 {
				p.Sprite.Flipped = false
			}
			actions = append(actions, ActionMove{
				X: p.shape.X + x,
				Y: p.shape.Y + y,
			})*/
		}
		if p.momentumX != 0 || p.momentumY != 0 {
			actions = append(actions, ActionMove{
				X: p.shape.X + p.momentumX,
				Y: p.shape.Y + p.momentumY,
			})
		}
		p.momentumX *= 0.3
		p.momentumY *= 0.3
		if math.Abs(p.momentumX) < 0.1 {
			p.momentumX = 0
		}
		if math.Abs(p.momentumY) < 0.1 {
			p.momentumY = 0
		}

		if p.Hand.Shape.X < p.shape.X {
			p.Sprite.Flipped = true
		} else {
			p.Sprite.Flipped = false
		}
	}

	p.previousInteraction = nil
	if p.impulses.Interaction != nil {
		switch imp := p.impulses.Interaction.(type) {
		case ImpulseReverse:
			if p.HasEnergyFor(imp) {
				p.Energy -= imp.Cost()
				p.TicksSinceLastInteraction = 0
				actions = append(actions, ActionReverse{
					X: imp.X,
					Y: imp.Y,
				})
				p.previousInteraction = ActionReverse{}
			}
		case ImpulseDeflect:
			if p.HasDeflect && p.HasEnergyFor(imp) {
				p.Energy -= imp.Cost()
				p.TicksSinceLastInteraction = 0
				actions = append(actions, ActionDeflect{
					Direction: math.Atan2(imp.Y-p.shape.Y, imp.X-p.shape.X),
					X:         imp.X,
					Y:         imp.Y,
				})
				p.previousInteraction = ActionDeflect{}
			}
		case ImpulseShield:
			if p.HasShield && p.HasEnergyFor(imp) {
				p.Energy -= imp.Cost()
				p.TicksSinceLastInteraction = 0
				actions = append(actions, ActionShield{})
				p.previousInteraction = ActionShield{}
			}
		default:
			// Do nothing.
		}
	}

	return actions
}

func (p *PC) Dead() bool {
	return p.Lives < 0
}

func (p *PC) HasEnergyFor(imp Impulse) bool {
	return p.Energy-imp.Cost() >= 0
}

func (p *PC) SetImpulses(impulses ImpulseSet) {
	p.impulses = impulses
}

// DrawHat draws the player's dumb hat.
func (p *PC) DrawHat(screen *ebiten.Image, x, y float64) {
	opts := &ebiten.DrawImageOptions{}
	if p.Sprite.Flipped {
		opts.GeoM.Scale(-1, 1)
		opts.GeoM.Translate(p.Hat.Width(), 0)
	}
	opts.GeoM.Translate(p.shape.X-float64(int(p.Hat.Width())/2)+x, p.Sprite.Y-p.Sprite.Height()/2-p.Hat.Height()+y)
	screen.DrawImage(p.Hat.Image(), opts)
}

func (p *PC) DrawHand(ctx states.DrawContext) {
	p.Hand.Sprite.Draw(ctx)
	p.Hand.HoverSprite.Draw(ctx)
}

func (p *PC) Draw(ctx states.DrawContext) {
	if !p.resurrected {
		// TODO: Actually play the death sprite backwards.
		p.DeathSprite.SetFrame(99)
		p.DeathSprite.SetXY(p.Sprite.X, p.Sprite.Y)
		p.DeathSprite.Draw(ctx)

		p.DrawHat(ctx.Screen, 0, 3+8)
	} else if p.Dead() {
		// Hackiness ahoy. This is hard coded to the exact death sprite frames and positions.
		p.DeathSprite.SetXY(p.Sprite.X, p.Sprite.Y)
		p.DeathSprite.Draw(ctx)
		y := 0.0
		switch p.DeathSprite.Frame() {
		case 9:
			y = 1.0
		case 10:
			y = 8
		}
		p.DrawHat(ctx.Screen, 0, 3+y)
		return
	}

	if _, ok := p.previousInteraction.(ActionDeflect); ok {
		vector.DrawFilledCircle(ctx.Screen, float32(p.Hand.Shape.X), float32(p.Hand.Shape.Y), 20, color.NRGBA{0xff, 0x66, 0x99, 0x33}, false)
	} else if _, ok := p.previousInteraction.(ActionReverse); ok {
		vector.DrawFilledCircle(ctx.Screen, float32(p.Hand.Shape.X), float32(p.Hand.Shape.Y), 20, color.NRGBA{0x66, 0x99, 0xff, 0x33}, false)
	} else if _, ok := p.previousInteraction.(ActionShield); ok {
		vector.DrawFilledCircle(ctx.Screen, float32(p.shape.X), float32(p.shape.Y), 20, color.NRGBA{0x66, 0xff, 0x99, 0x33}, false)
	}

	opts := &ebiten.DrawImageOptions{}

	if (p.InvulnerableTicks <= 0 || p.InvulnerableTicks%6 >= 3) && p.resurrected {
		p.Sprite.Draw(ctx)

		// Draw the player's phylactery (hit box representation). If the player has 0 lives, hide it, since it "broke"
		if p.Lives > 0 {
			opts.GeoM.Translate(p.shape.X-float64(int(p.Phylactery.Width())/2), p.shape.Y-float64(int(p.Phylactery.Height())/2))
			opts.ColorScale.Scale(0.5, 0.5, 1.0, 1.0)
			ctx.Screen.DrawImage(p.Phylactery.Image(), opts)
		}

		p.DrawHat(ctx.Screen, 0, 3)
	}

	if !p.resurrected {
		return
	}

	// Draw lives?
	for i := 0; i < p.Lives; i++ {
		opts = &ebiten.DrawImageOptions{}
		x := -(float64(p.Lives) * (p.Life.Width() + 1)) / 2
		x += p.Hand.Shape.X + float64(i)*(p.Life.Width()+1)
		y := p.Hand.Shape.Y + p.Hand.Sprite.Height()/2 + 3
		opts.GeoM.Translate(x, y)
		ctx.Screen.DrawImage(p.Life.Image(), opts)
	}

	r := math.Atan2(p.shape.Y-p.Hand.Shape.Y, p.shape.X-p.Hand.Shape.X)

	// TODO: The arrow image should change based on if we're reversing or deflecting.
	// Draw direction arrow
	opts = &ebiten.DrawImageOptions{}
	// Rotate about its center.
	opts.GeoM.Translate(-p.Arrow.Width()/2, -p.Arrow.Height()/2)
	opts.GeoM.Rotate(r)
	// Position.
	//opts.GeoM.Translate(p.Arrow.Width()/2, p.Arrow.Height()/2)
	opts.GeoM.Translate(p.shape.X, p.shape.Y)
	// Draw from center.
	ctx.Screen.DrawImage(p.Arrow.Image(), opts)
}

func (p *PC) Bounds() (x, y, w, h float64) {
	// Return the radius of the shape as width and height in diameter.
	return p.shape.X, p.shape.Y, p.shape.Radius * 2, p.shape.Radius * 2
}

func (p *PC) Shape() Shape {
	return &p.shape
}

func (p *PC) SetXY(x, y float64) {
	p.shape.X = x
	p.shape.Y = y
	p.Sprite.X = x
	p.Sprite.Y = y + 2 // We lightly offset the sprite so the phylactery is in a nicer visual position.
}

func (p *PC) SetSize(r float64) {
	p.shape.Radius = r
}

func (p *PC) Destroyed() bool {
	return false
}

func (p *PC) Hurtie() {
	if p.InvulnerableTicks <= 0 {
		p.Lives--
		p.InvulnerableTicks = 40
		p.hurtSfx.Play(0.5)
	}
}

func (p *PC) PlayAudio(deflecting, reversing, shielding bool) {
	// Stop playing sfx
	if !deflecting && !reversing && !shielding {
		if p.audioPlayer != nil {
			p.audioPlayer.Pause()
			p.audioPlayer.Close()
			p.audioPlayer = nil
			p.currentSfx = nil
		}
		return
	}

	// Set reverse sfx
	if reversing && p.currentSfx != p.reverseSfx {
		if p.audioPlayer != nil {
			p.audioPlayer.Pause()
			p.audioPlayer.Close()
		}
		p.audioPlayer = p.reverseSfx.Play(1.0)
		p.currentSfx = p.reverseSfx
	}

	// Set deflect sfx
	if deflecting && p.currentSfx != p.deflectSfx {
		if p.audioPlayer != nil {
			p.audioPlayer.Pause()
			p.audioPlayer.Close()
		}
		p.audioPlayer = p.deflectSfx.Play(1.0)
		p.currentSfx = p.deflectSfx
	}

	// Set shield sfx
	if shielding && p.currentSfx != p.shieldSfx {
		if p.audioPlayer != nil {
			p.audioPlayer.Pause()
			p.audioPlayer.Close()
		}
		p.audioPlayer = p.shieldSfx.Play(1.0)
		p.currentSfx = p.shieldSfx
	}
}
