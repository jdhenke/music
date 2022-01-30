package note

import (
	"embed"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
)

//go:embed c4.wav
var fs embed.FS

func Piano(format beep.Format, frequency, amplitude float64, dur time.Duration) beep.Streamer {
	f, err := fs.Open("c4.wav")
	if err != nil {
		panic(err)
	}
	s, nativeFormat, err := wav.Decode(f)
	if err != nil {
		panic(err)
	}

	a4 := beep.ResampleRatio(4, frequency/261.63, s)
	return beep.Take(nativeFormat.SampleRate.N(dur), a4)
}
