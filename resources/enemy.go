package resources

type Enemy struct {
	Sprite   string         `yaml:"sprite"`
	Health   int            `yaml:"health"`
	Speed    int            `yaml:"speed"`
	Behavior string         `yaml:"behavior"`
	Bullets  []*BulletGroup `yaml:"bullets"`
	Phases   []*Enemy       `yaml:"phases"`
}

func (e *Enemy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type enemy Enemy
	if err := unmarshal((*enemy)(e)); err != nil {
		return err
	}

	return nil
}
