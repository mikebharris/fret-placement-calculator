package handler

import (
	"fmt"
	"math"

	"github.com/mikebharris/music"
)

type Fret struct {
	Label    string  `json:"label"`
	Position float64 `json:"position"`
	Comment  string  `json:"comment,omitempty"`
	Interval string  `json:"interval,omitempty"`
}

type Fretboard struct {
	System      string  `json:"system"`
	Description string  `json:"description,omitempty"`
	ScaleLength float64 `json:"scaleLength"`
	Frets       []Fret  `json:"frets"`
}

func newFretboardFromJustScale(length float64, s music.JustScale) Fretboard {
	fretboard := Fretboard{
		System:      s.System(),
		Description: fmt.Sprintf("Fret positions based on %s", s.Description()),
		ScaleLength: length,
	}
	fretboard.makeJustFrets(s.Intervals())
	return fretboard
}

func newFretboardFromTemperedScale(length float64, s music.TemperedScale) Fretboard {
	fretboard := Fretboard{
		System:      s.System(),
		Description: fmt.Sprintf("Fret positions based on %s", s.Description()),
		ScaleLength: length,
	}
	fretboard.makeTemperedFrets(s.Intervals())
	return fretboard
}

func (f *Fretboard) makeTemperedFrets(intervals []music.TemperedInterval) {
	for _, interval := range intervals {
		f.Frets = append(f.Frets, Fret{
			Label:    interval.String(),
			Position: math.Round(f.ScaleLength-(f.ScaleLength/interval.Value())*100) / 100,
		})
	}
}

func (f *Fretboard) makeJustFrets(intervals []music.JustInterval) {
	var previousInterval = music.Unison
	for _, interval := range intervals {
		f.Frets = append(f.Frets, Fret{
			Label:    interval.String(),
			Position: math.Round((f.ScaleLength-(f.ScaleLength/float64(interval.Numerator()))*float64(interval.Denominator()))*100) / 100,
			Comment:  interval.Name(),
			Interval: interval.Subtract(previousInterval).String(),
		})
		previousInterval = interval
	}
}
