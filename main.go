package main

import (
	"ebijam23/states"
	"ebijam23/states/menu"
	"embed"
	"io/fs"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kettek/go-multipath/v2"
)

//go:embed assets/*
var embedFS embed.FS

func main() {
	game := &Game{}

	// Allow loading from filesystem.
	game.Manager.files.InsertFS(os.DirFS("assets"), multipath.FirstPriority)

	// Also allow loading from embedded filesystem.
	sub, err := fs.Sub(embedFS, "assets")
	if err != nil {
		panic(err)
	}
	game.Manager.files.InsertFS(sub, multipath.LastPriority)

	if err := game.Manager.Setup(); err != nil {
		panic(err)
	}

	// Might as well load all assets up front (for now -- might not want to with music later).
	if err := game.Manager.LoadAll(); err != nil {
		panic(err)
	}

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("ebijam23")

	// Push the state.
	game.PushState(&menu.Menu{})

	if err := ebiten.RunGame(game); err != nil {
		if err == states.ErrQuitGame {
			return
		}
		panic(err)
	}
}
