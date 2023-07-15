package game

import (
	"ebijam23/resources"
	"ebijam23/states"
)

// This can probably be attached to an actor instead being its own actor
type Spawner struct {
	shape        CircleShape
	bulletGroups []*BulletGroup
}

func CreateSpawner(ctx states.Context, bulletGroupDefs []*resources.BulletGroup) *Spawner {
	bulletGroups := make([]*BulletGroup, 0)

	// If we have bullet groups defined for the spawner, create them.
	if len(bulletGroupDefs) > 0 {
		for _, bg := range bulletGroupDefs {
			bulletAlias := (*resources.BulletGroup)(nil)
			if bg.Alias != nil {
				bulletAlias = ctx.R.GetAs("bullets", *bg.Alias, (*resources.BulletGroup)(nil)).(*resources.BulletGroup)
			}
			bulletGroups = append(bulletGroups, CreateBulletGroupFromDef(bg, bulletAlias))
		}
	}
	return &Spawner{
		shape:        CircleShape{Radius: 0},
		bulletGroups: bulletGroups,
	}
}

func (s *Spawner) SetXY(x, y float64) {
	s.shape.X = x
	s.shape.Y = y
	for _, bg := range s.bulletGroups {
		bg.SetXY(x, y)
	}
}

func (s *Spawner) Update() (actions []Action) {
	// Update the bullet groups
	for _, bg := range s.bulletGroups {
		// Add the actions from the bullet group to the list of actions
		actions = append(actions, bg.Update()...)
	}
	return actions
}

func (s *Spawner) Destroyed() bool                 { return false }
func (s *Spawner) Shape() Shape                    { return &s.shape }
func (s *Spawner) Save()                           {}
func (s *Spawner) Restore()                        {}
func (s *Spawner) Player() Player                  { return nil }
func (s *Spawner) SetPlayer(p Player)              {}
func (s *Spawner) SetImpulses(impulses ImpulseSet) {}
func (s *Spawner) Draw(states.DrawContext)         {}
func (s *Spawner) Bounds() (x, y, w, h float64)    { return 0, 0, 0, 0 }
func (s *Spawner) SetSize(r float64)               {}
func (s *Spawner) Dead() bool                      { return false }
