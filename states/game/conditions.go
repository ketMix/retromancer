package game

func CheckActiveCondition(actorsToCheck []string, mapActors []*Interactive) bool {
	checkNum := len(actorsToCheck)
	checkedNum := 0
	// Check all actor ids in args are active
	for _, a := range mapActors {
		for _, arg := range actorsToCheck {
			// Find actor by id within actor array
			if a.ID() == arg && a.Active() {
				checkedNum++
			}
		}
	}
	if checkNum == checkedNum {
		return true
	}
	return false
}
