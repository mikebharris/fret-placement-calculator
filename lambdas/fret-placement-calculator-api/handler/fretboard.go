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

func newFretboardFromJustScale(length float64, octaves int, scale music.JustScale) Fretboard {
	fretboard := Fretboard{
		System:      scale.System(),
		Description: fmt.Sprintf("Fret positions based on %s", scale.Description()),
		ScaleLength: length,
	}
	fretboard.makeJustFrets(scale.Intervals(), octaves)
	return fretboard
}

func newFretboardFromTemperedScale(length float64, octaves int, scale music.TemperedScale) Fretboard {
	fretboard := Fretboard{
		System:      scale.System(),
		Description: fmt.Sprintf("Fret positions based on %s", scale.Description()),
		ScaleLength: length,
	}
	fretboard.makeTemperedFrets(scale.Intervals(), octaves)
	return fretboard
}

func (f *Fretboard) makeTemperedFrets(intervals []music.TemperedInterval, octaves int) {
	for octave := 0; octave < octaves; octave++ {
		for i, interval := range intervals {
			if octave > 0 && i == 0 {
				continue // skip unison at octave 0
			}
			f.Frets = append(f.Frets, Fret{
				Label:    fmt.Sprintf("%.2f cents", interval.ToCents()+float64(octave)*1200),
				Position: math.Round((f.ScaleLength-(f.ScaleLength/interval.Value())/math.Pow(2, float64(octave)))*100) / 100,
			})
		}
	}
}

func (f *Fretboard) makeJustFrets(intervals []music.JustInterval, octaves int) {
	var previousInterval = music.Unison()
	for octave := 0; octave < octaves; octave++ {
		for _, interval := range intervals {
			if octave > 0 && interval == music.Unison() {
				previousInterval = music.Unison()
				continue // skip unison at octave 0
			}
			f.Frets = append(f.Frets, Fret{
				Label:    interval.String(),
				Position: math.Round((f.ScaleLength-(f.ScaleLength/float64(interval.Numerator()))*float64(interval.Denominator())/math.Pow(2, float64(octave)))*100) / 100,
				Comment:  interval.Name(),
				Interval: interval.Subtract(previousInterval).String(),
			})
			previousInterval = interval
		}
	}
}
