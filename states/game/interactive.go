package game

import (
	"ebijam23/resources"

	"github.com/hajimehoshi/ebiten/v2"
)

type Interactive struct {
	ID               string // ID of the interactive object, used to identify it in the game for condition triggering
	Active           bool
	ActiveSprite     *resources.Sprite
	InactiveSprite   *resources.Sprite
	shape            RectangleShape
	activationIdx    int // Holds the degree of activation
	activateCooldown int // Holds the cooldown for activation, can only decrement activation when this is 0
}

func CreateInteractive(x, y float64, id string, active bool, activeImages, inactiveImages []*ebiten.Image) *Interactive {
	// Set activation index to fully inactive
	activationIdx := len(inactiveImages) - 1
	if activationIdx < 0 {
		activationIdx = 0
	}

	// Create the active sprite
	activeSprite := resources.NewAnimatedSprite(activeImages)
	activeSprite.X = x
	activeSprite.Y = y
	activeSprite.Framerate = 5
	activeSprite.Loop = true

	// Create the inactive sprite
	inactiveSprite := resources.NewAnimatedSprite(inactiveImages)
	inactiveSprite.X = x
	inactiveSprite.Y = y
	inactiveSprite.Framerate = 0
	inactiveSprite.SetFrame(activationIdx)

	return &Interactive{
		ID:             id,
		Active:         active,
		ActiveSprite:   activeSprite,
		InactiveSprite: inactiveSprite,
		shape: RectangleShape{
			X:      x,
			Y:      y,
			Width:  activeSprite.Width(),
			Height: activeSprite.Height(),
		},
		activationIdx: activationIdx,
	}
}

func (i *Interactive) Update() []Action {
	if !i.Active {
		if i.activateCooldown > 0 {
			i.activateCooldown--
		}
		if i.InactiveSprite.Frame() != i.activationIdx {
			i.InactiveSprite.SetFrame(i.activationIdx)
		}
		i.InactiveSprite.Update()
	} else {
		i.ActiveSprite.Update()
	}
	return nil
}

func (i *Interactive) Draw(screen *ebiten.Image) {
	if i.Active {
		i.ActiveSprite.Draw(screen)
	} else {
		i.InactiveSprite.Draw(screen)
	}
}

// Reverse the interactive object
func (i *Interactive) Reverse() {
	// If already active or on cooldown, do nothing
	if i.Active || i.activateCooldown > 0 {
		return
	}
	// Set the cooldown, decrement the activation index
	i.activateCooldown = 30
	i.activationIdx--

	// If the activation index is now negative, set it to 0 and activate the object
	if i.activationIdx < 0 {
		i.activationIdx = 0
		i.Active = true
	}
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
