package audio

import (
	"bytes"
	"io"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"

	"go-inovation/ino/internal/assets"
)

const sampleRate = 44100

var (
	audioContext   = audio.NewContext(sampleRate)
	soundFilenames = []string{
		"damage.wav",
		"heal.wav",
		"ino1.ogg",
		"ino2.ogg",
		"itemget.wav",
		"itemget2.wav",
		"jump.wav",
	}
	soundPlayers = map[string]*audio.Player{}
	mute         = false
)

func Mute() {
	mute = true
}

func Load() error {
	for _, n := range soundFilenames {
		b, err := assets.Asset("resources/sound/" + n)
		if err != nil {
			return err
		}
		f := bytes.NewReader(b)
		var s io.ReadSeeker
		switch {
		case strings.HasSuffix(n, ".ogg"):
			stream, err := vorbis.Decode(audioContext, f)
			if err != nil {
				return err
			}
			s = audio.NewInfiniteLoop(stream, stream.Length())
		case strings.HasSuffix(n, ".wav"):
			stream, err := wav.Decode(audioContext, f)
			if err != nil {
				return err
			}
			s = stream
		default:
			panic("invalid file name")
		}
		p, err := audio.NewPlayer(audioContext, s)
		if err != nil {
			return err
		}
		soundPlayers[n] = p
	}
	return nil
}

func Finalize() error {
	for _, p := range soundPlayers {
		if err := p.Close(); err != nil {
			return err
		}
	}
	return nil
}

type BGM string

const (
	BGM0 BGM = "ino1.ogg"
	BGM1 BGM = "ino2.ogg"
)

func SetBGMVolume(volume float64) {
	if mute {
		return
	}
	for _, b := range []BGM{BGM0, BGM1} {
		p := soundPlayers[string(b)]
		if !p.IsPlaying() {
			continue
		}
		p.SetVolume(volume)
		return
	}
}

func PauseBGM() {
	if mute {
		return
	}
	for _, b := range []BGM{BGM0, BGM1} {
		p := soundPlayers[string(b)]
		p.Pause()
	}
}

func ResumeBGM(bgm BGM) {
	if mute {
		return
	}
	PauseBGM()
	p := soundPlayers[string(bgm)]
	p.SetVolume(1)
	p.Play()
}

func PlayBGM(bgm BGM) error {
	if mute {
		return nil
	}
	PauseBGM()
	p := soundPlayers[string(bgm)]
	p.SetVolume(1)
	if err := p.Rewind(); err != nil {
		return err
	}
	p.Play()
	return nil
}

type SE string

const (
	SE_DAMAGE   SE = "damage.wav"
	SE_HEAL     SE = "heal.wav"
	SE_ITEMGET  SE = "itemget.wav"
	SE_ITEMGET2 SE = "itemget2.wav"
	SE_JUMP     SE = "jump.wav"
)

func PlaySE(se SE) {
	if mute {
		return
	}
	p := soundPlayers[string(se)]
	p.Rewind()
	p.Play()
}
