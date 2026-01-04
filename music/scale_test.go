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
	assert.Equal(t, Interval{Numerator: 1, Denominator: 1, Name: "Perfect Unison"}, intervals[0])
	assert.Equal(t, Interval{Numerator: 256, Denominator: 243, Name: "Pythagorean Minor Second"}, intervals[1])
	assert.Equal(t, Interval{Numerator: 9, Denominator: 8, Name: "Pythagorean (Greater) Major Second"}, intervals[2])
	assert.Equal(t, Interval{Numerator: 32, Denominator: 27, Name: "Pythagorean Minor Third"}, intervals[3])
	assert.Equal(t, Interval{Numerator: 81, Denominator: 64, Name: "Pythagorean Major Third"}, intervals[4])
	assert.Equal(t, Interval{Numerator: 4, Denominator: 3, Name: "Perfect Fourth"}, intervals[5])
	assert.Equal(t, Interval{Numerator: 1024, Denominator: 729, Name: "Pythagorean Diminished Fifth"}, intervals[6])
	assert.Equal(t, Interval{Numerator: 729, Denominator: 512, Name: "Pythagorean Augmented Fourth"}, intervals[7])
	assert.Equal(t, Interval{Numerator: 3, Denominator: 2, Name: "Perfect Fifth"}, intervals[8])
	assert.Equal(t, Interval{Numerator: 128, Denominator: 81, Name: "Pythagorean Minor Sixth"}, intervals[9])
	assert.Equal(t, Interval{Numerator: 27, Denominator: 16, Name: "Pythagorean Major Sixth"}, intervals[10])
	assert.Equal(t, Interval{Numerator: 16, Denominator: 9, Name: "Pythagorean (Lesser) Minor Seventh"}, intervals[11])
	assert.Equal(t, Interval{Numerator: 243, Denominator: 128, Name: "Pythagorean Major Seventh"}, intervals[12])
	assert.Equal(t, Interval{Numerator: 2, Denominator: 1, Name: "Perfect Octave"}, intervals[13])
}

func Test_ShouldReturn5LimitScaleWithExpectedScaleDegrees(t *testing.T) {
	// Given
	scale := New5LimitPythagoreanScale()

	// When
	intervals := scale.Intervals()

	// Then
	assert.Equal(t, "5-limit Pythagorean", scale.System())
	assert.Equal(t, "5-limit just intonation pure ratios chromatic scale derived from applying syntonic comma to Pythagorean ratios.", scale.Description())
	assert.Equal(t, 13, len(intervals))
	assert.Equal(t, Interval{Numerator: 1, Denominator: 1, Name: "Perfect Unison"}, intervals[0])
	assert.Equal(t, Interval{Numerator: 16, Denominator: 15, Name: "Minor Second"}, intervals[1])
	assert.Equal(t, Interval{Numerator: 10, Denominator: 9, Name: "Greater Major Second"}, intervals[2])
	assert.Equal(t, Interval{Numerator: 6, Denominator: 5, Name: "Minor Third"}, intervals[3])
	assert.Equal(t, Interval{Numerator: 5, Denominator: 4, Name: "Major Third"}, intervals[4])
	assert.Equal(t, Interval{Numerator: 4, Denominator: 3, Name: "Perfect Fourth"}, intervals[5])
	assert.Equal(t, Interval{Numerator: 45, Denominator: 32, Name: "Augmented Fourth (Tritone)"}, intervals[6])
	assert.Equal(t, Interval{Numerator: 3, Denominator: 2, Name: "Perfect Fifth"}, intervals[7])
	assert.Equal(t, Interval{Numerator: 8, Denominator: 5}, intervals[8])
	assert.Equal(t, Interval{Numerator: 5, Denominator: 3}, intervals[9])
	assert.Equal(t, Interval{Numerator: 9, Denominator: 5}, intervals[10])
	assert.Equal(t, Interval{Numerator: 15, Denominator: 8}, intervals[11])
	assert.Equal(t, Interval{Numerator: 2, Denominator: 1}, intervals[12])
}

//assert.Equal(t, "16:15", fretboard.Frets[0].Label)
//assert.Equal(t, 33.75, fretboard.Frets[0].Position)
//assert.Equal(t, "Minor Second", fretboard.Frets[0].Comment)
//assert.Equal(t, "16:15", fretboard.Frets[0].Interval)
//assert.Equal(t, "10:9", fretboard.Frets[1].Label) // not 9:8
//assert.Equal(t, 54.0, fretboard.Frets[1].Position)
//assert.Equal(t, "25:24", fretboard.Frets[1].Interval)
//assert.Equal(t, "Just (Lesser) Major Second", fretboard.Frets[1].Comment)
//assert.Equal(t, "6:5", fretboard.Frets[2].Label)
//assert.Equal(t, "27:25", fretboard.Frets[2].Interval)
//assert.Equal(t, "5:4", fretboard.Frets[3].Label)
//assert.Equal(t, "25:24", fretboard.Frets[3].Interval)
//assert.Equal(t, 108.0, fretboard.Frets[3].Position)
//assert.Equal(t, "4:3", fretboard.Frets[4].Label)
//assert.Equal(t, "16:15", fretboard.Frets[4].Interval)
//assert.Equal(t, 135.0, fretboard.Frets[4].Position)
//assert.Equal(t, "64:45", fretboard.Frets[5].Label)
//assert.Equal(t, "16:15", fretboard.Frets[5].Interval)
//assert.Equal(t, "45:32", fretboard.Frets[6].Label)
//assert.Equal(t, "3:2", fretboard.Frets[7].Label)
//assert.Equal(t, "16:15", fretboard.Frets[7].Interval)
//assert.Equal(t, 180.0, fretboard.Frets[7].Position)
//assert.Equal(t, "8:5", fretboard.Frets[8].Label)
//assert.Equal(t, "16:15", fretboard.Frets[8].Interval)
//assert.Equal(t, "5:3", fretboard.Frets[9].Label)
//assert.Equal(t, "25:24", fretboard.Frets[9].Interval)
//assert.Equal(t, 216.0, fretboard.Frets[9].Position)
//assert.Equal(t, "9:5", fretboard.Frets[10].Label) // not 16:9
//assert.Equal(t, "27:25", fretboard.Frets[10].Interval)
//assert.Equal(t, "Just (Greater) Minor Seventh", fretboard.Frets[10].Comment)
//assert.Equal(t, "15:8", fretboard.Frets[11].Label)
//assert.Equal(t, "25:24", fretboard.Frets[11].Interval)
//assert.Equal(t, 252.0, fretboard.Frets[11].Position)
//assert.Equal(t, "2:1", fretboard.Frets[12].Label)
//assert.Equal(t, "16:15", fretboard.Frets[12].Interval)
//assert.Equal(t, 270.0, fretboard.Frets[12].Position)

func Test_ShouldReturnSazScaleWithExpectedScaleDegrees(t *testing.T) {
	// Given
	scale := NewSazScale()

	// When
	intervals := scale.Intervals()

	// Then
	assert.Equal(t, "Saz", scale.System())
	assert.Equal(t, "Turkish Saz tuning ratios.", scale.Description())
	assert.Equal(t, 18, len(intervals))
	assert.Equal(t, Interval{Numerator: 1, Denominator: 1}, intervals[0])
	assert.Equal(t, Interval{Numerator: 18, Denominator: 17}, intervals[1])
	assert.Equal(t, Interval{Numerator: 12, Denominator: 11}, intervals[2])
	assert.Equal(t, Interval{Numerator: 9, Denominator: 8}, intervals[3])
	assert.Equal(t, Interval{Numerator: 81, Denominator: 68}, intervals[4])
	assert.Equal(t, Interval{Numerator: 27, Denominator: 22}, intervals[5])
	assert.Equal(t, Interval{Numerator: 81, Denominator: 64}, intervals[6])
	assert.Equal(t, Interval{Numerator: 4, Denominator: 3}, intervals[7])
	assert.Equal(t, Interval{Numerator: 24, Denominator: 17}, intervals[8])
	assert.Equal(t, Interval{Numerator: 16, Denominator: 11}, intervals[9])
	assert.Equal(t, Interval{Numerator: 3, Denominator: 2}, intervals[10])
	assert.Equal(t, Interval{Numerator: 27, Denominator: 17}, intervals[11])
	assert.Equal(t, Interval{Numerator: 18, Denominator: 11}, intervals[12])
	assert.Equal(t, Interval{Numerator: 27, Denominator: 16}, intervals[13])
	assert.Equal(t, Interval{Numerator: 16, Denominator: 9}, intervals[14])
	assert.Equal(t, Interval{Numerator: 32, Denominator: 17}, intervals[15])
	assert.Equal(t, Interval{Numerator: 64, Denominator: 33}, intervals[16])
	assert.Equal(t, Interval{Numerator: 2, Denominator: 1}, intervals[17])

}
