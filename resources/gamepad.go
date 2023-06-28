package resources

import (
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

type GamepadMap map[int]int

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
)

var GamepadMaps = map[string]GamepadMap{
	"standard": {
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
		ButtonX:            int(ebiten.StandardGamepadButtonRightLeft),
	},
	"Xbox360": {
		AxisLeftX:  0,
		AxisLeftY:  1,
		AxisRightX: 3,
		AxisRightY: 4,
		//
		ButtonBumperLeft:  102, // 100 means it's technically an axis, but with -1=0 and 1=1
		ButtonBumperRight: 105,
		//
		ButtonTriggerLeft:  4,
		ButtonTriggerRight: 5,
		ButtonBack:         6,
		ButtonStart:        7,
		ButtonA:            1,
		ButtonB:            0,
		ButtonX:            2,
	},
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
	name := ebiten.GamepadName(ebiten.GamepadID(id))
	if strings.Contains(name, "Xbox") || strings.Contains(name, "X-Box") || strings.Contains(name, "360") {
		return "Xbox360"
	}
	// TODO: Add PS5 or something.
	return ""
}

func GetAxis(gamemap string, id int, which int) float64 {
	a, ok := GamepadMaps[gamemap][which]
	if !ok {
		return 0
	}
	if gamemap == "standard" {
		return ebiten.StandardGamepadAxisValue(ebiten.GamepadID(id), ebiten.StandardGamepadAxis(a))
	}
	return ebiten.GamepadAxisValue(ebiten.GamepadID(id), a)
}

func GetButton(gamemap string, id int, which int) bool {
	b, ok := GamepadMaps[gamemap][which]
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
