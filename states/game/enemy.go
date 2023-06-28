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
	behavior          string
	nextPhase         string
	spawner           *Spawner
	invulnerableTicks int // Ticks the enemy should be invulnerable for
	hitAccumulator    int // Hits accumulated.
}

func CreateEnemy(ctx states.Context, id, enemyName string) *Enemy {
	// Get the enemy definition using enemy name
	enemyDef := ctx.Manager.GetAs("enemies", enemyName, (*resources.Enemy)(nil)).(*resources.Enemy)

	// Get the alive and dead sprites
	aliveImageNames := ctx.Manager.GetNamesWithPrefix("images", enemyDef.Sprite+"-alive")
	aliveImages := make([]*ebiten.Image, 0)
	for _, s := range aliveImageNames {
		aliveImages = append(aliveImages, ctx.Manager.GetAs("images", s, (*ebiten.Image)(nil)).(*ebiten.Image))
	}
	aliveSprite := resources.NewAnimatedSprite(aliveImages)
	aliveSprite.Framerate = enemyDef.Framerate / 2
	aliveSprite.Loop = true

	deadImageNames := ctx.Manager.GetNamesWithPrefix("images", enemyDef.Sprite+"-dead")
	deadImages := make([]*ebiten.Image, 0)
	for _, s := range deadImageNames {
		deadImages = append(deadImages, ctx.Manager.GetAs("images", s, (*ebiten.Image)(nil)).(*ebiten.Image))
	}
	deadSprite := resources.NewAnimatedSprite(deadImages)
	deadSprite.Framerate = enemyDef.Framerate
	deadSprite.Loop = false

	// Get the hit and dead sounds
	hitSfx := ctx.Manager.GetAs("sounds", enemyDef.Sprite+"-hit", (*resources.Sound)(nil)).(*resources.Sound)
	deadSfx := ctx.Manager.GetAs("sounds", enemyDef.Sprite+"-dead", (*resources.Sound)(nil)).(*resources.Sound)

	// Create the spawner
	var spawner *Spawner
	if enemyDef.Bullets != nil {
		spawner = CreateSpawner(ctx, enemyDef.Bullets)
	}

	firstState := EnemyStateHunt
	if enemyDef.Wander {
		firstState = EnemyStateWander
	}

	return &Enemy{
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
		health:      enemyDef.Health,
		speed:       enemyDef.Speed,
		behavior:    enemyDef.Behavior,
		alwaysShoot: enemyDef.AlwaysShoot,
		spawner:     spawner,
		nextPhase:   enemyDef.NextPhase,
	}
}

func (e *Enemy) ID() string {
	return e.id
}

func (e *Enemy) SetXY(x, y float64) {
	e.sprite.SetXY(x, y)
	e.deadSprite.SetXY(x, y)
	e.spawner.SetXY(x+e.sprite.Width()/2, y+e.sprite.Height()/2)
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
		e.state = EnemyStateChase
	}
}

func (e *Enemy) Update() (a []Action) {
	if e.invulnerableTicks > 0 {
		e.invulnerableTicks--
	}

	if e.health <= 0 {
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
	return e.health > 0
}

func (e *Enemy) Damage(amount int) bool {
	if e.invulnerableTicks > 0 {
		return false
	}

	e.health -= amount
	if e.health <= 0 {
		e.deadSfx.Play(0.5) // TODO: use global volume setting?
		if e.nextPhase != "" {
			nextPhase := CreateEnemy(*e.ctx, e.id, e.nextPhase)
			e.health = nextPhase.health
			e.speed = nextPhase.speed
			e.sprite = nextPhase.sprite
			e.deadSprite = nextPhase.deadSprite
			e.spawner = nextPhase.spawner
			e.nextPhase = nextPhase.nextPhase
		}
	} else {
		e.hitSfx.Play(0.5) // TODO: use global volume setting?
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
