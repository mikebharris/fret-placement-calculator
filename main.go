package main

import (
	"flag"
	"fmt"
	"math"
	"slices"
)

// inspired by https://liutaiomottola.com/formulae/fret.htm
var temper = flag.String("temper", "", "Temper the intonation (equal, meantone, pythagorean, saz)")
var meantoneFifthTemperedBy = flag.Float64("temper-by", 0.25, "Meantone fifths tempered by (fraction less than one)")
var extendMeantone = flag.Bool("extend", false, "Extend meantone scale")
var scaleLength = flag.Float64("length", 540.0, "length of scale")
var divisionsOfOctave = flag.Uint("divisions", 12, "equal tempered divisions of octave (12, 19, 23, 31, 53, etc)")
var numberOfFrets = flag.Uint("frets", 22, "number of frets for equal temperament scale")
var diatonic = flag.Bool("diatonic", false, "produce diatonic scale (not yet implemented)")
var diatonicMode = flag.String("mode", "ionian", "mode for the diatonic scale (not yet implemented)")

//var modeMap []bool = []bool{false, true, false, true}
//var intervalMap [][]string = [][]string{{"lydian", "TTTSTTS"}, {"ionian", "TTSTTTS"}}

func main() {
	flag.Parse()

	switch *temper {
	case "equal":
		equalTemperament(*scaleLength, *divisionsOfOctave, *numberOfFrets)
	case "pythagorean":
		pythagorean(*scaleLength)
	case "meantone":
		meantone(*scaleLength, *meantoneFifthTemperedBy, *extendMeantone)
	case "saz":
		saz(*scaleLength)
	default:
		justIntonation(*scaleLength)
	}
}

func saz(scaleLength float64) {
	// as per https://en.wikipedia.org/wiki/Ba%C4%9Flama and the cura that I have
	var ratios = [][]uint{{18, 17}, {12, 11}, {9, 8}, {81, 68}, {27, 22}, {81, 64}, {4, 3}, {24, 17}, {16, 11}, {3, 2}, {27, 17}, {18, 11}, {27, 16}, {16, 9}, {32, 17}, {64, 33}, {2, 1}}

	fmt.Println("Calculating fret positions on the saz cura....")
	for _, ratio := range ratios {
		distanceFromNut := scaleLength - (scaleLength/float64(ratio[0]))*float64(ratio[1])
		fmt.Printf("Place %d:%d fret at %.3f\n", ratio[0], ratio[1], distanceFromNut)
	}
}

func equalTemperament(scaleLength float64, divisionsOfOctave uint, numberOfFrets uint) {
	fmt.Printf("Calculating for %d equal temperament, scale length %.3f, frets %d:\n", divisionsOfOctave, scaleLength, numberOfFrets)

	if numberOfFrets < divisionsOfOctave {
		numberOfFrets = divisionsOfOctave
	}

	for f := uint(1); f <= numberOfFrets; f++ {
		distanceFromNut := scaleLength - (scaleLength / math.Exp2(float64(f)/float64(divisionsOfOctave)))
		fmt.Printf("Place fret # %d at %.3f\n", int(f), distanceFromNut)
	}
}

func justIntonation(scaleLength float64) {
	var ratios = [][]uint{{16, 15}, {10, 9}, {9, 8}, {6, 5}, {5, 4}, {35, 25}, {4, 3}, {45, 32}, {3, 2}, {8, 5}, {5, 3}, {16, 9}, {9, 5}, {15, 8}, {2, 1}}

	fmt.Println("Calculating based on just intonation pure ratios....")
	for _, ratio := range ratios {
		distanceFromNut := scaleLength - (scaleLength/float64(ratio[0]))*float64(ratio[1])

		// frequency is the top of the ratio divided by the bottom of the ratio
		// need therefore the starting frequency of the string?

		// 440 * 3 / 2 = 660Hz - 1/3 way along the string
		// 440 * 2 / 1 = 880Hz - 1/2 way along the string
		// wavelength = 1/freq
		// 440 * 16 / 15 = 469.33Hz - 1 / 16 way along the string (rest of string is 15/16)
		// placement = string length - (string length / 16 * 15)

		fmt.Printf("Place %d:%d fret around %.3f\n", ratio[0], ratio[1], distanceFromNut)
	}
}

func pythagorean(scaleLength float64) {
	var ratios = [][]uint{{256, 243}, {9, 8}, {32, 27}, {81, 64}, {4, 3}, {1024, 729}, {729, 512}, {3, 2}, {128, 81}, {27, 16}, {16, 9}, {243, 128}, {2, 1}}

	fmt.Println("Calculating based on Pythagorean ratios....")
	for _, ratio := range ratios {
		distanceFromNut := scaleLength - (scaleLength/float64(ratio[0]))*float64(ratio[1])
		fmt.Printf("Place %d:%d fret around %.3f\n", ratio[0], ratio[1], distanceFromNut)
	}
}

func meantone(scaleLength float64, fifthTemperedByFractionOfSyntonicComma float64, extend bool) {

	// thanks to https://johncarlosbaez.wordpress.com/2023/12/13/quarter-comma-meantone-part-1/
	//
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

	fmt.Printf("Calculating based on extended meantone based on narrowing of fifths by %.2f of a syntonic comma (81/80)....\nNominal note names used based on a tonic of D:\n", fifthTemperedByFractionOfSyntonicComma)
	syntonicComma := 81.0 / 80.0
	quarterCommaFifth := 3.0 / 2.0 * math.Pow(syntonicComma, -fifthTemperedByFractionOfSyntonicComma)

	var fifthsFromTonic int
	var noteNames []string

	if extend {
		fifthsFromTonic = 9
		noteNames = []string{"D", "D#", "Eb", "E", "Fb", "F", "F#", "Gb", "G", "G#", "Ab", "A", "A#", "Bb", "B", "Cb", "C", "C#", "Db"}
	} else {
		fifthsFromTonic = 6
		noteNames = []string{"D", "Eb", "E", "F", "F#", "G", "G#", "Ab", "A", "Bb", "B", "C", "C#"}
	}

	var ratiosOfNotesToFundamental []float64
	for i := -fifthsFromTonic; i <= fifthsFromTonic; i++ {
		ratiosOfNotesToFundamental = append(ratiosOfNotesToFundamental, octaveReduce(math.Pow(quarterCommaFifth, float64(i))))
	}

	slices.Sort(ratiosOfNotesToFundamental)

	prevRatio := 1.0
	for fretNumber, ratio := range ratiosOfNotesToFundamental {
		if fretNumber == 0 {
			fmt.Printf("Open string at %.1f is D (ratio %.3f; interval %.6f)\n", 0.0, 1.0/1.0, 0.0)
			continue
		}
		distanceFromNut := scaleLength - (scaleLength / ratio)
		interval := ratio / prevRatio
		fmt.Printf("Place %d (%s) fret at %.1f (ratio %.3f; interval %.6f)\n", fretNumber, noteNames[fretNumber], distanceFromNut, ratio, interval)
		prevRatio = ratio
	}

	fmt.Printf("Place octave fret at %.1f (ratio %.3f; interval %.6f)\n", scaleLength/2, 2.0/1.0, 2.0/prevRatio)
}

func octaveReduce(ratio float64) float64 {
	for ratio >= 2.0 || ratio < 1.0 {
		if ratio >= 2.0 {
			ratio = ratio / 2.0
		}
		if ratio < 1.0 {
			ratio = ratio * 2.0
		}
	}
	return ratio
}
