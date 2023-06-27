package resources

// Omits empty to allow for overriding from enemy definition
type Bullet struct {
	BulletType      *string  `yaml:"bulletType,omitempty"`
	Color           *[]int   `yaml:"color,omitempty"`
	Radius          *int     `yaml:"radius,omitempty"`
	Speed           *float64 `yaml:"speed,omitempty"`
	Acceleration    *float64 `yaml:"acceleration,omitempty"`
	AccelAccel      *float64 `yaml:"accelAccel,omitempty"`
	MaxSpeed        *float64 `yaml:"maxSpeed,omitempty"`
	AngularVelocity *float64 `yaml:"angularVelocity,omitempty"`
	AimTime         *int     `yaml:"aimTime,omitempty"`
	AimDelay        *int     `yaml:"aimDelay,omitempty"`
	Damage          *int     `yaml:"damage,omitempty"`
}

type BulletGroup struct {
	Alias         *string `yaml:"alias,omitempty"`
	Angle         *string `yaml:"angle,omitempty"`
	FixedAngle    *int    `yaml:"fixedAngle,omitempty"`
	BulletCount   *int    `yaml:"bulletCount,omitempty"`
	LastSpawnedAt *int    `yaml:"lastSpawnedAt,omitempty"`
	SpawnRate     *int    `yaml:"spawnRate,omitempty"`
	LoopCount     *int    `yaml:"loopCount,omitempty"`
	Bullet        *Bullet `yaml:"bullet,omitempty"`
}

func (m *BulletGroup) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type bulletGroup BulletGroup
	if err := unmarshal((*bulletGroup)(m)); err != nil {
		return err
	}

	return nil
}
