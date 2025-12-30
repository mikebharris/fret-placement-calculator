package handler

import "math"

type EquallyDividedOctave struct {
	NumberOfDivisions uint
}

type Division struct {
	Ratio float64
	Cents float64
}

func (o EquallyDividedOctave) divisions() []Division {
	var divisions []Division
	for i := 1; i <= int(o.NumberOfDivisions); i++ {
		divisions = append(divisions, o.division(i))
	}
	return divisions
}

func (o EquallyDividedOctave) division(i int) Division {
	ratio := math.Exp2(float64(i) / float64(o.NumberOfDivisions))
	return Division{ratio, o.divisionInCents(ratio)}
}

func (o EquallyDividedOctave) divisionInCents(i float64) float64 {
	return math.Round(math.Log10(i) / math.Log10(2) * 1200)
}
