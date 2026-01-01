package handler

import "cmp"

type Scale struct {
	System      string
	Description string
	Notes       []Note
	Algorithm   computeNotesFn
}

type Note struct {
	DistanceFromTonic Interval
	Position          uint
	Name              string
}

func (n Note) sortWith(m Note) int {
	return cmp.Compare(n.DistanceFromTonic.Numerator*m.DistanceFromTonic.Denominator, m.DistanceFromTonic.Numerator*n.DistanceFromTonic.Denominator)
}

type computeNotesFn func() []Note
