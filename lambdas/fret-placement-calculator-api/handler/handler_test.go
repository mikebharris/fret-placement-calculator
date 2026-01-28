package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/mikebharris/music"
	"github.com/stretchr/testify/assert"
)

func Test_parseIntegerQueryParameter(t *testing.T) {
	type args struct {
		q        map[string]string
		key      string
		fallback int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "valid integer parameter",
			args: args{
				q:        map[string]string{"octaves": "3"},
				key:      "octaves",
				fallback: 1,
			},
			want: 3,
		},
		{
			name: "missing parameter uses fallback",
			args: args{
				q:        map[string]string{},
				key:      "octaves",
				fallback: 2,
			},
			want: 2,
		},
		{
			name: "invalid integer parameter uses fallback",
			args: args{
				q:        map[string]string{"octaves": "invalid"},
				key:      "octaves",
				fallback: 4,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseIntegerQueryParameter(tt.args.q, tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("parseIntegerQueryParameter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shouldReturnErrorWhenScaleLengthIsNotProvided(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"a numeric scaleLength greater than zero is required"}`}, response)
}

func Test_shouldReturnErrorWhenScaleLengthIsZero(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "0"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"a numeric scaleLength greater than zero is required"}`}, response)
}

func Test_shouldReturnErrorWhenScaleLengthIsLessThanZero(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "-100"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"a numeric scaleLength greater than zero is required"}`}, response)
}

func Test_shouldReturnErrorWhenScaleLengthIsNotANumber(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "three"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"a numeric scaleLength greater than zero is required"}`}, response)
}

func Test_ShouldReturnErrorWhenTuningSystemIsNotProvided(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"please provide a valid tuning system"}`}, response)
}

func Test_ShouldReturnErrorWhenTuningSystemIsInvalid(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "invalid"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, events.LambdaFunctionURLResponse{StatusCode: http.StatusUnprocessableEntity, Headers: headers, Body: `{"error":"please provide a valid tuning system"}`}, response)
}

func Test_ShouldDefaultToIonianIfNonSensicalPtolemyDiatonicScaleIsProvided(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "ptolemy", "diatonicMode": "Athenian"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard music.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "Ptolemy Intense Diatonic", fretboard.System)
	assert.Equal(t, "Fret positions based on Ptolemy's 5-limit intense diatonic scale in Ionian mode.", fretboard.Description)
	assert.Equal(t, 8, len(fretboard.Frets))
	assert.Equal(t, "9:8", fretboard.Frets[1].Label)
	assert.Equal(t, 60.0, fretboard.Frets[1].Position)
	assert.Equal(t, "Pythagorean (Greater) Major Second", fretboard.Frets[1].Comment)
	assert.Equal(t, "9:8", fretboard.Frets[1].Interval)
	assert.Equal(t, "5:4", fretboard.Frets[2].Label)
	assert.Equal(t, 108.0, fretboard.Frets[2].Position)
	assert.Equal(t, "4:3", fretboard.Frets[3].Label)
	assert.Equal(t, 135.0, fretboard.Frets[3].Position)
	assert.Equal(t, "3:2", fretboard.Frets[4].Label)
	assert.Equal(t, 180.0, fretboard.Frets[4].Position)
	assert.Equal(t, "5:3", fretboard.Frets[5].Label)
	assert.Equal(t, 216.0, fretboard.Frets[5].Position)
	assert.Equal(t, "15:8", fretboard.Frets[6].Label)
	assert.Equal(t, 252.0, fretboard.Frets[6].Position)
	assert.Equal(t, "2:1", fretboard.Frets[7].Label)
	assert.Equal(t, 270.0, fretboard.Frets[7].Position)
	assert.Equal(t, "Perfect Octave", fretboard.Frets[7].Comment)
	assert.Equal(t, "16:15", fretboard.Frets[7].Interval)
}

func Test_ShouldReturnFretPlacementsForPtolemyDiatonicScaleForProvidedMusicalMode(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "ptolemy", "diatonicMode": music.LydianMode.String()},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard music.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "Ptolemy Intense Diatonic", fretboard.System)
	assert.Equal(t, "Fret positions based on Ptolemy's 5-limit intense diatonic scale in Lydian mode.", fretboard.Description)
	assert.Equal(t, 8, len(fretboard.Frets))
	assert.Equal(t, "9:8", fretboard.Frets[1].Label)
	assert.Equal(t, 60.0, fretboard.Frets[1].Position)
	assert.Equal(t, "Pythagorean (Greater) Major Second", fretboard.Frets[1].Comment)
	assert.Equal(t, "9:8", fretboard.Frets[1].Interval)
	assert.Equal(t, "5:4", fretboard.Frets[2].Label)
	assert.Equal(t, 108.0, fretboard.Frets[2].Position)
	assert.Equal(t, "45:32", fretboard.Frets[3].Label)
	assert.Equal(t, 156.0, fretboard.Frets[3].Position)
	assert.Equal(t, "3:2", fretboard.Frets[4].Label)
	assert.Equal(t, 180.0, fretboard.Frets[4].Position)
	assert.Equal(t, "5:3", fretboard.Frets[5].Label)
	assert.Equal(t, 216.0, fretboard.Frets[5].Position)
	assert.Equal(t, "15:8", fretboard.Frets[6].Label)
	assert.Equal(t, 252.0, fretboard.Frets[6].Position)
	assert.Equal(t, "2:1", fretboard.Frets[7].Label)
	assert.Equal(t, 270.0, fretboard.Frets[7].Position)
	assert.Equal(t, "Perfect Octave", fretboard.Frets[7].Comment)
	assert.Equal(t, "16:15", fretboard.Frets[7].Interval)
}

func Test_ShouldReturnFretPlacementsForPythagoreanScale(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "pythagorean"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard music.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "Pythagorean", fretboard.System)
	assert.Equal(t, "Fret positions based on 3-limit Pythagorean ratios.", fretboard.Description)
	assert.Equal(t, 14, len(fretboard.Frets))
}

func Test_ShouldReturnFretPlacementsForFiveLimitPythagoreanScale(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "just5limitFromPythagorean"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard music.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "5-limit Pythagorean", fretboard.System)
	assert.Equal(t, "Fret positions based on 5-limit just intonation pure ratios chromatic scale derived from applying syntonic comma to Pythagorean ratios.", fretboard.Description)
	assert.Equal(t, 14, len(fretboard.Frets))
}

func Test_ShouldReturnFretPlacementsForFiveLimitJustScaleFromPureRatios(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "justFromRatios"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard music.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "5-limit Just Intonation", fretboard.System)
	assert.Equal(t, "Fret positions based on Just Intonation chromatic scale based on 5-limit pure ratios.", fretboard.Description)
	assert.Equal(t, 13, len(fretboard.Frets))
}

func Test_ShouldReturnFretPlacementsForFiveLimitJustScaleFromPureRatiosOverTwoOctaves(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "justFromRatios", "octaves": "2"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard music.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "5-limit Just Intonation", fretboard.System)
	assert.Equal(t, "Fret positions based on Just Intonation chromatic scale based on 5-limit pure ratios.", fretboard.Description)
	assert.Equal(t, 25, len(fretboard.Frets))

	assert.Equal(t, music.Fret{Label: "1:1", Position: 0, Comment: "Perfect Unison", Interval: "1:1"}, fretboard.Frets[0])
	assert.Equal(t, music.Fret{Label: "16:15", Position: 33.75, Comment: "Minor Second", Interval: "16:15"}, fretboard.Frets[1])
	assert.Equal(t, music.Fret{Label: "9:8", Position: 60.0, Comment: "Pythagorean (Greater) Major Second", Interval: "135:128"}, fretboard.Frets[2])
	assert.Equal(t, music.Fret{Label: "6:5", Position: 90.0, Comment: "Minor Third", Interval: "16:15"}, fretboard.Frets[3])
	assert.Equal(t, music.Fret{Label: "5:4", Position: 108.0, Comment: "Major Third", Interval: "25:24"}, fretboard.Frets[4])
	assert.Equal(t, music.Fret{Label: "4:3", Position: 135.0, Comment: "Perfect Fourth", Interval: "16:15"}, fretboard.Frets[5])
	assert.Equal(t, music.Fret{Label: "45:32", Position: 156.0, Comment: "Augmented Fourth", Interval: "135:128"}, fretboard.Frets[6])
	assert.Equal(t, music.Fret{Label: "3:2", Position: 180.0, Comment: "Perfect Fifth", Interval: "16:15"}, fretboard.Frets[7])
	assert.Equal(t, music.Fret{Label: "8:5", Position: 202.5, Comment: "Just Minor Sixth", Interval: "16:15"}, fretboard.Frets[8])
	assert.Equal(t, music.Fret{Label: "5:3", Position: 216.0, Comment: "Major Sixth", Interval: "25:24"}, fretboard.Frets[9])
	assert.Equal(t, music.Fret{Label: "9:5", Position: 240.0, Comment: "Just (Greater) Minor Seventh", Interval: "27:25"}, fretboard.Frets[10])
	assert.Equal(t, music.Fret{Label: "15:8", Position: 252.0, Comment: "Just Major Seventh", Interval: "25:24"}, fretboard.Frets[11])
	assert.Equal(t, music.Fret{Label: "2:1", Position: 270.0, Comment: "Perfect Octave", Interval: "16:15"}, fretboard.Frets[12])
	assert.Equal(t, music.Fret{Label: "16:15", Position: 286.88, Comment: "Minor Second", Interval: "16:15"}, fretboard.Frets[13])
	assert.Equal(t, music.Fret{Label: "9:8", Position: 300.0, Comment: "Pythagorean (Greater) Major Second", Interval: "135:128"}, fretboard.Frets[14])
	assert.Equal(t, music.Fret{Label: "6:5", Position: 315.0, Comment: "Minor Third", Interval: "16:15"}, fretboard.Frets[15])
	assert.Equal(t, music.Fret{Label: "5:4", Position: 324.0, Comment: "Major Third", Interval: "25:24"}, fretboard.Frets[16])
	assert.Equal(t, music.Fret{Label: "4:3", Position: 337.5, Comment: "Perfect Fourth", Interval: "16:15"}, fretboard.Frets[17])
	assert.Equal(t, music.Fret{Label: "45:32", Position: 348.0, Comment: "Augmented Fourth", Interval: "135:128"}, fretboard.Frets[18])
	assert.Equal(t, music.Fret{Label: "3:2", Position: 360.0, Comment: "Perfect Fifth", Interval: "16:15"}, fretboard.Frets[19])
	assert.Equal(t, music.Fret{Label: "8:5", Position: 371.25, Comment: "Just Minor Sixth", Interval: "16:15"}, fretboard.Frets[20])
	assert.Equal(t, music.Fret{Label: "5:3", Position: 378.0, Comment: "Major Sixth", Interval: "25:24"}, fretboard.Frets[21])
	assert.Equal(t, music.Fret{Label: "9:5", Position: 390.0, Comment: "Just (Greater) Minor Seventh", Interval: "27:25"}, fretboard.Frets[22])
	assert.Equal(t, music.Fret{Label: "15:8", Position: 396.0, Comment: "Just Major Seventh", Interval: "25:24"}, fretboard.Frets[23])
	assert.Equal(t, music.Fret{Label: "2:1", Position: 405.0, Comment: "Perfect Octave", Interval: "16:15"}, fretboard.Frets[24])
}

func Test_ShouldReturnFretPlacementsForThirteenLimitJustScaleFromPureRatios(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "justFromRatios", "limit": "13"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard music.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "13-limit Just Intonation", fretboard.System)
	assert.Equal(t, "Fret positions based on Just Intonation chromatic scale based on 13-limit pure ratios.", fretboard.Description)
	assert.Equal(t, 13, len(fretboard.Frets))
}

func Test_ShouldReturnFretPlacementsForSaz(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "saz"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard music.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "Saz", fretboard.System)
	assert.Equal(t, "Fret positions based on Turkish Saz tuning ratios.", fretboard.Description)
	assert.Equal(t, 18, len(fretboard.Frets))
}

func Test_ShouldReturnFretPlacementsForBachWellTemperament(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "bachWellTemperament"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard music.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "Bach's Well-Tempered Tuning", fretboard.System)
	assert.Equal(t, "Fret positions based on Derived from Lehman's decoding of Bach's Well-Tempered tuning, using sixth-comma, twelfth-comma, and pure fifths.", fretboard.Description)
	assert.Equal(t, 13, len(fretboard.Frets))
}

func Test_ShouldReturnFretPlacementsForQuarterCommaMeantoneScale(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "meantone"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard music.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "Quarter-Comma Meantone", fretboard.System)
	assert.Equal(t, "Fret positions based on Meantone temperament achieved by narrowing of fifths by 0.25 of a syntonic comma (81/80).", fretboard.Description)
	assert.Equal(t, 14, len(fretboard.Frets))
	assert.Equal(t, music.Fret{Label: "117.13 cents", Position: 35.33}, fretboard.Frets[1])
}

func Test_ShouldReturnFretPlacementsForExtendedQuarterCommaMeantoneScale(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "extendedMeantone"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, headers, response.Headers)

	var fretboard music.Fretboard
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, 540.0, fretboard.ScaleLength)
	assert.Equal(t, "Extended Quarter-Comma Meantone", fretboard.System)
	assert.Equal(t, "Fret positions based on Meantone temperament achieved by narrowing of fifths by 0.25 of a syntonic comma (81/80).", fretboard.Description)
	assert.Equal(t, 20, len(fretboard.Frets))
	assert.Equal(t, music.Fret{Label: "76.20 cents", Position: 23.25}, fretboard.Frets[1])
	assert.Equal(t, music.Fret{Label: "117.13 cents", Position: 35.33}, fretboard.Frets[2])
}

func Test_ShouldDefaultTo31EqualTemperamentByDefault(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "equal"}})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	fretboard := music.Fretboard{}
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, float64(540), fretboard.ScaleLength)
	assert.Equal(t, "Equal Temperament", fretboard.System)
	assert.Equal(t, "Fret positions based on 31-tone equal temperament.", fretboard.Description)
	assert.Equal(t, 32, len(fretboard.Frets))
}

func Test_ShouldReturnEqualTemperamentPlacementsWithCustomDivisions(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "600", "tuningSystem": "equal", "divisions": "12"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	fretboard := music.Fretboard{}
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, float64(600), fretboard.ScaleLength)
	assert.Equal(t, "Equal Temperament", fretboard.System)
	assert.Equal(t, "Fret positions based on 12-tone equal temperament.", fretboard.Description)
	assert.Equal(t, 13, len(fretboard.Frets))
	assert.Equal(t, music.Fret{Label: "0.00 cents", Position: 0}, fretboard.Frets[0])
	assert.Equal(t, music.Fret{Label: "100.00 cents", Position: 33.68}, fretboard.Frets[1])
	assert.Equal(t, music.Fret{Label: "200.00 cents", Position: 65.46}, fretboard.Frets[2])
	assert.Equal(t, music.Fret{Label: "300.00 cents", Position: 95.46}, fretboard.Frets[3])
	assert.Equal(t, music.Fret{Label: "400.00 cents", Position: 123.78}, fretboard.Frets[4])
	assert.Equal(t, music.Fret{Label: "500.00 cents", Position: 150.51}, fretboard.Frets[5])
	assert.Equal(t, music.Fret{Label: "600.00 cents", Position: 175.74}, fretboard.Frets[6])
	assert.Equal(t, music.Fret{Label: "700.00 cents", Position: 199.55}, fretboard.Frets[7])
	assert.Equal(t, music.Fret{Label: "800.00 cents", Position: 222.02}, fretboard.Frets[8])
	assert.Equal(t, music.Fret{Label: "900.00 cents", Position: 243.24}, fretboard.Frets[9])
	assert.Equal(t, music.Fret{Label: "1000.00 cents", Position: 263.26}, fretboard.Frets[10])
	assert.Equal(t, music.Fret{Label: "1100.00 cents", Position: 282.16}, fretboard.Frets[11])
	assert.Equal(t, music.Fret{Label: "1200.00 cents", Position: 300.0}, fretboard.Frets[12])
}

func Test_ShouldReturnMultipleOctavesOfFretsIfSpecified(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "600", "tuningSystem": "equal", "divisions": "12", "octaves": "2"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	fretboard := music.Fretboard{}
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, float64(600), fretboard.ScaleLength)
	assert.Equal(t, "Equal Temperament", fretboard.System)
	assert.Equal(t, "Fret positions based on 12-tone equal temperament.", fretboard.Description)
	assert.Equal(t, 25, len(fretboard.Frets))
	assert.Equal(t, music.Fret{Label: "0.00 cents", Position: 0}, fretboard.Frets[0])
	assert.Equal(t, music.Fret{Label: "100.00 cents", Position: 33.68}, fretboard.Frets[1])
	assert.Equal(t, music.Fret{Label: "200.00 cents", Position: 65.46}, fretboard.Frets[2])
	assert.Equal(t, music.Fret{Label: "300.00 cents", Position: 95.46}, fretboard.Frets[3])
	assert.Equal(t, music.Fret{Label: "400.00 cents", Position: 123.78}, fretboard.Frets[4])
	assert.Equal(t, music.Fret{Label: "500.00 cents", Position: 150.51}, fretboard.Frets[5])
	assert.Equal(t, music.Fret{Label: "600.00 cents", Position: 175.74}, fretboard.Frets[6])
	assert.Equal(t, music.Fret{Label: "700.00 cents", Position: 199.55}, fretboard.Frets[7])
	assert.Equal(t, music.Fret{Label: "800.00 cents", Position: 222.02}, fretboard.Frets[8])
	assert.Equal(t, music.Fret{Label: "900.00 cents", Position: 243.24}, fretboard.Frets[9])
	assert.Equal(t, music.Fret{Label: "1000.00 cents", Position: 263.26}, fretboard.Frets[10])
	assert.Equal(t, music.Fret{Label: "1100.00 cents", Position: 282.16}, fretboard.Frets[11])
	assert.Equal(t, music.Fret{Label: "1200.00 cents", Position: 300.0}, fretboard.Frets[12])
	assert.Equal(t, music.Fret{Label: "1300.00 cents", Position: 316.84}, fretboard.Frets[13])
	assert.Equal(t, music.Fret{Label: "2400.00 cents", Position: 450.0}, fretboard.Frets[24])
}

func Test_ShouldReturnSazPlacementsWithProvidedScaleLength(t *testing.T) {
	// Given
	// When
	response, err := Handler{}.HandleRequest(context.Background(), events.LambdaFunctionURLRequest{
		QueryStringParameters: map[string]string{"scaleLength": "540", "tuningSystem": "saz"},
	})

	// Then
	assert.Nil(t, err)
	assert.Equal(t, response.StatusCode, http.StatusOK)
	assert.Equal(t, response.Headers, headers)

	fretboard := music.Fretboard{}
	_ = json.Unmarshal([]byte(response.Body), &fretboard)
	assert.Equal(t, float64(540), fretboard.ScaleLength)
	assert.Equal(t, "Saz", fretboard.System)
	assert.Equal(t, "Fret positions based on Turkish Saz tuning ratios.", fretboard.Description)
	assert.Equal(t, 18, len(fretboard.Frets))
}
