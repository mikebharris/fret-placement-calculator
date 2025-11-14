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

var headers = map[string]string{
	"Content-Type": "application/json",
}

type Handler struct {
}

type Fret struct {
	Label    string  `json:"label"`
	Position float64 `json:"position"`
	Comment  string  `json:"comment,omitempty"`
}

type Fretting struct {
	System      string  `json:"system"`
	Description string  `json:"description, omitempty"`
	ScaleLength float64 `json:"scaleLength"`
	Frets       []Fret  `json:"frets"`
}

func (h Handler) HandleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	if request.QueryStringParameters["scaleLength"] == "" {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"scaleLength query parameter is required"}`}, nil
	}

	scaleLength, _ := strconv.ParseFloat(request.QueryStringParameters["scaleLength"], 64)

	var fretting = Fretting{
		ScaleLength: scaleLength,
	}
	var frets []Fret
	if request.QueryStringParameters["temper"] == "" {
		var ratios = [][]uint{{16, 15}, {10, 9}, {9, 8}, {6, 5}, {5, 4}, {35, 25}, {4, 3}, {45, 32}, {3, 2}, {8, 5}, {5, 3}, {16, 9}, {9, 5}, {15, 8}, {2, 1}}

		for _, ratio := range ratios {
			distanceFromNut := math.Round((scaleLength-(scaleLength/float64(ratio[0]))*float64(ratio[1]))*100) / 100
			frets = append(frets, Fret{
				Label:    fmt.Sprintf("%d:%d", ratio[0], ratio[1]),
				Position: distanceFromNut,
			})
		}
		fretting.Frets = frets
		fretting.System = "ji"
	}

	if request.QueryStringParameters["temper"] == "meantone" {
		fretting.System = "meantone"

		var fifthTemperedByFractionOfSyntonicComma float64
		if request.QueryStringParameters["temper-by"] != "" {
			fifthTemperedByFractionOfSyntonicComma, _ = strconv.ParseFloat(request.QueryStringParameters["temper-by"], 64)
		} else {
			fifthTemperedByFractionOfSyntonicComma = 0.25
		}

		extended := false
		if request.QueryStringParameters["extended"] != "" {
			extended, _ = strconv.ParseBool(request.QueryStringParameters["extended"])
		}

		if extended {
			fretting.Description = fmt.Sprintf("Calculating extended meantone based on narrowing of fifths by %.2f of a syntonic comma (81/80).  Nominal note names used based on a tonic of D.", fifthTemperedByFractionOfSyntonicComma)
		} else {
			fretting.Description = fmt.Sprintf("Calculating meantone based on narrowing of fifths by %.2f of a syntonic comma (81/80).  Nominal note names used based on a tonic of D.", fifthTemperedByFractionOfSyntonicComma)
		}

		syntonicComma := 81.0 / 80.0
		temperedFifth := 3.0 / 2.0 * math.Pow(syntonicComma, -fifthTemperedByFractionOfSyntonicComma)

		var fifthsFromTonic int
		var noteNames []string

		if extended {
			fifthsFromTonic = 9
			noteNames = []string{"D", "D#", "Eb", "E", "Fb", "F", "F#", "Gb", "G", "G#", "Ab", "A", "A#", "Bb", "B", "Cb", "C", "C#", "Db"}
		} else {
			fifthsFromTonic = 6
			noteNames = []string{"D", "Eb", "E", "F", "F#", "G", "G#", "Ab", "A", "Bb", "B", "C", "C#"}
		}

		var ratiosOfNotesToFundamental []float64
		for i := -fifthsFromTonic; i <= fifthsFromTonic; i++ {
			ratiosOfNotesToFundamental = append(ratiosOfNotesToFundamental, octaveReduce(math.Pow(temperedFifth, float64(i))))
		}

		slices.Sort(ratiosOfNotesToFundamental)

		prevRatio := 1.0
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
	}

	if request.QueryStringParameters["temper"] != "pythagorean" {

	}

	if request.QueryStringParameters["temper"] == "equal" {
		divisionsOfOctave := 12
		if request.QueryStringParameters["divisions"] != "" {
			divisionsOfOctave, _ = strconv.Atoi(request.QueryStringParameters["divisions"])
		}
		fretting.System = fmt.Sprintf("%d-TET", divisionsOfOctave)
		fretting.Description = fmt.Sprintf("Calculating equal temperament (%s).", fretting.System)

		for fretNumber := 1; fretNumber <= divisionsOfOctave; fretNumber++ {
			distanceFromNut, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", scaleLength-(scaleLength/math.Exp2(float64(fretNumber)/float64(divisionsOfOctave)))), 64)
			frets = append(frets, Fret{
				Label:    fmt.Sprintf("Fret %d", fretNumber),
				Position: distanceFromNut,
			})
		}
	}

	fretting.Frets = frets
	body, _ := json.Marshal(fretting)

	return events.LambdaFunctionURLResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(body)}, nil

}

func octaveReduce(ratio float64) float64 {
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
