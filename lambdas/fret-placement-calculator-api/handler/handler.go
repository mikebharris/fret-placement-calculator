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

const (
	defaultEqualTemperamentDivisions = 31
	defaultJustSymmetry              = "asymmetric"
	defaultOctavesToCompute          = 1
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

func (h Handler) HandleRequest(_ context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	q := request.QueryStringParameters

	scaleLength, err := strconv.ParseFloat(q["scaleLength"], 64)
	if err != nil || scaleLength <= 0 {
		return errorResponse(http.StatusUnprocessableEntity, `{"error":"a numeric scaleLength greater than zero is required"}`), nil
	}

	var fretPlacements FretPlacements

	switch q["tuningSystem"] {
	case "equal":
		fretPlacements = h.fretPlacementsForEqualTemperamentTuning(scaleLength, parseIntegerQueryParameter(q, "divisions", defaultEqualTemperamentDivisions))

	case "saz":
		fretPlacements = h.fretPlacementsForSazTuning(scaleLength)

	case "pythagorean":
		fretPlacements = h.fretPlacementsForPythagoreanTuning(scaleLength)

	case "meantone":
		fretPlacements = h.fretPlacementsForQuarterCommaMeantoneTuning(scaleLength, false)
	case "extendedMeantone":
		fretPlacements = h.fretPlacementsForQuarterCommaMeantoneTuning(scaleLength, true)

	case "":
		fallthrough
	case "ptolemy":
		fretPlacements = h.fretPlacementsForPtolemysIntenseDiatonicTuning(scaleLength, parseIntegerQueryParameter(q, "octaves", defaultOctavesToCompute), validDiatonicModeOrDefault(q["diatonicMode"]))

	case "just5limitFromPythagorean":
		fretPlacements = h.fretPlacementsFor5LimitJustChromaticTuningBuiltFromAdjustingPythagoreanScale(scaleLength)

	case "just5limitFromRatios":
		fretPlacements = h.fretPlacementsFor5LimitJustChromaticScaleBasedOnPureRatios(scaleLength, parseStringQueryParameter(q, "justSymmetry", defaultJustSymmetry))

	case "just7limitFromRatios":
		fretPlacements = h.fretPlacementsFor7LimitJustChromaticScaleBasedOnPureRatios(scaleLength, parseStringQueryParameter(q, "justSymmetry", defaultJustSymmetry))

	case "bachWellTemperament":
		fretPlacements = h.fretPlacementsForBachWohltemperierteKlavier(scaleLength)

	default:
		return errorResponse(http.StatusUnprocessableEntity, `{"error":"please provide a valid tuning system"}`), nil
	}

	fretPlacements.ScaleLength = scaleLength
	body, _ := json.Marshal(fretPlacements)
	return events.LambdaFunctionURLResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(body)}, nil
}

func validDiatonicModeOrDefault(mode string) string {
	if mode == "" || !isValidDiatonicMode(mode) {
		mode = "Ionian"
	}
	return mode
}

func errorResponse(status int, body string) events.LambdaFunctionURLResponse {
	return events.LambdaFunctionURLResponse{StatusCode: status, Headers: headers, Body: body}
}

func parseStringQueryParameter(q map[string]string, key, fallback string) string {
	if v := q[key]; v != "" {
		return v
	}
	return fallback
}

func parseIntegerQueryParameter(q map[string]string, key string, fallback int) int {
	if q[key] == "" {
		return fallback
	}
	atoi, err := strconv.Atoi(q[key])
	if err != nil || atoi <= 0 {
		return fallback
	}
	return atoi
}

func isValidDiatonicMode(mode string) bool {
	switch mode {
	case "Lydian", "Ionian", "Mixolydian", "Dorian", "Aeolian", "Phrygian", "Locrian":
		return true
	default:
		return false
	}
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

func (h Handler) fretPlacementsFor5LimitJustChromaticTuningBuiltFromAdjustingPythagoreanScale(scaleLength float64) FretPlacements {
	// Derive scale by adjusting Pythagorean scale by syntonic comma (81/80)
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
		Description: "Fret positions for chromatic scale based on 5-limit just intonation pure ratios derived from applying syntonic comma to Pythagorean ratios.",
		Frets:       h.ratiosToFretPlacements(scaleLength, ratios),
	}
}

// noteFilterFunction defines a function type for excluding certain ratios based on scale symmetry.
type noteFilterFunction func([]uint) bool

func (h Handler) fretPlacementsFor5LimitJustChromaticScaleBasedOnPureRatios(scaleLength float64, symmetry string) FretPlacements {
	return FretPlacements{
		System:      "5-limit Just Intonation",
		Description: "Fret positions for chromatic scale based on 5-limit just intonation pure ratios derived from third- and fifth-partial ratios.",
		Frets:       h.ratiosToFretPlacements(scaleLength, computeFiveLimitJustRatios(symmetry)),
	}
}

func computeFiveLimitJustRatios(symmetry string) [][]uint {
	thirdPartialMultipliers := [][]uint{{1, 9}, {1, 3}, {1, 1}, {3, 1}, {9, 1}}
	fifthPartialMultipliers := [][]uint{{5, 1}, {1, 1}, {1, 5}}
	return computeJustRatios(fifthPartialMultipliers, thirdPartialMultipliers, fiveLimitScaleFilter(symmetry))
}

func fiveLimitScaleFilter(symmetry string) func(r []uint) bool {
	return func(r []uint) bool {
		if symmetry == "asymmetric" && (isLesserMajorSecond(r) || isLesserMinorSeventh(r)) {
			return true
		}
		if symmetry == "symmetric1" && (isLesserMajorSecond(r) || isGreaterMinorSeventh(r)) {
			return true
		}
		if symmetry == "symmetric2" && (isGreaterMajorSecond(r) || isLesserMinorSeventh(r)) {
			return true
		}
		return false
	}
}

func computeJustRatios(multiplier1, multiplier2 [][]uint, filter noteFilterFunction) [][]uint {
	var ratios [][]uint
	for _, tpm := range multiplier1 {
		for _, fpm := range multiplier2 {
			numerator := tpm[0] * fpm[0]
			denominator := tpm[1] * fpm[1]
			ratio := octaveReduceIntegerRatio(fractionToLowestDenominator([]uint{numerator, denominator}))
			if isFundamental(ratio) || isDiminishedFifth(ratio) {
				continue
			}
			if filter(ratio) {
				continue
			}
			ratios = append(ratios, ratio)
		}
	}

	ratios = append(ratios, []uint{2, 1})

	slices.SortFunc(ratios, func(x, y []uint) int {
		return cmp.Compare(float64(x[0])/float64(x[1]), float64(y[0])/float64(y[1]))
	})
	return ratios
}

func isFundamental(ratio []uint) bool {
	return ratio[0] == 1 && ratio[1] == 1
}

func isLesserMajorSecond(ratio []uint) bool {
	return ratio[0] == 10 && ratio[1] == 9
}

func isGreaterMajorSecond(ratio []uint) bool {
	return ratio[0] == 9 && ratio[1] == 8
}

func isDiminishedFifth(ratio []uint) bool {
	return ratio[0] == 64 && ratio[1] == 45
}

func isLesserMinorSeventh(ratio []uint) bool {
	return ratio[0] == 16 && ratio[1] == 9
}

func isGreaterMinorSeventh(ratio []uint) bool {
	return ratio[0] == 9 && ratio[1] == 5
}

func ratioIsPerfect(ratio []uint) bool {
	return (ratio[0] == 1 && ratio[1] == 1) || (ratio[0] == 4 && ratio[1] == 3) || (ratio[0] == 3 && ratio[1] == 2) || (ratio[0] == 2 && ratio[1] == 1)
}

// Modified function to produce 7-limit just intonation frets.
// It reuses computeJustRatios by building multiplier lists that include 5 and 7 factors.
func (h Handler) fretPlacementsFor7LimitJustChromaticScaleBasedOnPureRatios(scaleLength float64, symmetry string) FretPlacements {
	return FretPlacements{
		System:      "7-limit Just Intonation",
		Description: "Fret positions for chromatic scale based on 7-limit just intonation pure ratios derived from third-, fifth- and seventh-partial ratios.",
		Frets:       h.ratiosToFretPlacements(scaleLength, computeSevenLimitJustRatios(symmetry)),
	}
}

// computeSevenLimitJustRatios builds multipliers combining small powers of 5 and 7 (exponents -1..1)
// and combines them with powers of 3 (as in the previous thirdPartialMultipliers).
func computeSevenLimitJustRatios(symmetry string) [][]uint {
	// powers of 3 used previously: 3^-2 .. 3^2
	thirdPartialMultipliers := [][]uint{{1, 9}, {1, 3}, {1, 1}, {3, 1}, {9, 1}}

	// build multipliers for combinations of 5^e5 * 7^e7 for e in {-1,0,1}
	var fiveSevenMultipliers [][]uint
	for e5 := -1; e5 <= 1; e5++ {
		for e7 := -1; e7 <= 1; e7++ {
			num, den := uint(1), uint(1)

			if e5 > 0 {
				num *= uint(intPow(5, e5))
			} else if e5 < 0 {
				den *= uint(intPow(5, -e5))
			}

			if e7 > 0 {
				num *= uint(intPow(7, e7))
			} else if e7 < 0 {
				den *= uint(intPow(7, -e7))
			}

			fiveSevenMultipliers = append(fiveSevenMultipliers, []uint{num, den})
		}
	}

	pool := computeJustRatios(fiveSevenMultipliers, thirdPartialMultipliers, sevenLimitScaleFilter(symmetry))
	return selectRatiosFromPool(pool, sevenLimitWantedRatios())
}

func sevenLimitWantedRatios() [][]uint {
	return [][]uint{
		{15, 14},
		{8, 7},
		{6, 5},
		{5, 4},
		{4, 3},
		{7, 5},
		{3, 2},
		{8, 5},
		{5, 3},
		{7, 4},
		{15, 8},
		{2, 1},
	}
}

func selectRatiosFromPool(pool, wanted [][]uint) [][]uint {
	present := make(map[[2]uint]struct{}, len(pool))
	for _, r := range pool {
		present[[2]uint{r[0], r[1]}] = struct{}{}
	}

	out := make([][]uint, 0, len(wanted))
	for _, w := range wanted {
		if _, ok := present[[2]uint{w[0], w[1]}]; ok {
			out = append(out, w)
		}
	}
	return out
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

		interval := ratio / prevRatio
		frets = append(frets, Fret{
			Label:    fmt.Sprintf("%d (%s)", fretNumber, noteNames[fretNumber]),
			Position: math.Round((scaleLength-(scaleLength/ratio))*10) / 10,
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
	for _, justRatio := range JustRatioNoteNames {
		if justRatio.Numerator == ratio[0] && justRatio.Denominator == ratio[1] {
			return justRatio.Name
		}
	}
	return ""
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

// Bach's Well Temperament as per
// Bach's Extraordinary Temperament: Our Rosetta Stone - Bradley Lehman
// Early Music Volume 33, Number 1, February 2005, pp. 3-23 (Article)
// Reference: https://academic.oup.com/em/article-abstract/33/1/3/535436?redirectedFrom=fulltext
func (h Handler) fretPlacementsForBachWohltemperierteKlavier(scaleLength float64) FretPlacements {
	// Narrowing of the fifths as outlined by Lehman
	syntonicComma := 81.0 / 80.0
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
		temperedFifths[i] = 3.0 / 2.0 * math.Pow(syntonicComma, temperingFractions[i])
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

	intervalNames := []string{"Unison", "Minor Second", "Major Second", "Minor Third", "Major Third", "Fourth", "Augmented Fourth", "Fifth", "Augmented Fifth", "Major Sixth", "Minor Seventh", "Major Seventh"}
	prevRatio := 1.0
	var frets []Fret
	for fretNumber, ratio := range ratios {
		if fretNumber == 0 {
			continue // Skip the tonic (unfretted)
		}
		distanceFromNut := scaleLength - (scaleLength / ratio)
		frets = append(frets, Fret{
			Label:    fmt.Sprintf("%d (%s)", fretNumber, intervalNames[fretNumber%len(intervalNames)]),
			Position: math.Round(distanceFromNut*1000) / 1000, // Round to 3 decimal places
			Comment:  fmt.Sprintf("ratio: %.6f; interval: %.6f", ratio, ratio/prevRatio),
		})
		prevRatio = ratio
	}

	// Add the octave fret (2nd partial)
	frets = append(frets, Fret{
		Label:    fmt.Sprintf("%d (Octave)", len(frets)+1),
		Position: scaleLength / 2,
		Comment:  fmt.Sprintf("ratio: %.1f; interval: %.6f", 2.0, 2.0/prevRatio),
	})

	description := "Fret positions derived from Lehman's decoding of Bach's Well-Tempered tuning, using sixth-comma, twelfth-comma, and pure fifths."
	return FretPlacements{
		System:      "Bach's Well-Tempered Tuning",
		Description: description,
		ScaleLength: scaleLength,
		Frets:       frets,
	}
}

// simple integer power for small positive exponents
func intPow(base, exp int) int {
	if exp <= 0 {
		return 1
	}
	res := 1
	for i := 0; i < exp; i++ {
		res *= base
	}
	return res
}

// sevenLimitScaleFilter excludes ratios to enforce a particular "symmetry" choice.
// For the current API/tests, "asymmetric" prefers septimal notes where there are competing 5-limit options.
func sevenLimitScaleFilter(symmetry string) func([]uint) bool {
	return func(r []uint) bool {
		if symmetry != "asymmetric" {
			return false
		}

		// Prefer septimal major second (8:7) over Pythagorean major second (9:8).
		if isGreaterMajorSecond(r) {
			return true
		}
		// Prefer septimal minor seventh (7:4) over 5-limit greater minor seventh (9:5)
		// and Pythagorean lesser minor seventh (16:9).
		if isLesserMinorSeventh(r) || isGreaterMinorSeventh(r) {
			return true
		}

		return false
	}
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
