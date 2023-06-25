package resources

type ConditionType string

const (
	KilledEnemies ConditionType = "killedEnemies"
	Active                      = "active"
	Inactive                    = "inactive"
)

type ConditionDef struct {
	Type ConditionType `yaml:"type"`
	Args []string      `yaml:"args"`
}
