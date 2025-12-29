package handler

import (
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
	intervals := computePythagoreanIntervals()
	return FretPlacements{
		System:      "Pythagorean",
		Description: "Fret positions based on 3-limit Pythagorean ratios.",
		Frets:       h.intervalsToFretPlacements(scaleLength, intervals),
	}
}

func computePythagoreanIntervals() []Interval {
	// divide by 3:2 = 4/9 * 2/3 = 16/27 = 2^3 : 3^3 = 16/27
	// divide by 3:2 = 2/3 * 2/3 = 4/9 - octave reduce to 16:9
	// divide by 3:2 = 2:3 - octave reduce to 4:3
	// fundamental = 1,1
	// multiply by 3:2 = 3,2
	// multiply by 3:2 = 9,4 - octave reduce to 9:8
	// multiply by 3:2 = 27:8 = 3^3:2^3 - octave reduce to 27:16
	// multiply by 3:2 = 81:16 - octave reduce to 81:64 - 3^4:2^4 = 81:16

	var thirdPartialsFromTonicToCompute = 6
	var intervals []Interval

	for i := -thirdPartialsFromTonicToCompute; i <= thirdPartialsFromTonicToCompute; i++ {
		if i < 0 {
			intervals = append(intervals, perfectFifth.toPowerOf(i).reciprocal().octaveReduce())
		} else if i > 0 {
			intervals = append(intervals, perfectFifth.toPowerOf(i).octaveReduce())
		}
	}

	intervals = append(intervals, octave)
	sortIntervals(intervals)
	return intervals
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

	return FretPlacements{
		System:      "5-limit Just Intonation",
		Description: "Fret positions for chromatic scale based on 5-limit just intonation pure ratios derived from applying syntonic comma to Pythagorean ratios.",
		Frets:       h.intervalsToFretPlacements(scaleLength, intervals),
	}
}

// intervalFilterFunction defines a function type for excluding certain ratios based on scale symmetry.
type intervalFilterFunction func(ratio Interval) bool

func (h Handler) fretPlacementsFor5LimitJustChromaticScaleBasedOnPureRatios(scaleLength float64, symmetry string) FretPlacements {
	return FretPlacements{
		System:      "5-limit Just Intonation",
		Description: "Fret positions for chromatic scale based on 5-limit just intonation pure ratios derived from third- and fifth-partial ratios.",
		Frets:       h.intervalsToFretPlacements(scaleLength, computeFiveLimitJustIntervals(symmetry)),
	}
}

func computeFiveLimitJustIntervals(symmetry string) []Interval {
	thirdPartialMultipliers := [][]uint{{1, 9}, {1, 3}, {1, 1}, {3, 1}, {9, 1}}
	fifthPartialMultipliers := [][]uint{{5, 1}, {1, 1}, {1, 5}}
	return computeJustIntervals(createMultiplierTableOf(thirdPartialMultipliers, fifthPartialMultipliers), fiveLimitScaleFilter(symmetry))
}

func fiveLimitScaleFilter(symmetry string) func(interval Interval) bool {
	return func(interval Interval) bool {
		if symmetry == "asymmetric" && (interval.isLesserMajorSecond() || interval.isLesserMinorSeventh()) {
			return true
		}
		if symmetry == "symmetric1" && (interval.isLesserMajorSecond() || interval.isGreaterMinorSeventh()) {
			return true
		}
		if symmetry == "symmetric2" && (interval.isGreaterMajorSecond() || interval.isLesserMinorSeventh()) {
			return true
		}
		return false
	}
}

func computeJustIntervals(multiplierList [][]uint, filter intervalFilterFunction) []Interval {
	var intervals []Interval
	for _, multiplier := range multiplierList {
		interval := Interval{Numerator: multiplier[0], Denominator: multiplier[1]}.octaveReduce()
		if interval.isUnison() || interval.isDiminishedFifth() {
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
	var bar [][]uint
	for _, multiplierA := range multiplierListA {
		for _, multiplierB := range multiplierListB {
			bar = append(bar, []uint{multiplierA[0] * multiplierB[0], multiplierA[1] * multiplierB[1]})
		}
	}
	return bar
}

// Modified function to produce 7-limit just intonation frets.
// It reuses computeJustIntervals by building multiplier lists that include 5 and 7 factors.
func (h Handler) fretPlacementsFor7LimitJustChromaticScaleBasedOnPureRatios(scaleLength float64, symmetry string) FretPlacements {
	return FretPlacements{
		System:      "7-limit Just Intonation",
		Description: "Fret positions for chromatic scale based on 7-limit just intonation pure ratios derived from third-, fifth- and seventh-partial ratios.",
		Frets:       h.intervalsToFretPlacements(scaleLength, computeSevenLimitJustScale()),
	}
}

func computeSevenLimitJustScale() []Interval {
	thirdPartialMultipliers := [][]uint{{1, 9}, {1, 3}, {1, 1}, {3, 1}, {9, 1}}
	fifthPartialMultipliers := [][]uint{{1, 5}, {1, 1}, {5, 1}}
	seventhPartialMultipliers := [][]uint{{1, 7}, {1, 1}, {7, 1}}

	sevenLimitMultipliers := createMultiplierTableOf(createMultiplierTableOf(seventhPartialMultipliers, fifthPartialMultipliers), thirdPartialMultipliers)
	sevenLimitIntervalPool := computeJustIntervals(sevenLimitMultipliers, func(interval Interval) bool { return false })

	var preferredIntervals []Interval

	for r := float64(50); r <= 1200; r += 100 {
		var intervalsInRange []Interval
		for _, interval := range sevenLimitIntervalPool {
			cents := interval.toCents()
			if cents >= r && cents < r+100 {
				intervalsInRange = append(intervalsInRange, interval)
			}
		}
		//   chosen interval is the simplest integer ratio
		var chosenInterval Interval
		chosenInterval = intervalsInRange[0]
		for _, interval := range intervalsInRange {
			if interval.Numerator < chosenInterval.Numerator && interval.Denominator < chosenInterval.Denominator {
				chosenInterval = interval
			}
		}
		preferredIntervals = append(preferredIntervals, chosenInterval)

	}

	return preferredIntervals
}

func (h Handler) fretPlacementsForPtolemysIntenseDiatonicTuning(scaleLength float64, octaves int, mode string) FretPlacements {
	var intervalMap = map[string][]Interval{
		"Lydian":     {greaterMajorSecond, lesserMajorSecond, greaterMajorSecond, diatonicSemitone, lesserMajorSecond, greaterMajorSecond, diatonicSemitone},
		"Ionian":     {greaterMajorSecond, lesserMajorSecond, diatonicSemitone, greaterMajorSecond, lesserMajorSecond, greaterMajorSecond, diatonicSemitone},
		"Mixolydian": {greaterMajorSecond, lesserMajorSecond, diatonicSemitone, greaterMajorSecond, lesserMajorSecond, diatonicSemitone, greaterMajorSecond},
		"Dorian":     {greaterMajorSecond, diatonicSemitone, lesserMajorSecond, greaterMajorSecond, lesserMajorSecond, diatonicSemitone, greaterMajorSecond},
		"Aeolian":    {greaterMajorSecond, diatonicSemitone, lesserMajorSecond, greaterMajorSecond, diatonicSemitone, greaterMajorSecond, lesserMajorSecond},
		"Phrygian":   {diatonicSemitone, greaterMajorSecond, lesserMajorSecond, greaterMajorSecond, diatonicSemitone, lesserMajorSecond, greaterMajorSecond},
		"Locrian":    {diatonicSemitone, greaterMajorSecond, lesserMajorSecond, diatonicSemitone, greaterMajorSecond, lesserMajorSecond, greaterMajorSecond},
	}

	var intervals []Interval
	var interval = unison

	for i := 0; i < octaves; i++ {
		for _, v := range intervalMap[mode] {
			interval = Interval{Numerator: interval.Numerator * v.Numerator, Denominator: interval.Denominator * v.Denominator}.simplify()
			intervals = append(intervals, interval)
		}
	}

	return FretPlacements{
		System:      "Ptolemy",
		Description: fmt.Sprintf("Fret positions for Ptolemy's 5-limit intense diatonic scale in %s mode.", mode),
		Frets:       h.intervalsToFretPlacements(scaleLength, intervals),
	}
}

func (h Handler) fretPlacementsForSazTuning(scaleLength float64) FretPlacements {
	return FretPlacements{
		System:      "saz",
		Description: "Fret positions for traditional Turkish Saz tuning ratios.",
		ScaleLength: scaleLength,
		Frets:       h.intervalsToFretPlacements(scaleLength, sazIntervals),
	}
}

func (h Handler) fretPlacementsForQuarterCommaMeantoneTuning(scaleLength float64, extendScale bool) FretPlacements {
	fractionOfSyntonicCommaToTemperFifthsBy := 0.25
	temperedFifth := perfectFifth.toFloat() * math.Pow(syntonicComma.toFloat(), -fractionOfSyntonicCommaToTemperFifthsBy)

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
		ratiosOfNotesToFundamental = append(ratiosOfNotesToFundamental, h.octaveReduceFloat(math.Pow(temperedFifth, float64(i))))
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

func (h Handler) intervalsToFretPlacements(scaleLength float64, intervals []Interval) []Fret {
	var frets []Fret
	var previousInterval = unison
	for _, intervalOfNote := range intervals {
		frets = append(frets, Fret{
			Label:    intervalOfNote.String(),
			Position: math.Round((scaleLength-(scaleLength/float64(intervalOfNote.Numerator))*float64(intervalOfNote.Denominator))*100) / 100,
			Comment:  intervalOfNote.name(),
			Interval: intervalOfNote.subtract(previousInterval).String(),
		})
		previousInterval = intervalOfNote
	}
	return frets
}

func (h Handler) fretPlacementsForEqualTemperamentTuning(scaleLength float64, divisionsOfOctave int) FretPlacements {
	fretPlacements := FretPlacements{
		System:      fmt.Sprintf("%d-TET", divisionsOfOctave),
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

func (h Handler) octaveReduceFloat(ratio float64) float64 {
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
		temperedFifths[i] = 3.0 / 2.0 * math.Pow(syntonicComma.toFloat(), temperingFractions[i])
	}

	// Derive ratios by walking the circle of fifths
	ratios := make([]float64, 12)
	ratios[0] = 1.0 // Start with the tonic
	for i := 1; i < 12; i++ {
		ratios[i] = ratios[i-1] * temperedFifths[(i-1)%12]
	}

	// Reduce ratios to within the octave [1.0, 2.0)
	for i := range ratios {
		ratios[i] = h.octaveReduceFloat(ratios[i])
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
