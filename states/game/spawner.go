package game

import (
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

func (s *Spawner) Shape() Shape                    { return &s.shape }
func (s *Spawner) Player() Player                  { return nil }
func (s *Spawner) SetPlayer(p Player)              {}
func (s *Spawner) SetImpulses(impulses ImpulseSet) {}
func (s *Spawner) Draw(screen *ebiten.Image)       {}
func (s *Spawner) Bounds() (x, y, w, h float64)    { return 0, 0, 0, 0 }
func (s *Spawner) SetXY(x, y float64)              {}
func (s *Spawner) SetSize(r float64)               {}
