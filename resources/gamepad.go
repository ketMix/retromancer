package resources

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type GamepadMap map[int]int

type GamepadDefinition struct {
	Name  string     `yaml:"Name"`
	Match []string   `yaml:"Match"`
	Map   GamepadMap `yaml:"Map"`
}

func (g *GamepadMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var m map[string]int
	if err := unmarshal(&m); err != nil {
		return err
	}

	*g = make(GamepadMap)
	for k, v := range m {
		id, ok := idToInt[strings.TrimSpace(k)]
		if ok {
			(*g)[id] = v
		}
	}
	return nil
}

const (
	AxisLeftX = iota
	AxisLeftY
	AxisRightX
	AxisRightY
	//
	ButtonBumperLeft
	ButtonBumperRight
	//
	ButtonTriggerLeft
	ButtonTriggerRight
	ButtonBack
	ButtonStart
	ButtonA
	ButtonB
	ButtonX
	ButtonY
)

var idToInt = map[string]int{
	"AxisLeftX":          AxisLeftX,
	"AxisLeftY":          AxisLeftY,
	"AxisRightX":         AxisRightX,
	"AxisRightY":         AxisRightY,
	"ButtonBumperLeft":   ButtonBumperLeft,
	"ButtonBumperRight":  ButtonBumperRight,
	"ButtonTriggerLeft":  ButtonTriggerLeft,
	"ButtonTriggerRight": ButtonTriggerRight,
	"ButtonBack":         ButtonBack,
	"ButtonStart":        ButtonStart,
	"ButtonA":            ButtonA,
	"ButtonB":            ButtonB,
	"ButtonX":            ButtonX,
	"ButtonY":            ButtonY,
}

var GamepadDefinitions = []GamepadDefinition{
	{
		Name:  "standard",
		Match: []string{},
		Map: GamepadMap{
			AxisLeftX:          int(ebiten.StandardGamepadAxisLeftStickHorizontal),
			AxisLeftY:          int(ebiten.StandardGamepadAxisLeftStickVertical),
			AxisRightX:         int(ebiten.StandardGamepadAxisRightStickHorizontal),
			AxisRightY:         int(ebiten.StandardGamepadAxisRightStickVertical),
			ButtonBumperLeft:   int(ebiten.StandardGamepadButtonFrontBottomLeft),
			ButtonBumperRight:  int(ebiten.StandardGamepadButtonFrontBottomRight),
			ButtonTriggerLeft:  int(ebiten.StandardGamepadButtonFrontTopLeft),
			ButtonTriggerRight: int(ebiten.StandardGamepadButtonFrontTopRight),
			ButtonBack:         int(ebiten.StandardGamepadButtonCenterLeft),
			ButtonStart:        int(ebiten.StandardGamepadButtonCenterRight),
			ButtonA:            int(ebiten.StandardGamepadButtonRightBottom),
			ButtonB:            int(ebiten.StandardGamepadButtonRightRight),
			ButtonX:            int(ebiten.StandardGamepadButtonRightTop),
			ButtonY:            int(ebiten.StandardGamepadButtonRightLeft),
		},
	},
}

func AddGamepadDefinition(g GamepadDefinition) {
	GamepadDefinitions = append(GamepadDefinitions, g)
}

func GetFunctionalGamepads() (c []int) {
	controllers := ebiten.AppendGamepadIDs(nil)
	// Filter controllers that are not standard layout.
	for i := 0; i < len(controllers); {
		bestMap := GetBestGamemap(int(controllers[i]))
		if bestMap == "" {
			controllers = append(controllers[:i], controllers[i+1:]...)
		} else {
			c = append(c, int(controllers[i]))
			i++
		}
	}
	return c
}

func GetBestGamemap(id int) string {
	if ebiten.IsStandardGamepadLayoutAvailable(ebiten.GamepadID(id)) {
		return "standard"
	}
	name := strings.ToLower(ebiten.GamepadName(ebiten.GamepadID(id)))
	for _, g := range GamepadDefinitions {
		for _, m := range g.Match {
			if strings.Contains(name, strings.ToLower(m)) {
				return g.Name
			}
		}
	}
	return ""
}

func GetGamemap(name string) GamepadMap {
	for _, g := range GamepadDefinitions {
		if g.Name == name {
			return g.Map
		}
	}
	return nil
}

func GetAxis(gamemap string, id int, which int) float64 {
	a, ok := GetGamemap(gamemap)[which]
	if !ok {
		return 0
	}
	deadzone := 0.1
	value := ebiten.GamepadAxisValue(ebiten.GamepadID(id), a)
	if gamemap == "standard" {
		value = ebiten.StandardGamepadAxisValue(ebiten.GamepadID(id), ebiten.StandardGamepadAxis(a))
	}
	if value < deadzone && value > -deadzone {
		return 0
	}
	return value
}

func GetButton(gamemap string, id int, which int) bool {
	b, ok := GetGamemap(gamemap)[which]
	if !ok {
		return false
	}
	if gamemap == "standard" {
		return ebiten.IsStandardGamepadButtonPressed(ebiten.GamepadID(id), ebiten.StandardGamepadButton(b))
	}

	if b > 100 {
		a := ebiten.GamepadAxisValue(ebiten.GamepadID(id), b-100)
		if a < -0.5 {
			return false
		} else if a > 0.5 {
			return true
		}
	}
	return ebiten.IsGamepadButtonPressed(ebiten.GamepadID(id), ebiten.GamepadButton(b))
}
