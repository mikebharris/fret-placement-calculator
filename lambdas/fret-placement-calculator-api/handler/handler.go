package handler

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"slices"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

var headers = map[string]string{
	"Content-Type": "application/json",
}

type Handler struct {
}

type Fret struct {
	Label    string  `json:"label"`
	Position float64 `json:"position"`
	Comment  string  `json:"comment,omitempty"`
	Interval string  `json:"interval,omitempty"`
}

type FretPlacements struct {
	System      string  `json:"system"`
	Description string  `json:"description,omitempty"`
	ScaleLength float64 `json:"scaleLength"`
	Frets       []Fret  `json:"frets"`
}

type Ratio struct {
	Numerator   uint
	Denominator uint
	Name        string
}

var JustRatios = []Ratio{
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
	{729, 512, "Pythagorean Augmented Forth"},
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

func (h Handler) HandleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	scaleLength, err := strconv.ParseFloat(request.QueryStringParameters["scaleLength"], 64)
	if err != nil || scaleLength <= 0 {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"a numeric scaleLength greater than zero is required"}`}, nil
	}

	var fretPlacements FretPlacements

	switch request.QueryStringParameters["tuningSystem"] {
	case "equal":
		divisionsOfOctave := 31
		if request.QueryStringParameters["divisions"] != "" {
			divisionsOfOctave, err = strconv.Atoi(request.QueryStringParameters["divisions"])
		}
		if err != nil {
			return events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"please provide number of divisions for equal temperament"}`}, nil
		}
		fretPlacements = h.fretPlacementsForEqualTemperamentTuning(scaleLength, divisionsOfOctave)
	case "saz":
		fretPlacements = h.fretPlacementsForSazTuning(scaleLength)
	case "pythagorean":
		fretPlacements = h.fretPlacementsForPythagoreanTuning(scaleLength)
	case "meantone":
		extended := false
		if request.QueryStringParameters["extended"] != "" {
			extended, _ = strconv.ParseBool(request.QueryStringParameters["extended"])
		}
		fretPlacements = h.fretPlacementsForQuarterCommaMeantoneTuning(scaleLength, extended)
	case "":
		fallthrough
	case "just":
		var octaves = 1
		if request.QueryStringParameters["octaves"] != "" {
			octaves, err = strconv.Atoi(request.QueryStringParameters["octaves"])
			if err != nil {
				return events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"please provide a valid positive number for number of octaves worth of frets"}`}, nil
			}
		}
		if request.QueryStringParameters["diatonicMode"] != "" {
			mode := "Ionian"
			mode = request.QueryStringParameters["diatonicMode"]
			if mode != "Lydian" && mode != "Ionian" && mode != "Mixolydian" && mode != "Dorian" && mode != "Aeolian" && mode != "Phrygian" && mode != "Locrian" {
				return events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"please provide a valid mode for the diatonic scale"}`}, nil
			}
			fretPlacements = h.fretPlacementsForPtolemysIntenseDiatonicTuning(scaleLength, octaves, mode)
		} else {
			fretPlacements = h.fretPlacementsFor5LimitJustChromaticTuning(scaleLength, octaves)
		}
	default:
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"invalid temper parameter"}`}, nil
	}

	fretPlacements.ScaleLength = scaleLength
	body, _ := json.Marshal(fretPlacements)

	return events.LambdaFunctionURLResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(body)}, nil

}

func (h Handler) fretPlacementsForPythagoreanTuning(scaleLength float64) FretPlacements {
	return FretPlacements{
		System:      "Pythagorean",
		Description: "Fret positions based on 3-limit Pythagorean ratios.",
		Frets:       h.ratiosToFretPlacements(scaleLength, computePythagoreanRatios()),
	}
}

func computePythagoreanRatios() [][]uint {
	// divide by 3:2 = 4/9 * 2/3 = 16/27 = 2^3 : 3^3 = 16/27
	// divide by 3:2 = 2/3 * 2/3 = 4/9 - octave reduce to 16:9
	// divide by 3:2 = 2:3 - octave reduce to 4:3
	// fundamental = 1,1
	// multiply by 3:2 = 3,2
	// multiply by 3:2 = 9,4 - octave reduce to 9:8
	// multiply by 3:2 = 27:8 = 3^3:2^3 - octave reduce to 27:16
	// multiply by 3:2 = 81:16 - octave reduce to 81:64 - 3^4:2^4 = 81:16

	var thirdPartial = []uint{3, 2}
	var thirdPartialsFromTonicToCompute = 6
	var ratiosFromFundamental [][]uint

	for i := -thirdPartialsFromTonicToCompute; i <= thirdPartialsFromTonicToCompute; i++ {
		t := math.Pow(float64(thirdPartial[0]), math.Abs(float64(i)))
		b := math.Pow(float64(thirdPartial[1]), math.Abs(float64(i)))
		if i < 0 {
			ratiosFromFundamental = append(ratiosFromFundamental, octaveReduceIntegerRatio([]uint{uint(b), uint(t)}))
		} else if i > 0 {
			ratiosFromFundamental = append(ratiosFromFundamental, octaveReduceIntegerRatio([]uint{uint(t), uint(b)}))
		}
	}

	ratiosFromFundamental = append(ratiosFromFundamental, []uint{2, 1})

	slices.SortFunc(ratiosFromFundamental, func(x, y []uint) int {
		return cmp.Compare(float64(x[0])/float64(x[1]), float64(y[0])/float64(y[1]))
	})
	return ratiosFromFundamental
}

func octaveReduceIntegerRatio(ratio []uint) []uint {
	for ratio[0]/ratio[1] >= 2.0 || ratio[0]/ratio[1] < 1.0 {
		if ratio[0]/ratio[1] < 1.0 {
			ratio[0] *= 2
		}
		if ratio[0]/ratio[1] >= 2.0 {
			ratio[1] *= 2
		}
	}
	return ratio
}

func (h Handler) fretPlacementsFor5LimitJustChromaticTuning(scaleLength float64, octaves int) FretPlacements {
	// 	m2 : 256/243 → 16/15
	//	M2 : 9/8 → 10/9
	//	m3 : 32/27 → 6/5
	//	M3 : 81/64 → 5/4
	//	m6 : 128/81 → 8/5
	//	M6 : 27/16 → 5/3
	//	m7 : 16/9 → 9/5
	//	M7 : 243/128 → 15/8

	var acuteUnison = []uint{81, 80}
	var graveUnison = []uint{80, 81}

	var ratios [][]uint

	for _, ratio := range computePythagoreanRatios() {
		if ratioIsPerfect(ratio) {
			ratios = append(ratios, ratio)
			continue
		}

		graveRatio := octaveReduceIntegerRatio(fractionToLowestDenominator([]uint{ratio[0] * acuteUnison[0], ratio[1] * acuteUnison[1]}))
		acuteRatio := octaveReduceIntegerRatio(fractionToLowestDenominator([]uint{ratio[0] * graveUnison[0], ratio[1] * graveUnison[1]}))

		if graveRatio[1] < acuteRatio[1] {
			ratios = append(ratios, graveRatio)
		} else {
			ratios = append(ratios, acuteRatio)
		}
	}

	return FretPlacements{
		System:      "5-limit Just Intonation",
		Description: fmt.Sprintf("Fret positions for chromatic scale based on 5-limit just intonation pure ratios derived from applying syntonic comma to Pythagorean ratios."),
		Frets:       h.ratiosToFretPlacements(scaleLength, ratios),
	}
}

func (h Handler) fretPlacementsFor7LimitJustChromaticTuning(scaleLength float64, octaves int) FretPlacements {
	// 	m2 : 256/243 → 16/15
	//	M2 : 9/8 → 10/9
	//	m3 : 32/27 → 6/5
	//	M3 : 81/64 → 5/4
	//	m6 : 128/81 → 8/5
	//	M6 : 27/16 → 5/3
	//	m7 : 16/9 → 9/5
	//	M7 : 243/128 → 15/8

	var acuteUnison = []uint{81, 80}
	var graveUnison = []uint{80, 81}

	var septimalKleisma = []uint{225, 224}

	var ratios [][]uint

	for _, ratio := range computePythagoreanRatios() {
		if ratioIsPerfect(ratio) {
			ratios = append(ratios, ratio)
			continue
		}

		graveRatio := octaveReduceIntegerRatio(fractionToLowestDenominator([]uint{ratio[0] * acuteUnison[0], ratio[1] * acuteUnison[1]}))
		acuteRatio := octaveReduceIntegerRatio(fractionToLowestDenominator([]uint{ratio[0] * graveUnison[0], ratio[1] * graveUnison[1]}))

		if graveRatio[1] < acuteRatio[1] {
			ratios = append(ratios, graveRatio)
		} else {
			ratios = append(ratios, acuteRatio)
		}
	}

	for i, ratio := range ratios {
		graveKleismaRatio := octaveReduceIntegerRatio(fractionToLowestDenominator([]uint{ratio[0] * septimalKleisma[0], ratio[1] * septimalKleisma[1]}))
		acuteKleismaRatio := octaveReduceIntegerRatio(fractionToLowestDenominator([]uint{ratio[0] * septimalKleisma[1], ratio[1] * septimalKleisma[0]}))
		fmt.Println(acuteKleismaRatio)
		fmt.Println(graveKleismaRatio)
		if graveKleismaRatio[1] < acuteKleismaRatio[1] {
			if graveKleismaRatio[0] < ratios[i][0] {
				ratios[i][0] = graveKleismaRatio[0]
				ratios[i][1] = graveKleismaRatio[1]
			}
		} else {
			if acuteKleismaRatio[0] < ratios[i][0] {
				ratios[i][0] = acuteKleismaRatio[0]
				ratios[i][1] = acuteKleismaRatio[1]
			}
		}
	}

	fmt.Println(ratios)

	return FretPlacements{
		System:      "7-limit Just Intonation",
		Description: fmt.Sprintf("Fret positions for chromatic scale based on 7-limit just intonation pure ratios derived from applying syntonic comma and septimal kleisma to Pythagorean ratios."),
		Frets:       h.ratiosToFretPlacements(scaleLength, ratios),
	}
}

func ratioIsPerfect(ratio []uint) bool {
	return (ratio[0] == 1 && ratio[1] == 1) || (ratio[0] == 4 && ratio[1] == 3) || (ratio[0] == 3 && ratio[1] == 2) || (ratio[0] == 2 && ratio[1] == 1)
}

func fractionToLowestDenominator(fraction []uint) []uint {
	gcd := func(a, b uint) uint {
		for b != 0 {
			a, b = b, a%b
		}
		return a
	}(fraction[0], fraction[1])
	fraction[0] = fraction[0] / gcd
	fraction[1] = fraction[1] / gcd
	return fraction
}

func (h Handler) fretPlacementsForPtolemysIntenseDiatonicTuning(scaleLength float64, octaves int, mode string) FretPlacements {
	var intervalMap = map[string][][]uint{
		"Lydian":     {{9, 8}, {10, 9}, {9, 8}, {16, 15}, {10, 9}, {9, 8}, {16, 15}},
		"Ionian":     {{9, 8}, {10, 9}, {16, 15}, {9, 8}, {10, 9}, {9, 8}, {16, 15}},
		"Mixolydian": {{9, 8}, {10, 9}, {16, 15}, {9, 8}, {10, 9}, {16, 15}, {9, 8}},
		"Dorian":     {{9, 8}, {16, 15}, {10, 9}, {9, 8}, {10, 9}, {16, 15}, {9, 8}},
		"Aeolian":    {{9, 8}, {16, 15}, {10, 9}, {9, 8}, {16, 15}, {9, 8}, {10, 9}},
		"Phrygian":   {{16, 15}, {9, 8}, {10, 9}, {9, 8}, {16, 15}, {10, 9}, {9, 8}},
		"Locrian":    {{16, 15}, {9, 8}, {10, 9}, {16, 15}, {9, 8}, {10, 9}, {9, 8}},
	}

	var ratios = make([][]uint, 0)
	var ratio = []uint{1, 1}

	for i := 0; i < octaves; i++ {
		for _, v := range intervalMap[mode] {
			ratio = fractionToLowestDenominator(
				[]uint{
					ratio[0] * v[0], ratio[1] * v[1],
				})

			ratios = append(ratios, ratio)
		}
	}

	return FretPlacements{
		System:      "Ptolemy",
		Description: fmt.Sprintf("Fret positions for Ptolemy's 5-limit intense diatonic scale in %s mode.", mode),
		Frets:       h.ratiosToFretPlacements(scaleLength, ratios),
	}
}

func (h Handler) fretPlacementsForSazTuning(scaleLength float64) FretPlacements {
	// as per https://en.wikipedia.org/wiki/Ba%C4%9Flama and the cura that I have
	return FretPlacements{
		System:      "saz",
		Description: "Fret positions for traditional Turkish Saz tuning ratios.",
		ScaleLength: scaleLength,
		Frets:       h.ratiosToFretPlacements(scaleLength, [][]uint{{18, 17}, {12, 11}, {9, 8}, {81, 68}, {27, 22}, {81, 64}, {4, 3}, {24, 17}, {16, 11}, {3, 2}, {27, 17}, {18, 11}, {27, 16}, {16, 9}, {32, 17}, {64, 33}, {2, 1}}),
	}
}

func (h Handler) fretPlacementsForQuarterCommaMeantoneTuning(scaleLength float64, extendScale bool) FretPlacements {
	syntonicComma := 81.0 / 80.0
	fractionOfSyntonicCommaToTemperFifthsBy := 0.25
	temperedFifth := 3.0 / 2.0 * math.Pow(syntonicComma, -fractionOfSyntonicCommaToTemperFifthsBy)

	var fifthsFromTonic int
	var noteNames []string

	if extendScale {
		fifthsFromTonic = 9
		noteNames = []string{"D", "D#", "Eb", "E", "Fb", "F", "F#", "Gb", "G", "G#", "Ab", "A", "A#", "Bb", "B", "Cb", "C", "C#", "Db"}
	} else {
		fifthsFromTonic = 6
		noteNames = []string{"D", "Eb", "E", "F", "F#", "G", "G#", "Ab", "A", "Bb", "B", "C", "C#"}
	}

	var ratiosOfNotesToFundamental []float64
	for i := -fifthsFromTonic; i <= fifthsFromTonic; i++ {
		ratiosOfNotesToFundamental = append(ratiosOfNotesToFundamental, octaveReduceFloat(math.Pow(temperedFifth, float64(i))))
	}

	slices.Sort(ratiosOfNotesToFundamental)

	prevRatio := 1.0
	var frets []Fret
	for fretNumber, ratio := range ratiosOfNotesToFundamental {
		if fretNumber == 0 {
			continue
		}
		distanceFromNut, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", scaleLength-(scaleLength/ratio)), 64)
		interval := ratio / prevRatio
		frets = append(frets, Fret{
			Label:    fmt.Sprintf("%d (%s)", fretNumber, noteNames[fretNumber]),
			Position: distanceFromNut,
			Comment:  fmt.Sprintf("ratio: %.3f; interval: %.6f", ratio, interval),
		})
		prevRatio = ratio
	}

	frets = append(frets, Fret{
		Label:    fmt.Sprintf("%d (Octave)", len(frets)+1),
		Position: scaleLength / 2,
		Comment:  fmt.Sprintf("ratio: %.1f; interval: %.6f", 2.0, 2.0/prevRatio),
	})

	var description string
	if extendScale {
		description = fmt.Sprintf("Fret positions for extended meantone computed by narrowing of fifths by %.2f of a syntonic comma (81/80).  Nominal note names used given a tonic of D.", fractionOfSyntonicCommaToTemperFifthsBy)
	} else {
		description = fmt.Sprintf("Fret positions for meantone computed by narrowing of fifths by %.2f of a syntonic comma (81/80).  Nominal note names used given a tonic of D.", fractionOfSyntonicCommaToTemperFifthsBy)
	}

	return FretPlacements{
		System:      "meantone",
		Description: description,
		ScaleLength: scaleLength,
		Frets:       frets,
	}
}

func (h Handler) ratiosToFretPlacements(scaleLength float64, ratios [][]uint) []Fret {
	var frets []Fret
	var previousRatio = []uint{1, 1}
	for _, ratio := range ratios {
		distanceFromNut := math.Round((scaleLength-(scaleLength/float64(ratio[0]))*float64(ratio[1]))*100) / 100
		interval := intervalBetween(ratio, previousRatio)
		frets = append(frets, Fret{
			Label:    fmt.Sprintf("%d:%d", ratio[0], ratio[1]),
			Position: distanceFromNut,
			Comment:  intervalNameFromRatio(ratio),
			Interval: fmt.Sprintf("%d:%d", interval[0], interval[1]),
		})
		previousRatio = ratio
	}
	return frets
}

func intervalBetween(ratio []uint, ratio2 []uint) []uint {
	t := ratio[0] * ratio2[1]
	b := ratio[1] * ratio2[0]
	i := cmp.Compare(float64(t)/float64(b), 1.0)
	if i < 0 {
		return fractionToLowestDenominator([]uint{b, t})
	} else if i > 0 {
		return fractionToLowestDenominator([]uint{t, b})
	}
	return ratio
}

func intervalNameFromRatio(ratio []uint) string {
	return func() string {
		for _, justRatio := range JustRatios {
			if justRatio.Numerator == ratio[0] && justRatio.Denominator == ratio[1] {
				return justRatio.Name
			}
		}
		return ""
	}()
}

func (h Handler) fretPlacementsForEqualTemperamentTuning(scaleLength float64, divisionsOfOctave int) FretPlacements {
	system := fmt.Sprintf("%d-TET", divisionsOfOctave)
	fretPlacements := FretPlacements{
		System:      system,
		Description: fmt.Sprintf("Fret positions for %d-tone equal temperament.", divisionsOfOctave),
		ScaleLength: scaleLength,
		Frets:       nil,
	}

	for fretNumber := 1; fretNumber <= divisionsOfOctave; fretNumber++ {
		distanceFromNut, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", scaleLength-(scaleLength/math.Exp2(float64(fretNumber)/float64(divisionsOfOctave)))), 64)
		fretPlacements.Frets = append(fretPlacements.Frets, Fret{
			Label:    fmt.Sprintf("Fret %d", fretNumber),
			Position: distanceFromNut,
		})
	}
	return fretPlacements
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
