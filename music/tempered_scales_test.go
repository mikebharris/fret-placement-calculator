package music

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ShouldReturnScaleForBachWohltemperierteKlavierTuning(t *testing.T) {
	// Given
	scale := NewBachWohltemperierteKlavierScale()

	// When
	intervals := scale.Intervals()
	assert.Equal(t, "Bach's Well-Tempered Tuning", scale.System())
	assert.Equal(t, "Derived from Lehman's decoding of Bach's Well-Tempered tuning, using sixth-comma, twelfth-comma, and pure fifths.", scale.Description())
	assert.Equal(t, 13, len(intervals))

	assert.Equal(t, TemperedInterval(1.0), intervals[0])
	assert.Equal(t, TemperedInterval(1.061), intervals[1])
	assert.Equal(t, TemperedInterval(1.124), intervals[2])
	assert.Equal(t, TemperedInterval(1.19), intervals[3])
	assert.Equal(t, TemperedInterval(1.262), intervals[4])
	assert.Equal(t, TemperedInterval(1.336), intervals[5])
	assert.Equal(t, TemperedInterval(1.415), intervals[6])
	assert.Equal(t, TemperedInterval(1.5), intervals[7])
	assert.Equal(t, TemperedInterval(1.589), intervals[8])
	assert.Equal(t, TemperedInterval(1.682), intervals[9])
	assert.Equal(t, TemperedInterval(1.785), intervals[10])
	assert.Equal(t, TemperedInterval(1.889), intervals[11])
	assert.Equal(t, TemperedInterval(2.0), intervals[12])
}

func Test_ShouldReturnScaleForQuarterCommaMeantone(t *testing.T) {
	// Given
	scale := NewQuarterCommaMeantoneScale()

	// When
	intervals := scale.Intervals()
	assert.Equal(t, "Quarter-Comma Meantone", scale.System())
	assert.Equal(t, "Meantone temperament achieved by narrowing of fifths by 0.25 of a syntonic comma (81/80).", scale.Description())
	assert.Equal(t, 13, len(intervals))

	//greaterSemitone := 1.069984
	//lesserSemitone := 1.044907
	//dieses := 1.024

	assert.Equal(t, TemperedInterval(1.0), intervals[0])
	assert.Equal(t, TemperedInterval(1.044), intervals[1])
	assert.Equal(t, TemperedInterval(1.118), intervals[2])
	assert.Equal(t, TemperedInterval(1.189), intervals[3])
	assert.Equal(t, TemperedInterval(1.26), intervals[4])
	assert.Equal(t, TemperedInterval(1.334), intervals[5])
	assert.Equal(t, TemperedInterval(1.414), intervals[6])
	assert.Equal(t, TemperedInterval(1.5), intervals[7])
	assert.Equal(t, TemperedInterval(1.587), intervals[8])
	assert.Equal(t, TemperedInterval(1.682), intervals[9])
	assert.Equal(t, TemperedInterval(1.782), intervals[10])
	assert.Equal(t, TemperedInterval(1.888), intervals[11])
	assert.Equal(t, TemperedInterval(2.0), intervals[12])
}

func Test_ShouldReturnScaleForExtendedQuarterCommaMeantone(t *testing.T) {
	// Given
	scale := NewQuarterCommaMeantoneScale()

	// When
	intervals := scale.Intervals()
	assert.Equal(t, "Extended Quarter-Comma Meantone", scale.System())
	assert.Equal(t, "Meantone temperament achieved by narrowing of fifths by 0.25 of a syntonic comma (81/80).", scale.Description())
	assert.Equal(t, 19, len(intervals))

	//lesserSemitone := 1.044907
	//dieses := 1.024

	assert.Equal(t, TemperedInterval(1.0), intervals[0])
	assert.Equal(t, TemperedInterval(1.044), intervals[1])
	assert.Equal(t, TemperedInterval(1.118), intervals[2])
	assert.Equal(t, TemperedInterval(1.189), intervals[3])
	assert.Equal(t, TemperedInterval(1.26), intervals[4])
	assert.Equal(t, TemperedInterval(1.334), intervals[5])
	assert.Equal(t, TemperedInterval(1.414), intervals[6])
	assert.Equal(t, TemperedInterval(1.5), intervals[7])
	assert.Equal(t, TemperedInterval(1.587), intervals[8])
	assert.Equal(t, TemperedInterval(1.682), intervals[9])
	assert.Equal(t, TemperedInterval(1.782), intervals[10])
	assert.Equal(t, TemperedInterval(1.888), intervals[11])
	assert.Equal(t, TemperedInterval(2.0), intervals[12])
}

func Test_ShouldReturnScaleFor12ToneEqualTemperament(t *testing.T) {
	// Given
	scale := NewEqualTemperamentScale(12)

	// When
	intervals := scale.Intervals()
	assert.Equal(t, "Equal Temperament", scale.System())
	assert.Equal(t, "12-tone equal temperament.", scale.Description())
	assert.Equal(t, 13, len(intervals))

	assert.Equal(t, 0.0, intervals[0].toCents())
	assert.Equal(t, 100.0, intervals[1].toCents())
	assert.Equal(t, 200.0, intervals[2].toCents())
	assert.Equal(t, 300.0, intervals[3].toCents())
	assert.Equal(t, 400.0, intervals[4].toCents())
	assert.Equal(t, 500.0, intervals[5].toCents())
	assert.Equal(t, 600.0, intervals[7].toCents())
	assert.Equal(t, 700.0, intervals[8].toCents())
	assert.Equal(t, 800.0, intervals[9].toCents())
	assert.Equal(t, 900.0, intervals[10].toCents())
	assert.Equal(t, 1000.0, intervals[11].toCents())
	assert.Equal(t, 1100.0, intervals[12].toCents())
	assert.Equal(t, 1200.0, intervals[13].toCents())
}
