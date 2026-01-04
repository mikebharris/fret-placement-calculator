package music

import (
	"fmt"
	"slices"
)

type JustScale struct {
	system      string
	description string
	algorithm   computeJustIntervalsFn
}

type Symmetry string

const (
	Asymmetric Symmetry = "asymmetric"
	Symmetric1 Symmetry = "symmetric1"
	Symmetric2 Symmetry = "symmetric2"
)

func (s Symmetry) String() string {
	return string(s)
}

func NewPythagoreanScale() JustScale {
	return JustScale{
		system:      "Pythagorean",
		description: "3-limit Pythagorean ratios.",
		algorithm:   computePythagoreanIntervals,
	}
}

func New5LimitPythagoreanScale() JustScale {
	return JustScale{
		system:      "5-limit Pythagorean",
		description: "5-limit just intonation pure ratios chromatic scale derived from applying syntonic comma to Pythagorean ratios.",
		algorithm:   compute5LimitPythagoreanIntervals,
	}
}

func New5LimitJustIntonationChromaticScale(symmetry Symmetry) JustScale {
	return JustScale{
		system:      "5-limit Just Intonation",
		description: "5-limit just intonation pure ratios derived from third- and fifth-partial ratios.",
		algorithm: func() []JustInterval {
			return computeJustScale(buildMultiplierTablesFrom(multipliers(3), multipliers(5), multipliers(9)), fiveLimitScaleFilter(symmetry))
		},
	}
}

func New7LimitJustIntonationChromaticScale() JustScale {
	return JustScale{
		system:      "7-limit Just Intonation",
		description: "7-limit just intonation pure ratios derived from third-, fifth- and seventh-partial ratios.",
		algorithm: func() []JustInterval {
			return computeJustScale(buildMultiplierTablesFrom(multipliers(3), multipliers(5), multipliers(9), multipliers(7)), nullScaleFilter())
		},
	}
}

func New13LimitJustIntonationChromaticScale() JustScale {
	return JustScale{
		system:      "13-limit Just Intonation",
		description: "Just Intonation chromatic scale based on 13-limit pure ratios.",
		algorithm: func() []JustInterval {
			return computeJustScale(buildMultiplierTablesFrom(multipliers(3), multipliers(5), multipliers(9), multipliers(7), multipliers(11), multipliers(13)), nullScaleFilter())
		},
	}
}

func NewIntenseDiatonicScale(mode MusicalMode) JustScale {
	return JustScale{
		system:      "Ptolemy Intense Diatonic",
		description: fmt.Sprintf("Ptolemy's 5-limit intense diatonic scale in %s mode.", mode),
		algorithm: func() []JustInterval {
			return computePtolemeicIntenseDiatonicScale(mode)
		},
	}
}

func NewSazScale() JustScale {
	// as per https://en.wikipedia.org/wiki/Ba%C4%9Flama and the cura that I have
	return JustScale{
		system:      "Saz",
		description: "Turkish Saz tuning ratios.",
		algorithm:   computeSazScale,
	}
}

func (s JustScale) System() string {
	return s.system
}

func (s JustScale) Description() string {
	return s.description
}

func (s JustScale) Intervals() []JustInterval {
	return s.algorithm()
}

type computeJustIntervalsFn func() []JustInterval

func computePtolemeicIntenseDiatonicScale(mode MusicalMode) []JustInterval {
	var intervalMap = map[MusicalMode][]JustInterval{
		Lydian:     {greaterMajorSecond, lesserMajorSecond, greaterMajorSecond, diatonicSemitone, lesserMajorSecond, greaterMajorSecond, diatonicSemitone},
		Ionian:     {greaterMajorSecond, lesserMajorSecond, diatonicSemitone, greaterMajorSecond, lesserMajorSecond, greaterMajorSecond, diatonicSemitone},
		Mixolydian: {greaterMajorSecond, lesserMajorSecond, diatonicSemitone, greaterMajorSecond, lesserMajorSecond, diatonicSemitone, greaterMajorSecond},
		Dorian:     {greaterMajorSecond, diatonicSemitone, lesserMajorSecond, greaterMajorSecond, lesserMajorSecond, diatonicSemitone, greaterMajorSecond},
		Aeolian:    {greaterMajorSecond, diatonicSemitone, lesserMajorSecond, greaterMajorSecond, diatonicSemitone, greaterMajorSecond, lesserMajorSecond},
		Phrygian:   {diatonicSemitone, greaterMajorSecond, lesserMajorSecond, greaterMajorSecond, diatonicSemitone, greaterMajorSecond, lesserMajorSecond},
		Locrian:    {diatonicSemitone, greaterMajorSecond, lesserMajorSecond, diatonicSemitone, greaterMajorSecond, lesserMajorSecond, greaterMajorSecond},
	}

	var intervals = []JustInterval{Unison}
	var interval = Unison

	for _, v := range intervalMap[mode] {
		interval = JustInterval{numerator: interval.numerator * v.numerator, denominator: interval.denominator * v.denominator}.simplify()
		intervals = append(intervals, interval)
	}

	return intervals
}

func computePythagoreanIntervals() []JustInterval {
	var fifthsFromTonicToCompute = 6
	var intervals []JustInterval
	for i := -fifthsFromTonicToCompute; i <= fifthsFromTonicToCompute; i++ {
		if i < 0 {
			intervals = append(intervals, PerfectFifth.ToPowerOf(i).reciprocal().octaveReduce())
		} else {
			intervals = append(intervals, PerfectFifth.ToPowerOf(i).octaveReduce())
		}
	}

	intervals = append(intervals, octave)
	slices.SortFunc(intervals, func(i, j JustInterval) int {
		return i.sortWith(j)
	})
	return intervals
}

func compute5LimitPythagoreanIntervals() []JustInterval {
	var intervals []JustInterval
	for _, interval := range computePythagoreanIntervals() {
		if interval.isPerfect() {
			intervals = append(intervals, interval)
			continue
		}

		graveRatio := interval.add(acuteUnison)
		acuteRatio := interval.add(graveUnison)

		if graveRatio.denominator < acuteRatio.denominator {
			intervals = append(intervals, graveRatio)
		} else {
			intervals = append(intervals, acuteRatio)
		}
	}
	return intervals
}

func computeJustScale(multipliers [][]uint, filter intervalFilterFunction) []JustInterval {
	poolOfPotentialIntervals := justIntervalsFromMultipliers(multipliers, filter)
	var preferredIntervals = []JustInterval{Unison}
	centsInOctave := 1200.0
	for r := 50.0; r <= centsInOctave; r += 100 {
		var intervalsInNoteRange []JustInterval
		for _, interval := range poolOfPotentialIntervals {
			cents := interval.toCents()
			if cents >= r && cents < r+100 {
				intervalsInNoteRange = append(intervalsInNoteRange, interval)
			}
		}

		//   chosen interval is the simplest integer ratio
		var chosenInterval JustInterval
		for i, interval := range intervalsInNoteRange {
			if i == 0 || (interval.numerator < chosenInterval.numerator && interval.denominator < chosenInterval.denominator) {
				chosenInterval = interval
				continue
			}
		}
		preferredIntervals = append(preferredIntervals, chosenInterval)
	}

	return preferredIntervals
}

func fiveLimitScaleFilter(symmetry Symmetry) func(interval JustInterval) bool {
	return func(interval JustInterval) bool {
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

func nullScaleFilter() func(interval JustInterval) bool {
	return func(interval JustInterval) bool {
		return false
	}
}

func buildMultiplierTablesFrom(multipliers ...[][]uint) [][]uint {
	if len(multipliers) == 1 {
		return multipliers[0]
	}
	return createMultiplierTableOf(multipliers[0], buildMultiplierTablesFrom(multipliers[1:]...))
}

func computeSazScale() []JustInterval {
	return intervalsFromIntegers([][]uint{{1, 1}, {18, 17}, {12, 11}, {9, 8}, {81, 68}, {27, 22}, {81, 64}, {4, 3}, {24, 17}, {16, 11}, {3, 2}, {27, 17}, {18, 11}, {27, 16}, {16, 9}, {32, 17}, {64, 33}, {2, 1}})
}
