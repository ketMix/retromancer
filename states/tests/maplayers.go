package tests

import (
	"ebijam23/resources"
	"ebijam23/states"

	"github.com/hajimehoshi/ebiten/v2"
)

type MapLayer struct {
	Map *resources.Map
}

func (m *MapLayer) Init(ctx states.Context) error {
	m.Map = ctx.Manager.GetAs("maps", "test", (*resources.Map)(nil)).(*resources.Map)

	for i, l := range m.Map.Layers {
		for j, row := range l.Cells {
			for k, cell := range row {
				switch cell.Type {
				case '#':
					cell.Sprite = resources.NewSprite(ctx.Manager.GetAs("images", "wall", (*ebiten.Image)(nil)).(*ebiten.Image))
				case 'f':
					cell.Sprite = resources.NewSprite(ctx.Manager.GetAs("images", "fire", (*ebiten.Image)(nil)).(*ebiten.Image))
				case '.':
					cell.Sprite = resources.NewSprite(ctx.Manager.GetAs("images", "floor", (*ebiten.Image)(nil)).(*ebiten.Image))
				default:
					cell.Sprite = resources.NewSprite(ctx.Manager.GetAs("images", "empty", (*ebiten.Image)(nil)).(*ebiten.Image))
				}
				cell.Sprite.SetXY(
					cell.Sprite.Width()*float64(k),
					cell.Sprite.Height()*float64(j),
				)
				row[k] = cell
			}
			l.Cells[j] = row
		}
		m.Map.Layers[i] = l
	}

	return nil
}

func (m *MapLayer) Update(ctx states.Context) error {
	return nil
}

func (m *MapLayer) Draw(screen *ebiten.Image) {
	for i := len(m.Map.Layers) - 1; i >= 0; i-- {
		l := m.Map.Layers[i]
		for j, row := range l.Cells {
			for k := range row {
				m.DrawCoord(k, j, i, screen)
			}
		}
	}
}

func (m *MapLayer) DrawCoord(x, y, z int, screen *ebiten.Image) {
	if z >= len(m.Map.Layers) {
		return
	}
	if y >= len(m.Map.Layers[z].Cells) {
		return
	}
	if x >= len(m.Map.Layers[z].Cells[y]) {
		return
	}
	l := m.Map.Layers[z]
	cell := l.Cells[y][x]
	opts := &ebiten.DrawImageOptions{}

	opts.GeoM.Translate(0, cell.Sprite.Height()/2*float64(len(m.Map.Layers)))

	zv := 1 - float64(z)/float64(len(m.Map.Layers))
	//zv := 1

	cellH := 4

	opts.GeoM.Translate(0, float64(z)*float64(cellH))

	if cell.Type == '.' {
		//opts.GeoM.Translate(1, 2)
	}
	if cell.Type == 'f' {
		cell.Sprite.DrawWithOptions(screen, opts)
	} else {
		for i := 0; i < cellH; i++ {
			v := float32(i) / float32(cellH) * float32(zv)
			opts.ColorScale.SetR(v)
			opts.ColorScale.SetG(v)
			opts.ColorScale.SetB(v)
			cell.Sprite.DrawWithOptions(screen, opts)
			opts.GeoM.Translate(-1, -2)
		}
	}
	/*} else if cell.Type == '.' {
		opts.GeoM.Translate(0, -2)
		cell.Sprite.DrawWithOptions(screen, opts)
	}*/
}
