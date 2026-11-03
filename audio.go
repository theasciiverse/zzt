package main

// The PC speaker, in software: ZZT's sound code turns a square-wave tone on
// and off (Sound, NoSound) from the 18.2 Hz timer interrupt, which here is a
// goroutine. The tone is synthesised and played through oto; with no audio
// device the game runs silent.

import (
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ebitengine/oto/v3"
)

const speakerRate = 44100

// speakerFreq is the tone now sounding in Hz, 0 for silence.
var speakerFreq atomic.Uint32

// soundMu guards the sound queue between the game and the timer.
var soundMu sync.Mutex

var speakerPlayer *oto.Player

type speakerWave struct {
	phase float64
	level float64 // eased toward on/off, so tones start and stop without clicks
}

func (w *speakerWave) Read(p []byte) (int, error) {
	n := len(p) / 4 * 4
	for i := 0; i < n; i += 4 {
		f := float64(speakerFreq.Load())
		target := 0.0
		if f > 0 {
			target = 1
			w.phase += f / speakerRate
			if w.phase >= 1 {
				w.phase -= math.Floor(w.phase)
			}
		}
		w.level += (target - w.level) * 0.02
		v := -1.0
		if w.phase < 0.5 {
			v = 1
		}
		s := int16(v * w.level * 0.12 * 32767)
		p[i], p[i+1], p[i+2], p[i+3] = byte(s), byte(s>>8), byte(s), byte(s>>8)
	}
	return n, nil
}

// AudioInstall opens the sound device and starts the timer. Call it before
// the screen is set up, so any message from the sound system lands outside it.
func AudioInstall(muted bool) {
	if !muted {
		ctx, ready, err := oto.NewContext(&oto.NewContextOptions{
			SampleRate:   speakerRate,
			ChannelCount: 2,
			Format:       oto.FormatSignedInt16LE,
			BufferSize:   30 * time.Millisecond,
		})
		if err == nil {
			<-ready
			speakerPlayer = ctx.NewPlayer(&speakerWave{})
			speakerPlayer.Play()
		} else {
			logf("no sound: %v", err)
		}
	}
	go func() {
		// the PC timer: 1193182 / 65536 Hz
		t := time.NewTicker(65536 * time.Second / 1193182)
		for range t.C {
			soundMu.Lock()
			SoundTimerHandler()
			soundMu.Unlock()
		}
	}()
}
