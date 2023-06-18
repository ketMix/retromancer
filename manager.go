package main

import (
	"ebijam23/resources"
	"errors"
	"fmt"
	"image/color"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/kettek/go-multipath/v2"
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
	}

	return nil, ErrNoSuchCategory
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
	}
	return nil
}

func (m *ResourceManager) LoadAll() error {
	m.files.Walk("images/", func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			m.Load("images", entry.Name())
		}
		return nil
	})
	m.files.Walk("maps/", func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			m.Load("maps", entry.Name())
		}
		return nil
	})
	return nil
}
