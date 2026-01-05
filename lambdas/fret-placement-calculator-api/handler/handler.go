package handler

import (
	"context"
	"encoding/json"
	"main/music"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

const (
	defaultEqualTemperamentDivisions = 31
	defaultJustSymmetry              = music.Asymmetric
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
		fretboard = h.fretPlacementsFor5LimitJustChromaticScaleBasedOnPureRatios(scaleLength, music.Symmetry(parseStringQueryParameter(q, "justSymmetry", defaultJustSymmetry.String())))

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
	if mode == "" || !music.MusicalMode(mode).IsDiatonic() {
		mode = music.Ionian.String()
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

func (h Handler) fretPlacementsForPythagoreanTuning(scaleLength float64) Fretboard {
	return newFretboardFromScale(scaleLength, music.NewPythagoreanScale())
}

func (h Handler) fretPlacementsFor5LimitJustChromaticTuningBuiltFromAdjustingPythagoreanScale(scaleLength float64) Fretboard {
	// Derive scale by adjusting Pythagorean scale by syntonic comma (81/80)
	return newFretboardFromScale(scaleLength, music.New5LimitPythagoreanScale())
}

func (h Handler) fretPlacementsFor5LimitJustChromaticScaleBasedOnPureRatios(scaleLength float64, symmetry music.Symmetry) Fretboard {
	return newFretboardFromScale(scaleLength, music.New5LimitJustIntonationChromaticScale(symmetry))
}

func (h Handler) fretPlacementsFor7LimitJustChromaticScaleBasedOnPureRatios(scaleLength float64) Fretboard {
	return newFretboardFromScale(scaleLength, music.New7LimitJustIntonationChromaticScale())
}

func (h Handler) fretPlacementsFor13LimitJustChromaticScaleBasedOnPureRatios(scaleLength float64, limit int) Fretboard {
	return newFretboardFromScale(scaleLength, music.New13LimitJustIntonationChromaticScale())
}

func (h Handler) fretPlacementsForPtolemysIntenseDiatonicTuning(scaleLength float64, octaves int, mode string) Fretboard {
	return newFretboardFromScale(scaleLength, music.NewIntenseDiatonicScale(music.MusicalMode(mode)))
}

func (h Handler) fretPlacementsForSazTuning(scaleLength float64) Fretboard {
	return newFretboardFromScale(scaleLength, music.NewSazScale())
}

func (h Handler) fretPlacementsForQuarterCommaMeantoneTuning(scaleLength float64, extendScale bool) Fretboard {
	if extendScale {
		return newFretboardFromScale(scaleLength, music.NewExtendedQuarterCommaMeantoneScale())
	}

	return newFretboardFromScale(scaleLength, music.NewQuarterCommaMeantoneScale())
}

func (h Handler) fretPlacementsForEqualTemperamentTuning(scaleLength float64, divisionsOfOctave int) Fretboard {
	return newFretboardFromScale(scaleLength, music.NewEqualTemperamentScale(uint(divisionsOfOctave)))
}

func (h Handler) fretPlacementsForBachWohltemperierteKlavier(scaleLength float64) Fretboard {
	return newFretboardFromScale(scaleLength, music.NewBachWohltemperierteKlavierScale())
}
