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
	assert.Equal(t, 14, len(intervals))

	assert.Equal(t, TemperedInterval(1.000), intervals[0])
	assert.Equal(t, TemperedInterval(1.070), intervals[1])
	assert.Equal(t, TemperedInterval(1.118), intervals[2])
	assert.Equal(t, TemperedInterval(1.196), intervals[3])
	assert.Equal(t, TemperedInterval(1.250), intervals[4])
	assert.Equal(t, TemperedInterval(1.337), intervals[5])
	assert.Equal(t, TemperedInterval(1.398), intervals[6])
	assert.Equal(t, TemperedInterval(1.431), intervals[7])
	assert.Equal(t, TemperedInterval(1.495), intervals[8])
	assert.Equal(t, TemperedInterval(1.600), intervals[9])
	assert.Equal(t, TemperedInterval(1.672), intervals[10])
	assert.Equal(t, TemperedInterval(1.789), intervals[11])
	assert.Equal(t, TemperedInterval(1.869), intervals[12])
	assert.Equal(t, TemperedInterval(2.000), intervals[13])
}

func Test_ShouldReturnScaleForExtendedQuarterCommaMeantone(t *testing.T) {
	// Given
	scale := NewExtendedQuarterCommaMeantoneScale()

	// When
	intervals := scale.Intervals()
	assert.Equal(t, "Extended Quarter-Comma Meantone", scale.System())
	assert.Equal(t, "Meantone temperament achieved by narrowing of fifths by 0.25 of a syntonic comma (81/80).", scale.Description())
	assert.Equal(t, 20, len(intervals))

	assert.Equal(t, TemperedInterval(1.0), intervals[0])
	assert.Equal(t, TemperedInterval(1.045), intervals[1])
	assert.Equal(t, TemperedInterval(1.070), intervals[2])
	assert.Equal(t, TemperedInterval(1.118), intervals[3])
	assert.Equal(t, TemperedInterval(1.168), intervals[4])
	assert.Equal(t, TemperedInterval(1.196), intervals[5])
	assert.Equal(t, TemperedInterval(1.250), intervals[6])
	assert.Equal(t, TemperedInterval(1.280), intervals[7])
	assert.Equal(t, TemperedInterval(1.337), intervals[8])
	assert.Equal(t, TemperedInterval(1.398), intervals[9])
	assert.Equal(t, TemperedInterval(1.431), intervals[10])
	assert.Equal(t, TemperedInterval(1.495), intervals[11])
	assert.Equal(t, TemperedInterval(1.562), intervals[12])
	assert.Equal(t, TemperedInterval(1.600), intervals[13])
	assert.Equal(t, TemperedInterval(1.672), intervals[14])
	assert.Equal(t, TemperedInterval(1.712), intervals[15])
	assert.Equal(t, TemperedInterval(1.789), intervals[16])
	assert.Equal(t, TemperedInterval(1.869), intervals[17])
	assert.Equal(t, TemperedInterval(1.914), intervals[18])
	assert.Equal(t, TemperedInterval(2.000), intervals[19])
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
	assert.Equal(t, 600.0, intervals[6].toCents())
	assert.Equal(t, 700.0, intervals[7].toCents())
	assert.Equal(t, 800.0, intervals[8].toCents())
	assert.Equal(t, 900.0, intervals[9].toCents())
	assert.Equal(t, 1000.0, intervals[10].toCents())
	assert.Equal(t, 1100.0, intervals[11].toCents())
	assert.Equal(t, 1200.0, intervals[12].toCents())
}

func Test_ShouldReturnScaleFor19ToneEqualTemperament(t *testing.T) {
	// Given
	scale := NewEqualTemperamentScale(19)

	// When
	intervals := scale.Intervals()
	assert.Equal(t, "Equal Temperament", scale.System())
	assert.Equal(t, "19-tone equal temperament.", scale.Description())
	assert.Equal(t, 20, len(intervals))

	assert.Equal(t, 0.0, intervals[0].toCents())
	assert.Equal(t, 63.16, intervals[1].toCents())
	assert.Equal(t, 126.32, intervals[2].toCents())
	assert.Equal(t, 189.47, intervals[3].toCents())
	assert.Equal(t, 252.63, intervals[4].toCents())
	assert.Equal(t, 315.79, intervals[5].toCents())
	assert.Equal(t, 378.95, intervals[6].toCents())
	assert.Equal(t, 442.11, intervals[7].toCents())
	assert.Equal(t, 505.26, intervals[8].toCents())
	assert.Equal(t, 568.42, intervals[9].toCents())
	assert.Equal(t, 631.58, intervals[10].toCents())
	assert.Equal(t, 694.74, intervals[11].toCents())
	assert.Equal(t, 757.89, intervals[12].toCents())
	assert.Equal(t, 821.05, intervals[13].toCents())
	assert.Equal(t, 884.21, intervals[14].toCents())
	assert.Equal(t, 947.37, intervals[15].toCents())
	assert.Equal(t, 1010.53, intervals[16].toCents())
	assert.Equal(t, 1073.68, intervals[17].toCents())
	assert.Equal(t, 1136.84, intervals[18].toCents())
	assert.Equal(t, 1200.00, intervals[19].toCents())
}
