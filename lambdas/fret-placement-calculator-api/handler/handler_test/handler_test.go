package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"main/lambdas/fret-placement-calculator-api/handler"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

var headers = map[string]string{"Content-Type": "application/json"}

func Test_shouldReturnFretPlacementsForAChromaticFiveLimitJustIntonationScaleBuiltFromPythagoreanRatios(t *testing.T) {
	// Given
	// Mike read the page at https://www.microtonaltheory.com/tuning-theory/five-limit-just-intonation

	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "just5limitFromPythagorean"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard handler.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "5-limit Just Intonation", fretboard.System)
	assert.Equal(t, "Fret positions based on 5-limit just intonation pure ratios chromatic scale derived from applying syntonic comma to Pythagorean ratios.", fretboard.Description)
	assert.Equal(t, 13, len(fretboard.Frets))

	assert.Equal(t, "16:15", fretboard.Frets[0].Label)
	assert.Equal(t, 33.75, fretboard.Frets[0].Position)
	assert.Equal(t, "Minor Second", fretboard.Frets[0].Comment)
	assert.Equal(t, "16:15", fretboard.Frets[0].Interval)
	assert.Equal(t, "10:9", fretboard.Frets[1].Label) // not 9:8
	assert.Equal(t, 54.0, fretboard.Frets[1].Position)
	assert.Equal(t, "25:24", fretboard.Frets[1].Interval)
	assert.Equal(t, "Just (Lesser) Major Second", fretboard.Frets[1].Comment)
	assert.Equal(t, "6:5", fretboard.Frets[2].Label)
	assert.Equal(t, "27:25", fretboard.Frets[2].Interval)
	assert.Equal(t, "5:4", fretboard.Frets[3].Label)
	assert.Equal(t, "25:24", fretboard.Frets[3].Interval)
	assert.Equal(t, 108.0, fretboard.Frets[3].Position)
	assert.Equal(t, "4:3", fretboard.Frets[4].Label)
	assert.Equal(t, "16:15", fretboard.Frets[4].Interval)
	assert.Equal(t, 135.0, fretboard.Frets[4].Position)
	assert.Equal(t, "64:45", fretboard.Frets[5].Label)
	assert.Equal(t, "16:15", fretboard.Frets[5].Interval)
	assert.Equal(t, "45:32", fretboard.Frets[6].Label)
	assert.Equal(t, "3:2", fretboard.Frets[7].Label)
	assert.Equal(t, "16:15", fretboard.Frets[7].Interval)
	assert.Equal(t, 180.0, fretboard.Frets[7].Position)
	assert.Equal(t, "8:5", fretboard.Frets[8].Label)
	assert.Equal(t, "16:15", fretboard.Frets[8].Interval)
	assert.Equal(t, "5:3", fretboard.Frets[9].Label)
	assert.Equal(t, "25:24", fretboard.Frets[9].Interval)
	assert.Equal(t, 216.0, fretboard.Frets[9].Position)
	assert.Equal(t, "9:5", fretboard.Frets[10].Label) // not 16:9
	assert.Equal(t, "27:25", fretboard.Frets[10].Interval)
	assert.Equal(t, "Just (Greater) Minor Seventh", fretboard.Frets[10].Comment)
	assert.Equal(t, "15:8", fretboard.Frets[11].Label)
	assert.Equal(t, "25:24", fretboard.Frets[11].Interval)
	assert.Equal(t, 252.0, fretboard.Frets[11].Position)
	assert.Equal(t, "2:1", fretboard.Frets[12].Label)
	assert.Equal(t, "16:15", fretboard.Frets[12].Interval)
	assert.Equal(t, 270.0, fretboard.Frets[12].Position)
}

func Test_shouldReturnFretPlacementsForAsymmetricJustChromaticScaleBasedOnPureRatios(t *testing.T) {
	// Given
	// Mike read the page at https://en.wikipedia.org/wiki/Five-limit_tuning

	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "just5limitFromRatios", "justSymmetry": "asymmetric"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard handler.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "5-limit Just Intonation", fretboard.System)
	assert.Equal(t, "Fret positions for chromatic scale based on 5-limit just intonation pure ratios derived from third- and fifth-partial ratios.", fretboard.Description)
	assert.Equal(t, 12, len(fretboard.Frets))

	assert.Equal(t, "16:15", fretboard.Frets[0].Label)
	assert.Equal(t, "9:8", fretboard.Frets[1].Label)
	assert.Equal(t, "6:5", fretboard.Frets[2].Label)
	assert.Equal(t, "5:4", fretboard.Frets[3].Label)
	assert.Equal(t, "4:3", fretboard.Frets[4].Label)
	assert.Equal(t, "45:32", fretboard.Frets[5].Label)
	assert.Equal(t, "3:2", fretboard.Frets[6].Label)
	assert.Equal(t, "8:5", fretboard.Frets[7].Label)
	assert.Equal(t, "5:3", fretboard.Frets[8].Label)
	assert.Equal(t, "9:5", fretboard.Frets[9].Label)
	assert.Equal(t, "15:8", fretboard.Frets[10].Label)
	assert.Equal(t, "2:1", fretboard.Frets[11].Label)
}

func Test_shouldReturnFretPlacementsForSymmetricJustChromaticScaleWithLesserMajorSecondBasedOnPureRatios(t *testing.T) {
	// Given
	// Mike read the page at https://en.wikipedia.org/wiki/Five-limit_tuning

	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "just5limitFromRatios", "justSymmetry": "symmetric1"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard handler.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "5-limit Just Intonation", fretboard.System)
	assert.Equal(t, "Fret positions for chromatic scale based on 5-limit just intonation pure ratios derived from third- and fifth-partial ratios.", fretboard.Description)
	assert.Equal(t, 12, len(fretboard.Frets))

	assert.Equal(t, "9:8", fretboard.Frets[1].Label)
	assert.Equal(t, "16:9", fretboard.Frets[9].Label)
}

func Test_shouldReturnFretPlacementsForSymmetricJustChromaticScaleWithGreaterMajorSecondBasedOnPureRatios(t *testing.T) {
	// Given
	// Mike read the page at https://en.wikipedia.org/wiki/Five-limit_tuning

	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "just5limitFromRatios", "justSymmetry": "symmetric2"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard handler.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "5-limit Just Intonation", fretboard.System)
	assert.Equal(t, "Fret positions for chromatic scale based on 5-limit just intonation pure ratios derived from third- and fifth-partial ratios.", fretboard.Description)
	assert.Equal(t, 12, len(fretboard.Frets))

	assert.Equal(t, "10:9", fretboard.Frets[1].Label)
	assert.Equal(t, "9:5", fretboard.Frets[9].Label)
}

func Test_shouldReturnFretPlacementsForAsymmetric7LimitJustChromaticScaleBasedOnPureRatios(t *testing.T) {
	// Given
	// Mike read the page at https://en.wikipedia.org/wiki/Five-limit_tuning

	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "just7limitFromRatios"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard handler.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "7-limit Just Intonation", fretboard.System)
	assert.Equal(t, "Fret positions for chromatic scale based on 7-limit just intonation pure ratios derived from third-, fifth- and seventh-partial ratios.", fretboard.Description)
	assert.Equal(t, 12, len(fretboard.Frets))

	assert.Equal(t, "15:14", fretboard.Frets[0].Label)
	assert.Equal(t, "8:7", fretboard.Frets[1].Label)
	assert.Equal(t, "6:5", fretboard.Frets[2].Label)
	assert.Equal(t, "5:4", fretboard.Frets[3].Label)
	assert.Equal(t, "4:3", fretboard.Frets[4].Label)
	assert.Equal(t, "7:5", fretboard.Frets[5].Label)
	assert.Equal(t, "3:2", fretboard.Frets[6].Label)
	assert.Equal(t, "8:5", fretboard.Frets[7].Label)
	assert.Equal(t, "5:3", fretboard.Frets[8].Label)
	assert.Equal(t, "7:4", fretboard.Frets[9].Label)
	assert.Equal(t, "15:8", fretboard.Frets[10].Label)
	assert.Equal(t, "2:1", fretboard.Frets[11].Label)
}

func Test_shouldReturnFretPlacementsForAsymmetric13LimitJustChromaticScaleBasedOnPureRatios(t *testing.T) {
	// Given
	// Mike read the page at https://en.wikipedia.org/wiki/Five-limit_tuning

	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "just13limitFromRatios", "limit": "13"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard handler.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "13-limit Just Intonation", fretboard.System)
	assert.Equal(t, "Fret positions for just intonation chromatic scale based on 13-limit pure ratios.", fretboard.Description)
	assert.Equal(t, 12, len(fretboard.Frets))

	assert.Equal(t, "13:12", fretboard.Frets[0].Label)
	assert.Equal(t, "8:7", fretboard.Frets[1].Label)
	assert.Equal(t, "6:5", fretboard.Frets[2].Label)
	assert.Equal(t, "5:4", fretboard.Frets[3].Label)
	assert.Equal(t, "4:3", fretboard.Frets[4].Label)
	assert.Equal(t, "7:5", fretboard.Frets[5].Label)
	assert.Equal(t, "3:2", fretboard.Frets[6].Label)
	assert.Equal(t, "8:5", fretboard.Frets[7].Label)
	assert.Equal(t, "5:3", fretboard.Frets[8].Label)
	assert.Equal(t, "7:4", fretboard.Frets[9].Label)
	assert.Equal(t, "13:7", fretboard.Frets[10].Label)
	assert.Equal(t, "2:1", fretboard.Frets[11].Label)
}

func Test_shouldReturnPtolemysIntenseDiatonicScaleInTheIonianModeWhenOnlyScaleLengthIsProvided(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "ptolemy", "octaves": "2"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard handler.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "Ptolemy", fretboard.System)
	assert.Equal(t, "Fret positions for Ptolemy's 5-limit intense diatonic scale in Ionian mode.", fretboard.Description)
	assert.Equal(t, 14, len(fretboard.Frets))

	assert.Equal(t, "9:8", fretboard.Frets[0].Label)
	assert.Equal(t, "9:8", fretboard.Frets[0].Interval)
	assert.Equal(t, 60.0, fretboard.Frets[0].Position)
	assert.Equal(t, "5:4", fretboard.Frets[1].Label)
	assert.Equal(t, "10:9", fretboard.Frets[1].Interval)
	assert.Equal(t, 108.0, fretboard.Frets[1].Position)
	assert.Equal(t, "4:3", fretboard.Frets[2].Label)
	assert.Equal(t, 135.0, fretboard.Frets[2].Position)
	assert.Equal(t, "3:2", fretboard.Frets[3].Label)
	assert.Equal(t, 180.0, fretboard.Frets[3].Position)
	assert.Equal(t, "5:3", fretboard.Frets[4].Label)
	assert.Equal(t, 216.0, fretboard.Frets[4].Position)
	assert.Equal(t, "15:8", fretboard.Frets[5].Label)
	assert.Equal(t, 252.0, fretboard.Frets[5].Position)
	assert.Equal(t, "2:1", fretboard.Frets[6].Label)
	assert.Equal(t, 270.0, fretboard.Frets[6].Position)
}

func Test_shouldReturnDiatonicDorianJustIntonationPlacements(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "ptolemy", "diatonicMode": "Dorian"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard handler.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "Ptolemy", fretboard.System)
	assert.Equal(t, "Fret positions for Ptolemy's 5-limit intense diatonic scale in Dorian mode.", fretboard.Description)
	assert.Equal(t, 7, len(fretboard.Frets))

	assert.Equal(t, "9:8", fretboard.Frets[0].Label)
	assert.Equal(t, "6:5", fretboard.Frets[1].Label)
	assert.Equal(t, "4:3", fretboard.Frets[2].Label)
	assert.Equal(t, "3:2", fretboard.Frets[3].Label)
	assert.Equal(t, "5:3", fretboard.Frets[4].Label)
	assert.Equal(t, "16:9", fretboard.Frets[5].Label)
	assert.Equal(t, "2:1", fretboard.Frets[6].Label)
}

func Test_shouldReturnDiatonicPhrygianJustIntonationPlacements(t *testing.T) {
	// Given Phrygian
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "ptolemy", "diatonicMode": "Phrygian"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard handler.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "Ptolemy", fretboard.System)
	assert.Equal(t, "Fret positions for Ptolemy's 5-limit intense diatonic scale in Phrygian mode.", fretboard.Description)
	assert.Equal(t, 7, len(fretboard.Frets))

	assert.Equal(t, "16:15", fretboard.Frets[0].Label)
	assert.Equal(t, "6:5", fretboard.Frets[1].Label)
	assert.Equal(t, "4:3", fretboard.Frets[2].Label)
	assert.Equal(t, "3:2", fretboard.Frets[3].Label)
	assert.Equal(t, "8:5", fretboard.Frets[4].Label)
	assert.Equal(t, "16:9", fretboard.Frets[5].Label)
	assert.Equal(t, "2:1", fretboard.Frets[6].Label)
}

func Test_shouldReturnDiatonicLydianJustIntonationPlacements(t *testing.T) {
	// Given Phrygian
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "ptolemy", "diatonicMode": "Lydian"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard handler.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "Ptolemy", fretboard.System)
	assert.Equal(t, "Fret positions for Ptolemy's 5-limit intense diatonic scale in Lydian mode.", fretboard.Description)
	assert.Equal(t, 7, len(fretboard.Frets))

	assert.Equal(t, "9:8", fretboard.Frets[0].Label)
	assert.Equal(t, "5:4", fretboard.Frets[1].Label)
	assert.Equal(t, "45:32", fretboard.Frets[2].Label)
	assert.Equal(t, "3:2", fretboard.Frets[3].Label)
	assert.Equal(t, "5:3", fretboard.Frets[4].Label)
	assert.Equal(t, "15:8", fretboard.Frets[5].Label)
	assert.Equal(t, "2:1", fretboard.Frets[6].Label)
}

func Test_shouldReturnDiatonicMixolydianJustIntonationPlacements(t *testing.T) {
	// Given Phrygian
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "ptolemy", "diatonicMode": "Mixolydian"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard handler.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "Ptolemy", fretboard.System)
	assert.Equal(t, "Fret positions for Ptolemy's 5-limit intense diatonic scale in Mixolydian mode.", fretboard.Description)
	assert.Equal(t, 7, len(fretboard.Frets))

	assert.Equal(t, "9:8", fretboard.Frets[0].Label)
	assert.Equal(t, "5:4", fretboard.Frets[1].Label)
	assert.Equal(t, "4:3", fretboard.Frets[2].Label)
	assert.Equal(t, "3:2", fretboard.Frets[3].Label)
	assert.Equal(t, "5:3", fretboard.Frets[4].Label)
	assert.Equal(t, "16:9", fretboard.Frets[5].Label)
	assert.Equal(t, "2:1", fretboard.Frets[6].Label)
}

func Test_shouldReturnDiatonicAeolianJustIntonationPlacements(t *testing.T) {
	// Given Phrygian
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "ptolemy", "diatonicMode": "Aeolian"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard handler.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "Ptolemy", fretboard.System)
	assert.Equal(t, "Fret positions for Ptolemy's 5-limit intense diatonic scale in Aeolian mode.", fretboard.Description)
	assert.Equal(t, 7, len(fretboard.Frets))

	assert.Equal(t, "9:8", fretboard.Frets[0].Label)
	assert.Equal(t, "6:5", fretboard.Frets[1].Label)
	assert.Equal(t, "4:3", fretboard.Frets[2].Label)
	assert.Equal(t, "3:2", fretboard.Frets[3].Label)
	assert.Equal(t, "8:5", fretboard.Frets[4].Label)
	assert.Equal(t, "9:5", fretboard.Frets[5].Label)
	assert.Equal(t, "2:1", fretboard.Frets[6].Label)
}

func Test_shouldReturnDiatonicLocrianJustIntonationPlacements(t *testing.T) {
	// Given Phrygian
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "ptolemy", "diatonicMode": "Locrian"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard handler.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "Ptolemy", fretboard.System)
	assert.Equal(t, "Fret positions for Ptolemy's 5-limit intense diatonic scale in Locrian mode.", fretboard.Description)
	assert.Equal(t, 7, len(fretboard.Frets))

	assert.Equal(t, "16:15", fretboard.Frets[0].Label)
	assert.Equal(t, "6:5", fretboard.Frets[1].Label)
	assert.Equal(t, "4:3", fretboard.Frets[2].Label)
	assert.Equal(t, "64:45", fretboard.Frets[3].Label)
	assert.Equal(t, "8:5", fretboard.Frets[4].Label)
	assert.Equal(t, "16:9", fretboard.Frets[5].Label)
	assert.Equal(t, "2:1", fretboard.Frets[6].Label)
}

func Test_shouldReturnErrorWhenScaleLengthIsNotProvided(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"a numeric scaleLength greater than zero is required"}`}, response)
}

func Test_shouldReturnErrorWhenScaleLengthIsZero(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "0"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"a numeric scaleLength greater than zero is required"}`}, response)
}

func Test_shouldReturnErrorWhenScaleLengthIsLessThanZero(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "-100"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"a numeric scaleLength greater than zero is required"}`}, response)
}

func Test_shouldReturnErrorWhenScaleLengthIsNotANumber(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "three"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"a numeric scaleLength greater than zero is required"}`}, response)
}

func Test_ShouldReturnErrorWhenTuningSystemIsInvalid(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "invalid"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"please provide a valid tuning system"}`}, response)
}

func Test_ShouldDefaultToIonianIfNonSensicalPtolemicDiatonicScaleIsProvided(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "ptolemy", "diatonicMode": "Athenian"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard handler.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "Ptolemy", fretboard.System)
	assert.Equal(t, "Fret positions for Ptolemy's 5-limit intense diatonic scale in Ionian mode.", fretboard.Description)
	assert.Equal(t, 7, len(fretboard.Frets))
}

func Test_shouldReturnQuarterCommaMeantonePlacementsWithProvidedScaleLength(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "meantone"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	fretboard := handler.Fretboard{}
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, float64(540), fretboard.ScaleLength)
	assert.Equal(t, "meantone", fretboard.System)
	assert.Equal(t, "Fret positions for meantone computed by narrowing of fifths by 0.25 of a syntonic comma (81/80).  Nominal note names used given a tonic of D.", fretboard.Description)
	assert.Equal(t, 13, len(fretboard.Frets))

	greaterSemitone := 1.069984
	lesserSemitone := 1.044907
	dieses := 1.024

	assert.Equal(t, handler.Fret{
		Label:    "1 (Eb)",
		Position: 35.3,
		Comment:  fmt.Sprintf("ratio: 1.070; interval: %f", greaterSemitone),
	}, fretboard.Frets[0])

	assert.Equal(t, handler.Fret{
		Label:    "2 (E)",
		Position: 57.0,
		Comment:  fmt.Sprintf("ratio: 1.118; interval: %f", lesserSemitone),
	}, fretboard.Frets[1])

	assert.Equal(t, handler.Fret{
		Label:    "7 (Ab)",
		Position: 162.7,
		Comment:  fmt.Sprintf("ratio: 1.431; interval: %f", dieses),
	}, fretboard.Frets[6])

	assert.Equal(t, handler.Fret{
		Label:    "13 (Octave)",
		Position: 270.0,
		Comment:  fmt.Sprintf("ratio: 2.0; interval: %f", greaterSemitone),
	}, fretboard.Frets[12])
}

func Test_shouldReturnExtendedQuarterCommaMeantonePlacementsWithProvidedScaleLength(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "extendedMeantone"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	fretboard := handler.Fretboard{}
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, float64(540), fretboard.ScaleLength)
	assert.Equal(t, "meantone", fretboard.System)
	assert.Equal(t, "Fret positions for extended meantone computed by narrowing of fifths by 0.25 of a syntonic comma (81/80).  Nominal note names used given a tonic of D.", fretboard.Description)
	assert.Equal(t, 19, len(fretboard.Frets))

	lesserSemitone := 1.044907
	dieses := 1.024

	assert.Equal(t, handler.Fret{
		Label:    "1 (D#)",
		Position: 23.2,
		Comment:  fmt.Sprintf("ratio: 1.045; interval: %f", lesserSemitone),
	}, fretboard.Frets[0])

	assert.Equal(t, handler.Fret{
		Label:    "2 (Eb)",
		Position: 35.3,
		Comment:  fmt.Sprintf("ratio: 1.070; interval: %f", dieses),
	}, fretboard.Frets[1])

	assert.Equal(t, handler.Fret{
		Label:    "10 (Ab)",
		Position: 162.7,
		Comment:  fmt.Sprintf("ratio: 1.431; interval: %f", dieses),
	}, fretboard.Frets[9])

	assert.Equal(t, handler.Fret{
		Label:    "19 (Octave)",
		Position: 270.0,
		Comment:  fmt.Sprintf("ratio: 2.0; interval: %f", lesserSemitone),
	}, fretboard.Frets[18])
}

func Test_shouldReturnPythagoreanPlacementsWithProvidedScaleLength(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "pythagorean"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	fretboard := handler.Fretboard{}
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, float64(540), fretboard.ScaleLength)
	assert.Equal(t, "Pythagorean", fretboard.System)
	assert.Equal(t, "Fret positions based on 3-limit Pythagorean ratios.", fretboard.Description)
	assert.Equal(t, 13, len(fretboard.Frets))

	assert.Equal(t, "256:243", fretboard.Frets[0].Label)  // minor second
	assert.Equal(t, "9:8", fretboard.Frets[1].Label)      // major second
	assert.Equal(t, "32:27", fretboard.Frets[2].Label)    // minor third
	assert.Equal(t, "81:64", fretboard.Frets[3].Label)    // major third
	assert.Equal(t, "4:3", fretboard.Frets[4].Label)      // perfect fourth
	assert.Equal(t, "1024:729", fretboard.Frets[5].Label) // augmented fourth
	assert.Equal(t, "729:512", fretboard.Frets[6].Label)  // diminished fifth
	assert.Equal(t, "3:2", fretboard.Frets[7].Label)      // perfect fifth
	assert.Equal(t, "128:81", fretboard.Frets[8].Label)   // minor sixth
	assert.Equal(t, "27:16", fretboard.Frets[9].Label)    // major sixth
	assert.Equal(t, "16:9", fretboard.Frets[10].Label)    // minor seventh
	assert.Equal(t, "243:128", fretboard.Frets[11].Label) // major seventh
	assert.Equal(t, "2:1", fretboard.Frets[12].Label)     // octave

	assert.Equal(t, "27.42", fmt.Sprintf("%.2f", fretboard.Frets[0].Position))
	assert.Equal(t, "60.00", fmt.Sprintf("%.2f", fretboard.Frets[1].Position))
	assert.Equal(t, "84.38", fmt.Sprintf("%.2f", fretboard.Frets[2].Position))
}

func Test_ShouldDefaultTo31EqualTemperamentByDefault(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "equal"}})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	fretboard := handler.Fretboard{}
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, float64(540), fretboard.ScaleLength)
	assert.Equal(t, "31-TET", fretboard.System)
	assert.Equal(t, "Fret positions for 31-tone equal temperament.", fretboard.Description)
	assert.Equal(t, 31, len(fretboard.Frets))
}

func Test_ShouldReturnEqualTemperamentPlacementsWithCustomDivisions(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "equal", "divisions": "19"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	fretboard := handler.Fretboard{}
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, float64(540), fretboard.ScaleLength)
	assert.Equal(t, "19-TET", fretboard.System)
	assert.Equal(t, "Fret positions for 19-tone equal temperament.", fretboard.Description)
	assert.Equal(t, 19, len(fretboard.Frets))
	assert.Equal(t, handler.Fret{Label: "Fret 1", Position: 19.345, Comment: "63.00 cents"}, fretboard.Frets[0])
	assert.Equal(t, handler.Fret{Label: "Fret 19", Position: 270.0, Comment: "1200.00 cents"}, fretboard.Frets[18])
}

func Test_ShouldReturnSazPlacementsWithProvidedScaleLength(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "saz"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	fretboard := handler.Fretboard{}
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, float64(540), fretboard.ScaleLength)
	assert.Equal(t, "saz", fretboard.System)
	assert.Equal(t, "Fret positions for traditional Turkish Saz tuning ratios.", fretboard.Description)
	assert.Equal(t, 17, len(fretboard.Frets))
}

func Test_ShouldReturnFretPlacementsForBachWohltemperierteKlavierTuning(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "bachWellTemperament"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	fretboard := handler.Fretboard{}
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, float64(540), fretboard.ScaleLength)
	assert.Equal(t, "Bach's Well-Tempered Tuning", fretboard.System)
	assert.Equal(t, "Fret positions derived from Lehman's decoding of Bach's Well-Tempered tuning, using sixth-comma, twelfth-comma, and pure fifths.", fretboard.Description)
	assert.Equal(t, 12, len(fretboard.Frets))
}
