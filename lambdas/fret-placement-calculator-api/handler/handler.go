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
	defaultNumberOfOctaves           = 1
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

	var fretboard music.Fretboard
	octaves := parseIntegerQueryParameter(q, "octaves", defaultNumberOfOctaves)

	switch q["tuningSystem"] {
	case "equal":
		fretboard = music.NewFretboardFromTemperedScale(scaleLength, octaves, music.NewEqualTemperamentScale(uint(parseIntegerQueryParameter(q, "divisions", defaultEqualTemperamentDivisions))))
	case "saz":
		fretboard = music.NewFretboardFromJustScale(scaleLength, octaves, music.NewSazScale())
	case "pythagorean":
		fretboard = music.NewFretboardFromJustScale(scaleLength, octaves, music.NewPythagoreanScale())
	case "meantone":
		fretboard = music.NewFretboardFromTemperedScale(scaleLength, octaves, music.NewQuarterCommaMeantoneScale())
	case "extendedMeantone":
		fretboard = music.NewFretboardFromTemperedScale(scaleLength, octaves, music.NewExtendedQuarterCommaMeantoneScale())
	case "ptolemy":
		fretboard = music.NewFretboardFromJustScale(scaleLength, octaves, music.NewIntenseDiatonicScale(music.MusicalMode(validDiatonicModeOrDefault(q["diatonicMode"]))))
	case "just5limitFromPythagorean":
		fretboard = music.NewFretboardFromJustScale(scaleLength, octaves, music.New5LimitPythagoreanScale())
	case "justFromRatios":
		fretboard = music.NewFretboardFromJustScale(scaleLength, octaves, music.NewJustIntonationChromaticScaleWithLimit(parseIntegerQueryParameter(q, "limit", defaultJustLimit)))
	case "bachWellTemperament":
		fretboard = music.NewFretboardFromTemperedScale(scaleLength, octaves, music.NewBachWohltemperierteKlavierScale())
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
