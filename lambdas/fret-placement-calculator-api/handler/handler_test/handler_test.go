package handler_test

import (
	"context"
	"encoding/json"
	"main/lambdas/fret-placement-calculator-api/handler"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

var headers = map[string]string{"Content-Type": "application/json"}

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
