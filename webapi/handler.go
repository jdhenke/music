package webapi

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
	"github.com/jdhenke/music/note"
	"github.com/orcaman/writerseeker"
)

func NoteHandler() http.Handler {
	return http.HandlerFunc(handleNote)
}

func handleNote(rw http.ResponseWriter, r *http.Request) {
	noteReq, err := parseNoteRequest(r.URL.Query())
	if err != nil {
		log.Printf("Error parsing note request: %v", err)
		http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	ws := &writerseeker.WriterSeeker{}
	format := beep.Format{
		SampleRate:  44100,
		NumChannels: 2,
		Precision:   2,
	}
	source := sources[noteReq.source]
	streamer := source(format, noteReq.freq, 0.05, noteReq.dur)
	if err := wav.Encode(ws, streamer, format); err != nil {
		log.Println(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	reader := ws.BytesReader()
	rw.Header().Add("Content-Length", fmt.Sprint(reader.Len()))
	if _, err := io.Copy(rw, reader); err != nil {
		log.Println(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

var sources = map[string]func(format beep.Format, freq, amp float64, dur time.Duration) beep.Streamer{
	"square": note.Square,
	"sine":   note.Sine,
	"piano":  note.Piano,
}

func parseNoteRequest(q url.Values) (noteRequest, error) {
	sourceStr := q.Get("source")
	freqStr := q.Get("freq")
	durStr := q.Get("dur")
	if sourceStr == "" {
		return noteRequest{}, fmt.Errorf("no source provided")
	}
	if _, ok := sources[sourceStr]; !ok {
		return noteRequest{}, fmt.Errorf("unknown source: %v", sourceStr)
	}
	if freqStr == "" {
		return noteRequest{}, fmt.Errorf("no freq provided")
	}
	freq, err := strconv.ParseFloat(freqStr, 64)
	if err != nil {
		return noteRequest{}, fmt.Errorf("bad frequency: %v: %v", freqStr, err)
	}
	if durStr == "" {
		return noteRequest{}, fmt.Errorf("no dur provided")
	}
	dur, err := time.ParseDuration(durStr)
	if err != nil {
		return noteRequest{}, fmt.Errorf("bad duration: %v: %v", durStr, err)
	}
	return noteRequest{
		source: sourceStr,
		freq:   freq,
		dur:    dur,
	}, nil
}

type noteRequest struct {
	source string
	freq   float64
	dur    time.Duration
}
