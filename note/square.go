package note

import (
	"time"

	"github.com/faiface/beep"
)

// Square returns a beep.Streamer that streams the sound of a square wave with the given frequency and amplitude for the
// given duration.
func Square(format beep.Format, frequency, amplitude float64, dur time.Duration) beep.Streamer {
	s := 0
	end := format.SampleRate.N(dur)
	// Get samples / period = (samples / second) * (seconds / period), where (seconds / period) is the inverse of
	// (periods / second) = frequency, so this becomes (samples / second) / frequency.
	samplesPerSecond := format.SampleRate.N(time.Second)
	samplesPerPeriod := samplesPerSecond / int(frequency)
	halfPeriod := samplesPerPeriod / 2
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		if s == end {
			return 0, false
		}
		start := s
		next := min(s+len(samples), end)
		for ; s < next; s++ {
			val := amplitude
			if s%samplesPerPeriod < halfPeriod {
				val = -amplitude
			}
			samples[s-start][0] = val
			samples[s-start][1] = val
		}
		return next - start, true
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
