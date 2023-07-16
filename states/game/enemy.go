package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type EnemyState int

const (
	EnemyStateHunt EnemyState = iota
	EnemyStateWander
	EnemyStateChase
	EnemyStateFriendly
)

type Enemy struct {
	ctx               *states.Context
	id                string
	sprite            *resources.Sprite
	deadSprite        *resources.Sprite
	hitSfx            *resources.Sound
	deadSfx           *resources.Sound
	shape             RectangleShape
	target            Actor
	state             EnemyState
	alwaysShoot       bool
	wanderDir         float64
	rethinkTime       int
	health            int
	speed             int
	friendly          bool
	behavior          string
	nextPhase         string
	hasDied           bool // Set to true during Update when health < 0
	spawnOnDeath      []string
	spawner           *Spawner
	invulnerableTicks int // Ticks the enemy should be invulnerable for
	hitAccumulator    int // Hits accumulated.
	ticksUntilSfx     int // Ticks since last hit sound effect.
}

func CreateEnemy(ctx states.Context, id, enemyName string) *Enemy {
	// Get the enemy definition using enemy name and difficulty
	// If difficulty definition doesn't exist, use base definition
	enemyDef := ctx.R.GetAs("enemies", enemyName+"-"+string(ctx.Difficulty), (*resources.Enemy)(nil)).(*resources.Enemy)
	if enemyDef.Sprite == "" {
		enemyDef = ctx.R.GetAs("enemies", enemyName, (*resources.Enemy)(nil)).(*resources.Enemy)
	}

	// Get the alive and dead sprites
	aliveImageNames := ctx.R.GetNamesWithPrefix("images", enemyDef.Sprite+"-alive")
	aliveImages := make([]*ebiten.Image, 0)
	for _, s := range aliveImageNames {
		aliveImages = append(aliveImages, ctx.R.GetAs("images", s, (*ebiten.Image)(nil)).(*ebiten.Image))
	}
	aliveSprite := resources.NewAnimatedSprite(aliveImages)
	aliveSprite.Framerate = enemyDef.Framerate / 2
	aliveSprite.Loop = true

	deadImageNames := ctx.R.GetNamesWithPrefix("images", enemyDef.Sprite+"-dead")
	deadImages := make([]*ebiten.Image, 0)
	for _, s := range deadImageNames {
		deadImages = append(deadImages, ctx.R.GetAs("images", s, (*ebiten.Image)(nil)).(*ebiten.Image))
	}
	deadSprite := resources.NewAnimatedSprite(deadImages)
	deadSprite.Framerate = enemyDef.Framerate
	deadSprite.Loop = false

	// Get the hit and dead sounds
	hitSfx := ctx.R.GetAs("sounds", enemyDef.Sprite+"-hit", (*resources.Sound)(nil)).(*resources.Sound)
	deadSfx := ctx.R.GetAs("sounds", enemyDef.Sprite+"-dead", (*resources.Sound)(nil)).(*resources.Sound)

	// Create the spawner
	var spawner *Spawner
	if enemyDef.Bullets != nil {
		spawner = CreateSpawner(ctx, enemyDef.Bullets)
	}

	firstState := EnemyStateHunt
	if enemyDef.Friendly {
		//
	} else if enemyDef.Wander {
		firstState = EnemyStateWander
	}

	return &Enemy{
		id:         id,
		ctx:        &ctx,
		state:      firstState,
		sprite:     aliveSprite,
		deadSprite: deadSprite,
		hitSfx:     hitSfx,
		deadSfx:    deadSfx,
		shape: RectangleShape{
			Width:  aliveSprite.Width(),
			Height: aliveSprite.Height(),
		},
		friendly:     enemyDef.Friendly,
		health:       enemyDef.Health,
		speed:        enemyDef.Speed,
		behavior:     enemyDef.Behavior,
		alwaysShoot:  enemyDef.AlwaysShoot,
		spawner:      spawner,
		nextPhase:    enemyDef.NextPhase,
		spawnOnDeath: enemyDef.SpawnOnDeath,
	}
}

func (e *Enemy) ID() string {
	return e.id
}

func (e *Enemy) SetXY(x, y float64) {
	e.sprite.SetXY(x, y)
	e.deadSprite.SetXY(x, y)
	if e.spawner != nil {
		e.spawner.SetXY(x+e.sprite.Width()/2, y+e.sprite.Height()/2)
	}
	e.shape.X = x
	e.shape.Y = y
}

func (e *Enemy) Draw(ctx states.DrawContext) {
	if e.health <= 0 {
		e.deadSprite.Draw(ctx)
	} else if e.invulnerableTicks <= 0 || e.invulnerableTicks%4 < 2 {
		e.sprite.Draw(ctx)
		if e.spawner != nil {
			e.spawner.Draw(ctx)
		}
	}
}

func (e *Enemy) SetTarget(a Actor) {
	e.target = a
	if a == nil {
		e.state = EnemyStateWander
	} else {
		if e.friendly {
			e.state = EnemyStateFriendly
		} else {
			e.state = EnemyStateChase
		}
	}
}

func (e *Enemy) Update() (a []Action) {
	if e.ticksUntilSfx > 0 {
		e.ticksUntilSfx--
	}

	if e.invulnerableTicks > 0 {
		e.invulnerableTicks--
	}

	if e.health <= 0 {
		if !e.hasDied {
			e.hasDied = true
			e.deadSfx.Play(0.5) // TODO: use global volume setting?
			for _, spawn := range e.spawnOnDeath {
				a = append(a, ActionSpawnEnemy{
					ID:   spawn,
					Name: spawn,
					X:    e.shape.X,
					Y:    e.shape.Y,
				})
			}
			if e.nextPhase != "" {
				a = append(a, ActionSpawnEnemy{
					ID:   e.nextPhase,
					Name: e.nextPhase,
					X:    e.shape.X,
					Y:    e.shape.Y,
				})
			}
		}
		e.deadSprite.Update()
	} else {
		e.sprite.Update()
		// FIXME: Add a flag for some enemies to fire even without a target.
		if e.spawner != nil && (e.target != nil || e.alwaysShoot) {
			a = append(a, e.spawner.Update()...)
		}

		switch e.state {
		case EnemyStateHunt:
			a = append(a, ActionFindNearestActor{Actor: (*PC)(nil)})
		case EnemyStateFriendly:
			if e.target == nil {
				e.state = EnemyStateHunt
			} else {
				tx, ty, _, _ := e.target.Shape().Bounds()
				d := math.Sqrt(math.Pow(tx-e.shape.X, 2) + math.Pow(ty-e.shape.Y, 2))
				r := math.Atan2(ty-e.shape.Y, tx-e.shape.X)
				if d > 5 {
					a = append(a, ActionMove{X: e.shape.X + math.Cos(r)*float64(e.speed)*0.5, Y: e.shape.Y + math.Sin(r)*float64(e.speed)*0.25})
				}
			}
		case EnemyStateWander:
			e.rethinkTime++
			if e.rethinkTime > 0 {
				e.rethinkTime = -(30 + rng.Intn(20))
				e.wanderDir = math.Pi * 2 * rng.Float64()
				if e.target == nil {
					a = append(a, ActionFindNearestActor{Actor: (*PC)(nil)})
				}
			}
			a = append(a, ActionMove{X: e.shape.X + math.Cos(e.wanderDir)*float64(e.speed)*0.5, Y: e.shape.Y + math.Sin(e.wanderDir)*float64(e.speed)*0.25})
		case EnemyStateChase:
			if e.target == nil || e.target.Dead() {
				e.state = EnemyStateWander
			} else {
				tx, ty, _, _ := e.target.Shape().Bounds()
				r := math.Atan2(ty-e.shape.Y, tx-e.shape.X)
				a = append(a, ActionMove{X: e.shape.X + math.Cos(r)*float64(e.speed)*0.5, Y: e.shape.Y + math.Sin(r)*float64(e.speed)*0.25})
			}
		}
	}
	return a
}

func (e *Enemy) IsAlive() bool {
	return e.health > 0 && !e.hasDied
}

func (e *Enemy) Damage(amount int) bool {
	if e.invulnerableTicks > 0 {
		return false
	}

	e.health -= amount
	if e.health > 0 {
		if e.ticksUntilSfx <= 0 {
			e.ticksUntilSfx = 10
			e.hitSfx.Play(0.5) // TODO: use global volume setting?
		}
		e.hitAccumulator += amount
		if e.hitAccumulator >= e.health/4 {
			e.invulnerableTicks = 10
			e.hitAccumulator = 0
		}
	}
	return true
}

func (e *Enemy) Shape() Shape                    { return &e.shape }
func (e *Enemy) Save()                           {}
func (e *Enemy) Restore()                        {}
func (e *Enemy) Player() Player                  { return nil }
func (e *Enemy) SetPlayer(p Player)              {}
func (e *Enemy) SetImpulses(impulses ImpulseSet) {}
func (e *Enemy) Bounds() (x, y, w, h float64)    { return 0, 0, 0, 0 }
func (e *Enemy) SetSize(r float64)               {}
func (e *Enemy) Dead() bool                      { return false }
func (e *Enemy) Destroyed() bool                 { return false }
