package handler

import (
	"fmt"
	"math"
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

func newFretboardFromScale(length float64, s Scale) Fretboard {
	fretboard := Fretboard{
		System:      s.System,
		Description: fmt.Sprintf("Fret positions based on %s", s.Description),
		ScaleLength: length,
	}
	fretboard.makeFrets(s.Algorithm())
	return fretboard
}

func (f *Fretboard) makeFrets(notes []Note) {
	var previousInterval = unison
	for _, note := range notes {
		f.Frets = append(f.Frets, Fret{
			Label:    note.DistanceFromTonic.String(),
			Position: math.Round((f.ScaleLength-(f.ScaleLength/float64(note.DistanceFromTonic.Numerator))*float64(note.DistanceFromTonic.Denominator))*100) / 100,
			Comment:  note.DistanceFromTonic.name(),
			Interval: note.DistanceFromTonic.subtract(previousInterval).String(),
		})
		previousInterval = note.DistanceFromTonic
	}
}
