package resources

type ConditionType string

const (
	Cleared  ConditionType = "cleared"
	Active                 = "active"
	Inactive               = "inactive"
)

type ConditionDef struct {
	Type string `yaml:"type"`
	Args []int  `yaml:"args"`
}
