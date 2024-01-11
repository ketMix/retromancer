package game

import (
	"image/color"
	"math"

	"github.com/ketMix/retromancer/states"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/ketMix/retromancer/resources"
)

type BulletType string

// Defines the different types of rendered bullets
const (
	Circular    BulletType = "circular"    // Default (●)
	Directional            = "directional" // Bi-directional <●>
	Vector                 = "vector"      // Vector ●>
)

type Bullet struct {
	Shape           CircleShape
	bulletType      BulletType
	TargetActor     Actor       // Target actor to head towards
	Speed           float64     // How fastum the bullet goes
	Angle           float64     // What angle the bullet has
	Acceleration    float64     // How fast the bullet accelerates
	AccelAccel      float64     // How fast the bullet accelerates its acceleration
	MinSpeed        float64     // Minimum speed of the bullet
	MaxSpeed        float64     // Maximum speed of the bullet
	AngularVelocity float64     // How fast the bullet rotates
	Color           color.Color // Color of the bullet
	borderColor     color.Color // Color of the border
	aimDelay        int         // How long the bullet should wait before aiming at player
	aimTime         int         // How long the bullet should aim at player
	reversed        bool        // If the bullet has been reversed.
	deflected       bool        // If the bullet has been deflected.
	friendly        bool
	holdFor         int       // An amount of time to hold the bullet in place.
	timeLine        []*Bullet // Positions the bullet has been in
	nextParticle    int       // Next particle to spawn (negative upwards)
	sprite          *resources.Sprite
	Lifetime        int
	Deathtime       int // Maximum lifetime of the bullet
	Destroyed       bool
	Damage          int
}

// TODO: do this differently, hard to read and write arguments
func CreateBullet(
	bulletType BulletType,
	clr color.Color,
	radius, speed, angle, acceleration, accelAccel, minSpeed, maxSpeed, angularVelocity float64,
	aimTime, aimDelay int,
) *Bullet {
	b := &Bullet{
		Shape:           CircleShape{Radius: radius},
		bulletType:      bulletType,
		Speed:           speed,
		Acceleration:    acceleration,
		AccelAccel:      accelAccel,
		Angle:           angle,
		MinSpeed:        minSpeed,
		MaxSpeed:        maxSpeed,
		AngularVelocity: angularVelocity,
		Color:           clr,
		borderColor:     color.White,
		aimTime:         aimTime,
		aimDelay:        aimDelay,
		timeLine:        make([]*Bullet, 0),
		Damage:          5,
		Deathtime:       500,
	}
	b.sprite = resources.NewSprite(ebiten.NewImage(int(radius*2), int(radius*2)))
	return b
}

// Copy a bullet
func BulletFromExisting(b *Bullet, angle float64) *Bullet {
	bullet := CreateBullet(
		b.bulletType,
		b.Color,
		b.Shape.Radius,
		b.Speed,
		angle,
		b.Acceleration,
		b.AccelAccel,
		b.MinSpeed,
		b.MaxSpeed,
		b.AngularVelocity,
		b.aimTime,
		b.aimDelay,
	)
	bullet.Damage = b.Damage
	bullet.SetXY(b.Shape.X, b.Shape.Y)
	return bullet
}

func CreateBulletFromDef(override, alias *resources.Bullet) *Bullet {
	// Create a bullet group from a bullet group definition
	// Use override values if they exist
	// TODO: maybe have default values if properties aren't present in alias or override

	bulletType := *alias.BulletType
	c := *alias.Color
	radius := float64(*alias.Radius)
	speed := float64(*alias.Speed)
	acceleration := float64(*alias.Acceleration)
	accelAccel := float64(*alias.AccelAccel)
	minSpeed := float64(*alias.MinSpeed)
	maxSpeed := float64(*alias.MaxSpeed)
	angularVelocity := float64(*alias.AngularVelocity)
	aimTime := *alias.AimTime
	aimDelay := *alias.AimDelay
	damage := 5
	if alias.Damage != nil {
		damage = *alias.Damage
	}

	if override != nil {
		if override.BulletType != nil {
			bulletType = *override.BulletType
		}
		if override.Color != nil {
			c = *override.Color
		}
		if override.Radius != nil {
			radius = float64(*override.Radius)
		}
		if override.Speed != nil {
			speed = float64(*override.Speed)
		}
		if override.Acceleration != nil {
			acceleration = float64(*override.Acceleration)
		}
		if override.AccelAccel != nil {
			accelAccel = float64(*override.AccelAccel)
		}
		if override.MinSpeed != nil {
			minSpeed = float64(*override.MinSpeed)
		}
		if override.MaxSpeed != nil {
			maxSpeed = float64(*override.MaxSpeed)
		}
		if override.AngularVelocity != nil {
			angularVelocity = float64(*override.AngularVelocity)
		}
		if override.AimTime != nil {
			aimTime = *override.AimTime
		}
		if override.AimDelay != nil {
			aimDelay = *override.AimDelay
		}
		if override.Damage != nil {
			damage = *override.Damage
		}
	}
	color := color.RGBA{uint8(c[0]), uint8(c[1]), uint8(c[2]), uint8(c[3])}

	if damage == 0 {
		damage = 5
	}

	bullet := CreateBullet(
		BulletType(bulletType),
		color,
		radius,
		speed,
		0,
		acceleration,
		accelAccel,
		minSpeed,
		maxSpeed,
		angularVelocity,
		aimTime,
		aimDelay,
	)
	bullet.Damage = damage
	return bullet
}

func (b *Bullet) SetXY(x, y float64) {
	b.Shape.X = x
	b.Shape.Y = y
	b.sprite.SetXY(x, y)
}

// Update the bullet's position and speed
func (b *Bullet) Update() (actions []Action) {
	// Only do this if the bullet is not deflected
	if !b.deflected && b.Deathtime > 0 {
		b.Lifetime++
		if b.Lifetime > b.Deathtime {
			b.Destroyed = true
			return
		}
	}

	if len(b.timeLine) == 1 && b.reversed {
		// if we're at the first point in timeLine, use the bullet as current bullet
		prevBullet := b.timeLine[0]
		b.timeLine = b.timeLine[:0]
		b.borderColor = prevBullet.borderColor
		b.Speed = prevBullet.Speed
		b.Angle = prevBullet.Angle
		b.Acceleration = prevBullet.Acceleration
		b.AngularVelocity = prevBullet.AngularVelocity
		b.reversed = false
	}

	if b.holdFor > 0 {
		b.holdFor--
		return actions
	}

	if b.reversed && len(b.timeLine) > 0 {
		// Get previous bullet and remove it from the timeline
		prevBullet := b.timeLine[len(b.timeLine)-1]
		b.timeLine = b.timeLine[:len(b.timeLine)-1]

		// Set properties of the bullet
		b.Speed = prevBullet.Speed

		// Move bullet towards previous position, but keep it facing the same direction as previous bullet
		movementAngle := math.Atan2(prevBullet.Shape.Y-b.Shape.Y, prevBullet.Shape.X-b.Shape.X)
		b.Angle = prevBullet.Angle
		b.aimTime = prevBullet.aimTime
		b.aimDelay = prevBullet.aimDelay
		b.Shape.X += b.Speed * math.Cos(movementAngle)
		b.sprite.X = b.Shape.X
		b.Shape.Y += b.Speed * math.Sin(movementAngle)
		b.sprite.Y = b.Shape.Y
		return actions
	}

	b.Speed += b.Acceleration
	b.Acceleration += b.AccelAccel

	if b.Speed > b.MaxSpeed {
		b.Speed = b.MaxSpeed
	} else if b.Speed < b.MinSpeed {
		b.Speed = b.MinSpeed
	}
	b.Shape.X += b.Speed * math.Cos(b.Angle)
	b.sprite.X = b.Shape.X
	b.Shape.Y += b.Speed * math.Sin(b.Angle)
	b.sprite.Y = b.Shape.Y

	// Spawn a lil particle.
	b.nextParticle++
	if b.nextParticle >= 0 {
		actions = append(actions, ActionSpawnParticle{
			X:     b.Shape.X,
			Y:     b.Shape.Y,
			Angle: b.Angle + math.Pi,
			Speed: 0.05,
			Img:   "bullet",
			Life:  20,
		})

		b.nextParticle = -2
	}

	// Decrement delay
	if b.aimDelay > 0 {
		b.aimDelay--
	}

	// Add bullet to timeline if not deflected
	if !b.deflected {
		b.timeLine = append(b.timeLine, BulletFromExisting(b, b.Angle))
	}

	// If we're not aiming at the player yet, adjust angle by angular velocity.
	if b.aimDelay > 0 || b.aimTime <= 0 {
		b.Angle += b.AngularVelocity
		return actions
	}

	if b.aimTime > 0 {
		// Disable angular velocity.
		b.AngularVelocity = 0
		// Request closest player actor for next tick.
		if b.TargetActor == nil {
			actions = append(actions, ActionFindNearestActor{Actor: (*PC)(nil)})
		} else {
			// Aim at closest actor.
			x, y, _, _ := b.TargetActor.Bounds()
			b.Angle = math.Atan2(y-b.Shape.Y, x-b.Shape.X)
		}
		b.aimTime--
	}

	return actions
}

func (b *Bullet) Reverse() {
	if b.reversed {
		return
	}

	// Reset lifetime
	b.Lifetime = 0

	// If the bullet is friendly, empower it.
	if b.friendly {
		b.Damage = 2
		b.Shape.Radius = 2
	}

	// Turn the bullet border REVERSE color
	b.borderColor = color.NRGBA{0x66, 0x99, 0xff, 0xff}

	// Stop aiming the bullet if it was aimed. Perhaps this should deflect the bullet towards the spawner that created it.
	b.reversed = true
}

func (b *Bullet) Deflect(angle float64) {
	if b.deflected {
		return
	}

	// If the bullet is friendly, empower it.
	if b.friendly {
		b.Damage = 2
		b.Shape.Radius = 2
	}

	// Turn the bullet border DEFLECT color
	b.borderColor = color.NRGBA{0xff, 0x66, 0x99, 0xff}

	// Stop aiming the bullet if it was aimed. Perhaps this should deflect the bullet towards the spawner that created it.
	b.aimTime = 0
	// FIXME: Deflect should take into account the bullet's angle relative to the deflection angle and use that for the final angle.
	b.Angle = angle
	b.AngularVelocity = 0
	b.deflected = true
}

// Draw the bullet
func (b *Bullet) Draw(ctx states.DrawContext) {
	// Draw base bullet
	vector.DrawFilledCircle(ctx.Screen, float32(b.sprite.X), float32(b.sprite.Y), float32(b.Shape.Radius), b.Color, false)

	// Draw the border depending on its type
	switch b.bulletType {
	case Circular:
		// Draw circle border? Bit too visually noisy.
		vector.StrokeCircle(ctx.Screen, float32(b.sprite.X), float32(b.sprite.Y), float32(b.Shape.Radius)*1.2, 1, b.borderColor, false)
	case Directional:
		// Draw V shape on both ends
		vector.StrokeLine(
			ctx.Screen,
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)*2),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)*2),
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)+b.Shape.Radius*math.Cos(b.Angle+math.Pi/2)),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)+b.Shape.Radius*math.Sin(b.Angle+math.Pi/2)),
			1,
			b.borderColor,
			false,
		)
		vector.StrokeLine(
			ctx.Screen,
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)*2),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)*2),
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)-b.Shape.Radius*math.Cos(b.Angle+math.Pi/2)),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)-b.Shape.Radius*math.Sin(b.Angle+math.Pi/2)),
			1,
			b.borderColor,
			false,
		)
		vector.StrokeLine(
			ctx.Screen,
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)*-2),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)*-2),
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)+b.Shape.Radius*math.Cos(b.Angle+math.Pi/2)),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)+b.Shape.Radius*math.Sin(b.Angle+math.Pi/2)),
			1,
			b.borderColor,
			false,
		)
		vector.StrokeLine(
			ctx.Screen,
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)*-2),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)*-2),
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)-b.Shape.Radius*math.Cos(b.Angle+math.Pi/2)),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)-b.Shape.Radius*math.Sin(b.Angle+math.Pi/2)),
			1,
			b.borderColor,
			false,
		)
	case Vector:
		// Circular border + V shape
		// Should be drawn on the edge in the direction of the bullet's angle
		vector.StrokeCircle(ctx.Screen, float32(b.sprite.X), float32(b.sprite.Y), float32(b.Shape.Radius)*1.2, 1, b.borderColor, false)

		// TODO: replace with triangle in borderColor
		vector.StrokeLine(
			ctx.Screen,
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)*3),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)*3),
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)+b.Shape.Radius*math.Cos(b.Angle+math.Pi/2)),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)+b.Shape.Radius*math.Sin(b.Angle+math.Pi/2)),
			1,
			b.borderColor,
			false,
		)
		vector.StrokeLine(
			ctx.Screen,
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)*3),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)*3),
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)-b.Shape.Radius*math.Cos(b.Angle+math.Pi/2)),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)-b.Shape.Radius*math.Sin(b.Angle+math.Pi/2)),
			1,
			b.borderColor,
			false,
		)
	}
}
