package main

import (
	"ebijam23/resources"
	"errors"
	"fmt"
	"image/color"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	ErrNoSuchCategory = errors.New("no such category")
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
		var bg *resources.BulletGroupDef
		if err := yaml.Unmarshal(bytes, &bg); err != nil {
			return nil, err
		}
		group.data[strings.TrimSuffix(name, filepath.Ext(name))] = bg
		return bg, nil
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
	case *resources.BulletGroupDef:
		d := m.Get(category, name)
		if d == nil {
			return &resources.BulletGroupDef{} // FIXME: Use an actual fallback bullet group.
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
	}

	return nil
}

func (m *ResourceManager) LoadAll() error {
	m.files.Walk("images/", func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			if _, err := m.Load("images", entry.Name()); err != nil {
				return err
			}
		}
		return nil
	})
	fmt.Println("loaded", len(m.groups["images"].data), "images")
	m.files.Walk("maps/", func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			if _, err := m.Load("maps", entry.Name()); err != nil {
				return err
			}
		}
		return nil
	})
	fmt.Println("loaded", len(m.groups["maps"].data), "maps")
	m.files.Walk("bullets/", func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			if _, err := m.Load("bullets", entry.Name()); err != nil {
				return err
			}
		}
		return nil
	})
	fmt.Println("loaded", len(m.groups["bullets"].data), "bullet groups")
	m.files.Walk("fonts/", func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			if _, err := m.Load("fonts", entry.Name()); err != nil {
				return err
			}
		}
		return nil
	})
	fmt.Println("loaded", len(m.groups["fonts"].data), "fonts")
	m.files.Walk("sounds/", func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			if _, err := m.Load("sounds", entry.Name()); err != nil {
				return err
			}
		}
		return nil
	})
	fmt.Println("loaded", len(m.groups["sounds"].data), "sounds")
	return nil
}
