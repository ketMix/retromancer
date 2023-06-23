package main

import (
	"ebijam23/states"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

type MusicPlayer struct {
	player *audio.Player
	song   states.Song
	loop   bool
	volume float64
}

func (m *MusicPlayer) Play(s states.Song) (err error) {
	if m.player != nil {
		m.player.Pause()
		m.player.Close()
	}

	m.player, err = audio.CurrentContext().NewPlayer(s.Stream())
	if err != nil {
		return
	}
	m.player.SetVolume(m.volume)
	m.player.Play()
	m.song = s
	return
}

func (m *MusicPlayer) Update() {
	if m.player != nil && m.song != nil && !m.player.IsPlaying() {
		m.player.Seek(0)
		m.player.Play()
	}

}

func (m *MusicPlayer) Pause() {
	if m.player != nil {
		m.player.Pause()
	}
}

func (m *MusicPlayer) Resume() {
	if m.player != nil {
		m.player.Play()
	}
}

func (m *MusicPlayer) Loop() bool {
	return m.loop
}

func (m *MusicPlayer) SetLoop(loop bool) {
	m.loop = loop
}

func (m *MusicPlayer) Volume() float64 {
	return m.volume
}

func (m *MusicPlayer) SetVolume(vol float64) {
	m.volume = vol
	if m.player != nil {
		m.player.SetVolume(m.volume)
	}
}
