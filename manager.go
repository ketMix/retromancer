package main

import (
	"bytes"
	"errors"
	"fmt"
	"image/color"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/ketMix/retromancer/resources"
	"github.com/kettek/go-multipath/v2"
	"github.com/tinne26/etxt/font"
	"golang.org/x/image/font/sfnt"
	"gopkg.in/yaml.v2"
)

type ResourceGroup struct {
	data map[string]interface{}
}

type ResourceManager struct {
	files         multipath.FS
	groups        map[string]ResourceGroup
	imageFallback *ebiten.Image
}

var (
	ErrNoSuchCategory   = errors.New("no such category")
	ErrMissingDirectory = errors.New("missing directory")
)

func (m *ResourceManager) Setup() error {
	m.groups = make(map[string]ResourceGroup)
	// Create a default image.
	m.imageFallback = ebiten.NewImage(16, 16)
	m.imageFallback.Fill(color.NRGBA{0xff, 0x00, 0x00, 0xff})
	return nil
}

func (m *ResourceManager) Load(category string, name string) (interface{}, error) {
	if _, ok := m.groups[category]; !ok {
		m.groups[category] = ResourceGroup{
			data: make(map[string]interface{}),
		}
	}

	group := m.groups[category]

	if data, ok := group.data[name]; ok {
		return data, nil
	}

	file, err := m.files.Open(fmt.Sprintf("%s/%s", category, name))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if category == "images" {
		img, _, err := ebitenutil.NewImageFromFileSystem(m.files, fmt.Sprintf("%s/%s", category, name))
		if err != nil {
			return nil, err
		}
		group.data[strings.TrimSuffix(name, filepath.Ext(name))] = img
		return img, nil
	} else if category == "maps" {
		bytes, err := m.files.ReadFile(fmt.Sprintf("%s/%s", category, name))
		if err != nil {
			return nil, err
		}
		var m *resources.Map
		if err := yaml.Unmarshal(bytes, &m); err != nil {
			return nil, err
		}
		group.data[strings.TrimSuffix(name, filepath.Ext(name))] = m
		return m, nil
	} else if category == "bullets" {
		bytes, err := m.files.ReadFile(fmt.Sprintf("%s/%s", category, name))
		if err != nil {
			return nil, err
		}
		var bg *resources.BulletGroup
		if err := yaml.Unmarshal(bytes, &bg); err != nil {
			return nil, err
		}
		group.data[strings.TrimSuffix(name, filepath.Ext(name))] = bg
		return bg, nil
	} else if category == "enemies" {
		bytes, err := m.files.ReadFile(fmt.Sprintf("%s/%s", category, name))
		if err != nil {
			return nil, err
		}
		var e *resources.Enemy
		if err := yaml.Unmarshal(bytes, &e); err != nil {
			return nil, err
		}
		group.data[strings.TrimSuffix(name, filepath.Ext(name))] = e
		return e, nil
	} else if category == "fonts" {
		if strings.HasSuffix(name, ".ttf") {
			bytes, err := m.files.ReadFile(fmt.Sprintf("%s/%s", category, name))
			if err != nil {
				return nil, err
			}

			f, _, err := font.ParseFromBytes(bytes)
			if err != nil {
				return nil, err
			}
			group.data[strings.TrimSuffix(name, filepath.Ext(name))] = f
			return f, nil
		} else {
			return nil, nil
		}
	} else if category == "sounds" {
		if strings.HasSuffix(name, ".ogg") {
			bytes, err := m.files.ReadFile(fmt.Sprintf("%s/%s", category, name))
			if err != nil {
				return nil, err
			}

			snd, err := resources.NewSound(bytes)
			if err != nil {
				return nil, err
			}
			group.data[strings.TrimSuffix(name, filepath.Ext(name))] = snd
			return snd, nil
		} else {
			return nil, nil
		}
	} else if category == "locales" {
		if strings.HasSuffix(name, ".yaml") {
			bytes, err := m.files.ReadFile(fmt.Sprintf("%s/%s", category, name))
			if err != nil {
				return nil, err
			}

			var l *resources.Locale
			if err := yaml.Unmarshal(bytes, &l); err != nil {
				return nil, err
			}
			group.data[strings.TrimSuffix(name, filepath.Ext(name))] = l
			return l, nil
		} else {
			return nil, nil
		}
	} else if category == "songs" {
		if strings.HasSuffix(name, ".ogg") {
			b, err := m.files.ReadFile(fmt.Sprintf("%s/%s", category, name))
			if err != nil {
				return nil, err
			}
			song, err := resources.NewSong(bytes.NewReader(b))
			if err != nil {
				return nil, err
			}
			group.data[strings.TrimSuffix(name, filepath.Ext(name))] = song
			return song, nil
		} else {
			return nil, nil
		}
	}

	return nil, ErrNoSuchCategory
}

func (m *ResourceManager) GetNamesWithPrefix(category string, prefix string) []string {
	if c, ok := m.groups[category]; !ok {
		return nil
	} else {
		var names []string
		for k := range c.data {
			if strings.HasPrefix(k, prefix) {
				names = append(names, k)
			}
		}
		sort.Slice(names, func(i, j int) bool {
			return strings.Compare(names[i], names[j]) < 0
		})
		return names
	}
}

func (m *ResourceManager) GetWithPrefix(category string, prefix string) (items []interface{}) {
	names := m.GetNamesWithPrefix(category, prefix)
	for _, name := range names {
		items = append(items, m.Get(category, name))
	}
	return
}

func (m *ResourceManager) GetAsWithPrefix(category string, prefix string, target interface{}) (items []interface{}) {
	names := m.GetNamesWithPrefix(category, prefix)
	for _, name := range names {
		items = append(items, m.GetAs(category, name, target))
	}
	return
}

func (m *ResourceManager) Get(category string, name string) interface{} {
	if c, ok := m.groups[category]; !ok {
		return nil
	} else {
		return c.data[name]
	}
}

func (m *ResourceManager) GetAs(category string, name string, target interface{}) interface{} {
	switch target.(type) {
	case *ebiten.Image:
		d := m.Get(category, name)
		if d == nil {
			return m.imageFallback
		}
		return d
	case *resources.Map:
		d := m.Get(category, name)
		if d == nil {
			return &resources.Map{} // FIXME: Use an actual fallback map.
		}
		return d
	case *resources.BulletGroup:
		d := m.Get(category, name)
		if d == nil {
			return &resources.BulletGroup{} // FIXME: Use an actual fallback bullet group.
		}
		return d
	case *resources.Enemy:
		d := m.Get(category, name)
		if d == nil {
			return &resources.Enemy{} // FIXME: Use an actual fallback enemy.
		}
		return d
	case *sfnt.Font:
		d := m.Get(category, name)
		if d == nil {
			return &sfnt.Font{} // FIXME: Use an actual fallback font.
		}
		return d
	case *resources.Sound:
		d := m.Get(category, name)
		if d == nil {
			return &resources.Sound{} // FIXME: Use an actual fallback sound.
		}
		return d
	case *resources.Locale:
		d := m.Get(category, name)
		if d == nil {
			d = m.Get("locales", "en.yaml")
			if d == nil {
				return &resources.Locale{} // This shouldn't be reached.
			}
		}
		return d
	case *resources.Song:
		d := m.Get(category, name)
		if d == nil {
			return &resources.Song{} // FIXME: Use an actual fallback song.
		}
		return d
	}

	return nil
}

func (m *ResourceManager) LoadDir(category, path string) error {
	err := m.files.Walk(path, func(path string, entry fs.DirEntry, err error) error {
		if entry == nil {
			return ErrMissingDirectory
		}
		if !entry.IsDir() {
			if _, err := m.Load(category, entry.Name()); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	fmt.Println("loaded", len(m.groups[category].data), category)
	return nil
}

func (m *ResourceManager) LoadAll() error {
	if err := m.LoadDir("images", "images/"); err != nil {
		return err
	}
	if err := m.LoadDir("maps", "maps/"); err != nil {
		return err
	}
	if err := m.LoadDir("bullets", "bullets/"); err != nil {
		return err
	}
	if err := m.LoadDir("enemies", "enemies/"); err != nil {
		return err
	}
	if err := m.LoadDir("fonts", "fonts/"); err != nil {
		return err
	}
	if err := m.LoadDir("sounds", "sounds/"); err != nil {
		return err
	}
	if err := m.LoadDir("locales", "locales/"); err != nil {
		return err
	}
	if err := m.LoadDir("songs", "songs/"); err != nil {
		return err
	}
	return nil
}
