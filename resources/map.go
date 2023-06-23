package resources

import (
	"strings"
)

type RuneDef struct {
	Sprite    string `yaml:"sprite"`
	Blocks    bool   `yaml:"blocks"`
	Wall      bool   `yaml:"wall"`
	Floor     bool   `yaml:"floor"`
	Isometric bool   `yaml:"isometric"`
	ID        string `yaml:"id,omitempty"`
}

type Cell struct {
	Sprite    *Sprite `yaml:"-"`
	Type      rune    `yaml:"-"` // I guess using runes is okay.
	Blocks    bool    `yaml:"-"`
	Wall      bool    `yaml:"-"`
	Floor     bool    `yaml:"-"`
	Isometric bool    `yaml:"-"`
	ID        string  `yaml:"-"`
}

type Layer struct {
	Cells [][]Cell `yaml:"-"`
}

type Map struct {
	Title        string
	RuneMap      map[string]RuneDef `yaml:"runes"`
	Layers       []Layer            `yaml:"-"`
	SourceLayers []string           `yaml:"layers"`
	Start        [3]int             `yaml:"start"`
	Actors       []ActorSpawn       `yaml:"actors"`
}

type DoorDef struct {
	Map        string          `yaml:"map"`
	Open       bool            `yaml:"open"`
	Conditions []*ConditionDef `yaml:"conditions,omitempty"`
}

type ActorSpawn struct {
	ID           string            `yaml:"id"`
	Spawn        [3]int            `yaml:"spawn,omitempty"`
	Type         string            `yaml:"type"`
	Active       bool              `yaml:"active"`
	Door         *DoorDef          `yaml:"door,omitempty"`
	BulletGroups []*BulletGroupDef `yaml:"bullets,omitempty"`
}

func (m *Map) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type mapAlias Map
	if err := unmarshal((*mapAlias)(m)); err != nil {
		return err
	}

	for _, layer := range m.SourceLayers {
		var l Layer

		rows := strings.Split(layer, "\n")
		for _, row := range rows {
			l.Cells = append(l.Cells, []Cell{})
			for _, cell := range row {
				l.Cells[len(l.Cells)-1] = append(l.Cells[len(l.Cells)-1], Cell{
					Type: rune(cell),
				})
			}
		}

		m.Layers = append(m.Layers, l)
	}

	return nil
}
