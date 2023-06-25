package game

import (
	"ebijam23/resources"
	"ebijam23/states"

	"github.com/hajimehoshi/ebiten/v2"
)

// Default cooldown for interactive activation + degradation
const ACTIVATE_COOLDOWN = 5
const DEGRADE_COOLDOWN = ACTIVATE_COOLDOWN * 10

type Interactive struct {
	id                 string         // ID of the interactive object, used to identify it in the game for condition triggering
	maxHp              int            // Hit points of the interactive object, if it reaches 0, it is activated (only relevant for shooting)
	hp                 int            // Current hit points of the interactive object, also determines if it is shootable
	text               string         // Text to display when the interactive is active and touched
	linkedInteractives []*Interactive // Holds a list of interactives that have their activation state linked to this one
	active             bool
	activeSprite       *resources.Sprite
	inactiveSprite     *resources.Sprite
	conditions         []*resources.ConditionDef
	shape              RectangleShape
	nextMap            *string
	reversable         bool
	touchable          bool // Whether or not it can be reversed by touching
	degrade            bool // Whether or not the activation can be degraded, should this always be true?
	activationIdx      int  // Holds the degree of activation
	activateCooldown   int  // Holds the cooldown for activation, can only decrement activation when this is 0
}

func CreateInteractive(ctx states.Context, actorDef resources.ActorSpawn) *Interactive {
	// Set up sprites
	spritePrefix := actorDef.Sprite

	// Create the active sprite
	activeImageNames := ctx.Manager.GetNamesWithPrefix("images", spritePrefix+"-active")
	activeImages := make([]*ebiten.Image, 0)
	for _, s := range activeImageNames {
		activeImages = append(activeImages, ctx.Manager.GetAs("images", s, (*ebiten.Image)(nil)).(*ebiten.Image))
	}
	activeSprite := resources.NewAnimatedSprite(activeImages)
	if activeSprite != nil {
		activeSprite.Framerate = 5
		activeSprite.Loop = true
	}

	// Create the inactive sprite
	inactiveImageNames := ctx.Manager.GetNamesWithPrefix("images", spritePrefix+"-inactive")
	inactiveImages := make([]*ebiten.Image, 0)
	for _, s := range inactiveImageNames {
		inactiveImages = append(inactiveImages, ctx.Manager.GetAs("images", s, (*ebiten.Image)(nil)).(*ebiten.Image))
	}
	inactiveSprite := resources.NewAnimatedSprite(inactiveImages)

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

	interactive := &Interactive{
		id:             actorDef.ID,
		activeSprite:   activeSprite,
		inactiveSprite: inactiveSprite,
		shape: RectangleShape{
			Width:  width,
			Height: height,
		},
		activationIdx: activationIdx,
	}

	// If there are interactive definitions, set them
	i := actorDef.Interactive
	if i != nil {
		interactive.active = i.Active
		interactive.degrade = i.Degrade
		interactive.conditions = i.Conditions
		interactive.nextMap = i.Map
		interactive.maxHp = i.Health
		interactive.hp = i.Health
		interactive.reversable = i.Reversable
		interactive.touchable = i.Touchable
		interactive.text = i.Text
	}
	return interactive
}

func (i *Interactive) SetXY(x, y float64) {
	i.shape.X = x
	i.shape.Y = y
	if i.activeSprite != nil {
		i.activeSprite.SetXY(x, y)
	}
	if i.inactiveSprite != nil {
		i.inactiveSprite.SetXY(x, y)
	}
}

func (i *Interactive) Bounds() (x, y, w, h float64) {
	return i.shape.X, i.shape.Y, i.shape.Width, i.shape.Height
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
		// If the interactive is shootable, color it a bit red depending on the hp
		if i.hp > 0 {
			colorScale := float32(i.hp) / float32(i.maxHp)
			opts := &ebiten.DrawImageOptions{}
			opts.ColorScale.Scale(1, colorScale, colorScale, 1)
			i.inactiveSprite.DrawWithOptions(ctx, opts)
		} else {
			i.inactiveSprite.Draw(ctx)
		}
	}
}

func (i *Interactive) Reverseable() bool {
	if i.active {
		return false
	}

	return i.reversable
}

func (i *Interactive) Shootable() bool {
	if i.active {
		return false
	}
	return i.hp > 0
}

func (i *Interactive) Touchable() bool {
	if !i.active {
		return false
	}
	return i.touchable
}

func (i *Interactive) Hit() {
	if i.hp == 0 || i.active {
		return
	}
	i.hp--
	if i.hp <= 0 {
		i.IncreaseActivation(nil)
		i.hp = i.maxHp
	}
}

// Reverse the interactive object
func (i *Interactive) Reverse() {
	if !i.reversable && !i.touchable && i.linkedInteractives == nil {
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
	if i.activationIdx < 0 {
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
func (i *Interactive) SetSize(r float64)               {}
func (i *Interactive) Dead() bool                      { return false }
func (i *Interactive) Destroyed() bool                 { return false }
