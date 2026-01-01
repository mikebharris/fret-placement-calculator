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
	defaultJustLimit                 = 5
)

var headers = map[string]string{
	"Content-Type": "application/json",
}

type Handler struct {
}

func (h Handler) HandleRequest(_ context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	q := request.QueryStringParameters

	scaleLength, err := strconv.ParseFloat(q["scaleLength"], 64)
	if err != nil || scaleLength <= 0 {
		return errorResponse(http.StatusUnprocessableEntity, `{"error":"a numeric scaleLength greater than zero is required"}`), nil
	}

	var fretboard Fretboard

	switch q["tuningSystem"] {
	case "equal":
		fretboard = h.fretPlacementsForEqualTemperamentTuning(scaleLength, parseIntegerQueryParameter(q, "divisions", defaultEqualTemperamentDivisions))

	case "saz":
		fretboard = h.fretPlacementsForSazTuning(scaleLength)

	case "pythagorean":
		fretboard = h.fretPlacementsForPythagoreanTuning(scaleLength)

	case "meantone":
		fretboard = h.fretPlacementsForQuarterCommaMeantoneTuning(scaleLength, false)
	case "extendedMeantone":
		fretboard = h.fretPlacementsForQuarterCommaMeantoneTuning(scaleLength, true)

	case "":
		fallthrough
	case "ptolemy":
		fretboard = h.fretPlacementsForPtolemysIntenseDiatonicTuning(scaleLength, parseIntegerQueryParameter(q, "octaves", defaultOctavesToCompute), validDiatonicModeOrDefault(q["diatonicMode"]))

	case "just5limitFromPythagorean":
		fretboard = h.fretPlacementsFor5LimitJustChromaticTuningBuiltFromAdjustingPythagoreanScale(scaleLength)

	case "just5limitFromRatios":
		fretboard = h.fretPlacementsFor5LimitJustChromaticScaleBasedOnPureRatios(scaleLength, parseStringQueryParameter(q, "justSymmetry", defaultJustSymmetry))

	case "just7limitFromRatios":
		fretboard = h.fretPlacementsFor7LimitJustChromaticScaleBasedOnPureRatios(scaleLength)

	case "just13limitFromRatios":
		fretboard = h.fretPlacementsFor13LimitJustChromaticScaleBasedOnPureRatios(scaleLength, parseIntegerQueryParameter(q, "limit", defaultJustLimit))

	case "bachWellTemperament":
		fretboard = h.fretPlacementsForBachWohltemperierteKlavier(scaleLength)

	default:
		return errorResponse(http.StatusUnprocessableEntity, `{"error":"please provide a valid tuning system"}`), nil
	}

	fretboard.ScaleLength = scaleLength
	body, _ := json.Marshal(fretboard)
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

func computePythagoreanNotes() []Note {
	var thirdPartialsFromTonicToCompute = 6
	var notes []Note
	for i := -thirdPartialsFromTonicToCompute; i <= thirdPartialsFromTonicToCompute; i++ {
		if i < 0 {
			var note = Note{DistanceFromTonic: perfectFifth.toPowerOf(i).reciprocal().octaveReduce()}
			notes = append(notes, note)
		} else if i > 0 {
			var note = Note{DistanceFromTonic: perfectFifth.toPowerOf(i).octaveReduce()}
			notes = append(notes, note)
		}
	}

	notes = append(notes, Note{DistanceFromTonic: octave})
	slices.SortFunc(notes, func(i, j Note) int {
		return i.sortWith(j)
	})
	return notes
}

func (h Handler) fretPlacementsForPythagoreanTuning(scaleLength float64) Fretboard {
	return newFretboardFromScale(scaleLength, Scale{System: "Pythagorean", Description: "3-limit Pythagorean ratios.", Algorithm: computePythagoreanNotes})
}

func (h Handler) fretPlacementsFor5LimitJustChromaticTuningBuiltFromAdjustingPythagoreanScale(scaleLength float64) Fretboard {
	// Derive scale by adjusting Pythagorean scale by syntonic comma (81/80)
	return newFretboardFromScale(scaleLength, Scale{
		System:      "5-limit Just Intonation",
		Description: "5-limit just intonation pure ratios chromatic scale derived from applying syntonic comma to Pythagorean ratios.",
		Algorithm:   computeFiveLimitNotesFromPythagorean,
	})
}

func computeFiveLimitNotesFromPythagorean() []Note {
	var notes []Note
	for _, note := range computePythagoreanNotes() {
		if note.DistanceFromTonic.isPerfect() {
			notes = append(notes, note)
			continue
		}

		graveRatio := note.DistanceFromTonic.add(acuteUnison)
		acuteRatio := note.DistanceFromTonic.add(graveUnison)

		if graveRatio.Denominator < acuteRatio.Denominator {
			notes = append(notes, Note{DistanceFromTonic: graveRatio})
		} else {
			notes = append(notes, Note{DistanceFromTonic: acuteRatio})
		}
	}
	return notes
}

func (h Handler) fretPlacementsFor5LimitJustChromaticScaleBasedOnPureRatios(scaleLength float64, symmetry string) Fretboard {
	return Fretboard{
		System:      "5-limit Just Intonation",
		Description: "Fret positions for chromatic scale based on 5-limit just intonation pure ratios derived from third- and fifth-partial ratios.",
		Frets:       h.intervalsToFretPlacements(scaleLength, computeJustScale(buildMultiplierTablesFrom(multipliers(3), multipliers(5), multipliers(9)), fiveLimitScaleFilter(symmetry))),
	}
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

func nullScaleFilter() func(interval Interval) bool {
	return func(interval Interval) bool {
		return false
	}
}

func (h Handler) fretPlacementsFor7LimitJustChromaticScaleBasedOnPureRatios(scaleLength float64) Fretboard {
	return Fretboard{
		System:      "7-limit Just Intonation",
		Description: "Fret positions for chromatic scale based on 7-limit just intonation pure ratios derived from third-, fifth- and seventh-partial ratios.",
		Frets:       h.intervalsToFretPlacements(scaleLength, computeJustScale(buildMultiplierTablesFrom(multipliers(3), multipliers(5), multipliers(9), multipliers(7)), nullScaleFilter())),
	}
}

func (h Handler) fretPlacementsFor13LimitJustChromaticScaleBasedOnPureRatios(scaleLength float64, limit int) Fretboard {
	return Fretboard{
		System:      fmt.Sprintf("%d-limit Just Intonation", limit),
		Description: fmt.Sprintf("Fret positions for just intonation chromatic scale based on %d-limit pure ratios.", limit),
		Frets:       h.intervalsToFretPlacements(scaleLength, computeJustScale(buildMultiplierTablesFrom(multipliers(3), multipliers(5), multipliers(9), multipliers(7), multipliers(13)), nullScaleFilter())),
	}
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

func (h Handler) fretPlacementsForPtolemysIntenseDiatonicTuning(scaleLength float64, octaves int, mode string) Fretboard {
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

	return Fretboard{
		System:      "Ptolemy",
		Description: fmt.Sprintf("Fret positions for Ptolemy's 5-limit intense diatonic scale in %s mode.", mode),
		Frets:       h.intervalsToFretPlacements(scaleLength, intervals),
	}
}

func (h Handler) fretPlacementsForSazTuning(scaleLength float64) Fretboard {
	return Fretboard{
		System:      "saz",
		Description: "Fret positions for traditional Turkish Saz tuning ratios.",
		ScaleLength: scaleLength,
		Frets:       h.intervalsToFretPlacements(scaleLength, sazIntervals),
	}
}

func (h Handler) fretPlacementsForQuarterCommaMeantoneTuning(scaleLength float64, extendScale bool) Fretboard {
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

	return Fretboard{
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

func (h Handler) fretPlacementsForEqualTemperamentTuning(scaleLength float64, divisionsOfOctave int) Fretboard {
	fretPlacements := Fretboard{
		System:      fmt.Sprintf("%d-TET", divisionsOfOctave),
		Description: fmt.Sprintf("Fret positions for %d-tone equal temperament.", divisionsOfOctave),
		ScaleLength: scaleLength,
		Frets:       nil,
	}

	edo := EquallyDividedOctave{
		NumberOfDivisions: uint(divisionsOfOctave),
	}.divisions()

	for f, d := range edo {
		distanceFromNut, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", scaleLength-(scaleLength/d.Ratio)), 64)
		fretPlacements.Frets = append(fretPlacements.Frets, Fret{
			Label:    fmt.Sprintf("Fret %d", f+1),
			Position: distanceFromNut,
			Comment:  fmt.Sprintf("%.2f cents", d.Cents),
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
func (h Handler) fretPlacementsForBachWohltemperierteKlavier(scaleLength float64) Fretboard {
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
	return Fretboard{
		System:      "Bach's Well-Tempered Tuning",
		Description: description,
		ScaleLength: scaleLength,
		Frets:       frets,
	}
}

func buildMultiplierTablesFrom(multipliers ...[][]uint) [][]uint {
	if len(multipliers) == 1 {
		return multipliers[0]
	}
	return createMultiplierTableOf(multipliers[0], buildMultiplierTablesFrom(multipliers[1:]...))
}
