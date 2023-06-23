package states

import "github.com/hajimehoshi/ebiten/v2/audio/vorbis"

type Song interface {
	Stream() *vorbis.Stream
}

type MusicPlayer interface {
	Play(s Song) error
	Resume()
	Pause()
	Loop() bool
	SetLoop(loop bool)
	Volume() float64
	SetVolume(volume float64)
}
