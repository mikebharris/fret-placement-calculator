package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"main/fret-placement-calculator-api/handler"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func Test_shouldReturnJustIntonationPlacementsWithProvidedScaleLengthWhenOnlyScaleLengthIsProvided(t *testing.T) {
	// Given
	ctx := context.Background()
	h := handler.Handler{}

	// When
	response, err := h.HandleRequest(ctx, events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"scaleLength": "540"}})

	// Then
	headers := map[string]string{"Content-Type": "application/json"}

	frets := []handler.Fret{
		{Label: "16:15", Position: 33.75},
		{Label: "10:9", Position: 54.0},
		{Label: "9:8", Position: 60.0},
		{Label: "6:5", Position: 90.0},
		{Label: "5:4", Position: 108.0},
		{Label: "35:25", Position: 154.29},
		{Label: "4:3", Position: 135.0},
		{Label: "45:32", Position: 156.0},
		{Label: "3:2", Position: 180.0},
		{Label: "8:5", Position: 202.5},
		{Label: "5:3", Position: 216.0},
		{Label: "16:9", Position: 236.25},
		{Label: "9:5", Position: 240.0},
		{Label: "15:8", Position: 252.0},
		{Label: "2:1", Position: 270.0},
	}

	fretting := handler.Fretting{
		System:      "ji",
		ScaleLength: 540,
		Frets:       frets,
	}

	m, _ := json.Marshal(fretting)

	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusOK, Headers: headers, Body: string(m)}, response)
	assert.Nil(t, err)
}

func Test_shouldReturnErrorWhenScaleLengthIsNotProvided(t *testing.T) {
	// Given
	ctx := context.Background()
	h := handler.Handler{}

	// When
	response, err := h.HandleRequest(ctx, events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{}})

	// Then
	headers := map[string]string{"Content-Type": "application/json"}

	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"scaleLength query parameter is required"}`}, response)
	assert.Nil(t, err)
}

func Test_shouldReturnQuarterCommaMeantonePlacementsWithProvidedScaleLength(t *testing.T) {
	// Given
	ctx := context.Background()
	h := handler.Handler{}
	headers := map[string]string{"Content-Type": "application/json"}
	greaterSemitone := 1.069984
	lesserSemitone := 1.044907
	dieses := 1.024

	// When
	response, err := h.HandleRequest(ctx, events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"scaleLength": "540", "temper": "meantone"}})

	// Then
	fretting := handler.Fretting{}

	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	_ = json.Unmarshal([]byte(response.Body), &fretting)

	assert.Equal(t, float64(540), fretting.ScaleLength)
	assert.Equal(t, "meantone", fretting.System)
	assert.Equal(t, "Calculating meantone based on narrowing of fifths by 0.25 of a syntonic comma (81/80).  Nominal note names used based on a tonic of D.", fretting.Description)
	assert.Equal(t, 13, len(fretting.Frets))

	assert.Equal(t, handler.Fret{
		Label:    "1 (Eb)",
		Position: 35.3,
		Comment:  fmt.Sprintf("ratio: 1.070; interval: %f", greaterSemitone),
	}, fretting.Frets[0])

	assert.Equal(t, handler.Fret{
		Label:    "2 (E)",
		Position: 57.0,
		Comment:  fmt.Sprintf("ratio: 1.118; interval: %f", lesserSemitone),
	}, fretting.Frets[1])

	assert.Equal(t, handler.Fret{
		Label:    "7 (Ab)",
		Position: 162.7,
		Comment:  fmt.Sprintf("ratio: 1.431; interval: %f", dieses),
	}, fretting.Frets[6])

	assert.Equal(t, handler.Fret{
		Label:    "13 (Octave)",
		Position: 270.0,
		Comment:  fmt.Sprintf("ratio: 2.0; interval: %f", greaterSemitone),
	}, fretting.Frets[12])

	assert.Nil(t, err)
}

func Test_shouldReturnExtendedQuarterCommaMeantonePlacementsWithProvidedScaleLength(t *testing.T) {
	// Given
	ctx := context.Background()
	h := handler.Handler{}

	// When
	response, err := h.HandleRequest(ctx, events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"scaleLength": "540", "temper": "meantone", "extended": "true"}})

	// Then
	headers := map[string]string{"Content-Type": "application/json"}

	fretting := handler.Fretting{}

	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	_ = json.Unmarshal([]byte(response.Body), &fretting)

	assert.Equal(t, float64(540), fretting.ScaleLength)
	assert.Equal(t, "meantone", fretting.System)
	assert.Equal(t, "Calculating extended meantone based on narrowing of fifths by 0.25 of a syntonic comma (81/80).  Nominal note names used based on a tonic of D.", fretting.Description)
	assert.Equal(t, 19, len(fretting.Frets))

	lesserSemitone := 1.044907
	dieses := 1.024

	assert.Equal(t, handler.Fret{
		Label:    "1 (D#)",
		Position: 23.2,
		Comment:  fmt.Sprintf("ratio: 1.045; interval: %f", lesserSemitone),
	}, fretting.Frets[0])

	assert.Equal(t, handler.Fret{
		Label:    "2 (Eb)",
		Position: 35.3,
		Comment:  fmt.Sprintf("ratio: 1.070; interval: %f", dieses),
	}, fretting.Frets[1])

	assert.Equal(t, handler.Fret{
		Label:    "10 (Ab)",
		Position: 162.7,
		Comment:  fmt.Sprintf("ratio: 1.431; interval: %f", dieses),
	}, fretting.Frets[9])

	assert.Equal(t, handler.Fret{
		Label:    "19 (Octave)",
		Position: 270.0,
		Comment:  fmt.Sprintf("ratio: 2.0; interval: %f", lesserSemitone),
	}, fretting.Frets[18])

	assert.Nil(t, err)
}

func Test_shouldReturnFifthCommaMeantonePlacementsWithProvidedScaleLength(t *testing.T) {
	// Given
	ctx := context.Background()
	h := handler.Handler{}
	headers := map[string]string{"Content-Type": "application/json"}
	greaterSemitone := 1.066667
	lesserSemitone := 1.049460
	dieses := 1.016396

	// When
	response, err := h.HandleRequest(ctx, events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"scaleLength": "540", "temper": "meantone", "temper-by": "0.2"}})

	// Then
	fretting := handler.Fretting{}

	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	_ = json.Unmarshal([]byte(response.Body), &fretting)

	assert.Equal(t, float64(540), fretting.ScaleLength)
	assert.Equal(t, "meantone", fretting.System)
	assert.Equal(t, "Calculating meantone based on narrowing of fifths by 0.20 of a syntonic comma (81/80).  Nominal note names used based on a tonic of D.", fretting.Description)
	assert.Equal(t, 13, len(fretting.Frets))

	assert.Equal(t, handler.Fret{
		Label:    "1 (Eb)",
		Position: 33.7,
		Comment:  fmt.Sprintf("ratio: 1.067; interval: %f", greaterSemitone),
	}, fretting.Frets[0])

	assert.Equal(t, handler.Fret{
		Label:    "2 (E)",
		Position: 57.6,
		Comment:  fmt.Sprintf("ratio: 1.119; interval: %f", lesserSemitone),
	}, fretting.Frets[1])

	assert.Equal(t, handler.Fret{
		Label:    "7 (Ab)",
		Position: 161.3,
		Comment:  fmt.Sprintf("ratio: 1.426; interval: %f", dieses),
	}, fretting.Frets[6])

	assert.Equal(t, handler.Fret{
		Label:    "13 (Octave)",
		Position: 270.0,
		Comment:  fmt.Sprintf("ratio: 2.0; interval: %f", greaterSemitone),
	}, fretting.Frets[12])

	assert.Nil(t, err)
}

func Test_shouldReturnPythagoreanPlacementsWithProvidedScaleLength(t *testing.T) {
	// Given
	ctx := context.Background()
	h := handler.Handler{}
	headers := map[string]string{"Content-Type": "application/json"}

	// When
	response, err := h.HandleRequest(ctx, events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"scaleLength": "540", "temper": "pythagorean"}})

	// Then
	fretting := handler.Fretting{}

	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	_ = json.Unmarshal([]byte(response.Body), &fretting)

	assert.Equal(t, float64(540), fretting.ScaleLength)
	assert.Equal(t, "pythagorean", fretting.System)
	assert.Equal(t, "Calculating based Pythagorean ratios.", fretting.Description)
	assert.Equal(t, 13, len(fretting.Frets))

	assert.Equal(t, handler.Fret{Label: "256:243", Position: 27.4}, fretting.Frets[0])

	assert.Nil(t, err)
}

func Test_ShouldReturnErrorWhenTemperParameterIsInvalid(t *testing.T) {
	// Given
	ctx := context.Background()
	h := handler.Handler{}

	// When
	response, err := h.HandleRequest(ctx, events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"scaleLength": "540", "temper": "invalid"}})

	// Then
	headers := map[string]string{"Content-Type": "application/json"}

	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"invalid temper parameter"}`}, response)
	assert.Nil(t, err)
}

func Test_ShouldReturnErrorWhenTemperByParameterIsOutOfRange(t *testing.T) {
	// Given
	ctx := context.Background()
	h := handler.Handler{}

	// When
	response, err := h.HandleRequest(ctx, events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"scaleLength": "540", "temper": "meantone", "temper-by": "1.5"}})

	// Then
	headers := map[string]string{"Content-Type": "application/json"}

	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"temper-by parameter must be less than 1"}`}, response)
	assert.Nil(t, err)
}

func Test_ShouldReturnEqualTemperamentPlacementsFor12EDOByDefault(t *testing.T) {
	// Given
	ctx := context.Background()
	h := handler.Handler{}

	// When
	response, err := h.HandleRequest(ctx, events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"scaleLength": "540", "temper": "equal"}})

	// Then
	headers := map[string]string{"Content-Type": "application/json"}

	fretting := handler.Fretting{}

	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	_ = json.Unmarshal([]byte(response.Body), &fretting)

	assert.Equal(t, float64(540), fretting.ScaleLength)
	assert.Equal(t, "12-TET", fretting.System)
	assert.Equal(t, 12, len(fretting.Frets))
	assert.Equal(t, handler.Fret{Label: "Fret 1", Position: 30.308}, fretting.Frets[0])
	assert.Equal(t, handler.Fret{Label: "Fret 12", Position: 270.0}, fretting.Frets[11])
	assert.Nil(t, err)
}

func Test_ShouldReturnEqualTemperamentPlacementsWithCustomDivisions(t *testing.T) {
	// Given
	ctx := context.Background()
	h := handler.Handler{}

	// When
	response, err := h.HandleRequest(ctx, events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"scaleLength": "540", "temper": "equal", "divisions": "19"}})

	// Then
	headers := map[string]string{"Content-Type": "application/json"}

	fretting := handler.Fretting{}

	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	_ = json.Unmarshal([]byte(response.Body), &fretting)

	assert.Equal(t, float64(540), fretting.ScaleLength)
	assert.Equal(t, "19-TET", fretting.System)
	assert.Equal(t, 19, len(fretting.Frets))
	assert.Equal(t, handler.Fret{Label: "Fret 1", Position: 19.345}, fretting.Frets[0])
	assert.Equal(t, handler.Fret{Label: "Fret 19", Position: 270.0}, fretting.Frets[18])
	assert.Nil(t, err)
}

func Test_ShouldReturnSazPlacementsWithProvidedScaleLength(t *testing.T) {
	// Given
	ctx := context.Background()
	h := handler.Handler{}

	// When
	response, err := h.HandleRequest(ctx, events.LambdaFunctionURLRequest{QueryStringParameters: map[string]string{"scaleLength": "540", "temper": "saz"}})

	// Then
	headers := map[string]string{"Content-Type": "application/json"}

	fretting := handler.Fretting{}

	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	_ = json.Unmarshal([]byte(response.Body), &fretting)

	assert.Equal(t, float64(540), fretting.ScaleLength)
	assert.Equal(t, "saz", fretting.System)
	assert.Equal(t, 17, len(fretting.Frets))
	assert.Nil(t, err)
}
