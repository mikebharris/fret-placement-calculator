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

type FretPlacements struct {
	System      string  `json:"system"`
	Description string  `json:"description,omitempty"`
	ScaleLength float64 `json:"scaleLength"`
	Frets       []Fret  `json:"frets"`
}

func (h Handler) HandleRequest(ctx context.Context, request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	scaleLength, err := strconv.ParseFloat(request.QueryStringParameters["scaleLength"], 64)
	if err != nil || scaleLength <= 0 {
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"a numeric scaleLength greater than zero is required"}`}, nil
	}

	var fretPlacements FretPlacements

	switch request.QueryStringParameters["temper"] {
	case "equal":
		divisionsOfOctave := 31
		if request.QueryStringParameters["divisions"] != "" {
			divisionsOfOctave, err = strconv.Atoi(request.QueryStringParameters["divisions"])
		}
		if err != nil {
			return events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"please provide number of divisions for equal temperament"}`}, nil
		}
		fretPlacements = h.equalTemperamentFretPlacements(scaleLength, divisionsOfOctave)
	case "saz":
		fretPlacements = h.sazFretPlacements(scaleLength)
	case "pythagorean":
		fretPlacements = h.pythagoreanFretPlacements(scaleLength)
	case "meantone":
		extended := false
		if request.QueryStringParameters["extended"] != "" {
			extended, _ = strconv.ParseBool(request.QueryStringParameters["extended"])
		}
		fretPlacements = h.quarterCommaMeantoneFretPlacements(scaleLength, extended)
	case "":
		var octaves = 1
		if request.QueryStringParameters["octaves"] != "" {
			octaves, err = strconv.Atoi(request.QueryStringParameters["octaves"])
			if err != nil {
				return events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"please provide a valid positive number for number of octaves worth of frets"}`}, nil
			}
		}
		var mode = "Ionian"
		if request.QueryStringParameters["mode"] != "" {
			mode = request.QueryStringParameters["mode"]
		}
		fretPlacements = h.justIntonationFretPlacements(scaleLength, octaves, mode)
	default:
		return events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"invalid temper parameter"}`}, nil
	}

	fretPlacements.ScaleLength = scaleLength
	body, _ := json.Marshal(fretPlacements)

	return events.LambdaFunctionURLResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(body)}, nil

}

func (h Handler) pythagoreanFretPlacements(scaleLength float64) FretPlacements {
	return FretPlacements{
		System:      "Pythagorean",
		Description: "Fret positions based on 3-limit Pythagorean ratios.",
		Frets:       h.ratiosToFretPlacements(scaleLength, [][]uint{{256, 243}, {9, 8}, {32, 27}, {81, 64}, {4, 3}, {1024, 729}, {729, 512}, {3, 2}, {128, 81}, {27, 16}, {16, 9}, {243, 128}, {2, 1}}),
	}
}

func (h Handler) justIntonationFretPlacements(scaleLength float64, octaves int, mode string) FretPlacements {
	var intervalMap = map[string][][]uint{
		"Lydian":     {{9, 8}, {10, 9}, {9, 8}, {16, 15}, {10, 9}, {9, 8}, {16, 15}},
		"Ionian":     {{9, 8}, {10, 9}, {16, 15}, {9, 8}, {10, 9}, {9, 8}, {16, 15}},
		"Mixolydian": {{9, 8}, {10, 9}, {16, 15}, {9, 8}, {10, 9}, {16, 15}, {9, 8}},
		"Dorian":     {{9, 8}, {16, 15}, {10, 9}, {9, 8}, {10, 9}, {16, 15}, {9, 8}},
		"Aeolian":    {{9, 8}, {16, 15}, {10, 9}, {9, 8}, {16, 15}, {10, 9}, {9, 8}},
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
		System:      "ji",
		Description: fmt.Sprintf("Fret positions based on 5-limit just intonation pure ratios and diatonic scale %s mode.", mode),
		Frets:       h.ratiosToFretPlacements(scaleLength, ratios),
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

func (h Handler) sazFretPlacements(scaleLength float64) FretPlacements {
	// as per https://en.wikipedia.org/wiki/Ba%C4%9Flama and the cura that I have
	return FretPlacements{
		System:      "saz",
		Description: "Fret positions for traditional Turkish Saz tuning ratios.",
		ScaleLength: scaleLength,
		Frets:       h.ratiosToFretPlacements(scaleLength, [][]uint{{18, 17}, {12, 11}, {9, 8}, {81, 68}, {27, 22}, {81, 64}, {4, 3}, {24, 17}, {16, 11}, {3, 2}, {27, 17}, {18, 11}, {27, 16}, {16, 9}, {32, 17}, {64, 33}, {2, 1}}),
	}
}

func (h Handler) quarterCommaMeantoneFretPlacements(scaleLength float64, extendScale bool) FretPlacements {
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
		ratiosOfNotesToFundamental = append(ratiosOfNotesToFundamental, octaveReduce(math.Pow(temperedFifth, float64(i))))
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
	for _, ratio := range ratios {
		distanceFromNut := math.Round((scaleLength-(scaleLength/float64(ratio[0]))*float64(ratio[1]))*100) / 100
		frets = append(frets, Fret{
			Label:    fmt.Sprintf("%d:%d", ratio[0], ratio[1]),
			Position: distanceFromNut,
		})
	}
	return frets
}

func (h Handler) equalTemperamentFretPlacements(scaleLength float64, divisionsOfOctave int) FretPlacements {
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
