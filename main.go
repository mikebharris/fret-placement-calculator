package main

import (
	"flag"
	"fmt"
	"math"
	"slices"
)

// inspired by https://liutaiomottola.com/formulae/fret.htm
var temper = flag.String("temper", "", "Temper the intonation (equal, quarter-comma, none")
var scaleLength = flag.Float64("length", 540.0, "length of scale (any unit, defaults to 540.0mm)")
var divisionsOfOctave = flag.Uint("divisions", 12, "equal tempered divisions of octave (12, 19, 23, 31, etc)")
var numberOfFrets = flag.Uint("frets", 22, "number of frets (defaults to 22)")
var diatonic = flag.Bool("diatonic", false, "produce diatonic scale (defaults to false)")
var diatonicMode = flag.String("mode", "ionian", "mode for the diatonic scale (ionian, dorian, etc)")

var modeMap []bool = []bool{false, true, false, true}
var intervalMap [][]string = [][]string{{"lydian", "TTTSTTS"}, {"ionian", "TTSTTTS"}}

func main() {
	flag.Parse()

	switch *temper {
	case "equal":
		equalTemperament()
	case "pyth":
		pythagorean()
	case "quarter-comma":
		quartercomma()
	default:
		justIntonation()
	}
}

func equalTemperament() {
	fmt.Printf("Calculating for %d equal temperament, scale length %.3f, frets %d:\n", *divisionsOfOctave, *scaleLength, *numberOfFrets)

	if *numberOfFrets < *divisionsOfOctave {
		*numberOfFrets = *divisionsOfOctave
	}

	for f := uint(1); f <= *numberOfFrets; f++ {
		distanceFromNut := *scaleLength - (*scaleLength / math.Exp2(float64(f)/float64(*divisionsOfOctave)))
		fmt.Printf("Place fret # %d at %.3f\n", int(f), distanceFromNut)
	}
}

func justIntonation() {
	var ratios = [][]uint{{16, 15}, {10, 9}, {9, 8}, {6, 5}, {5, 4}, {35, 25}, {4, 3}, {45, 32}, {3, 2}, {8, 5}, {5, 3}, {16, 9}, {9, 5}, {15, 8}, {2, 1}}

	fmt.Println("Calculating based on just intonation pure ratios....")
	for _, f := range ratios {
		distanceFromNut := *scaleLength - (*scaleLength/float64(f[0]))*float64(f[1])

		// frequency is the top of the ratio divided by the bottom of the ratio
		// need therefore the starting frequency of the string?

		// 440 * 3 / 2 = 660Hz - 1/3 way along the string
		// 440 * 2 / 1 = 880Hz - 1/2 way along the string
		// wavelength = 1/freq
		// 440 * 16 / 15 = 469.33Hz - 1 / 16 way along the string (rest of string is 15/16)
		// placement = string length - (string length / 16 * 15)

		fmt.Printf("Place %d:%d fret around %.3f\n", f[0], f[1], distanceFromNut)
	}
}

func pythagorean() {
	var ratios = [][]uint{{256, 243}, {9, 8}, {32, 27}, {81, 64}, {4, 3}, {1024, 729}, {729, 512}, {3, 2}, {128, 81}, {27, 16}, {16, 9}, {243, 128}, {2, 1}}

	fmt.Println("Calculating based on Pythagorean ratios....")
	for _, f := range ratios {
		distanceFromNut := *scaleLength - (*scaleLength/float64(f[0]))*float64(f[1])
		fmt.Printf("Place %d:%d fret around %.3f\n", f[0], f[1], distanceFromNut)
	}
}

func quartercomma() {

	// the scale has major (diatonic) and minor (chromatic) semitones
	// the chromatic is (5^7/4)/16 ~= 1.04491
	// the diatonic is 8/(5^5/4) ~= 1.06998
	// and the lesser dieses = 2^7 / 5^3 = 128/125 = 1.024
	// diatonic semitone = dieses x chromatic semitone
	// therefore in the scale: D Eb E F F# G G# Ab A Bb B C C# (13 notes)
	// we get rid of the dieses by ditching the diminished fifth (Ab) and keeping the augmented forth (G#)
	// this leaves D (d) Eb (c) E (d) F (c) F# (d) G (c) G# (d) A (d) Bb (c) B (d) C (c) C# (d) D
	// with two diatonic semitones either side of the tonic (D) and forth (G) or one back in the circle of ratiosOfNotesToFundamental
	//
	// there are two tritones: D -> G# and D -> Ab

	fmt.Println("Calculating based on quartercomma meantone based on narrowing of fifths by one quarter of a syntonic comma")
	syntonicComma := 81.0 / 80.0
	quarterCommaFifth := 3.0 / 2.0 * math.Pow(syntonicComma, -1.0/4.0)

	var ratiosOfNotesToFundamental []float64

	fifthsFromTonic := 6
	for i := -fifthsFromTonic; i <= fifthsFromTonic; i++ {
		ratiosOfNotesToFundamental = append(ratiosOfNotesToFundamental, octaveReduce(math.Pow(quarterCommaFifth, float64(i))))
	}

	slices.Sort(ratiosOfNotesToFundamental)

	noteNames := []string{"C", "Db", "D", "Eb", "E", "F", "F#", "Gb", "G", "Ab", "A", "Bb", "B"}

	prevRatio := 1.0
	for fretNumber, ratio := range ratiosOfNotesToFundamental {
		distanceFromNut := *scaleLength - (*scaleLength / ratio)
		interval := ratio / prevRatio
		fmt.Printf("Place %d (%s) fret at %.1f (ratio %.3f; interval %.6f)\n", fretNumber, noteNames[fretNumber], distanceFromNut, ratio, interval)
		prevRatio = ratio
	}

	octaveFretNumber := fifthsFromTonic*2 + 1
	fmt.Printf("Place %d (C) fret at %.1f (ratio 2.0; interval %.6f)\n", octaveFretNumber, *scaleLength/2.0, 2.0/prevRatio)

	chromaticSemitone := 1.044906727
	diatonicSemitone := 1.069984488

	intervals := []float64{diatonicSemitone, chromaticSemitone, diatonicSemitone, chromaticSemitone, diatonicSemitone, chromaticSemitone, diatonicSemitone, diatonicSemitone, chromaticSemitone, diatonicSemitone, chromaticSemitone, diatonicSemitone}

	fmt.Printf("Calculating based on quartercomma meantone using chromatic (lesser) semitone of %.6f and diatonic (greater) semitone of %.6f\n for a scale length of %.2f\n", chromaticSemitone, diatonicSemitone, *scaleLength)

	frequencyMultiplier := 1.0
	for n, interval := range intervals {
		frequencyMultiplier = frequencyMultiplier * interval
		distanceFromNut := *scaleLength - *scaleLength*1/frequencyMultiplier
		fmt.Printf("Place %d fret at %.3f\n", n+1, distanceFromNut)
	}
}

func octaveReduce(pow float64) float64 {
	for pow >= 2.0 || pow < 1.0 {
		if pow >= 2.0 {
			pow = pow / 2.0
		}
		if pow < 1.0 {
			pow = pow * 2.0
		}
	}
	return pow
}
