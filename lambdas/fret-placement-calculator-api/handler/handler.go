package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mikebharris/music"
)

const (
	defaultEqualTemperamentDivisions = 31
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
	octaves := parseIntegerQueryParameter(q, "octaves", 1)

	switch q["tuningSystem"] {
	case "equal":
		fretboard = newFretboardFromTemperedScale(scaleLength, octaves, music.NewEqualTemperamentScale(uint(parseIntegerQueryParameter(q, "divisions", defaultEqualTemperamentDivisions))))
	case "saz":
		fretboard = newFretboardFromJustScale(scaleLength, octaves, music.NewSazScale())
	case "pythagorean":
		fretboard = newFretboardFromJustScale(scaleLength, octaves, music.NewPythagoreanScale())
	case "meantone":
		fretboard = newFretboardFromTemperedScale(scaleLength, octaves, music.NewQuarterCommaMeantoneScale())
	case "extendedMeantone":
		fretboard = newFretboardFromTemperedScale(scaleLength, octaves, music.NewExtendedQuarterCommaMeantoneScale())
	case "ptolemy":
		fretboard = newFretboardFromJustScale(scaleLength, octaves, music.NewIntenseDiatonicScale(music.MusicalMode(validDiatonicModeOrDefault(q["diatonicMode"]))))
	case "just5limitFromPythagorean":
		fretboard = newFretboardFromJustScale(scaleLength, octaves, music.New5LimitPythagoreanScale())
	case "justFromRatios":
		fretboard = newFretboardFromJustScale(scaleLength, octaves, music.NewJustIntonationChromaticScaleWithLimit(parseIntegerQueryParameter(q, "limit", defaultJustLimit)))
	case "bachWellTemperament":
		fretboard = newFretboardFromTemperedScale(scaleLength, octaves, music.NewBachWohltemperierteKlavierScale())
	default:
		return errorResponse(http.StatusUnprocessableEntity, `{"error":"please provide a valid tuning system"}`), nil
	}

	fretboard.ScaleLength = scaleLength
	body, _ := json.Marshal(fretboard)
	return events.LambdaFunctionURLResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(body)}, nil
}

func validDiatonicModeOrDefault(mode string) string {
	if mode == "" || !music.MusicalMode(mode).IsDiatonic() {
		mode = music.IonianMode.String()
	}
	return mode
}

func errorResponse(status int, body string) events.LambdaFunctionURLResponse {
	return events.LambdaFunctionURLResponse{StatusCode: status, Headers: headers, Body: body}
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
