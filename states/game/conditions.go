package game

import (
	"ebijam23/resources"
)

// Check all provided conditions are true
// If no conditions, return false
func CheckConditions(conditions []*resources.ConditionDef, interactives []*Interactive, enemies []*Enemy) bool {
	if len(conditions) == 0 {
		return false
	}

	for _, condition := range conditions {
		args := condition.Args
		switch condition.Type {
		case resources.Active:
			return CheckActiveCondition(&args, interactives)
		case resources.KilledEnemies:
			return CheckKilledEnemiesCondition(&args, enemies)
		}
	}
	return true
}

// Check that all interactives in args are active
func CheckActiveCondition(check *[]string, interactives []*Interactive) bool {
	if len(*check) == 0 {
		return true
	}
	// Check all in args are active
	for _, a := range interactives {
		for _, arg := range *check {
			// Find actor by id within actor array
			if a.ID() == arg && !a.Active() {
				return false
			}
		}
	}
	return true
}

// Check that all enemies in args are dead
// If no args, check all enemies are dead
func CheckKilledEnemiesCondition(check *[]string, enemies []*Enemy) bool {
	checkAll := len(*check) == 0
	// Check all in args are active
	for _, e := range enemies {
		if checkAll && e.IsAlive() {
			return false
		} else {
			for i, arg := range *check {
				// Find actor by id within actor array
				if e.ID() == arg {
					if e.IsAlive() {
						return false
					} else {
						// remove it from the list safely
						if len(*check) == 1 {
							return true
						}
						*check = append((*check)[:i], (*check)[i+1:]...)
					}
				}
			}
		}
	}
	if !checkAll {
		return !(len(*check) > 0)
	}
	return true
}
