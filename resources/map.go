package resources

import (
	"strings"
)

type Cell struct {
	Sprite *Sprite `yaml:"-"`
	Type   rune    `yaml:"-"` // I guess using runes is okay.
}

type Layer struct {
	Cells [][]Cell `yaml:"-"`
}

type Map struct {
	Title        string
	RuneMap      map[string]string `yaml:"runes"`
	Layers       []Layer           `yaml:"-"`
	SourceLayers []string          `yaml:"layers"`
	Start        [3]int            `yaml:"start"`
	Actors       []ActorSpawn      `yaml:"actors"`
}

type ActorSpawn struct {
	Spawn [3]int `yaml:"spawn"`
	Type  string `yaml:"type"`
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
