package main

import (
	"flag"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/generators"
	"github.com/faiface/beep/wav"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	var (
		file string
		freq float64
		dur  time.Duration
	)
	flag.StringVar(&file, "file", "", "output file")
	flag.Float64Var(&freq, "frequency", 0, "frequency")
	flag.DurationVar(&dur, "duration", 0, "duration")
	flag.Parse()
	if file == "" {
		log.Fatal("Must set -file")
	}
	if freq == 0 {
		log.Fatal("Must set -frequency")
	}
	if dur == 0 {
		log.Fatal("Must set -duration")
	}
	format := beep.Format{
		SampleRate:  44100,
		NumChannels: 2,
		Precision:   2,
	}
	var all []beep.Streamer
	for _, t := range []struct {
		m    int // multiple of a third of a period
		gain float64
	}{
		{6, 16},
		{3, 10},
		{9, 8},
		{12, 7},
		{15, 5},
		{18, 4},
		{21, .9},
		{27, .7},
		{2, .6},
	} {
		s, err := generators.SinTone(format.SampleRate, int((freq/3)*float64(t.m)))
		if err != nil {
			log.Fatal(err)
		}
		s = beep.Seq(beep.Silence(rand.Intn(format.SampleRate.N(1*time.Second/262))), s)
		all = append(all, &effects.Gain{
			Streamer: s,
			Gain:     -1 + t.gain/100,
		})
	}
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}

	mix := beep.Mix(all...)
	sampled := beep.Take(format.SampleRate.N(dur), mix)
	if err := wav.Encode(f, sampled, format); err != nil {
		log.Fatal(err)
	}
	log.Println("Success.")
}
