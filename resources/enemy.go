package resources

type Enemy struct {
	Sprite      string         `yaml:"sprite"`
	Framerate   int            `yaml:"framerate"`
	Health      int            `yaml:"health"`
	Speed       int            `yaml:"speed"`
	Behavior    string         `yaml:"behavior"`
	Wander      bool           `yaml:"wander"`
	AlwaysShoot bool           `yaml:"alwaysShoot"`
	Bullets     []*BulletGroup `yaml:"bullets"`
	NextPhase   string         `yaml:"nextPhase"`
}

func (e *Enemy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type enemy Enemy
	if err := unmarshal((*enemy)(e)); err != nil {
		return err
	}

	return nil
}
