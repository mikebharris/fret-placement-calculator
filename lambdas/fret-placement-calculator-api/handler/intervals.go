package handler

import (
	"cmp"
	"fmt"
	"math"
	"slices"
)

type Interval struct {
	Numerator   uint   `json:"numerator"`
	Denominator uint   `json:"denominator"`
	Name        string `json:"name,omitempty"`
}

func newInterval(numerator, denominator uint) Interval {
	return Interval{Numerator: numerator, Denominator: denominator}.simplify()
}

func (i Interval) isUnison() bool {
	return i.Numerator == 1 && i.Denominator == 1
}

func (i Interval) isEqualTo(other Interval) bool {
	return i.Numerator == other.Numerator && i.Denominator == other.Denominator
}

func (i Interval) isDiminishedFifth() bool {
	return i.Numerator == 64 && i.Denominator == 45
}

func (i Interval) isLesserMajorSecond() bool {
	return i.Numerator == 10 && i.Denominator == 9
}

func (i Interval) isGreaterMajorSecond() bool {
	return i.Numerator == 9 && i.Denominator == 8
}

func (i Interval) isLesserMinorSeventh() bool {
	return i.Numerator == 16 && i.Denominator == 9
}

func (i Interval) isGreaterMinorSeventh() bool {
	return i.Numerator == 9 && i.Denominator == 5
}

func (i Interval) add(other Interval) Interval {
	return Interval{
		Numerator:   i.Numerator * other.Numerator,
		Denominator: i.Denominator * other.Denominator,
	}.simplify()
}

func (i Interval) isPerfectFourth() bool {
	return i.Numerator == 4 && i.Denominator == 3
}

func (i Interval) isPerfect() bool {
	return i.isUnison() || i.isPerfectFourth() || i.isPerfectFifth() || i.isOctave()
}

func (i Interval) isPerfectFifth() bool {
	return i.Numerator == 3 && i.Denominator == 2
}

func (i Interval) isOctave() bool {
	return i.Numerator == 2 && i.Denominator == 1
}

func (i Interval) simplify() Interval {
	if i.Denominator == 0 {
		return i
	}
	if i.Numerator == 0 {
		return Interval{Numerator: 0, Denominator: 1}
	}
	gcd := func(a, b uint) uint {
		for b != 0 {
			a, b = b, a%b
		}
		return a
	}(i.Numerator, i.Denominator)
	i.Numerator = i.Numerator / gcd
	i.Denominator = i.Denominator / gcd
	return i
}

func (i Interval) octaveReduce() Interval {
	for i.Numerator/i.Denominator >= 2.0 || i.Numerator/i.Denominator < 1.0 {
		if i.Numerator/i.Denominator < 1.0 {
			i.Numerator *= 2
		}
		if i.Numerator/i.Denominator >= 2.0 {
			i.Denominator *= 2
		}
	}
	return i
}

func (i Interval) lessThan(other Interval) bool {
	return i.Numerator*other.Denominator < other.Numerator*i.Denominator
}

func (i Interval) greaterThan(other Interval) bool {
	return !i.lessThan(other) && !i.isEqualTo(other)
}

func (i Interval) subtract(other Interval) Interval {
	if i.lessThan(other) {
		return Interval{Numerator: i.Denominator * other.Numerator, Denominator: i.Numerator * other.Denominator}.simplify()
	} else if i.greaterThan(other) {
		return Interval{Numerator: i.Numerator * other.Denominator, Denominator: i.Denominator * other.Numerator}.simplify()
	}
	return i
}

func (i Interval) name() string {
	for _, n := range intervalNames {
		if n.Numerator == i.Numerator && n.Denominator == i.Denominator {
			return n.Name
		}
	}
	return ""
}

func (i Interval) toFloat() float64 {
	return float64(i.Numerator) / float64(i.Denominator)
}

func (i Interval) toPowerOf(p int) Interval {
	return Interval{uint(math.Pow(float64(perfectFifth.Numerator), math.Abs(float64(p)))), uint(math.Pow(float64(perfectFifth.Denominator), math.Abs(float64(p)))), ""}
}

func (i Interval) reciprocal() Interval {
	return Interval{Denominator: i.Numerator, Numerator: i.Denominator}
}

var unison = Interval{1, 1, "Unison"}
var acuteUnison = Interval{Numerator: 81, Denominator: 80}
var syntonicComma = Interval{Numerator: 81, Denominator: 80}
var dieses = Interval{Numerator: 128, Denominator: 125}
var justChromaticSemitone = Interval{Numerator: 25, Denominator: 24}
var graveUnison = Interval{Numerator: 80, Denominator: 81}
var lesserMajorSecond = Interval{10, 9, "Lesser Major Second"}
var greaterMajorSecond = Interval{9, 8, "Greater Major Second"}
var diatonicSemitone = Interval{16, 15, "Diatonic Semitone"}
var perfectFifth = Interval{3, 2, "Perfect Fifth"}
var octave = Interval{2, 1, "Octave"}

var intervalNames = []Interval{
	{1, 1, "Perfect Unison"},
	{225, 224, "Septimal Kleisma"},
	{81, 80, "Grave Unison"},
	{128, 125, "Dieses (Diminished Second)"},
	{25, 24, "Just (Lesser) Chromatic Semitone"},
	{256, 243, "Pythagorean Minor Second"},
	{135, 128, "Greater Chromatic Semitone"},
	{27, 25, "Acute Minor Second"},
	{16, 15, "Minor Second"},
	{15, 14, "Septimal Minor Second"},
	{10, 9, "Just (Lesser) Major Second"},
	{9, 8, "Pythagorean (Greater) Major Second"},
	{8, 7, "Septimal Major Second"},
	{6, 5, "Minor Third"},
	{5, 4, "Major Third"},
	{32, 27, "Diminished Fourth"},
	{81, 64, "Pythagorean Major Third"},
	{4, 3, "Perfect Fourth"},
	{45, 32, "Augmented Fourth"},
	{7, 5, "Septimal Augmented Fourth"},
	{1024, 729, "Pythagorean Diminished Fifth"},
	{729, 512, "Pythagorean Augmented Fourth"},
	{64, 45, "Diminished Fifth"},
	{10, 7, "Septimal Diminished Fifth"},
	{40, 27, "Grave Fifth"},
	{3, 2, "Perfect Fifth"},
	{8, 5, "Just Minor Sixth"},
	{128, 81, "Pythagorean Minor Sixth"},
	{5, 3, "Major Sixth"},
	{27, 16, "Pythagorean Major Sixth"},
	{16, 9, "Pythagorean (Lesser) Minor Seventh"},
	{9, 5, "Just (Greater) Minor Seventh"},
	{7, 4, "Septimal (Harmonic) Minor Seventh"},
	{15, 8, "Just Major Seventh"},
	{243, 128, "Pythagorean Major Seventh"},
	{2, 1, "Perfect Octave"},
}

func (i Interval) sortWith(j Interval) int {
	return cmp.Compare(float64(i.Numerator)/float64(i.Denominator), float64(j.Numerator)/float64(j.Denominator))
}

func (i Interval) fretPosition(scaleLength float64) float64 {
	return math.Round((scaleLength-(scaleLength/float64(i.Numerator))*float64(i.Denominator))*100) / 100
}

func (i Interval) String() string {
	return fmt.Sprintf("%d:%d", i.Numerator, i.Denominator)
}

// as per https://en.wikipedia.org/wiki/Ba%C4%9Flama and the cura that I have
var sazIntervals = intervalsFromIntegers([][]uint{{18, 17}, {12, 11}, {9, 8}, {81, 68}, {27, 22}, {81, 64}, {4, 3}, {24, 17}, {16, 11}, {3, 2}, {27, 17}, {18, 11}, {27, 16}, {16, 9}, {32, 17}, {64, 33}, {2, 1}})

func intervalsFromIntegers(integers [][]uint) []Interval {
	var intervals []Interval
	for _, pair := range integers {
		intervals = append(intervals, fromIntArray(pair))
	}
	return intervals
}

func fromIntArray(i []uint) Interval {
	return Interval{Numerator: i[0], Denominator: i[1]}
}

func sortIntervals(intervals []Interval) {
	slices.SortFunc(intervals, func(i, j Interval) int {
		return i.sortWith(j)
	})
}
