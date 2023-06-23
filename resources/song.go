package resources

import (
	"io"

	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
)

type Song struct {
	stream *vorbis.Stream
}

func (s *Song) Stream() *vorbis.Stream {
	return s.stream
}

func NewSong(rs io.ReadSeeker) (*Song, error) {
	song := &Song{}

	s, err := vorbis.DecodeWithSampleRate(44100, rs)
	if err != nil {
		return nil, err
	}

	song.stream = s

	return song, nil
}
