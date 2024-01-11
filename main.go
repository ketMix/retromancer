package main

import (
	"embed"
	"flag"
	"io/fs"
	"os"

	"github.com/ketMix/retromancer/net"

	"github.com/ketMix/retromancer/states"
	"github.com/ketMix/retromancer/states/menu"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/ketMix/retromancer/resources"
	"github.com/kettek/go-multipath/v2"
	"golang.org/x/image/font/sfnt"
	"gopkg.in/yaml.v2"

	gaem "github.com/ketMix/retromancer/states/game"
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
	flag.IntVar(&net.NetBufferSize, "net-buffer-size", 1024, "network buffer size")
	flag.IntVar(&net.NetDataShards, "net-data-shards", 5, "network data shards")
	flag.IntVar(&net.NetParityShards, "net-parity-shards", 2, "network parity shards")
	flag.IntVar(&net.NetChannelSize, "net-channel-size", 30, "network channel size")
	flag.StringVar(&game.Flags.Difficulty, "difficulty", string(states.DifficultyNormal), "difficulty to play at")
	flag.Parse()

	// Allow loading from filesystem.
	game.Resources.files.InsertFS(os.DirFS("assets"), multipath.FirstPriority)

	// Also allow loading from embedded filesystem.
	sub, err := fs.Sub(embedFS, "assets")
	if err != nil {
		panic(err)
	}
	game.Resources.files.InsertFS(sub, multipath.LastPriority)

	if err := game.Resources.Setup(); err != nil {
		panic(err)
	}

	// Might as well load all assets up front (for now -- might not want to with music later).
	if err := game.Resources.LoadAll(); err != nil {
		panic(err)
	}

	// Load up our gamepad maps.
	if b, err := game.Resources.files.ReadFile("gamepad.yaml"); err != nil {
		panic(err)
	} else {
		var m []resources.GamepadDefinition
		if err := yaml.Unmarshal(b, &m); err != nil {
			panic(err)
		}
		for _, v := range m {
			resources.AddGamepadDefinition(v)
		}
	}

	// Set our locale.
	game.Localizer.resources = &game.Resources
	game.Localizer.SetLocale(game.Flags.Locale, false) // Start without GPT
	game.Localizer.InitGPT()

	// Initialize game fields as necessary.
	if err := game.Init(); err != nil {
		panic(err)
	}

	// Ensure we have our font.
	if f := game.Resources.GetAs("fonts", game.Flags.Font, (*sfnt.Font)(nil)).(*sfnt.Font); f == nil {
		panic("missing font")
	} else {
		game.Text.SetFont(f)
		game.Text.Utils().SetCache8MiB()
	}

	// Initialize audio.
	audio.NewContext(44100)
	game.MusicPlayer.SetVolume(0.5) // Default to 50% volume.

	// Set initial mute/volume
	if game.Flags.Muted {
		game.MusicPlayer.SetVolume(0)
		resources.Volume = 0
	} else {
		game.MusicPlayer.volume = game.Flags.MusicVolume
		resources.Volume = game.Flags.SoundVolume
	}

	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Retromancer")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if game.Flags.Fullscreen {
		ebiten.SetFullscreen(true)
	}

	// Set up menu.
	game.PushState(&menu.Menu{})

	// Quick skip for map testing.
	if game.Flags.Map != "" {
		var difficulty states.Difficulty
		if game.Flags.Difficulty == "hard" {
			difficulty = states.DifficultyHard
		} else if game.Flags.Difficulty == "easy" {
			difficulty = states.DifficultyEasy
		} else {
			difficulty = states.DifficultyNormal
		}

		game.Difficulty = difficulty
		game.PushState(&gaem.World{
			StartingMap: game.Flags.Map,
			Players: []gaem.Player{
				gaem.NewLocalPlayer(),
			},
		})
	} else {
		if !game.Flags.SkipIntro {
			game.PushState(&menu.Intro{})
			game.PushState(&menu.PreIntro{})
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
