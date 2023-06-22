package resources

import (
	"strings"
)

type RuneDef struct {
	Sprite    string `yaml:"sprite"`
	Blocks    bool   `yaml:"blocks"`
	Wall      bool   `yaml:"wall"`
	Floor     bool   `yaml:"floor"`
	Door      bool   `yaml:"door"`
	Isometric bool   `yaml:"isometric"`
	Map       string `yaml:"map"`
}

type Cell struct {
	Sprite    *Sprite `yaml:"-"`
	Type      rune    `yaml:"-"` // I guess using runes is okay.
	Blocks    bool    `yaml:"-"`
	Wall      bool    `yaml:"-"`
	Floor     bool    `yaml:"-"`
	Door      bool    `yaml:"-"`
	Map       string  `yaml:"-"`
	Isometric bool    `yaml:"-"`
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

type ActorSpawn struct {
	ID           string            `yaml:"id"`
	Spawn        [3]int            `yaml:"spawn"`
	Type         string            `yaml:"type"`
	Active       bool              `yaml:"active"`
	Conditions   []*ConditionDef   `yaml:"conditions,omitempty"`
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
