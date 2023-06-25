package resources

import (
	"fmt"
	"strings"
)

type RuneDef struct {
	Sprite    string `yaml:"sprite"`
	BlockView bool   `yaml:"blockView"`
	BlockMove bool   `yaml:"blockMove"`
	Wall      bool   `yaml:"wall"`
	Floor     bool   `yaml:"floor"`
	Isometric bool   `yaml:"isometric"`
	ID        string `yaml:"id,omitempty"`
}

type Cell struct {
	Sprite    *Sprite `yaml:"-"`
	Type      rune    `yaml:"-"` // I guess using runes is okay.
	BlockView bool    `yaml:"-"`
	BlockMove bool    `yaml:"-"`
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
	Conditions   []*ConditionDef    `yaml:"conditions"`
	RuneMap      map[string]RuneDef `yaml:"runes"`
	Layers       []Layer            `yaml:"-"`
	Width        int                `yaml:"width"`
	Height       int                `yaml:"height"`
	SourceLayers []string           `yaml:"layers"`
	Actors       []ActorSpawn       `yaml:"actors"`
}

type Interactive struct {
	Map        *string         `yaml:"map,omitempty"`
	Conditions []*ConditionDef `yaml:"conditions,omitempty"`
	Linked     []string        `yaml:"linked,omitempty"`
	Collidable bool            `yaml:"collidable,omitempty"`
	Reversable bool            `yaml:"reversable"`
	Degrade    bool            `yaml:"degrade"`
	Active     bool            `yaml:"active"`
	Touchable  bool            `yaml:"touchable"`
	Health     int             `yaml:"health"`
}

type ActorSpawn struct {
	ID           string         `yaml:"id"`
	Spawn        [3]int         `yaml:"spawn,omitempty"`
	Type         string         `yaml:"type"`
	Sprite       string         `yaml:"sprite"`
	BulletGroups []*BulletGroup `yaml:"bullets,omitempty"`
	Interactive  *Interactive   `yaml:"interactive,omitempty"`
}

func (m *Map) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type mapAlias Map
	if err := unmarshal((*mapAlias)(m)); err != nil {
		fmt.Printf("Error unmarshalling map: %v", err)
		return err
	}

	for _, layer := range m.SourceLayers {
		var l Layer

		rows := strings.Split(layer, "\n")
		if len(rows) > m.Height {
			m.Height = len(rows)
		}
		for _, row := range rows {
			l.Cells = append(l.Cells, []Cell{})
			if len(row) > m.Width {
				m.Width = len(row)
			}
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
