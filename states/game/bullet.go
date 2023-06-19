package game

import (
	"ebijam23/resources"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	Type            BulletType
	TargetActor     Actor   // Target actor to head towards
	Speed           float64 // How fastum the bullet goes
	Angle           float64 // What angle the bullet has
	Acceleration    float64 // How fast the bullet accelerates
	AccelAccel      float64 // How fast the bullet accelerates its acceleration
	MinSpeed        float64 // Minimum speed of the bullet
	MaxSpeed        float64 // Maximum speed of the bullet
	AngularVelocity float64 // How fast the bullet rotates
	Color           color.Color
	aimDelay        int  // How long the bullet should wait before aiming at player
	aimTime         int  // How long the bullet should aim at player
	reflected       bool // If the bullet has been reflected
	deflected       bool // If the bullet has been deflected.
	sprite          *resources.Sprite
}

// TODO: do this differently, hard to read and write arguments
func CreateBullet(
	bulletType BulletType,
	color color.Color,
	x, y, radius, speed, acceleration, accelAccel, minSpeed, maxSpeed, angularVelocity float64,
	aimTime, aimDelay int,
) *Bullet {
	b := &Bullet{
		Shape:        CircleShape{X: x, Y: y, Radius: radius},
		Type:         bulletType,
		Speed:        speed,
		Acceleration: acceleration,
		AccelAccel:   accelAccel,
		// Angle:           angle, Set by spawner
		MinSpeed:        minSpeed,
		MaxSpeed:        maxSpeed,
		AngularVelocity: angularVelocity,
		Color:           color,
		aimTime:         aimTime,
		aimDelay:        aimDelay,
	}
	b.sprite = resources.NewSprite(ebiten.NewImage(int(radius*2), int(radius*2)))
	b.sprite.X = x
	b.sprite.Y = y
	return b
}

// Copy a bullet
func BulletFromExisting(b *Bullet) *Bullet {
	bullet := CreateBullet(
		b.Type,
		b.Color,
		b.Shape.X,
		b.Shape.Y,
		b.Shape.Radius,
		b.Speed,
		b.Acceleration,
		b.AccelAccel,
		b.MinSpeed,
		b.MaxSpeed,
		b.AngularVelocity,
		b.aimTime,
		b.aimDelay,
	)
	return bullet
}

// Update the bullet's position and speed
func (b *Bullet) Update() (actions []Action) {
	b.Speed += b.Acceleration
	b.Acceleration += b.AccelAccel

	if b.Speed < b.MinSpeed {
		b.Speed = b.MinSpeed
	}
	if b.Speed > b.MaxSpeed {
		b.Speed = b.MaxSpeed
	}
	b.Shape.X += b.Speed * math.Cos(b.Angle)
	b.sprite.X = b.Shape.X
	b.Shape.Y += b.Speed * math.Sin(b.Angle)
	b.sprite.Y = b.Shape.Y

	// Decrement delay
	if b.aimDelay > 0 {
		b.aimDelay--
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

func (b *Bullet) Reflect() {
	if b.reflected {
		return
	}
	// Stop aiming the bullet if it was aimed. Perhaps this should deflect the bullet towards the spawner that created it.
	b.aimTime = 0
	b.Angle = math.Mod(b.Angle+math.Pi, 2*math.Pi)
	b.reflected = true
}

func (b *Bullet) Deflect(angle float64) {
	if b.deflected {
		return
	}
	// Stop aiming the bullet if it was aimed. Perhaps this should deflect the bullet towards the spawner that created it.
	b.aimTime = 0
	// FIXME: Deflect should take into account the bullet's angle relative to the deflection angle and use that for the final angle.
	b.Angle = angle
	b.deflected = true
}

// Draw the bullet
func (b *Bullet) Draw(screen *ebiten.Image) {
	// Draw base bullet
	vector.DrawFilledCircle(screen, float32(b.sprite.X), float32(b.sprite.Y), float32(b.Shape.Radius), b.Color, false)

	// Draw the border depending on its type
	switch b.Type {
	case Circular:
		// Draw circle border? Bit too visually noisy.
		// vector.StrokeCircle(screen, float32(b.sprite.X), float32(b.sprite.Y), float32(b.Shape.Radius)+2, 1, color.White, false)
		return
	case Directional:
		// Draw V shape on both ends
		vector.StrokeLine(
			screen,
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)*2),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)*2),
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)+b.Shape.Radius*math.Cos(b.Angle+math.Pi/2)),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)+b.Shape.Radius*math.Sin(b.Angle+math.Pi/2)),
			1,
			color.White,
			false,
		)
		vector.StrokeLine(
			screen,
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)*2),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)*2),
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)-b.Shape.Radius*math.Cos(b.Angle+math.Pi/2)),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)-b.Shape.Radius*math.Sin(b.Angle+math.Pi/2)),
			1,
			color.White,
			false,
		)
		vector.StrokeLine(
			screen,
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)*-2),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)*-2),
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)+b.Shape.Radius*math.Cos(b.Angle+math.Pi/2)),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)+b.Shape.Radius*math.Sin(b.Angle+math.Pi/2)),
			1,
			color.White,
			false,
		)
		vector.StrokeLine(
			screen,
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)*-2),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)*-2),
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)-b.Shape.Radius*math.Cos(b.Angle+math.Pi/2)),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)-b.Shape.Radius*math.Sin(b.Angle+math.Pi/2)),
			1,
			color.White,
			false,
		)
	case Vector:
		// Draw V shape
		// Should be drawn on the edge in the direction of the bullet's angle
		vector.StrokeLine(
			screen,
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)*2),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)*2),
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)+b.Shape.Radius*math.Cos(b.Angle+math.Pi/2)),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)+b.Shape.Radius*math.Sin(b.Angle+math.Pi/2)),
			1,
			color.White,
			false,
		)
		vector.StrokeLine(
			screen,
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)*2),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)*2),
			float32(b.sprite.X+b.Shape.Radius*math.Cos(b.Angle)-b.Shape.Radius*math.Cos(b.Angle+math.Pi/2)),
			float32(b.sprite.Y+b.Shape.Radius*math.Sin(b.Angle)-b.Shape.Radius*math.Sin(b.Angle+math.Pi/2)),
			1,
			color.White,
			false,
		)
	}
}

func (b *Bullet) OutOfBounds() bool {
	w, h := ebiten.WindowSize()
	return b.Shape.X < 0 || b.Shape.X > float64(w) || b.Shape.Y < 0 || b.Shape.Y > float64(h)
}
