package resources

import (
	"bytes"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

type Sound struct {
	bytes []byte
}

var Volume float64 = 1.0

func NewSound(data []byte) (*Sound, error) {
	stream, err := vorbis.DecodeWithSampleRate(44100, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	_, err = b.ReadFrom(stream)
	if err != nil {
		return nil, err
	}

	return &Sound{
		bytes: b.Bytes(),
	}, nil
}

func (s *Sound) Play(v float64) *audio.Player {
	player := audio.CurrentContext().NewPlayerFromBytes(s.bytes)
	player.SetVolume(v * Volume)
	player.Play()
	return player
}
