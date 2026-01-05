package music

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ShouldReturnPythagorean3LimitScaleWithExpectedScaleDegrees(t *testing.T) {
	// Given
	scale := NewPythagoreanScale()

	// When
	intervals := scale.Intervals()

	// Then
	assert.Equal(t, "Pythagorean", scale.System())
	assert.Equal(t, "3-limit Pythagorean ratios.", scale.Description())
	assert.Equal(t, 14, len(intervals))
	assert.Equal(t, JustInterval{numerator: 1, denominator: 1}, intervals[0])
	assert.Equal(t, JustInterval{numerator: 256, denominator: 243}, intervals[1])
	assert.Equal(t, JustInterval{numerator: 9, denominator: 8}, intervals[2])
	assert.Equal(t, JustInterval{numerator: 32, denominator: 27}, intervals[3])
	assert.Equal(t, JustInterval{numerator: 81, denominator: 64}, intervals[4])
	assert.Equal(t, JustInterval{numerator: 4, denominator: 3}, intervals[5])
	assert.Equal(t, JustInterval{numerator: 1024, denominator: 729}, intervals[6])
	assert.Equal(t, JustInterval{numerator: 729, denominator: 512}, intervals[7])
	assert.Equal(t, JustInterval{numerator: 3, denominator: 2}, intervals[8])
	assert.Equal(t, JustInterval{numerator: 128, denominator: 81}, intervals[9])
	assert.Equal(t, JustInterval{numerator: 27, denominator: 16}, intervals[10])
	assert.Equal(t, JustInterval{numerator: 16, denominator: 9}, intervals[11])
	assert.Equal(t, JustInterval{numerator: 243, denominator: 128}, intervals[12])
	assert.Equal(t, JustInterval{numerator: 2, denominator: 1}, intervals[13])
}

func Test_ShouldReturn5LimitScaleWithExpectedScaleDegrees(t *testing.T) {
	// Given
	scale := New5LimitPythagoreanScale()

	// When
	intervals := scale.Intervals()

	// Then
	assert.Equal(t, "5-limit Pythagorean", scale.System())
	assert.Equal(t, "5-limit just intonation pure ratios chromatic scale derived from applying syntonic comma to Pythagorean ratios.", scale.Description())
	assert.Equal(t, 14, len(intervals))
	assert.Equal(t, JustInterval{numerator: 1, denominator: 1}, intervals[0])
	assert.Equal(t, JustInterval{numerator: 16, denominator: 15}, intervals[1])
	assert.Equal(t, JustInterval{numerator: 10, denominator: 9}, intervals[2])
	assert.Equal(t, JustInterval{numerator: 6, denominator: 5}, intervals[3])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 4}, intervals[4])
	assert.Equal(t, JustInterval{numerator: 4, denominator: 3}, intervals[5])
	assert.Equal(t, JustInterval{numerator: 64, denominator: 45}, intervals[6])
	assert.Equal(t, JustInterval{numerator: 45, denominator: 32}, intervals[7])
	assert.Equal(t, JustInterval{numerator: 3, denominator: 2}, intervals[8])
	assert.Equal(t, JustInterval{numerator: 8, denominator: 5}, intervals[9])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 3}, intervals[10])
	assert.Equal(t, JustInterval{numerator: 9, denominator: 5}, intervals[11])
	assert.Equal(t, JustInterval{numerator: 15, denominator: 8}, intervals[12])
	assert.Equal(t, JustInterval{numerator: 2, denominator: 1}, intervals[13])
}

func Test_shouldReturnAsymmetricJustChromaticScaleBasedOnPureRatios(t *testing.T) {
	// Given
	// Mike read the page at https://en.wikipedia.org/wiki/Five-limit_tuning

	// When
	scale := New5LimitJustIntonationChromaticScale(Asymmetric)
	intervals := scale.Intervals()

	// Then
	assert.Equal(t, "5-limit Just Intonation", scale.System())
	assert.Equal(t, "5-limit just intonation pure ratios derived from third- and fifth-partial ratios.", scale.Description())

	assert.Equal(t, 13, len(intervals))
	assert.Equal(t, JustInterval{numerator: 1, denominator: 1}, intervals[0])
	assert.Equal(t, JustInterval{numerator: 16, denominator: 15}, intervals[1])
	assert.Equal(t, JustInterval{numerator: 9, denominator: 8}, intervals[2])
	assert.Equal(t, JustInterval{numerator: 6, denominator: 5}, intervals[3])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 4}, intervals[4])
	assert.Equal(t, JustInterval{numerator: 4, denominator: 3}, intervals[5])
	assert.Equal(t, JustInterval{numerator: 45, denominator: 32}, intervals[6])
	assert.Equal(t, JustInterval{numerator: 3, denominator: 2}, intervals[7])
	assert.Equal(t, JustInterval{numerator: 8, denominator: 5}, intervals[8])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 3}, intervals[9])
	assert.Equal(t, JustInterval{numerator: 9, denominator: 5}, intervals[10])
	assert.Equal(t, JustInterval{numerator: 15, denominator: 8}, intervals[11])
	assert.Equal(t, JustInterval{numerator: 2, denominator: 1}, intervals[12])
}

func Test_shouldReturnJustChromaticScaleWithLesserMinorSeventhBasedOnPureRatios(t *testing.T) {
	// Given
	// Mike read the page at https://en.wikipedia.org/wiki/Five-limit_tuning

	// When
	scale := New5LimitJustIntonationChromaticScale(Symmetric1)
	intervals := scale.Intervals()

	// Then
	assert.Equal(t, "5-limit Just Intonation", scale.System())
	assert.Equal(t, "5-limit just intonation pure ratios derived from third- and fifth-partial ratios.", scale.Description())

	assert.Equal(t, 13, len(intervals))
	assert.Equal(t, JustInterval{numerator: 1, denominator: 1}, intervals[0])
	assert.Equal(t, JustInterval{numerator: 16, denominator: 15}, intervals[1])
	assert.Equal(t, JustInterval{numerator: 9, denominator: 8}, intervals[2])
	assert.Equal(t, JustInterval{numerator: 6, denominator: 5}, intervals[3])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 4}, intervals[4])
	assert.Equal(t, JustInterval{numerator: 4, denominator: 3}, intervals[5])
	assert.Equal(t, JustInterval{numerator: 45, denominator: 32}, intervals[6])
	assert.Equal(t, JustInterval{numerator: 3, denominator: 2}, intervals[7])
	assert.Equal(t, JustInterval{numerator: 8, denominator: 5}, intervals[8])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 3}, intervals[9])
	assert.Equal(t, JustInterval{numerator: 16, denominator: 9}, intervals[10])
	assert.Equal(t, JustInterval{numerator: 15, denominator: 8}, intervals[11])
	assert.Equal(t, JustInterval{numerator: 2, denominator: 1}, intervals[12])
}

func Test_shouldReturnJustChromaticScaleWithLesserMajorSecondBasedOnPureRatios(t *testing.T) {
	// Given
	// Mike read the page at https://en.wikipedia.org/wiki/Five-limit_tuning

	// When
	scale := New5LimitJustIntonationChromaticScale(Symmetric2)
	intervals := scale.Intervals()

	// Then
	assert.Equal(t, "5-limit Just Intonation", scale.System())
	assert.Equal(t, "5-limit just intonation pure ratios derived from third- and fifth-partial ratios.", scale.Description())

	assert.Equal(t, 13, len(intervals))
	assert.Equal(t, JustInterval{numerator: 1, denominator: 1}, intervals[0])
	assert.Equal(t, JustInterval{numerator: 16, denominator: 15}, intervals[1])
	assert.Equal(t, JustInterval{numerator: 10, denominator: 9}, intervals[2])
	assert.Equal(t, JustInterval{numerator: 6, denominator: 5}, intervals[3])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 4}, intervals[4])
	assert.Equal(t, JustInterval{numerator: 4, denominator: 3}, intervals[5])
	assert.Equal(t, JustInterval{numerator: 45, denominator: 32}, intervals[6])
	assert.Equal(t, JustInterval{numerator: 3, denominator: 2}, intervals[7])
	assert.Equal(t, JustInterval{numerator: 8, denominator: 5}, intervals[8])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 3}, intervals[9])
	assert.Equal(t, JustInterval{numerator: 9, denominator: 5}, intervals[10])
	assert.Equal(t, JustInterval{numerator: 15, denominator: 8}, intervals[11])
	assert.Equal(t, JustInterval{numerator: 2, denominator: 1}, intervals[12])
}

func Test_shouldReturn7LimitJustChromaticScaleBasedOnPureRatios(t *testing.T) {
	// Given
	// Mike read the page at https://en.wikipedia.org/wiki/Five-limit_tuning

	// When
	scale := New7LimitJustIntonationChromaticScale()
	intervals := scale.Intervals()

	// Then
	assert.Equal(t, "7-limit Just Intonation", scale.System())
	assert.Equal(t, "7-limit just intonation pure ratios derived from third-, fifth- and seventh-partial ratios.", scale.Description())

	assert.Equal(t, 13, len(intervals))
	assert.Equal(t, JustInterval{numerator: 1, denominator: 1}, intervals[0])
	assert.Equal(t, JustInterval{numerator: 15, denominator: 14}, intervals[1])
	assert.Equal(t, JustInterval{numerator: 8, denominator: 7}, intervals[2])
	assert.Equal(t, JustInterval{numerator: 6, denominator: 5}, intervals[3])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 4}, intervals[4])
	assert.Equal(t, JustInterval{numerator: 4, denominator: 3}, intervals[5])
	assert.Equal(t, JustInterval{numerator: 7, denominator: 5}, intervals[6])
	assert.Equal(t, JustInterval{numerator: 3, denominator: 2}, intervals[7])
	assert.Equal(t, JustInterval{numerator: 8, denominator: 5}, intervals[8])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 3}, intervals[9])
	assert.Equal(t, JustInterval{numerator: 7, denominator: 4}, intervals[10])
	assert.Equal(t, JustInterval{numerator: 15, denominator: 8}, intervals[11])
	assert.Equal(t, JustInterval{numerator: 2, denominator: 1}, intervals[12])
}

func Test_shouldReturn13LimitJustChromaticScaleBasedOnPureRatios(t *testing.T) {
	// Given
	// Mike read the page at https://en.wikipedia.org/wiki/Five-limit_tuning

	// When
	scale := New13LimitJustIntonationChromaticScale()
	intervals := scale.Intervals()

	// Then
	assert.Equal(t, "13-limit Just Intonation", scale.System())
	assert.Equal(t, "Just Intonation chromatic scale based on 13-limit pure ratios.", scale.Description())

	assert.Equal(t, 13, len(intervals))
	assert.Equal(t, JustInterval{numerator: 1, denominator: 1}, intervals[0])
	assert.Equal(t, JustInterval{numerator: 13, denominator: 12}, intervals[1])
	assert.Equal(t, JustInterval{numerator: 8, denominator: 7}, intervals[2])
	assert.Equal(t, JustInterval{numerator: 6, denominator: 5}, intervals[3])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 4}, intervals[4])
	assert.Equal(t, JustInterval{numerator: 4, denominator: 3}, intervals[5])
	assert.Equal(t, JustInterval{numerator: 7, denominator: 5}, intervals[6])
	assert.Equal(t, JustInterval{numerator: 3, denominator: 2}, intervals[7])
	assert.Equal(t, JustInterval{numerator: 8, denominator: 5}, intervals[8])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 3}, intervals[9])
	assert.Equal(t, JustInterval{numerator: 7, denominator: 4}, intervals[10])
	assert.Equal(t, JustInterval{numerator: 13, denominator: 7}, intervals[11])
	assert.Equal(t, JustInterval{numerator: 2, denominator: 1}, intervals[12])
}

func Test_ShouldReturnIntenseDiatonicScaleInLydianModeWithScaleDegrees(t *testing.T) {
	// Given
	scale := NewIntenseDiatonicScale(Lydian)

	// When
	intervals := scale.Intervals()

	// Then
	assert.Equal(t, "Ptolemy Intense Diatonic", scale.System())
	assert.Equal(t, "Ptolemy's 5-limit intense diatonic scale in Lydian mode.", scale.Description())
	assert.Equal(t, 8, len(intervals))

	assert.Equal(t, JustInterval{numerator: 1, denominator: 1}, intervals[0])
	assert.Equal(t, JustInterval{numerator: 9, denominator: 8}, intervals[1])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 4}, intervals[2])
	assert.Equal(t, JustInterval{numerator: 45, denominator: 32}, intervals[3])
	assert.Equal(t, JustInterval{numerator: 3, denominator: 2}, intervals[4])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 3}, intervals[5])
	assert.Equal(t, JustInterval{numerator: 15, denominator: 8}, intervals[6])
	assert.Equal(t, JustInterval{numerator: 2, denominator: 1}, intervals[7])
}

func Test_ShouldReturnIntenseDiatonicScaleInIonianModeWithScaleDegrees(t *testing.T) {
	// Given
	scale := NewIntenseDiatonicScale(Ionian)

	// When
	intervals := scale.Intervals()

	// Then
	assert.Equal(t, "Ptolemy Intense Diatonic", scale.System())
	assert.Equal(t, "Ptolemy's 5-limit intense diatonic scale in Ionian mode.", scale.Description())
	assert.Equal(t, 8, len(intervals))

	assert.Equal(t, JustInterval{numerator: 1, denominator: 1}, intervals[0])
	assert.Equal(t, JustInterval{numerator: 9, denominator: 8}, intervals[1])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 4}, intervals[2])
	assert.Equal(t, JustInterval{numerator: 4, denominator: 3}, intervals[3])
	assert.Equal(t, JustInterval{numerator: 3, denominator: 2}, intervals[4])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 3}, intervals[5])
	assert.Equal(t, JustInterval{numerator: 15, denominator: 8}, intervals[6])
	assert.Equal(t, JustInterval{numerator: 2, denominator: 1}, intervals[7])
}

func Test_ShouldReturnIntenseDiatonicScaleInMixolydianModeWithScaleDegrees(t *testing.T) {
	// Given
	scale := NewIntenseDiatonicScale(Mixolydian)

	// When
	intervals := scale.Intervals()

	// Then
	assert.Equal(t, "Ptolemy Intense Diatonic", scale.System())
	assert.Equal(t, "Ptolemy's 5-limit intense diatonic scale in Mixolydian mode.", scale.Description())
	assert.Equal(t, 8, len(intervals))

	assert.Equal(t, JustInterval{numerator: 1, denominator: 1}, intervals[0])
	assert.Equal(t, JustInterval{numerator: 9, denominator: 8}, intervals[1])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 4}, intervals[2])
	assert.Equal(t, JustInterval{numerator: 4, denominator: 3}, intervals[3])
	assert.Equal(t, JustInterval{numerator: 3, denominator: 2}, intervals[4])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 3}, intervals[5])
	assert.Equal(t, JustInterval{numerator: 16, denominator: 9}, intervals[6])
	assert.Equal(t, JustInterval{numerator: 2, denominator: 1}, intervals[7])
}

func Test_ShouldReturnIntenseDiatonicScaleInDorianModeWithScaleDegrees(t *testing.T) {
	// Given
	scale := NewIntenseDiatonicScale(Dorian)

	// When
	intervals := scale.Intervals()

	// Then
	assert.Equal(t, "Ptolemy Intense Diatonic", scale.System())
	assert.Equal(t, "Ptolemy's 5-limit intense diatonic scale in Dorian mode.", scale.Description())
	assert.Equal(t, 8, len(intervals))

	assert.Equal(t, JustInterval{numerator: 1, denominator: 1}, intervals[0])
	assert.Equal(t, JustInterval{numerator: 9, denominator: 8}, intervals[1])
	assert.Equal(t, JustInterval{numerator: 6, denominator: 5}, intervals[2])
	assert.Equal(t, JustInterval{numerator: 4, denominator: 3}, intervals[3])
	assert.Equal(t, JustInterval{numerator: 3, denominator: 2}, intervals[4])
	assert.Equal(t, JustInterval{numerator: 5, denominator: 3}, intervals[5])
	assert.Equal(t, JustInterval{numerator: 16, denominator: 9}, intervals[6])
	assert.Equal(t, JustInterval{numerator: 2, denominator: 1}, intervals[7])
}

func Test_ShouldReturnIntenseDiatonicScaleInAeolianModeWithScaleDegrees(t *testing.T) {
	// Given
	scale := NewIntenseDiatonicScale(Aeolian)

	// When
	intervals := scale.Intervals()

	// Then
	assert.Equal(t, "Ptolemy Intense Diatonic", scale.System())
	assert.Equal(t, "Ptolemy's 5-limit intense diatonic scale in Aeolian mode.", scale.Description())
	assert.Equal(t, 8, len(intervals))

	assert.Equal(t, JustInterval{numerator: 1, denominator: 1}, intervals[0])
	assert.Equal(t, JustInterval{numerator: 9, denominator: 8}, intervals[1])
	assert.Equal(t, JustInterval{numerator: 6, denominator: 5}, intervals[2])
	assert.Equal(t, JustInterval{numerator: 4, denominator: 3}, intervals[3])
	assert.Equal(t, JustInterval{numerator: 3, denominator: 2}, intervals[4])
	assert.Equal(t, JustInterval{numerator: 8, denominator: 5}, intervals[5])
	assert.Equal(t, JustInterval{numerator: 9, denominator: 5}, intervals[6])
	assert.Equal(t, JustInterval{numerator: 2, denominator: 1}, intervals[7])
}

func Test_ShouldReturnIntenseDiatonicScaleInPhrygianModeWithScaleDegrees(t *testing.T) {
	// Given
	scale := NewIntenseDiatonicScale(Phrygian)

	// When
	intervals := scale.Intervals()

	// Then
	assert.Equal(t, "Ptolemy Intense Diatonic", scale.System())
	assert.Equal(t, "Ptolemy's 5-limit intense diatonic scale in Phrygian mode.", scale.Description())
	assert.Equal(t, 8, len(intervals))

	assert.Equal(t, JustInterval{numerator: 1, denominator: 1}, intervals[0])
	assert.Equal(t, JustInterval{numerator: 16, denominator: 15}, intervals[1])
	assert.Equal(t, JustInterval{numerator: 6, denominator: 5}, intervals[2])
	assert.Equal(t, JustInterval{numerator: 4, denominator: 3}, intervals[3])
	assert.Equal(t, JustInterval{numerator: 3, denominator: 2}, intervals[4])
	assert.Equal(t, JustInterval{numerator: 8, denominator: 5}, intervals[5])
	assert.Equal(t, JustInterval{numerator: 9, denominator: 5}, intervals[6])
	assert.Equal(t, JustInterval{numerator: 2, denominator: 1}, intervals[7])
}

func Test_ShouldReturnIntenseDiatonicScaleInLocrianModeWithScaleDegrees(t *testing.T) {
	// Given
	scale := NewIntenseDiatonicScale(Locrian)

	// When
	intervals := scale.Intervals()

	// Then
	assert.Equal(t, "Ptolemy Intense Diatonic", scale.System())
	assert.Equal(t, "Ptolemy's 5-limit intense diatonic scale in Locrian mode.", scale.Description())
	assert.Equal(t, 8, len(intervals))

	assert.Equal(t, JustInterval{numerator: 1, denominator: 1}, intervals[0])
	assert.Equal(t, JustInterval{numerator: 16, denominator: 15}, intervals[1])
	assert.Equal(t, JustInterval{numerator: 6, denominator: 5}, intervals[2])
	assert.Equal(t, JustInterval{numerator: 4, denominator: 3}, intervals[3])
	assert.Equal(t, JustInterval{numerator: 64, denominator: 45}, intervals[4])
	assert.Equal(t, JustInterval{numerator: 8, denominator: 5}, intervals[5])
	assert.Equal(t, JustInterval{numerator: 16, denominator: 9}, intervals[6])
	assert.Equal(t, JustInterval{numerator: 2, denominator: 1}, intervals[7])
}

func Test_ShouldReturnSazScaleWithExpectedScaleDegrees(t *testing.T) {
	// Given
	scale := NewSazScale()

	// When
	intervals := scale.Intervals()

	// Then
	assert.Equal(t, "Saz", scale.System())
	assert.Equal(t, "Turkish Saz tuning ratios.", scale.Description())
	assert.Equal(t, 18, len(intervals))
	assert.Equal(t, JustInterval{numerator: 1, denominator: 1}, intervals[0])
	assert.Equal(t, JustInterval{numerator: 18, denominator: 17}, intervals[1])
	assert.Equal(t, JustInterval{numerator: 12, denominator: 11}, intervals[2])
	assert.Equal(t, JustInterval{numerator: 9, denominator: 8}, intervals[3])
	assert.Equal(t, JustInterval{numerator: 81, denominator: 68}, intervals[4])
	assert.Equal(t, JustInterval{numerator: 27, denominator: 22}, intervals[5])
	assert.Equal(t, JustInterval{numerator: 81, denominator: 64}, intervals[6])
	assert.Equal(t, JustInterval{numerator: 4, denominator: 3}, intervals[7])
	assert.Equal(t, JustInterval{numerator: 24, denominator: 17}, intervals[8])
	assert.Equal(t, JustInterval{numerator: 16, denominator: 11}, intervals[9])
	assert.Equal(t, JustInterval{numerator: 3, denominator: 2}, intervals[10])
	assert.Equal(t, JustInterval{numerator: 27, denominator: 17}, intervals[11])
	assert.Equal(t, JustInterval{numerator: 18, denominator: 11}, intervals[12])
	assert.Equal(t, JustInterval{numerator: 27, denominator: 16}, intervals[13])
	assert.Equal(t, JustInterval{numerator: 16, denominator: 9}, intervals[14])
	assert.Equal(t, JustInterval{numerator: 32, denominator: 17}, intervals[15])
	assert.Equal(t, JustInterval{numerator: 64, denominator: 33}, intervals[16])
	assert.Equal(t, JustInterval{numerator: 2, denominator: 1}, intervals[17])
}
