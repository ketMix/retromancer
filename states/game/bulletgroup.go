package game

import (
	"ebijam23/resources"
	"math"
	"math/rand"
)

type GroupAngle string

const (
	Fixed  = "fixed"  // Fixed angle
	Radial = "radial" // Radial angle from spawner
	Random = "random" // Random angle
)

type BulletGroup struct {
	bullet        *Bullet    // What bullet comprises this group
	angle         GroupAngle // What angle to spawn bullets at
	spawnRate     int        // How often to spawn bullets
	lastSpawnedAt int        // How long since spawn
	bulletCount   int        // How many bullets to spawn
	loopCount     int        // How many times to loop
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
		angle := 0.0
		for i := 0; i < bg.bulletCount; i++ {
			// Set the bullet angle
			switch bg.angle {
			case Radial:
				// Spread each bullet evenly
				angle = float64(i) * 2 * math.Pi / float64(bg.bulletCount)
			case Random:
				// TODO: Random angle
				// Generate a random angle
				angle = rand.Float64() * 360
			case Fixed:
			}
			// Add the bullet to the array
			bullets[i] = BulletFromExisting(bg.bullet, angle)
		}
		// Create the action to spawn the bullets
		actions = append(actions, ActionSpawnBullets{bullets})
	}
	bg.lastSpawnedAt++
	return actions
}

func CreateBulletGroupFromDef(x, y float64, override, alias *resources.BulletGroup) *BulletGroup {
	// Create a bullet group from a bullet group definition
	// Use override values if they exist
	// TODO: maybe have default values if properties aren't present in alias or override

	angle := GroupAngle(*alias.Angle)
	spawnRate := *alias.SpawnRate
	bulletCount := *alias.BulletCount
	loopCount := *alias.LoopCount
	lastSpawnedAt := alias.LastSpawnedAt

	if override != nil {
		if override.Angle != nil {
			angle = GroupAngle(*override.Angle)
		}
		if override.SpawnRate != nil {
			spawnRate = *override.SpawnRate
		}
		if override.BulletCount != nil {
			bulletCount = *override.BulletCount
		}
		if override.LoopCount != nil {
			loopCount = *override.LoopCount
		}
		if override.LastSpawnedAt != nil {
			lastSpawnedAt = override.LastSpawnedAt
		}
	}

	// Default to spawn rate if last spawned at is nil
	spawnAt := spawnRate
	if lastSpawnedAt != nil {
		spawnAt = *lastSpawnedAt
	}
	return &BulletGroup{
		bullet:        CreateBulletFromDef(x, y, override.Bullet, alias.Bullet),
		angle:         angle,
		spawnRate:     spawnRate,
		lastSpawnedAt: spawnAt,
		bulletCount:   bulletCount,
		loopCount:     loopCount,
	}
}
