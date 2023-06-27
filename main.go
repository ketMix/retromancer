package main

import (
	"ebijam23/resources"
	"ebijam23/states"
	"ebijam23/states/menu"
	"embed"
	"flag"
	"io/fs"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/kettek/go-multipath/v2"
	"golang.org/x/image/font/sfnt"

	gaem "ebijam23/states/game"
)

//go:embed assets/*
var embedFS embed.FS

func main() {
	game := &Game{}

	// Parse in flags.
	flag.Float64Var(&game.Flags.MusicVolume, "volume", 1.0, "volume to play music at")
	flag.Float64Var(&game.Flags.SoundVolume, "sound", 1.0, "volume to play sound at")
	flag.BoolVar(&game.Flags.Muted, "mute", false, "whether to start muted")
	flag.BoolVar(&game.Flags.Fullscreen, "fullscreen", false, "whether to start in fullscreen")
	flag.BoolVar(&game.Flags.SkipIntro, "skip-intro", false, "whether to skip the intro")
	flag.StringVar(&game.Flags.Locale, "locale", "en", "locale to use")
	flag.StringVar(&game.Flags.Font, "font", "x12y16pxMaruMonica", "font to use")
	flag.StringVar(&game.Flags.Map, "map", "", "map to load")
	flag.Parse()

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

	// Set our locale.
	game.Localizer.manager = &game.Manager
	game.Localizer.SetLocale(game.Flags.Locale, false) // Start without GPT
	game.Localizer.InitGPT()

	// Initialize game fields as necessary.
	if err := game.Init(); err != nil {
		panic(err)
	}

	// Ensure we have our font.
	if f := game.Manager.GetAs("fonts", game.Flags.Font, (*sfnt.Font)(nil)).(*sfnt.Font); f == nil {
		panic("missing font")
	} else {
		game.Text.SetFont(f)
		game.Text.Utils().SetCache8MiB()
	}

	// Initialize audio.
	audio.NewContext(44100)

	// Set initial mute/volume
	if game.Flags.Muted {
		game.MusicPlayer.SetVolume(0)
		resources.Volume = 0
	} else {
		game.MusicPlayer.volume = game.Flags.MusicVolume
		resources.Volume = game.Flags.SoundVolume
	}

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("ebijam23")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if game.Flags.Fullscreen {
		ebiten.SetFullscreen(true)
	}

	// Set up menu.
	game.PushState(&menu.Menu{})

	// Quick skip for map testing.
	if game.Flags.Map != "" {
		game.PushState(&gaem.World{
			StartingMap: game.Flags.Map,
			Players: []gaem.Player{
				gaem.NewLocalPlayer(),
			},
		})
	} else {
		if !game.Flags.SkipIntro {
			game.PushState(&menu.Intro{})
			game.PushState(&menu.Loading{})
		}
	}

	// Set up loading screen.
	if err := ebiten.RunGame(game); err != nil {
		if err == states.ErrQuitGame {
			return
		}
		panic(err)
	}
}
