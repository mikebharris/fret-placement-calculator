package music

import (
	"math"
	"slices"
)

func computePtolemeicIntenseDiatonicScale(mode MusicalMode) []Interval {
	var intervalMap = map[MusicalMode][]Interval{
		Lydian:     {greaterMajorSecond, lesserMajorSecond, greaterMajorSecond, diatonicSemitone, lesserMajorSecond, greaterMajorSecond, diatonicSemitone},
		Ionian:     {greaterMajorSecond, lesserMajorSecond, diatonicSemitone, greaterMajorSecond, lesserMajorSecond, greaterMajorSecond, diatonicSemitone},
		Mixolydian: {greaterMajorSecond, lesserMajorSecond, diatonicSemitone, greaterMajorSecond, lesserMajorSecond, diatonicSemitone, greaterMajorSecond},
		Dorian:     {greaterMajorSecond, diatonicSemitone, lesserMajorSecond, greaterMajorSecond, lesserMajorSecond, diatonicSemitone, greaterMajorSecond},
		Aeolian:    {greaterMajorSecond, diatonicSemitone, lesserMajorSecond, greaterMajorSecond, diatonicSemitone, greaterMajorSecond, lesserMajorSecond},
		Phrygian:   {diatonicSemitone, greaterMajorSecond, lesserMajorSecond, greaterMajorSecond, diatonicSemitone, lesserMajorSecond, greaterMajorSecond},
		Locrian:    {diatonicSemitone, greaterMajorSecond, lesserMajorSecond, diatonicSemitone, greaterMajorSecond, lesserMajorSecond, greaterMajorSecond},
	}

	var intervals []Interval
	var interval = Unison

	for _, v := range intervalMap[mode] {
		interval = Interval{Numerator: interval.Numerator * v.Numerator, Denominator: interval.Denominator * v.Denominator}.simplify()
		intervals = append(intervals, interval)
	}

	return intervals
}

func computePythagoreanIntervals() []Interval {
	var fifthsFromTonicToCompute = 6
	var intervals []Interval
	for i := -fifthsFromTonicToCompute; i <= fifthsFromTonicToCompute; i++ {
		if i < 0 {
			intervals = append(intervals, PerfectFifth.ToPowerOf(i).reciprocal().octaveReduce())
		} else {
			intervals = append(intervals, PerfectFifth.ToPowerOf(i).octaveReduce())
		}
	}

	intervals = append(intervals, octave)
	slices.SortFunc(intervals, func(i, j Interval) int {
		return i.sortWith(j)
	})
	return intervals
}

func compute5LimitPythagoreanIntervals() []Interval {
	var intervals []Interval
	for _, interval := range computePythagoreanIntervals() {
		if interval.isPerfect() {
			intervals = append(intervals, interval)
			continue
		}

		graveRatio := interval.add(acuteUnison)
		acuteRatio := interval.add(graveUnison)

		if graveRatio.Denominator < acuteRatio.Denominator {
			intervals = append(intervals, graveRatio)
		} else {
			intervals = append(intervals, acuteRatio)
		}
	}
	return intervals
}

func computeJustScale(multipliers [][]uint, filter intervalFilterFunction) []Interval {
	poolOfPotentialIntervals := justIntervalsFromMultipliers(multipliers, filter)

	var preferredIntervals []Interval
	centsInOctave := 1200.0
	for r := 50.0; r <= centsInOctave; r += 100 {
		var intervalsInNoteRange []Interval
		for _, interval := range poolOfPotentialIntervals {
			cents := interval.toCents()
			if cents >= r && cents < r+100 {
				intervalsInNoteRange = append(intervalsInNoteRange, interval)
			}
		}

		//   chosen interval is the simplest integer ratio
		var chosenInterval Interval
		for i, interval := range intervalsInNoteRange {
			if i == 0 || (interval.Numerator < chosenInterval.Numerator && interval.Denominator < chosenInterval.Denominator) {
				chosenInterval = interval
				continue
			}
		}
		preferredIntervals = append(preferredIntervals, chosenInterval)
	}

	return preferredIntervals
}

func computeQuarterCommaMeantoneScale() []float64 {
	fractionOfSyntonicCommaToTemperFifthsBy := 0.25
	temperedFifth := PerfectFifth.ToFloat() * math.Pow(SyntonicComma.ToFloat(), -fractionOfSyntonicCommaToTemperFifthsBy)

	var fifthsFromTonic int

	var ratiosOfNotesToFundamental []float64
	for i := -fifthsFromTonic; i <= fifthsFromTonic; i++ {
		ratiosOfNotesToFundamental = append(ratiosOfNotesToFundamental, octaveReduceFloat(math.Pow(temperedFifth, float64(i))))
	}

	slices.Sort(ratiosOfNotesToFundamental)

	return ratiosOfNotesToFundamental
}

func octaveReduceFloat(ratio float64) float64 {
	for ratio >= 2.0 || ratio < 1.0 {
		if ratio >= 2.0 {
			ratio /= 2.0
		}
		if ratio < 1.0 {
			ratio *= 2.0
		}
	}
	return ratio
}

func computeBachScale() []Interval {
	// Narrowing of the fifths as outlined by Lehman
	temperingFractions := []float64{
		0.0,         // Pure fifth
		-1.0 / 12.0, // Twelfth-comma narrowed
		-1.0 / 6.0,  // Sixth-comma narrowed
		0.0,         // Pure fifth
		-1.0 / 6.0,  // Sixth-comma narrowed
		-1.0 / 12.0, // Twelfth-comma narrowed
		0.0,         // Pure fifth
		-1.0 / 6.0,  // Sixth-comma narrowed
		-1.0 / 12.0, // Twelfth-comma narrowed
		0.0,         // Pure fifth
		-1.0 / 6.0,  // Sixth-comma narrowed
		-1.0 / 12.0, // Twelfth-comma narrowed
	}

	// Calculate tempered fifths
	temperedFifths := make([]float64, 12)
	for i := 0; i < 12; i++ {
		temperedFifths[i] = 3.0 / 2.0 * math.Pow(SyntonicComma.ToFloat(), temperingFractions[i])
	}

	// Derive ratios by walking the circle of fifths
	ratios := make([]float64, 12)
	ratios[0] = 1.0 // Start with the tonic
	for i := 1; i < 12; i++ {
		ratios[i] = ratios[i-1] * temperedFifths[(i-1)%12]
	}

	// Reduce ratios to within the octave [1.0, 2.0)
	for i := range ratios {
		ratios[i] = octaveReduceFloat(ratios[i])
	}

	slices.Sort(ratios) // Sort the ratios in ascending order
	var intervals []Interval
	for _, ratio := range ratios {
		intervals = append(intervals, Interval{Numerator: uint(ratio), Denominator: 1}.octaveReduce())
	}
	return intervals
}

func fiveLimitScaleFilter(symmetry Symmetry) func(interval Interval) bool {
	return func(interval Interval) bool {
		if symmetry == Asymmetric && (interval.IsLesserMajorSecond() || interval.IsLesserMinorSeventh()) {
			return true
		}
		if symmetry == Symmetric1 && (interval.IsLesserMajorSecond() || interval.IsGreaterMinorSeventh()) {
			return true
		}
		if symmetry == Symmetric2 && (interval.IsGreaterMajorSecond() || interval.IsLesserMinorSeventh()) {
			return true
		}
		return false
	}
}

func buildMultiplierTablesFrom(multipliers ...[][]uint) [][]uint {
	if len(multipliers) == 1 {
		return multipliers[0]
	}
	return createMultiplierTableOf(multipliers[0], buildMultiplierTablesFrom(multipliers[1:]...))
}

func computeSazScale() []Interval {
	return intervalsFromIntegers([][]uint{{1, 1}, {18, 17}, {12, 11}, {9, 8}, {81, 68}, {27, 22}, {81, 64}, {4, 3}, {24, 17}, {16, 11}, {3, 2}, {27, 17}, {18, 11}, {27, 16}, {16, 9}, {32, 17}, {64, 33}, {2, 1}})
}
