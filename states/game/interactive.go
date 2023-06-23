package game

import (
	"ebijam23/resources"
	"ebijam23/states"
)

// Default cooldown for interactive activation + degradation
const ACTIVATE_COOLDOWN = 5
const DEGRADE_COOLDOWN = ACTIVATE_COOLDOWN * 10

type Interactive struct {
	id                 string         // ID of the interactive object, used to identify it in the game for condition triggering
	linkedInteractives []*Interactive // Holds a list of interactives that have their activation state linked to this one
	active             bool
	activeSprite       *resources.Sprite
	inactiveSprite     *resources.Sprite
	conditions         []*resources.ConditionDef
	shape              RectangleShape
	reversable         bool
	degrade            bool // Whether or not the activation can be degraded, should this always be true?
	activationIdx      int  // Holds the degree of activation
	activateCooldown   int  // Holds the cooldown for activation, can only decrement activation when this is 0
}

func CreateInteractive(x, y float64, id string, active, reversable bool, conditions []*resources.ConditionDef, activeSprite, inactiveSprite *resources.Sprite) *Interactive {
	// Set activation index to fully inactive
	activationIdx := len(inactiveSprite.Images()) - 1
	if activationIdx < 0 {
		activationIdx = 0
	}
	inactiveSprite.SetFrame(activationIdx)

	width := 0.0
	height := 0.0
	if activeSprite != nil {
		width = activeSprite.Width()
		height = activeSprite.Height()
	} else {
		width = inactiveSprite.Width()
		height = inactiveSprite.Height()
	}
	return &Interactive{
		id:             id,
		active:         active,
		activeSprite:   activeSprite,
		inactiveSprite: inactiveSprite,
		conditions:     conditions,
		shape: RectangleShape{
			X:      x,
			Y:      y,
			Width:  width,
			Height: height,
		},
		reversable:    reversable,
		activationIdx: activationIdx,
	}
}

func (i *Interactive) Update() []Action {
	if !i.active {
		i.activateCooldown--

		// If we haven't reversed, degrade the activation
		if i.activateCooldown <= -(DEGRADE_COOLDOWN) {
			if i.degrade {
				i.DecreaseActivation(nil)
			}
			i.activateCooldown = 0
		}
		i.inactiveSprite.SetFrame(i.activationIdx)
		i.inactiveSprite.Update()
	} else {
		i.activeSprite.Update()
	}
	return nil
}

// If the interactive object is active, draw the active sprite
// otherwise draw the inactive sprite
func (i *Interactive) Draw(ctx states.DrawContext) {
	if i.active {
		i.activeSprite.Draw(ctx)
	} else {
		i.inactiveSprite.Draw(ctx)
	}
}

func (i *Interactive) Reverseable() bool {
	if !i.reversable {
		return false
	}
	if i.active {
		return false
	}

	return true
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
	// Set the cooldown, increase the activation
	i.activateCooldown = ACTIVATE_COOLDOWN
	i.IncreaseActivation(nil)
}

func (i *Interactive) Conditions() []*resources.ConditionDef {
	return i.conditions
}

func (i *Interactive) Active() bool {
	return i.active
}

// Takes parent ID to prevent infinite loops
func (i *Interactive) IncreaseActivation(parentIds []string) {
	for _, parentId := range parentIds {
		if parentId == i.id {
			return
		}
	}
	i.activationIdx--

	// If the activation index is now negative, and we have active sprite defined
	// set it to 0 and activate the object
	// otherwise reset the idx
	if i.activationIdx <= 0 {
		if i.activeSprite != nil {
			i.activationIdx = 0
			i.active = true
		} else {
			i.activationIdx = len(i.inactiveSprite.Images()) - 1
		}
	}

	// Should we set linked to active if we're active?
	for _, linked := range i.linkedInteractives {
		linked.IncreaseActivation(append(parentIds, i.id))
		linked.activateCooldown = ACTIVATE_COOLDOWN
	}
}

// Takes parent ID to prevent infinite loops
func (i *Interactive) DecreaseActivation(parentIds []string) {
	for _, parentId := range parentIds {
		if parentId == i.id {
			return
		}
	}
	i.activationIdx++
	if i.activationIdx >= len(i.inactiveSprite.Images()) {
		i.activationIdx = len(i.inactiveSprite.Images()) - 1
	}
	for _, linked := range i.linkedInteractives {
		linked.DecreaseActivation(append(parentIds, i.id))
	}
}

func (i *Interactive) ID() string {
	return i.id
}

func (i *Interactive) AddLinkedInteractive(interactives *Interactive) {
	i.linkedInteractives = append(i.linkedInteractives, interactives)
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
func (i *Interactive) Destroyed() bool                 { return false }
