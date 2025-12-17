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

func Test_shouldReturnJustIntonationPlacementsWithProvidedScaleLengthWhenOnlyScaleLengthIsProvided(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540"},
	})

	// Then
	assert.Nil(t, err)

	frets := []handler.Fret{
		// c
		{Label: "9:8", Position: 60.0},    // d - greater major second
		{Label: "5:4", Position: 108.0},   // e - major third
		{Label: "4:3", Position: 135.0},   // f - fourth
		{Label: "45:32", Position: 156.0}, // f# - augmented fourth (tritone) - two major tones and a small tone
		{Label: "3:2", Position: 180.0},   // g - fifth
		{Label: "5:3", Position: 216.0},   // a - major sixth
		{Label: "16:9", Position: 236.25}, // bb (lesser) minor seventh
		{Label: "15:8", Position: 252.0},  // b - major 7th
		{Label: "2:1", Position: 270.0},   // octave
	}

	fretPlacements := handler.FretPlacements{
		System:      "ji",
		Description: "Fret positions based on 5-limit just intonation pure ratios and diatonic scale.",
		ScaleLength: 540,
		Frets:       frets,
	}

	m, _ := json.Marshal(fretPlacements)
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(m)}, response)
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

func Test_ShouldReturnErrorWhenTemperParameterIsInvalid(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "temper": "invalid"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"invalid temper parameter"}`}, response)
}

func Test_shouldReturnQuarterCommaMeantonePlacementsWithProvidedScaleLength(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "temper": "meantone"},
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
		QueryStringParameters: map[string]string{"scaleLength": "540", "temper": "meantone", "extended": "yes please"},
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
		QueryStringParameters: map[string]string{"scaleLength": "540", "temper": "meantone", "extended": "true"},
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
		QueryStringParameters: map[string]string{"scaleLength": "540", "temper": "pythagorean"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	fretPlacements := handler.FretPlacements{}
	_ = json.Unmarshal([]byte(response.Body), &fretPlacements)
	assert.Equal(t, float64(540), fretPlacements.ScaleLength)
	assert.Equal(t, "pythagorean", fretPlacements.System)
	assert.Equal(t, "Fret positions based on 3-limit Pythagorean ratios.", fretPlacements.Description)
	assert.Equal(t, 13, len(fretPlacements.Frets))
	assert.Equal(t, handler.Fret{Label: "256:243", Position: 27.42}, fretPlacements.Frets[0])
}

func Test_ShouldReturnErrorIfNumberOfDivisionsForEqualTemperamentIsNotProvided(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "temper": "equal"}})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"please provide number of divisions for equal temperament"}`}, response)
}

func Test_ShouldReturnEqualTemperamentPlacementsWithCustomDivisions(t *testing.T) {
	// Given
	// When
	response, err := handler.Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "temper": "equal", "divisions": "19"},
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
		QueryStringParameters: map[string]string{"scaleLength": "540", "temper": "saz"},
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
