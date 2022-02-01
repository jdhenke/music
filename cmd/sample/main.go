package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/cmplx"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
	"github.com/mjibson/go-dsp/fft"
)

// Sample -in starting at -start and write to -out for -duration by manually rebuilding the sine waves from the FFT.
func main() {
	var (
		in, out         string
		start, duration time.Duration
	)
	flag.StringVar(&in, "in", "", "input file")
	flag.StringVar(&out, "out", "", "output file")
	flag.DurationVar(&start, "start", 0, "starting point to sample in input file")
	flag.DurationVar(&duration, "duration", 0, "how long the output file should be")
	flag.Parse()
	if err := sample(in, out, start, duration); err != nil {
		log.Fatalf("Failed: %v", err)
	}
	log.Println("Success.")
}

func sample(in string, out string, start time.Duration, duration time.Duration) error {
	// get samples
	log.Println("Getting samples...")
	inFile, err := os.Open(in)
	if err != nil {
		return fmt.Errorf("open input file: %v", err)
	}
	defer func() {
		_ = inFile.Close()
	}()
	raw, format, err := wav.Decode(inFile)
	if err != nil {
		return fmt.Errorf("decoding input wav file: %v", err)
	}
	if err := raw.Seek(format.SampleRate.N(start)); err != nil {
		return fmt.Errorf("seeking to start position in input file: %v", err)
	}
	const numPeriods = 100
	const frequency = 262
	limited := beep.Take(format.SampleRate.N(time.Second*numPeriods/frequency), raw)
	p := make([][2]float64, 1024)
	var samples []float64
	for {
		n, ok := limited.Stream(p)
		if !ok {
			break
		}
		for i := 0; i < n; i++ {
			samples = append(samples, p[i][0])
		}
	}
	if err := limited.Err(); err != nil {
		return fmt.Errorf("streaming samples: %v", err)
	}

	// perform fft on samples to get frequencies
	log.Println("Performing FFT...")
	freqs := fft.FFTReal(samples)

	// manually construct output streamers from frequencies
	log.Println("Creating output streamer...")
	mix := &beep.Mixer{}
	for freq, x := range freqs {
		r, θ := cmplx.Polar(x)
		f := float64(freq) * frequency / numPeriods
		s := Cosine(format.SampleRate, f, r, len(samples), θ)
		mix.Add(s)
	}

	// write output file
	log.Println("Writing output file...")
	outFile, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("creating output file: %v", err)
	}
	defer func() {
		_ = outFile.Close()
	}()
	limitedOut := beep.Take(format.SampleRate.N(duration), mix)
	if err := wav.Encode(outFile, limitedOut, format); err != nil {
		return fmt.Errorf("encoding output wav file: %v", err)
	}
	if err := outFile.Close(); err != nil {
		return fmt.Errorf("closing output file: %v", err)
	}
	return nil
}

func Cosine(rate beep.SampleRate, freq, r float64, N int, θ float64) beep.Streamer {
	var i int
	return beep.StreamerFunc(func(samples [][2]float64) (n int, ok bool) {
		for start := i; i-start < len(samples); i++ {
			val := (r / float64(N)) * math.Cos(θ+2*math.Pi*float64(i)*freq/float64(rate.N(1*time.Second)))
			samples[i-start][0] = val
			samples[i-start][1] = val
		}
		return len(samples), true
	})
}
