package game

import (
	"ebijam23/resources"
	"ebijam23/states"
	"errors"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	ErrMissingMap  = errors.New("missing map")
	ErrNoActiveMap = errors.New("no active map")
)

const (
	cellW = 16
	cellH = 16
)

type Map struct {
	filename   string
	data       *resources.Map
	actors     []Actor
	bullets    []*Bullet
	conditions []*resources.ConditionDef
	cleared    bool
	currentZ   int // This isn't the right location for this, but we need to keep track of the current/active Z for rendering appropriate fading.
	vfx        resources.VFXList
}

func (s *World) TravelToMap(ctx states.Context, mapName string) error {
	mapData := ctx.Manager.GetAs("maps", mapName, (*resources.Map)(nil)).(*resources.Map)
	if mapData == nil {
		return ErrMissingMap
	}

	m := &Map{
		filename:   mapName,
		data:       mapData,
		conditions: mapData.Conditions,
	}

	//wallH := 6

	for i, l := range m.data.Layers {
		//xoffset := float64(wallH * len(m.data.Layers))
		//yoffset := float64(wallH * 2 * len(m.data.Layers))
		xoffset := 0.0
		yoffset := 0.0
		for j, row := range l.Cells {
			for k, cell := range row {
				if r, ok := m.data.RuneMap[string(cell.Type)]; ok {
					cell.ID = r.ID
					cell.BlockMove = r.BlockMove
					cell.BlockView = r.BlockView
					cell.Wall = r.Wall
					cell.Floor = r.Floor
					cell.Sprite = resources.NewSprite(ctx.Manager.GetAs("images", r.Sprite, (*ebiten.Image)(nil)).(*ebiten.Image))
					cell.Sprite.SetXY(
						cellW*float64(k)+xoffset,
						cellH*float64(j)+yoffset,
					)
					// This is gross, but it visually allows the cell to be "isometric" without going the standard walls path.
					if r.Isometric {
						cell.Sprite.SetXY(
							cell.Sprite.X-(cellW/4),
							cell.Sprite.Y,
						)
					}
				}
				//cell.Sprite.Centered = true
				row[k] = cell
			}
			l.Cells[j] = row
		}
		m.data.Layers[i] = l
	}

	// Create actors.
	// Create map of actor IDs to actors
	interactiveMap := make(map[string]*Interactive)
	for _, a := range m.data.Actors {
		cell := m.FindCellById(a.ID)
		// TODO: consolidate this junk elsewhere
		x := float64(a.Spawn[0]) * cellW
		y := float64(a.Spawn[1]) * cellH
		if cell != nil {
			x = float64(cell.Sprite.X)
			y = float64(cell.Sprite.Y)
		}
		switch a.Type {
		case "interactive":
			spritePrefix := a.Sprite
			// Create the active sprite
			activeImageNames := ctx.Manager.GetNamesWithPrefix("images", spritePrefix+"-active")
			activeImages := make([]*ebiten.Image, 0)
			for _, s := range activeImageNames {
				activeImages = append(activeImages, ctx.Manager.GetAs("images", s, (*ebiten.Image)(nil)).(*ebiten.Image))
			}
			activeSprite := resources.NewAnimatedSprite(activeImages)
			activeSprite.X = x
			activeSprite.Y = y
			activeSprite.Framerate = 5
			activeSprite.Loop = true

			// Create the inactive sprite
			inactiveImageNames := ctx.Manager.GetNamesWithPrefix("images", spritePrefix+"-inactive")
			inactiveImages := make([]*ebiten.Image, 0)
			for _, s := range inactiveImageNames {
				inactiveImages = append(inactiveImages, ctx.Manager.GetAs("images", s, (*ebiten.Image)(nil)).(*ebiten.Image))
			}
			inactiveSprite := resources.NewAnimatedSprite(inactiveImages)
			inactiveSprite.X = x
			inactiveSprite.Y = y

			interactive := CreateInteractive(
				x,
				y,
				a.ID,
				a.Active,
				true,
				a.Condtions,
				activeSprite,
				inactiveSprite,
			)
			interactive.degrade = a.Degrade
			m.actors = append(m.actors, interactive)
			interactiveMap[a.ID] = interactive
		case "spawner":
			bulletGroups := make([]*BulletGroup, 0)

			// If we have bullet groups defined for the spawner, create them.
			if len(a.BulletGroups) > 0 {
				for _, bg := range a.BulletGroups {
					bulletAlias := (*resources.BulletGroupDef)(nil)
					if bg.Alias != nil {
						bulletAlias = ctx.Manager.GetAs("bullets", *bg.Alias, (*resources.BulletGroupDef)(nil)).(*resources.BulletGroupDef)
					}
					bulletGroups = append(bulletGroups, CreateBulletGroupFromDef(x, y, bg, bulletAlias))
				}
			}
			spawner := CreateSpawner(x, y, bulletGroups)
			m.actors = append(m.actors, spawner)
		}
	}

	// After creating the actors, link them up
	for _, a := range m.data.Actors {
		if a.Linked != nil {
			// Find the actor in the map
			i := interactiveMap[a.ID]
			if i == nil {
				continue
			}

			// Find the linked actors in the map and link them up
			for _, b := range a.Linked {
				childActor := interactiveMap[b]
				if childActor == nil {
					continue
				}
				i.AddLinkedInteractive(childActor)
			}
		}
	}

	// Set proper active layer.
	m.currentZ = m.data.Start[2]

	// Only add fade and title VFX if this map is not the same as the previous one.
	if s.activeMap == nil || s.activeMap.data != m.data {
		// Add fade in VFX.
		m.vfx.Add(&resources.Fade{
			Alpha:        1,
			Duration:     1 * time.Second,
			ApplyToImage: true,
		})
		// Add map title VFX.
		m.vfx.Add(&resources.Text{
			Text:         m.data.Title,
			Scale:        2.0,
			X:            320,
			Y:            320,
			Delay:        400 * time.Millisecond,
			Outline:      true,
			OutlineColor: color.NRGBA{0x22, 0x8b, 0x22, 0xff},
			InDuration:   500 * time.Millisecond,
			HoldDuration: 2 * time.Second,
			OutDuration:  500 * time.Millisecond,
		})

	}

	s.activeMap = m

	// Move players over to new map.
	for _, p := range s.Players {
		// Save actor right before entry.
		p.Actor().Save()
		// Position the actor and place them in the map.
		p.Actor().SetXY(float64(m.data.Start[0]*cellW), float64(m.data.Start[1]*cellH))
		m.actors = append(m.actors, p.Actor())
	}

	return nil
}

func (s *World) ResetActiveMap(ctx states.Context) error {
	if s.activeMap == nil {
		return ErrNoActiveMap
	}

	for _, p := range s.Players {
		p.Actor().Restore()
	}

	s.activeMap.actors = make([]Actor, 0)
	s.activeMap.bullets = make([]*Bullet, 0)
	return s.TravelToMap(ctx, s.activeMap.filename)
}

func (m *Map) Draw(ctx states.DrawContext) {
	wallH := 6

	for z := len(m.data.Layers) - 1; z >= 0; z-- {
		layerOpts := &ebiten.DrawImageOptions{}
		// Offset the layer -- this makes the player's collision position look better.
		layerOpts.GeoM.Translate(2, 5)

		//zv := float64(z) / float64(len(m.data.Layers))
		dz := 1.0 - math.Abs(float64(z)-float64(m.currentZ))*0.5

		layerOpts.GeoM.Translate(float64(wallH*z), float64(wallH*z)*2)

		// TODO: Draw/render operations should probably be queued, sorted by z-index, then rendered in game.Draw.
		l := m.data.Layers[z]
		for _, row := range l.Cells {
			for _, cell := range row {
				if cell.Floor {
					opts := &ebiten.DrawImageOptions{}
					opts.GeoM.Concat(layerOpts.GeoM)
					for i := wallH / 3; i >= 0; i-- {
						ds := 1 - float32(i)/float32(wallH/3)*float32(dz)
						opts.ColorScale.Reset()
						opts.ColorScale.Scale(ds, ds, ds, 1.0)
						cell.Sprite.DrawWithOptions(ctx, opts)
						opts.GeoM.Translate(-1, -2)
					}
				} else if cell.Wall {
					opts := &ebiten.DrawImageOptions{}
					opts.GeoM.Concat(layerOpts.GeoM)
					for i := 0; i < wallH; i++ {
						ds := float32(i) / float32(wallH) * float32(dz)
						opts.ColorScale.Reset()
						opts.ColorScale.Scale(ds, ds, ds, 1.0)
						cell.Sprite.DrawWithOptions(ctx, opts)
						opts.GeoM.Translate(-1, -2)
					}
				} else {
					cell.Sprite.Draw(ctx)
				}
			}
		}
	}
	for _, a := range m.actors {
		a.Draw(ctx)
	}
	for _, b := range m.bullets {
		b.Draw(ctx)
	}

	m.vfx.Process(ctx, nil)
}

func (m *Map) GetCell(x, y, z int) (resources.Cell, error) {
	if z < 0 || z >= len(m.data.Layers) || y < 0 || y >= len(m.data.Layers[z].Cells) || x < 0 || x >= len(m.data.Layers[z].Cells[y]) {
		return resources.Cell{}, errors.New("no such cell")
	}
	return m.data.Layers[z].Cells[y][x], nil
}

func (m *Map) FindCellById(id string) *resources.Cell {
	for z, layer := range m.data.Layers {
		for y, row := range layer.Cells {
			for x, cell := range row {
				if cell.ID == id {
					return &m.data.Layers[z].Cells[y][x]
				}
			}
		}
	}
	return nil
}

type CellCollision struct {
	Cell resources.Cell
}

func (m *Map) Collides(s Shape) *CellCollision {
	// Get nearest cell to shape coordinates and check adjacent cells for collisions.
	x, y, _, _ := s.Bounds()
	x /= cellW
	y /= cellH
	x = math.Round(x)
	y = math.Round(y)
	z := m.currentZ

	check := func(x, y int) *CellCollision {
		if y >= 0 && int(y) < len(m.data.Layers[z].Cells) && x >= 0 && int(x) < len(m.data.Layers[z].Cells[int(y)]) {
			if m.data.Layers[z].Cells[int(y)][int(x)].BlockMove {
				cell := m.data.Layers[z].Cells[int(y)][int(x)]
				if s.Collides(&RectangleShape{
					X:      cell.Sprite.X,
					Y:      cell.Sprite.Y,
					Width:  cell.Sprite.Width(),
					Height: cell.Sprite.Height(),
				}) {
					return &CellCollision{
						Cell: cell,
					}
				}
			}
		}
		return nil
	}

	// TODO: Get whatever its called when you get the minimum distance to ensure contact but not intersection and return it so the caller can still potentially move.
	// lol, this is bad
	collision := check(int(x), int(y))
	if collision == nil {
		collision = check(int(x), int(y+1))
	}
	if collision == nil {
		collision = check(int(x), int(y-1))
	}
	if collision == nil {
		collision = check(int(x-1), int(y))
	}
	if collision == nil {
		collision = check(int(x-1), int(y+1))
	}
	if collision == nil {
		collision = check(int(x-1), int(y-1))
	}
	if collision == nil {
		collision = check(int(x+1), int(y))
	}
	if collision == nil {
		collision = check(int(x+1), int(y+1))
	}
	if collision == nil {
		collision = check(int(x+1), int(y-1))
	}

	return collision
}

func (m *Map) DoesLineCollide(fx1, fy1, fx2, fy2 float64, z int) bool {
	x1 := int(math.Round(fx1 / cellW))
	y1 := int(math.Round(fy1 / cellH))
	x2 := int(math.Round(fx2 / cellW))
	y2 := int(math.Round(fy2 / cellH))

	// Bresenham's line algorithm.
	dx := math.Abs(float64(x2 - x1))
	dy := math.Abs(float64(y2 - y1))
	var sx, sy int

	if x1 < x2 {
		sx = 1
	} else {
		sx = -1
	}
	if y1 < y2 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	for {
		if cell, er := m.GetCell(x1, y1, z); er == nil {
			if cell.BlockView {
				return true
			}
		}
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}

	return false
}

func (m *Map) GetInteractiveActors() []*Interactive {
	var actors []*Interactive
	for _, a := range m.actors {
		if a, ok := a.(*Interactive); ok {
			actors = append(actors, a)
		}
	}
	return actors
}
