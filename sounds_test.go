package main

import "testing"

// A queued tune plays its notes at ZZT's pitches, one per timer tick at
// duration T, and the speaker falls silent at the end.
func TestSoundQueuePlaysNotes(t *testing.T) {
	SoundEnabled = true
	SoundBlockQueueing = false
	SoundClearQueue()
	SoundQueue(1, SoundParse("tcegx"))
	var got []uint32
	for i := 0; i < 6; i++ {
		SoundTimerHandler()
		got = append(got, speakerFreq.Load())
	}
	// C, E, G in octave 3, then the rest and silence
	want := []uint32{255, 322, 383, 0, 0, 0} // Exp(3 ln 2) * 32 truncates to 255
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("tick %d: %d Hz, want %d (all %v)", i, got[i], want[i], got)
		}
	}
	if SoundIsPlaying {
		t.Fatal("still playing after the tune")
	}
}

// The synthesised wave is a square at the speaker's frequency.
func TestSpeakerWave(t *testing.T) {
	speakerFreq.Store(441)
	w := &speakerWave{level: 1}
	buf := make([]byte, 4*speakerRate/10)
	w.Read(buf)
	flips := 0
	prev := int16(buf[0]) | int16(buf[1])<<8
	for i := 4; i < len(buf); i += 4 {
		s := int16(buf[i]) | int16(buf[i+1])<<8
		if (s > 0) != (prev > 0) {
			flips++
		}
		prev = s
	}
	speakerFreq.Store(0)
	if flips < 86 || flips > 90 { // 441 Hz for 0.1 s: 88 half-waves
		t.Fatalf("%d sign changes, want about 88", flips)
	}
}
