package game

import (
	"ebijam23/resources"

	"github.com/hajimehoshi/ebiten/v2"
)

type Interactive struct {
	id               string // ID of the interactive object, used to identify it in the game for condition triggering
	active           bool
	activeSprite     *resources.Sprite
	inactiveSprite   *resources.Sprite
	conditions       []*resources.ConditionDef
	shape            RectangleShape
	reversable       bool
	activationIdx    int // Holds the degree of activation
	activateCooldown int // Holds the cooldown for activation, can only decrement activation when this is 0
}

func CreateInteractive(x, y float64, id string, active, reversable bool, conditions []*resources.ConditionDef, activeSprite, inactiveSprite *resources.Sprite) *Interactive {
	// Set activation index to fully inactive
	activationIdx := len(inactiveSprite.Images()) - 1
	if activationIdx < 0 {
		activationIdx = 0
	}
	inactiveSprite.SetFrame(activationIdx)

	return &Interactive{
		id:             id,
		active:         active,
		activeSprite:   activeSprite,
		inactiveSprite: inactiveSprite,
		conditions:     conditions,
		shape: RectangleShape{
			X:      x,
			Y:      y,
			Width:  activeSprite.Width(),
			Height: activeSprite.Height(),
		},
		reversable:    reversable,
		activationIdx: activationIdx,
	}
}

func (i *Interactive) Update() []Action {
	if !i.active {
		if i.activateCooldown > 0 {
			i.activateCooldown--
		}
		if i.inactiveSprite.Frame() != i.activationIdx {
			i.inactiveSprite.SetFrame(i.activationIdx)
		}
		i.inactiveSprite.Update()
	} else {
		i.activeSprite.Update()
	}
	return nil
}

func (i *Interactive) Draw(screen *ebiten.Image) {
	if i.active {
		i.activeSprite.Draw(screen)
	} else {
		i.inactiveSprite.Draw(screen)
	}
}

// Reverse the interactive object
func (i *Interactive) Reverse() {
	if !i.reversable {
		return
	}
	// If already active or on cooldown, do nothing
	if i.active || i.activateCooldown > 0 {
		return
	}
	// Set the cooldown, decrement the activation index
	i.activateCooldown = 30
	i.activationIdx--

	// If the activation index is now negative, set it to 0 and activate the object
	if i.activationIdx < 0 {
		i.activationIdx = 0
		i.active = true
	}
}

func (i *Interactive) Conditions() []*resources.ConditionDef {
	return i.conditions
}

func (i *Interactive) Active() bool {
	return i.active
}

func (i *Interactive) SetActive(active bool) {
	i.active = active
}

func (i *Interactive) ID() string {
	return i.id
}
func (i *Interactive) Shape() Shape                    { return &i.shape }
func (i *Interactive) Save()                           {}
func (i *Interactive) Restore()                        {}
func (i *Interactive) Player() Player                  { return nil }
func (i *Interactive) SetPlayer(p Player)              {}
func (i *Interactive) SetImpulses(impulses ImpulseSet) {}
func (i *Interactive) Bounds() (x, y, w, h float64)    { return 0, 0, 0, 0 }
func (i *Interactive) SetXY(x, y float64)              {}
func (i *Interactive) SetSize(r float64)               {}
func (i *Interactive) Dead() bool                      { return false }
