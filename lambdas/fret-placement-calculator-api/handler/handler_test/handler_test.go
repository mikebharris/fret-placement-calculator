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

func Test_shouldReturnFretPlacementsForAChromaticFiveLimitJustIntonationScale(t *testing.T) {
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

	var fretPlacements handler.FretPlacements
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, 540.0, fretPlacements.ScaleLength)
	assert.Equal(t, "5-limit Just Intonation", fretPlacements.System)
	assert.Equal(t, "Fret positions for chromatic scale based on 5-limit just intonation pure ratios derived from applying syntonic comma to Pythagorean ratios.", fretPlacements.Description)
	assert.Equal(t, 13, len(fretPlacements.Frets))

	assert.Equal(t, "16:15", fretPlacements.Frets[0].Label)
	assert.Equal(t, 33.75, fretPlacements.Frets[0].Position)
	assert.Equal(t, "Minor Second", fretPlacements.Frets[0].Comment)
	assert.Equal(t, "16:15", fretPlacements.Frets[0].Interval)
	assert.Equal(t, "10:9", fretPlacements.Frets[1].Label) // not 9:8
	assert.Equal(t, 54.0, fretPlacements.Frets[1].Position)
	assert.Equal(t, "25:24", fretPlacements.Frets[1].Interval)
	assert.Equal(t, "Just (Lesser) Major Second", fretPlacements.Frets[1].Comment)
	assert.Equal(t, "6:5", fretPlacements.Frets[2].Label)
	assert.Equal(t, "27:25", fretPlacements.Frets[2].Interval)
	assert.Equal(t, "5:4", fretPlacements.Frets[3].Label)
	assert.Equal(t, "25:24", fretPlacements.Frets[3].Interval)
	assert.Equal(t, 108.0, fretPlacements.Frets[3].Position)
	assert.Equal(t, "4:3", fretPlacements.Frets[4].Label)
	assert.Equal(t, "16:15", fretPlacements.Frets[4].Interval)
	assert.Equal(t, 135.0, fretPlacements.Frets[4].Position)
	assert.Equal(t, "64:45", fretPlacements.Frets[5].Label)
	assert.Equal(t, "16:15", fretPlacements.Frets[5].Interval)
	assert.Equal(t, "45:32", fretPlacements.Frets[6].Label)
	assert.Equal(t, "3:2", fretPlacements.Frets[7].Label)
	assert.Equal(t, "16:15", fretPlacements.Frets[7].Interval)
	assert.Equal(t, 180.0, fretPlacements.Frets[7].Position)
	assert.Equal(t, "8:5", fretPlacements.Frets[8].Label)
	assert.Equal(t, "16:15", fretPlacements.Frets[8].Interval)
	assert.Equal(t, "5:3", fretPlacements.Frets[9].Label)
	assert.Equal(t, "25:24", fretPlacements.Frets[9].Interval)
	assert.Equal(t, 216.0, fretPlacements.Frets[9].Position)
	assert.Equal(t, "9:5", fretPlacements.Frets[10].Label) // not 16:9
	assert.Equal(t, "27:25", fretPlacements.Frets[10].Interval)
	assert.Equal(t, "Just (Greater) Minor Seventh", fretPlacements.Frets[10].Comment)
	assert.Equal(t, "15:8", fretPlacements.Frets[11].Label)
	assert.Equal(t, "25:24", fretPlacements.Frets[11].Interval)
	assert.Equal(t, 252.0, fretPlacements.Frets[11].Position)
	assert.Equal(t, "2:1", fretPlacements.Frets[12].Label)
	assert.Equal(t, "16:15", fretPlacements.Frets[12].Interval)
	assert.Equal(t, 270.0, fretPlacements.Frets[12].Position)
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

	var fretPlacements handler.FretPlacements
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, 540.0, fretPlacements.ScaleLength)
	assert.Equal(t, "5-limit Just Intonation", fretPlacements.System)
	assert.Equal(t, "Fret positions for chromatic scale based on 5-limit just intonation pure ratios derived from third- and forth-partial ratios.", fretPlacements.Description)
	assert.Equal(t, 12, len(fretPlacements.Frets))

	assert.Equal(t, "16:15", fretPlacements.Frets[0].Label)
	assert.Equal(t, "9:8", fretPlacements.Frets[1].Label)
	assert.Equal(t, "6:5", fretPlacements.Frets[2].Label)
	assert.Equal(t, "5:4", fretPlacements.Frets[3].Label)
	assert.Equal(t, "4:3", fretPlacements.Frets[4].Label)
	assert.Equal(t, "45:32", fretPlacements.Frets[5].Label)
	assert.Equal(t, "3:2", fretPlacements.Frets[6].Label)
	assert.Equal(t, "8:5", fretPlacements.Frets[7].Label)
	assert.Equal(t, "5:3", fretPlacements.Frets[8].Label)
	assert.Equal(t, "9:5", fretPlacements.Frets[9].Label)
	assert.Equal(t, "15:8", fretPlacements.Frets[10].Label)
	assert.Equal(t, "2:1", fretPlacements.Frets[11].Label)
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

	var fretPlacements handler.FretPlacements
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, 540.0, fretPlacements.ScaleLength)
	assert.Equal(t, "5-limit Just Intonation", fretPlacements.System)
	assert.Equal(t, "Fret positions for chromatic scale based on 5-limit just intonation pure ratios derived from third- and forth-partial ratios.", fretPlacements.Description)
	assert.Equal(t, 12, len(fretPlacements.Frets))

	assert.Equal(t, "9:8", fretPlacements.Frets[1].Label)
	assert.Equal(t, "16:9", fretPlacements.Frets[9].Label)
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

	var fretPlacements handler.FretPlacements
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, 540.0, fretPlacements.ScaleLength)
	assert.Equal(t, "5-limit Just Intonation", fretPlacements.System)
	assert.Equal(t, "Fret positions for chromatic scale based on 5-limit just intonation pure ratios derived from third- and forth-partial ratios.", fretPlacements.Description)
	assert.Equal(t, 12, len(fretPlacements.Frets))

	assert.Equal(t, "10:9", fretPlacements.Frets[1].Label)
	assert.Equal(t, "9:5", fretPlacements.Frets[9].Label)
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

	var fretPlacements handler.FretPlacements
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, 540.0, fretPlacements.ScaleLength)
	assert.Equal(t, "Ptolemy", fretPlacements.System)
	assert.Equal(t, "Fret positions for Ptolemy's 5-limit intense diatonic scale in Ionian mode.", fretPlacements.Description)
	assert.Equal(t, 14, len(fretPlacements.Frets))

	assert.Equal(t, "9:8", fretPlacements.Frets[0].Label)
	assert.Equal(t, "9:8", fretPlacements.Frets[0].Interval)
	assert.Equal(t, 60.0, fretPlacements.Frets[0].Position)
	assert.Equal(t, "5:4", fretPlacements.Frets[1].Label)
	assert.Equal(t, "10:9", fretPlacements.Frets[1].Interval)
	assert.Equal(t, 108.0, fretPlacements.Frets[1].Position)
	assert.Equal(t, "4:3", fretPlacements.Frets[2].Label)
	assert.Equal(t, 135.0, fretPlacements.Frets[2].Position)
	assert.Equal(t, "3:2", fretPlacements.Frets[3].Label)
	assert.Equal(t, 180.0, fretPlacements.Frets[3].Position)
	assert.Equal(t, "5:3", fretPlacements.Frets[4].Label)
	assert.Equal(t, 216.0, fretPlacements.Frets[4].Position)
	assert.Equal(t, "15:8", fretPlacements.Frets[5].Label)
	assert.Equal(t, 252.0, fretPlacements.Frets[5].Position)
	assert.Equal(t, "2:1", fretPlacements.Frets[6].Label)
	assert.Equal(t, 270.0, fretPlacements.Frets[6].Position)
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

	var fretPlacements handler.FretPlacements
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, 540.0, fretPlacements.ScaleLength)
	assert.Equal(t, "Ptolemy", fretPlacements.System)
	assert.Equal(t, "Fret positions for Ptolemy's 5-limit intense diatonic scale in Dorian mode.", fretPlacements.Description)
	assert.Equal(t, 7, len(fretPlacements.Frets))

	assert.Equal(t, "9:8", fretPlacements.Frets[0].Label)
	assert.Equal(t, "6:5", fretPlacements.Frets[1].Label)
	assert.Equal(t, "4:3", fretPlacements.Frets[2].Label)
	assert.Equal(t, "3:2", fretPlacements.Frets[3].Label)
	assert.Equal(t, "5:3", fretPlacements.Frets[4].Label)
	assert.Equal(t, "16:9", fretPlacements.Frets[5].Label)
	assert.Equal(t, "2:1", fretPlacements.Frets[6].Label)
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

	var fretPlacements handler.FretPlacements
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, 540.0, fretPlacements.ScaleLength)
	assert.Equal(t, "Ptolemy", fretPlacements.System)
	assert.Equal(t, "Fret positions for Ptolemy's 5-limit intense diatonic scale in Phrygian mode.", fretPlacements.Description)
	assert.Equal(t, 7, len(fretPlacements.Frets))

	assert.Equal(t, "16:15", fretPlacements.Frets[0].Label)
	assert.Equal(t, "6:5", fretPlacements.Frets[1].Label)
	assert.Equal(t, "4:3", fretPlacements.Frets[2].Label)
	assert.Equal(t, "3:2", fretPlacements.Frets[3].Label)
	assert.Equal(t, "8:5", fretPlacements.Frets[4].Label)
	assert.Equal(t, "16:9", fretPlacements.Frets[5].Label)
	assert.Equal(t, "2:1", fretPlacements.Frets[6].Label)
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

	var fretPlacements handler.FretPlacements
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, 540.0, fretPlacements.ScaleLength)
	assert.Equal(t, "Ptolemy", fretPlacements.System)
	assert.Equal(t, "Fret positions for Ptolemy's 5-limit intense diatonic scale in Lydian mode.", fretPlacements.Description)
	assert.Equal(t, 7, len(fretPlacements.Frets))

	assert.Equal(t, "9:8", fretPlacements.Frets[0].Label)
	assert.Equal(t, "5:4", fretPlacements.Frets[1].Label)
	assert.Equal(t, "45:32", fretPlacements.Frets[2].Label)
	assert.Equal(t, "3:2", fretPlacements.Frets[3].Label)
	assert.Equal(t, "5:3", fretPlacements.Frets[4].Label)
	assert.Equal(t, "15:8", fretPlacements.Frets[5].Label)
	assert.Equal(t, "2:1", fretPlacements.Frets[6].Label)
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

	var fretPlacements handler.FretPlacements
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, 540.0, fretPlacements.ScaleLength)
	assert.Equal(t, "Ptolemy", fretPlacements.System)
	assert.Equal(t, "Fret positions for Ptolemy's 5-limit intense diatonic scale in Mixolydian mode.", fretPlacements.Description)
	assert.Equal(t, 7, len(fretPlacements.Frets))

	assert.Equal(t, "9:8", fretPlacements.Frets[0].Label)
	assert.Equal(t, "5:4", fretPlacements.Frets[1].Label)
	assert.Equal(t, "4:3", fretPlacements.Frets[2].Label)
	assert.Equal(t, "3:2", fretPlacements.Frets[3].Label)
	assert.Equal(t, "5:3", fretPlacements.Frets[4].Label)
	assert.Equal(t, "16:9", fretPlacements.Frets[5].Label)
	assert.Equal(t, "2:1", fretPlacements.Frets[6].Label)
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

	var fretPlacements handler.FretPlacements
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, 540.0, fretPlacements.ScaleLength)
	assert.Equal(t, "Ptolemy", fretPlacements.System)
	assert.Equal(t, "Fret positions for Ptolemy's 5-limit intense diatonic scale in Aeolian mode.", fretPlacements.Description)
	assert.Equal(t, 7, len(fretPlacements.Frets))

	assert.Equal(t, "9:8", fretPlacements.Frets[0].Label)
	assert.Equal(t, "6:5", fretPlacements.Frets[1].Label)
	assert.Equal(t, "4:3", fretPlacements.Frets[2].Label)
	assert.Equal(t, "3:2", fretPlacements.Frets[3].Label)
	assert.Equal(t, "8:5", fretPlacements.Frets[4].Label)
	assert.Equal(t, "9:5", fretPlacements.Frets[5].Label)
	assert.Equal(t, "2:1", fretPlacements.Frets[6].Label)
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

	var fretPlacements handler.FretPlacements
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, 540.0, fretPlacements.ScaleLength)
	assert.Equal(t, "Ptolemy", fretPlacements.System)
	assert.Equal(t, "Fret positions for Ptolemy's 5-limit intense diatonic scale in Locrian mode.", fretPlacements.Description)
	assert.Equal(t, 7, len(fretPlacements.Frets))

	assert.Equal(t, "16:15", fretPlacements.Frets[0].Label)
	assert.Equal(t, "6:5", fretPlacements.Frets[1].Label)
	assert.Equal(t, "4:3", fretPlacements.Frets[2].Label)
	assert.Equal(t, "64:45", fretPlacements.Frets[3].Label)
	assert.Equal(t, "8:5", fretPlacements.Frets[4].Label)
	assert.Equal(t, "16:9", fretPlacements.Frets[5].Label)
	assert.Equal(t, "2:1", fretPlacements.Frets[6].Label)
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

func Test_ShouldReturnErrorWhenPtolemicDiatonicScaleIsInvalid(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "ptolemy", "diatonicMode": "Athenian"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"please provide a valid mode for the diatonic scale"}`}, response)
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

	fretPlacements := handler.FretPlacements{}
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, float64(540), fretPlacements.ScaleLength)
	assert.Equal(t, "meantone", fretPlacements.System)
	assert.Equal(t, "Fret positions for meantone computed by narrowing of fifths by 0.25 of a syntonic comma (81/80).  Nominal note names used given a tonic of D.", fretPlacements.Description)
	assert.Equal(t, 13, len(fretPlacements.Frets))

	greaterSemitone := 1.069984
	lesserSemitone := 1.044907
	dieses := 1.024

	assert.Equal(t, handler.Fret{
		Label:    "1 (Eb)",
		Position: 35.3,
		Comment:  fmt.Sprintf("ratio: 1.070; interval: %f", greaterSemitone),
	}, fretPlacements.Frets[0])

	assert.Equal(t, handler.Fret{
		Label:    "2 (E)",
		Position: 57.0,
		Comment:  fmt.Sprintf("ratio: 1.118; interval: %f", lesserSemitone),
	}, fretPlacements.Frets[1])

	assert.Equal(t, handler.Fret{
		Label:    "7 (Ab)",
		Position: 162.7,
		Comment:  fmt.Sprintf("ratio: 1.431; interval: %f", dieses),
	}, fretPlacements.Frets[6])

	assert.Equal(t, handler.Fret{
		Label:    "13 (Octave)",
		Position: 270.0,
		Comment:  fmt.Sprintf("ratio: 2.0; interval: %f", greaterSemitone),
	}, fretPlacements.Frets[12])
}

func Test_shouldReturnQuarterCommaMeantonePlacementsWithProvidedScaleLengthWhenExtendedParameterIsInvalid(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "meantone", "extended": "yes please"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	fretPlacements := handler.FretPlacements{}
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, "Fret positions for meantone computed by narrowing of fifths by 0.25 of a syntonic comma (81/80).  Nominal note names used given a tonic of D.", fretPlacements.Description)
}

func Test_shouldReturnExtendedQuarterCommaMeantonePlacementsWithProvidedScaleLength(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "meantone", "extended": "true"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	fretPlacements := handler.FretPlacements{}
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, float64(540), fretPlacements.ScaleLength)
	assert.Equal(t, "meantone", fretPlacements.System)
	assert.Equal(t, "Fret positions for extended meantone computed by narrowing of fifths by 0.25 of a syntonic comma (81/80).  Nominal note names used given a tonic of D.", fretPlacements.Description)
	assert.Equal(t, 19, len(fretPlacements.Frets))

	lesserSemitone := 1.044907
	dieses := 1.024

	assert.Equal(t, handler.Fret{
		Label:    "1 (D#)",
		Position: 23.2,
		Comment:  fmt.Sprintf("ratio: 1.045; interval: %f", lesserSemitone),
	}, fretPlacements.Frets[0])

	assert.Equal(t, handler.Fret{
		Label:    "2 (Eb)",
		Position: 35.3,
		Comment:  fmt.Sprintf("ratio: 1.070; interval: %f", dieses),
	}, fretPlacements.Frets[1])

	assert.Equal(t, handler.Fret{
		Label:    "10 (Ab)",
		Position: 162.7,
		Comment:  fmt.Sprintf("ratio: 1.431; interval: %f", dieses),
	}, fretPlacements.Frets[9])

	assert.Equal(t, handler.Fret{
		Label:    "19 (Octave)",
		Position: 270.0,
		Comment:  fmt.Sprintf("ratio: 2.0; interval: %f", lesserSemitone),
	}, fretPlacements.Frets[18])
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

	fretPlacements := handler.FretPlacements{}
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, float64(540), fretPlacements.ScaleLength)
	assert.Equal(t, "Pythagorean", fretPlacements.System)
	assert.Equal(t, "Fret positions based on 3-limit Pythagorean ratios.", fretPlacements.Description)
	assert.Equal(t, 13, len(fretPlacements.Frets))

	assert.Equal(t, "256:243", fretPlacements.Frets[0].Label)  // minor second
	assert.Equal(t, "9:8", fretPlacements.Frets[1].Label)      // major second
	assert.Equal(t, "32:27", fretPlacements.Frets[2].Label)    // minor third
	assert.Equal(t, "81:64", fretPlacements.Frets[3].Label)    // major third
	assert.Equal(t, "4:3", fretPlacements.Frets[4].Label)      // perfect fourth
	assert.Equal(t, "1024:729", fretPlacements.Frets[5].Label) // augmented fourth
	assert.Equal(t, "729:512", fretPlacements.Frets[6].Label)  // diminished fifth
	assert.Equal(t, "3:2", fretPlacements.Frets[7].Label)      // perfect fifth
	assert.Equal(t, "128:81", fretPlacements.Frets[8].Label)   // minor sixth
	assert.Equal(t, "27:16", fretPlacements.Frets[9].Label)    // major sixth
	assert.Equal(t, "16:9", fretPlacements.Frets[10].Label)    // minor seventh
	assert.Equal(t, "243:128", fretPlacements.Frets[11].Label) // major seventh
	assert.Equal(t, "2:1", fretPlacements.Frets[12].Label)     // octave

	assert.Equal(t, "27.42", fmt.Sprintf("%.2f", fretPlacements.Frets[0].Position))
	assert.Equal(t, "60.00", fmt.Sprintf("%.2f", fretPlacements.Frets[1].Position))
	assert.Equal(t, "84.38", fmt.Sprintf("%.2f", fretPlacements.Frets[2].Position))
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

	fretPlacements := handler.FretPlacements{}
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, float64(540), fretPlacements.ScaleLength)
	assert.Equal(t, "31-TET", fretPlacements.System)
	assert.Equal(t, "Fret positions for 31-tone equal temperament.", fretPlacements.Description)
	assert.Equal(t, 31, len(fretPlacements.Frets))
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

	fretPlacements := handler.FretPlacements{}
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, float64(540), fretPlacements.ScaleLength)
	assert.Equal(t, "19-TET", fretPlacements.System)
	assert.Equal(t, "Fret positions for 19-tone equal temperament.", fretPlacements.Description)
	assert.Equal(t, 19, len(fretPlacements.Frets))
	assert.Equal(t, handler.Fret{Label: "Fret 1", Position: 19.345}, fretPlacements.Frets[0])
	assert.Equal(t, handler.Fret{Label: "Fret 19", Position: 270.0}, fretPlacements.Frets[18])
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

	fretPlacements := handler.FretPlacements{}
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, float64(540), fretPlacements.ScaleLength)
	assert.Equal(t, "saz", fretPlacements.System)
	assert.Equal(t, "Fret positions for traditional Turkish Saz tuning ratios.", fretPlacements.Description)
	assert.Equal(t, 17, len(fretPlacements.Frets))
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

	fretPlacements := handler.FretPlacements{}
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, float64(540), fretPlacements.ScaleLength)
	assert.Equal(t, "Bach's Well-Tempered Tuning", fretPlacements.System)
	assert.Equal(t, "Fret positions derived from Lehman's decoding of Bach's Well-Tempered tuning, using sixth-comma, twelfth-comma, and pure fifths.", fretPlacements.Description)
	assert.Equal(t, 12, len(fretPlacements.Frets))
}
