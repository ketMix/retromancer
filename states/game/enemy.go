package game

import (
	"ebijam23/resources"
	"ebijam23/states"

	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	id         string
	sprite     *resources.Sprite
	deadSprite *resources.Sprite
	hitSfx     *resources.Sound
	deadSfx    *resources.Sound
	shape      RectangleShape
	phases     []*resources.Enemy
	health     int
	speed      int
	behavior   string
	spawner    *Spawner
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

	deadImageNames := ctx.Manager.GetNamesWithPrefix("images", enemyDef.Sprite+"-dead")
	deadImages := make([]*ebiten.Image, 0)
	for _, s := range deadImageNames {
		deadImages = append(deadImages, ctx.Manager.GetAs("images", s, (*ebiten.Image)(nil)).(*ebiten.Image))
	}
	deadSprite := resources.NewAnimatedSprite(deadImages)

	// Get the hit and dead sounds
	hitSfx := ctx.Manager.GetAs("sounds", enemyDef.Sprite+"-hit", (*resources.Sound)(nil)).(*resources.Sound)
	deadSfx := ctx.Manager.GetAs("sounds", enemyDef.Sprite+"-dead", (*resources.Sound)(nil)).(*resources.Sound)

	// Create the spawner
	var spawner *Spawner
	if enemyDef.Bullets != nil {
		spawner = CreateSpawner(ctx, enemyDef.Bullets)
	}
	return &Enemy{
		sprite:     aliveSprite,
		deadSprite: deadSprite,
		hitSfx:     hitSfx,
		deadSfx:    deadSfx,
		shape: RectangleShape{
			Width:  aliveSprite.Width(),
			Height: aliveSprite.Height(),
		},
		phases:   enemyDef.Phases,
		health:   enemyDef.Health,
		speed:    enemyDef.Speed,
		behavior: enemyDef.Behavior,
		spawner:  spawner,
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
	} else {
		e.sprite.Draw(ctx)
		if e.spawner != nil {
			e.spawner.Draw(ctx)
		}
	}
}

func (e *Enemy) Update() (a []Action) {
	if e.health <= 0 {
		e.deadSprite.Update()
	} else {
		e.sprite.Update()
		if e.spawner != nil {
			a = append(a, e.spawner.Update()...)
		}
	}
	return a
}

func (e *Enemy) IsAlive() bool {
	return e.health > 0
}

func (e *Enemy) Damage(amount int) {
	e.health -= amount
	if e.health <= 0 {
		e.deadSfx.Play(0.5) // TODO: use global volume setting?
	} else {
		e.hitSfx.Play(0.5) // TODO: use global volume setting?
	}
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
