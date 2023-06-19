package tests

import (
	"ebijam23/resources"
	"ebijam23/states"
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type MapLayer struct {
	Map      *resources.Map
	Which    int
	currentZ int
}

func (m *MapLayer) Init(ctx states.Context) error {
	m.Map = ctx.Manager.GetAs("maps", "test2", (*resources.Map)(nil)).(*resources.Map)

	for i, l := range m.Map.Layers {
		for j, row := range l.Cells {
			for k, cell := range row {
				switch m.Map.RuneMap[string(cell.Type)] {
				case "wall":
					cell.Sprite = resources.NewSprite(ctx.Manager.GetAs("images", "wall", (*ebiten.Image)(nil)).(*ebiten.Image))
				case "fire":
					cell.Sprite = resources.NewSprite(ctx.Manager.GetAs("images", "fire", (*ebiten.Image)(nil)).(*ebiten.Image))
				case "floor":
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
	if ebiten.IsKeyPressed(ebiten.Key1) {
		m.Which = 0
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		m.Which = 1
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		m.currentZ++
		if m.currentZ >= len(m.Map.Layers) {
			m.currentZ = 0
		}
		fmt.Println(m.currentZ)
	}
	return nil
}

func (m *MapLayer) Draw(screen *ebiten.Image) {
	for i := len(m.Map.Layers) - 1; i >= 0; i-- {
		l := m.Map.Layers[i]
		for j, row := range l.Cells {
			for k := range row {
				if m.Which == 1 {
					m.DrawCoord(k, j, i, screen)
				} else {
					m.DrawCoord2(k, j, i, screen)
				}
			}
		}
	}
}

func (m *MapLayer) DrawCoord2(x, y, z int, screen *ebiten.Image) {
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

	opts.GeoM.Translate(0, cell.Sprite.Height()*float64(len(m.Map.Layers)))

	// Get the distance from our z to the current z:
	//dz := math.Abs(float64(m.currentZ-z)) / float64(len(m.Map.Layers))
	dz := math.Abs(float64(m.currentZ-z)) / 2
	//fmt.Println(currentZ, z, dz)
	zv := math.Max(0.0, 1.0-dz)
	fmt.Println(m.currentZ, z, zv)

	//zv := 1.0 - float64(z)/float64(len(m.Map.Layers))

	/*opts.ColorScale.SetR(float32(zv))
	opts.ColorScale.SetG(float32(zv))
	opts.ColorScale.SetB(float32(zv))*/
	sx := math.Max(1.0-dz, 0.8)
	sy := math.Max(1.0-dz, 0.8)
	if sx != 1 || sy != 1 {
		opts.GeoM.Scale(sx, sy)
		opts.GeoM.Translate(float64(len(m.Map.Layers[z].Cells[y]))*cell.Sprite.Width()/2, float64(len(m.Map.Layers[z].Cells))*cell.Sprite.Height()/2)
		opts.GeoM.Translate(-float64(len(m.Map.Layers[z].Cells[y]))*cell.Sprite.Width()/2*sx, -float64(len(m.Map.Layers[z].Cells))*cell.Sprite.Height()/2*sy)
	}
	opts.ColorScale.ScaleAlpha(float32(zv))
	cell.Sprite.DrawWithOptions(screen, opts)
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
	dz := 0.0
	if z >= m.currentZ {
		dz = 1
	} else if z < m.currentZ {
		dz = math.Max(0.0, 1-float64(m.currentZ-z))
	}

	cellH := 6

	opts.GeoM.Translate(0, float64(z)*float64(cellH))

	if cell.Type == '.' {
		opts.GeoM.Translate(1, 2)
	}
	if cell.Type == 'f' {
		opts.ColorScale.ScaleAlpha(float32(dz))
		cell.Sprite.DrawWithOptions(screen, opts)
	} else if cell.Type == '.' {
		for i := 0; i < 2; i++ {
			v := float32(i) / float32(cellH) * float32(zv)
			opts.ColorScale.SetR(v)
			opts.ColorScale.SetG(v)
			opts.ColorScale.SetB(v)
			//opts.ColorScale.SetA(1 - float32(dz))
			opts.ColorScale.ScaleAlpha(float32(dz))
			cell.Sprite.DrawWithOptions(screen, opts)
			opts.GeoM.Translate(-1, -2)
		}
	} else {
		for i := 0; i < cellH; i++ {
			v := float32(i) / float32(cellH) * float32(zv)
			opts.ColorScale.SetR(v)
			opts.ColorScale.SetG(v)
			opts.ColorScale.SetB(v)
			opts.ColorScale.ScaleAlpha(float32(dz))
			cell.Sprite.DrawWithOptions(screen, opts)
			opts.GeoM.Translate(-1, -2)
		}
	}
	/*} else if cell.Type == '.' {
		opts.GeoM.Translate(0, -2)
		cell.Sprite.DrawWithOptions(screen, opts)
	}*/
}
