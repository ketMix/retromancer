package main

type Flags struct {
	Locale          string
	Font            string
	MusicVolume     float64
	SoundVolume     float64
	Muted           bool
	Fullscreen      bool
	SkipIntro       bool
	Map             string
	NetBufferSize   int
	NetDataShards   int
	NetParityShards int
}
