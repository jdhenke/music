package note

import (
	"math"
	"time"

	"github.com/faiface/beep"
)

// Sine returns a beep.Streamer that streams the sound of a sine wave with the given frequency and amplitude for the
// given duration.
func Sine(format beep.Format, frequency, amplitude float64, dur time.Duration) beep.Streamer {
	s := 0
	end := format.SampleRate.N(dur)
	samplesPerSecond := format.SampleRate.N(time.Second)
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		if s == end {
			return 0, false
		}
		start := s
		next := min(s+len(samples), end)
		for ; s < next; s++ {
			val := amplitude * math.Sin(float64(s)*2*math.Pi*frequency/float64(samplesPerSecond))
			samples[s-start][0] = val
			samples[s-start][1] = val
		}
		return next - start, true
	})
}
