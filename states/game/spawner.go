package game

import (
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type GroupAngle string

const (
	Aimed  GroupAngle = "aimed"  // Aim at nearest player
	Fixed             = "fixed"  // Fixed angle
	Radial            = "radial" // Radial angle from spawner
	Random            = "random" // Random angle
)

type BulletGroup struct {
	bullet        *Bullet       // What bullet comprises this group
	angle         GroupAngle    // What angle to spawn bullets at
	spawnRate     int           // How often to spawn bullets
	lastSpawnedAt int           // How long since spawn
	bulletCount   int           // How many bullets to spawn
	loopCount     int           // How many times to loop
	subGroup      []BulletGroup // Subgroups TODO: implement
}

func (bg *BulletGroup) Update() (actions []Action) {
	// Update the bullet group
	// TODO: maybe handle infinite loop differently
	if bg.lastSpawnedAt >= bg.spawnRate && (bg.loopCount > 0 || bg.loopCount == -1) {
		if bg.loopCount > 0 {
			bg.loopCount--
		}
		bg.lastSpawnedAt = 0

		// Spawn new bullets
		// Init bullet array
		bullets := make([]*Bullet, bg.bulletCount)
		for i := 0; i < bg.bulletCount; i++ {
			// Create a new bullet
			bullet := BulletFromExisting(bg.bullet)
			// Set the bullet angle
			switch bg.angle {
			case Aimed:
				bullet.aimed = true
			case Fixed:
				bullet.Angle = 0
			case Radial:
				// Spread each bullet evenly
				bullet.Angle = float64(i) * (360 / float64(bg.bulletCount))
			case Random:
				// TODO: Random angle
				// Generate a random angle
				bullet.Angle = rand.Float64() * 360
			}
			// Add the bullet to the array
			bullets[i] = bullet
		}
		// Create the action to spawn the bullets
		actions = append(actions, ActionSpawnBullets{bullets})
	}
	bg.lastSpawnedAt++
	return actions
}

// This can probably be attached to an actor instead being its own actor
type Spawner struct {
	shape        Shape
	bulletGroups []*BulletGroup
}

// TODO:
// - these bullet groups should be created from external file for specific enemies
func CreateSpawner(x, y float64) *Spawner {
	return &Spawner{
		shape: Shape{X: x, Y: y, Radius: 0},
		bulletGroups: []*BulletGroup{
			// WHITE: Aim radially
			{
				angle:         Radial,
				bulletCount:   24,
				lastSpawnedAt: 10, // Spawn immediately
				spawnRate:     15,
				loopCount:     -1, // Loop forever
				bullet:        CreateBullet(x, y, 3, Circular, 3, 0, 0, 0, 100, 0, color.White),
			},
			// BLUE: Aim radially, but accelerated
			{
				angle:         Radial,
				bulletCount:   12,
				lastSpawnedAt: 10, // Spawn immediately
				spawnRate:     25,
				loopCount:     -1,
				bullet:        CreateBullet(x, y, 3, Circular, 1, 0, 1, 0, 3, 0, color.RGBA{0, 0, 255, 255}),
			},
			// GREEN: Radial with a bit of angular velocity
			{
				angle:         Radial,
				bulletCount:   5,
				lastSpawnedAt: 10, // Spawn immediately
				spawnRate:     20,
				loopCount:     -1,
				bullet:        CreateBullet(x, y, 3, Circular, 5, 0, 0, 0, 100, 0.02, color.RGBA{0, 255, 0, 255}),
			},
			// PURPLE: Cool lil' spiral thing
			{
				angle:         Radial,
				bulletCount:   5,
				lastSpawnedAt: 10, // Spawn immediately
				spawnRate:     20,
				loopCount:     -1,
				bullet:        CreateBullet(x, y, 3, Circular, 5, 0, 0, 0, 100, 0.1, color.RGBA{255, 0, 255, 255}),
			},
			// YELLOW: Random Angle
			{
				angle:         Random,
				bulletCount:   20,
				lastSpawnedAt: 10, // Spawn immediately
				spawnRate:     5,
				loopCount:     -1,
				bullet:        CreateBullet(x, y, 3, Circular, 5, 0, 0, 0, 100, 0, color.RGBA{255, 255, 0, 255}),
			},
			// RED: Aim at player
			{
				angle:       Aimed,
				bulletCount: 1,
				spawnRate:   35,
				loopCount:   -1,
				bullet:      CreateBullet(x, y, 5, Circular, 2, 0, 0, 0, 100, 0, color.RGBA{255, 0, 0, 255}),
			},
		},
	}
}

func (s *Spawner) Update() (actions []Action) {
	// Update the bullet groups
	for _, bg := range s.bulletGroups {
		// Add the actions from the bullet group to the list of actions
		bgActions := bg.Update()
		actions = append(actions, bgActions...)
	}
	return actions
}

func (s *Spawner) Player() Player                  { return nil }
func (s *Spawner) SetPlayer(p Player)              {}
func (s *Spawner) SetImpulses(impulses ImpulseSet) {}
func (s *Spawner) Draw(screen *ebiten.Image)       {}
func (s *Spawner) Bounds() (x, y, w, h float64)    { return 0, 0, 0, 0 }
func (s *Spawner) SetXY(x, y float64)              {}
func (s *Spawner) SetSize(r float64)               {}
