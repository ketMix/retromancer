package game

import (
	"ebijam23/resources"
	"ebijam23/states"

	"github.com/hajimehoshi/ebiten/v2"
)

type Enemy struct {
	sprite     *resources.Sprite
	deadSprite *resources.Sprite
	shape      RectangleShape
	phases     []*resources.Enemy
	health     int
	speed      int
	behavior   string
	spawner    Spawner
}

func CreateEnemy(ctx states.Context, enemyDef resources.Enemy) *Enemy {
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

	return &Enemy{
		sprite:     aliveSprite,
		deadSprite: deadSprite,
		shape: RectangleShape{
			Width:  aliveSprite.Width(),
			Height: aliveSprite.Height(),
		},
		phases:   enemyDef.Phases,
		health:   enemyDef.Health,
		speed:    enemyDef.Speed,
		behavior: enemyDef.Behavior,
		spawner:  *CreateSpawner(ctx, enemyDef.Bullets),
	}
}

func (e *Enemy) SetXY(x, y float64) {
	e.sprite.SetXY(x, y)
	e.deadSprite.SetXY(x, y)
	e.shape.X = x
	e.shape.Y = y
}

func (e *Enemy) Draw(ctx states.DrawContext) {
	if e.health <= 0 {
		e.deadSprite.Draw(ctx)
	} else {
		e.sprite.Draw(ctx)
		e.spawner.Draw(ctx)
	}
}

func (e *Enemy) Update() []Action {
	return nil
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
