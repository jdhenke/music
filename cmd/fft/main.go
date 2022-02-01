package main

import (
	"bufio"
	"fmt"
	"log"
	"math/cmplx"
	"os"
	"strconv"

	"github.com/mjibson/go-dsp/fft"
)

func main() {
	var nums []float64
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		x, err := strconv.ParseFloat(s.Text(), 64)
		if err != nil {
			log.Fatalf("Bad number '%s': %v", s.Text(), err)
		}
		nums = append(nums, x)
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	freqs := fft.FFTReal(nums)
	for _, c := range freqs {
		r, _ := cmplx.Polar(c)
		fmt.Println(r)
	}
	log.Println("Success.")
}
