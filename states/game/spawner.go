package game

import (
	"ebijam23/resources"

	"github.com/hajimehoshi/ebiten/v2"
)

// This can probably be attached to an actor instead being its own actor
type Spawner struct {
	shape        CircleShape
	bulletGroups []*BulletGroup
}

func CreateSpawner(x, y float64, bulletGroups []*BulletGroup) *Spawner {
	return &Spawner{
		shape:        CircleShape{X: x, Y: y, Radius: 0},
		bulletGroups: bulletGroups,
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

func (s *Spawner) Shape() Shape                          { return &s.shape }
func (s *Spawner) Save()                                 {}
func (s *Spawner) Restore()                              {}
func (s *Spawner) Player() Player                        { return nil }
func (s *Spawner) SetPlayer(p Player)                    {}
func (s *Spawner) SetImpulses(impulses ImpulseSet)       {}
func (s *Spawner) Draw(screen *ebiten.Image)             {}
func (s *Spawner) Bounds() (x, y, w, h float64)          { return 0, 0, 0, 0 }
func (s *Spawner) SetXY(x, y float64)                    {}
func (s *Spawner) SetSize(r float64)                     {}
func (s *Spawner) Dead() bool                            { return false }
func (s *Spawner) Reverse()                              {}
func (s *Spawner) Conditions() []*resources.ConditionDef { return nil }
func (s *Spawner) Active() bool                          { return false }
func (s *Spawner) SetActive(bool)                        {}
func (s *Spawner) ID() string                            { return "" }
