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
	Shape           Shape
	Type            BulletType
	Speed           float64 // How fastum the bullet goes
	Angle           float64 // What angle the bullet has
	Acceleration    float64 // How fast the bullet accelerates
	MinSpeed        float64 // Minimum speed of the bullet
	MaxSpeed        float64 // Maximum speed of the bullet
	AngularVelocity float64 // How fast the bullet rotates
	Color           color.Color
	aimed           bool // If the bullet is aimed at a player (TODO: maybe do differently)
	sprite          *resources.Sprite
}

// TODO: do this differently, hard to read and write arguments
func CreateBullet(x, y, radius float64, bulletType BulletType, speed, angle, acceleration, minSpeed, maxSpeed, angularVelocity float64, color color.Color) *Bullet {
	b := &Bullet{
		Shape:           Shape{X: x, Y: y, Radius: radius},
		Type:            bulletType,
		Speed:           speed,
		Angle:           angle,
		Acceleration:    acceleration,
		MinSpeed:        minSpeed,
		MaxSpeed:        maxSpeed,
		AngularVelocity: angularVelocity,
		Color:           color,
	}
	b.sprite = resources.NewSprite(ebiten.NewImage(int(radius*2), int(radius*2)))
	b.sprite.X = x
	b.sprite.Y = y
	return b
}

// Copy a bullet
func BulletFromExisting(b *Bullet) *Bullet {
	bullet := CreateBullet(
		b.Shape.X,
		b.Shape.Y,
		b.Shape.Radius,
		b.Type,
		b.Speed,
		b.Angle,
		b.Acceleration,
		b.MinSpeed,
		b.MaxSpeed,
		b.AngularVelocity,
		b.Color,
	)
	return bullet
}

// Update the bullet's position and speed
// Pass in a slice of players to find the closest player to aim at
// Maybe do this differently
func (b *Bullet) Update(p []Player) {
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

	b.Speed += b.Acceleration

	if !b.aimed {
		b.Angle += b.AngularVelocity
	} else {
		// Find closest player
		var closestPlayer *Player
		var closestDistance float64
		var playerX, playerY float64
		for _, player := range p {
			x, y, _, _ := player.Actor().Bounds()
			distance := math.Sqrt(math.Pow(b.Shape.X-x, 2) + math.Pow(b.Shape.Y-y, 2))
			if closestPlayer == nil || distance < closestDistance {
				closestPlayer = &player
				closestDistance = distance
				playerX = x
				playerY = y
			}
		}
		// Aim at closest player
		// Need to add some momentum so it doesn't just follow the player
		b.Angle = math.Atan2(playerY-b.Shape.Y, playerX-b.Shape.X)
	}
}

// Draw the bullet
func (b *Bullet) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, float32(b.sprite.X), float32(b.sprite.Y), float32(b.Shape.Radius), b.Color, false)
}

func (b *Bullet) OutOfBounds() bool {
	w, h := ebiten.WindowSize()
	return b.Shape.X < 0 || b.Shape.X > float64(w) || b.Shape.Y < 0 || b.Shape.Y > float64(h)
}
