package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
)

func main() {
	var (
		file  string
		start time.Duration = -1
		dur   time.Duration
	)
	flag.StringVar(&file, "file", "", "file")
	flag.DurationVar(&start, "start", 0, "duration to wait until sampling")
	flag.DurationVar(&dur, "duration", 0, "duration to sample")
	flag.Parse()
	if file == "" {
		log.Fatal("Must provide -file.")
	}
	if start == -1 {
		log.Fatal("Must provide -start.")
	}
	if dur == 0 {
		log.Fatal("Must provide -duration.")
	}
	if err := decode(file, start, dur); err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Println("Success.")
}

func decode(file string, start, dur time.Duration) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()
	orig, format, err := wav.Decode(f)
	if err != nil {
		return fmt.Errorf("decoding wave file: %v", err)
	}
	startSample := format.SampleRate.N(start)
	if err := orig.Seek(startSample); err != nil {
		return fmt.Errorf("seeking to start: %v", err)
	}
	s := beep.Take(format.SampleRate.N(dur), orig)
	samples := make([][2]float64, 1024)
	for {
		n, ok := s.Stream(samples)
		if n == 0 && !ok {
			break
		}
		if !ok {
			return fmt.Errorf("scanning: %v", s.Err())
		}
		for i := 0; i < n; i++ {
			fmt.Println(samples[i][0])
		}
	}
	return nil
}
