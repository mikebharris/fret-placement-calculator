package music

import (
	"cmp"
	"fmt"
	"math"
	"slices"
)

type JustInterval struct {
	numerator   uint
	denominator uint
	name        string
}

func newInterval(numerator, denominator uint) JustInterval {
	return JustInterval{numerator: numerator, denominator: denominator}.simplify()
}

func (i JustInterval) Numerator() uint {
	return i.numerator
}

func (i JustInterval) Denominator() uint {
	return i.denominator
}

func (i JustInterval) isUnison() bool {
	return i.numerator == 1 && i.denominator == 1
}

func (i JustInterval) isEqualTo(other JustInterval) bool {
	return i.numerator == other.numerator && i.denominator == other.denominator
}

func (i JustInterval) isDiminishedFifth() bool {
	return i.numerator == 64 && i.denominator == 45
}

func (i JustInterval) IsLesserMajorSecond() bool {
	return i.numerator == 10 && i.denominator == 9
}

func (i JustInterval) IsGreaterMajorSecond() bool {
	return i.numerator == 9 && i.denominator == 8
}

func (i JustInterval) IsLesserMinorSeventh() bool {
	return i.numerator == 16 && i.denominator == 9
}

func (i JustInterval) IsGreaterMinorSeventh() bool {
	return i.numerator == 9 && i.denominator == 5
}

func (i JustInterval) add(other JustInterval) JustInterval {
	interval := JustInterval{
		numerator:   i.numerator * other.numerator,
		denominator: i.denominator * other.denominator,
	}.simplify()

	return interval
}

func (i JustInterval) isPerfectFourth() bool {
	return i.numerator == 4 && i.denominator == 3
}

func (i JustInterval) isPerfect() bool {
	simplestForm := i.octaveReduce().simplify()
	return simplestForm.isUnison() || simplestForm.isPerfectFourth() || simplestForm.isPerfectFifth() || simplestForm.isOctave()
}

func (i JustInterval) isPerfectFifth() bool {
	return i.numerator == 3 && i.denominator == 2
}

func (i JustInterval) isOctave() bool {
	return i.numerator == 2 && i.denominator == 1
}

func (i JustInterval) simplify() JustInterval {
	if i.denominator == 0 {
		return i
	}
	if i.numerator == 0 {
		return JustInterval{numerator: 0, denominator: 1}
	}
	gcd := func(a, b uint) uint {
		for b != 0 {
			a, b = b, a%b
		}
		return a
	}(i.numerator, i.denominator)
	i.numerator = i.numerator / gcd
	i.denominator = i.denominator / gcd
	return i
}

func (i JustInterval) octaveReduce() JustInterval {
	for i.ToFloat() >= 2.0 || i.ToFloat() < 1.0 {
		if i.ToFloat() < 1.0 {
			i.numerator *= 2
		}
		if i.ToFloat() >= 2.0 {
			i.denominator *= 2
		}
	}
	return i
}

func (i JustInterval) lessThan(other JustInterval) bool {
	return i.numerator*other.denominator < other.numerator*i.denominator
}

func (i JustInterval) greaterThan(other JustInterval) bool {
	return !i.lessThan(other) && !i.isEqualTo(other)
}

func (i JustInterval) Subtract(other JustInterval) JustInterval {
	if i.lessThan(other) {
		return JustInterval{numerator: i.denominator * other.numerator, denominator: i.numerator * other.denominator}.simplify()
	} else if i.greaterThan(other) {
		return JustInterval{numerator: i.numerator * other.denominator, denominator: i.denominator * other.numerator}.simplify()
	}
	return JustInterval{numerator: 1, denominator: 1}
}

func (i JustInterval) Name() string {
	for _, n := range intervalNames {
		if n.numerator == i.numerator && n.denominator == i.denominator {
			return n.name
		}
	}
	return ""
}

func (i JustInterval) ToFloat() float64 {
	return float64(i.numerator) / float64(i.denominator)
}

func (i JustInterval) ToTemperedInterval() TemperedInterval {
	return TemperedInterval(float64(i.numerator) / float64(i.denominator))
}

func (i JustInterval) ToPowerOf(p int) JustInterval {
	return JustInterval{numerator: uint(math.Pow(float64(i.numerator), math.Abs(float64(p)))), denominator: uint(math.Pow(float64(i.denominator), math.Abs(float64(p))))}.simplify()
}

func (i JustInterval) reciprocal() JustInterval {
	interval := JustInterval{denominator: i.numerator, numerator: i.denominator}
	return interval
}

func perfectFourth() JustInterval {
	return JustInterval{4, 3, "Perfect Fourth"}
}

var Unison = JustInterval{numerator: 1, denominator: 1}
var acuteUnison = JustInterval{numerator: 81, denominator: 80}
var SyntonicComma = JustInterval{numerator: 81, denominator: 80}
var dieses = JustInterval{numerator: 128, denominator: 125}
var justChromaticSemitone = JustInterval{numerator: 25, denominator: 24}
var graveUnison = JustInterval{numerator: 80, denominator: 81}
var lesserMajorSecond = JustInterval{numerator: 10, denominator: 9}
var greaterMajorSecond = JustInterval{numerator: 9, denominator: 8}
var diatonicSemitone = JustInterval{numerator: 16, denominator: 15}
var PerfectFifth = JustInterval{numerator: 3, denominator: 2}
var octave = JustInterval{numerator: 2, denominator: 1}

var intervalNames = []JustInterval{
	{1, 1, "Perfect Unison"},
	{225, 224, "Septimal Kleisma"},
	{80, 81, "Syntonic Comma"},
	{81, 80, "Grave Unison"},
	{128, 125, "Dieses (Diminished Second)"},
	{25, 24, "Just (Lesser) Chromatic Semitone"},
	{256, 243, "Pythagorean Minor Second"},
	{135, 128, "Greater Chromatic Semitone"},
	{27, 25, "Acute Minor Second"},
	{16, 15, "Minor Second"},
	{13, 12, "Tridecimal Minor Second (Avicenna)"},
	{12, 11, "Undecimal Minor Second"},
	{15, 14, "Septimal Minor Second"},
	{10, 9, "Just (Lesser) Major Second"},
	{9, 8, "Pythagorean (Greater) Major Second"},
	{8, 7, "Septimal Major Second"},
	{6, 5, "Minor Third"},
	{5, 4, "Major Third"},
	{32, 27, "Pythagorean Minor Third"},
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
	{13, 7, "Tridecimal Major Seventh"},
	{243, 128, "Pythagorean Major Seventh"},
	{2, 1, "Perfect Octave"},
}

func (i JustInterval) sortWith(j JustInterval) int {
	return cmp.Compare(i.numerator*j.denominator, j.numerator*i.denominator)
}

func (i JustInterval) String() string {
	return fmt.Sprintf("%d:%d", i.numerator, i.denominator)
}

func (i JustInterval) toCents() float64 {
	return math.Log10(float64(i.numerator)/float64(i.denominator)) / math.Log10(2) * 1200
}

func intervalsFromIntegers(integers [][]uint) []JustInterval {
	var intervals []JustInterval
	for _, pair := range integers {
		intervals = append(intervals, fromIntArray(pair))
	}
	return intervals
}

func fromIntArray(i []uint) JustInterval {
	return JustInterval{numerator: i[0], denominator: i[1]}.simplify()
}

func sortIntervals(intervals []JustInterval) {
	slices.SortFunc(intervals, func(i, j JustInterval) int {
		return i.sortWith(j)
	})
}

// intervalFilterFunction defines a function type for excluding certain ratios based on scale symmetry.
type intervalFilterFunction func(ratio JustInterval) bool

func multipliers(base uint) [][]uint {
	return [][]uint{{base, 1}, {1, 1}, {1, base}}
}

func justIntervalsFromMultipliers(multiplierList [][]uint, filter intervalFilterFunction) []JustInterval {
	var intervals []JustInterval
	for _, multiplier := range multiplierList {
		interval := JustInterval{numerator: multiplier[0], denominator: multiplier[1]}.octaveReduce()
		if interval.isDiminishedFifth() {
			continue
		}
		if filter(interval) {
			continue
		}
		intervals = append(intervals, interval)
	}
	intervals = append(intervals, octave)
	sortIntervals(intervals)
	return intervals
}

func createMultiplierTableOf(multiplierListA, multiplierListB [][]uint) [][]uint {
	var multiplierTable [][]uint
	for _, multiplierA := range multiplierListA {
		for _, multiplierB := range multiplierListB {
			multiplierTable = append(multiplierTable, []uint{multiplierA[0] * multiplierB[0], multiplierA[1] * multiplierB[1]})
		}
	}
	return multiplierTable
}
